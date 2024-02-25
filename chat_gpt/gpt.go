package chat_gpt

import (
	"context"
	"fmt"
	cu "github.com/Davincible/chromedp-undetected"
	"github.com/go-rod/stealth"
	"log"
	"sync"
	"time"

	_ "github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

const (
	RandomShortNoFollowUpUrl = "https://chat.openai.com/c/d2bdaf60-5c6d-4681-8436-863f08e07b0f"
	NewsClassifierUrl        = "https://chat.openai.com/c/21525c55-a858-4086-9016-b2bceaf0d07d"
)

var (
	tabContexts = make(map[string]context.Context)
)

func createNewUndetectedContext(headless bool) context.Context {
	var ctx context.Context
	var err error
	if headless {
		ctx, _, err = cu.New(
			cu.NewConfig(
				cu.WithHeadless(),
			),
		)
	} else {
		ctx, _, err = cu.New(
			cu.NewConfig(),
		)
	}

	if err != nil {
		panic(err)
	}
	return ctx
}

func createNewUndetectedContextFromParentCtx(ctx context.Context, chatUrl string) context.Context {
	childCtx, _ := chromedp.NewContext(ctx)
	tabContexts[chatUrl] = childCtx
	return childCtx
}

func GetModelContextForGivenType(chatUrl string) context.Context {
	if _, ok := tabContexts[chatUrl]; !ok {
		panic("Tab context not found")
	}
	return tabContexts[chatUrl]
}

func openChatGptLoginPage(ctx context.Context) {
	if err := chromedp.Run(
		ctx,
		chromedp.Evaluate(stealth.JS, nil),
		chromedp.Navigate("https://chat.openai.com/auth/login"),
		chromedp.Sleep(2*time.Second),
	); err != nil {
		fmt.Println("Error navigating to login page:", err)
		return
	}
	log.Print("Chat-gpt login page opened.")
	if err := chromedp.Run(
		ctx,
		chromedp.Click("button", chromedp.ByQuery),
		chromedp.Sleep(10*time.Second),
	); err != nil {
		fmt.Println("Error clicking login button:", err)
		return
	}
	log.Print("Chat-gpt login button clicked.")
}

func handleTypingEmailAndSubmitForVarient1(ctx context.Context, email string) {
	if err := chromedp.Run(ctx, chromedp.Click(`div.c09f465e8.c2ee5f981.text.c64c67016.c4596eb33`)); err != nil {
		panic(err)
	}
	if err := chromedp.Run(
		ctx, chromedp.SendKeys("input.cc0c4f1dc.c15825b59", email),
	); err != nil {
		panic(err)
	}
	log.Print("email entered in variant 1")
	if err := chromedp.Run(
		ctx, chromedp.Click(`//button[contains(@class,'_button-login-id')]`),
		chromedp.Sleep(time.Second*3),
	); err != nil {
		fmt.Println("Error clicking submit button for email:", err)
		return
	}
	log.Print("email submitted in variant 1")
}

func handleTypingEmailAndSubmitForVarient2(ctx context.Context, email string) {
	fmt.Println("Input field does not exist on the page, trying fallback method ...")
	if err := chromedp.Run(
		ctx,
		chromedp.Click(`#email-input`, chromedp.ByID),
		chromedp.SetValue(`#email-input`, email, chromedp.ByID),
	); err != nil {
		fmt.Println("Error typing email:", err)
		return
	}
	log.Print("email entered in variant 2")
	if err := chromedp.Run(ctx, chromedp.Click(`button.continue-btn`), chromedp.Sleep(time.Second*3)); err != nil {
		panic(err)
	}
	log.Print("email submitted in variant 2")
}

func handleTypingEmailAndSubmit(ctx context.Context, email string) {
	var exists bool
	err := chromedp.Run(
		ctx, chromedp.EvaluateAsDevTools(`document.querySelector("input.cc0c4f1dc.c15825b59") !== null`, &exists),
	)
	if err != nil {
		panic(err)
	}
	if exists {
		log.Print("Chat-gpt email login variant 1 found")
		handleTypingEmailAndSubmitForVarient1(ctx, email)
	} else {
		log.Print("Chat-gpt email login variant 2 found")
		handleTypingEmailAndSubmitForVarient2(ctx, email)
		handleTypingEmailAndSubmitForVarient1(ctx, email)
	}
}

func handleTypingPasswordAndSubmit(ctx context.Context, password string) {
	log.Print("typing password for chatgpt")
	if err := chromedp.Run(
		ctx,
		chromedp.SendKeys(`input[type="password"]`, password, chromedp.ByQuery),
	); err != nil {
		fmt.Println("Error typing password:", err)
		return
	}
	if err := chromedp.Run(
		ctx,
		chromedp.Click(`button[name="action"][value="default"]`, chromedp.ByQuery), // Click the button
		chromedp.Sleep(10*time.Second),
	); err != nil {
		fmt.Println("Error clicking login with email/password button:", err)
		return
	}
	log.Print("password submitted.")
}

func openAllChatsInDifferentTabs(ctx context.Context) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		if err := chromedp.Run(
			createNewUndetectedContextFromParentCtx(ctx, NewsClassifierUrl),
			chromedp.Navigate(NewsClassifierUrl),
			chromedp.Sleep(15*time.Second),
		); err != nil {
			fmt.Println("Error navigating to relevant chat:", err)
			return
		}
		wg.Done()
		log.Print("NewsClassifier tab opened...")
	}()
	go func() {
		if err := chromedp.Run(
			createNewUndetectedContextFromParentCtx(ctx, RandomShortNoFollowUpUrl),
			chromedp.Navigate(RandomShortNoFollowUpUrl),
			chromedp.Sleep(15*time.Second),
		); err != nil {
			fmt.Println("Error navigating to relevant chat:", err)
			return
		}
		wg.Done()
		log.Print("RandomShortNoFollowUp tab opened...")
	}()
	wg.Wait()
}

func getLatestResponseFromChat(ctx context.Context) string {
	var text string
	err := chromedp.Run(
		ctx, chromedp.Evaluate(
			`
        var div = document.querySelector('.flex.flex-col.pb-9.text-sm > div:last-child > div:last-child > div:last-child > div:last-child > div:last-child > div:first-child');
        if (div) {
            div.textContent;
        }
    `, &text,
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	return text
}

func GetResponse(ctx context.Context, message string) string {
	log.Println("generating reponse for: " + message)
	if err := chromedp.Run(
		ctx,
		chromedp.WaitVisible(`#prompt-textarea`, chromedp.ByID),
		chromedp.SetValue(`#prompt-textarea`, message, chromedp.ByID),
		chromedp.Click("[data-testid='send-button']"),
		chromedp.Sleep(5*time.Second),
	); err != nil {
		log.Println("panic while generating response")
		panic("Error finding text area or sending message:" + err.Error())
	}
	return getLatestResponseFromChat(ctx)
}

func Run(email string, password string, headless bool) {
	ctx := createNewUndetectedContext(headless)
	openChatGptLoginPage(ctx)
	handleTypingEmailAndSubmit(ctx, email)
	handleTypingPasswordAndSubmit(ctx, password)
	openAllChatsInDifferentTabs(ctx)
}
