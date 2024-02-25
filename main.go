package main

import (
	"time"
	"web-scraper/chat_gpt"
)

func main() {
	go func() {
		chat_gpt.Run("", "")
	}()
	time.Sleep(time.Minute)
	chat_gpt.GetResponse(
		chat_gpt.GetModelContextForGivenType(chat_gpt.RandomShortNoFollowUpUrl), "what is the weather in toronto",
	)
	chat_gpt.GetResponse(
		chat_gpt.GetModelContextForGivenType(chat_gpt.NewsClassifierUrl), "reliance gonna pump tomorrow fosho",
	)
	chat_gpt.GetResponse(
		chat_gpt.GetModelContextForGivenType(chat_gpt.RandomShortNoFollowUpUrl), "what is 2+2",
	)
	chat_gpt.GetResponse(
		chat_gpt.GetModelContextForGivenType(chat_gpt.RandomShortNoFollowUpUrl), "write hello world in c++",
	)
	chat_gpt.GetResponse(
		chat_gpt.GetModelContextForGivenType(chat_gpt.NewsClassifierUrl),
		"🇺🇸CONGRESS ON WINTER BREAK\n\nSnow geese take flight...\n\nSource: ABC",
	)
	select {}
}
