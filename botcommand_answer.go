package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type AnswerCommand struct {
	BotCommandBase
}

func (h *AnswerCommand) Satisfies(context *MessageContext) bool {
	return strings.HasPrefix(context.Message.Content, "!answer")
}

func (h *AnswerCommand) Exec(context *MessageContext) {
	question := strings.Replace(context.Message.Content, "!answer ", "", -1)
	question = url.QueryEscape(question)

	apiCall := fmt.Sprintf("http://api.wolframalpha.com/v1/result?appid=%s&i=%s",
		context.BotRegistry.BotConfig.WolframAlphaKey,
		question)

	resp, _ := http.Get(apiCall)
	answerBytes, _ := ioutil.ReadAll(resp.Body)
	answer := string(answerBytes)

	if answer == "" {
		context.Session.ChannelMessageSend(context.Message.ChannelID,
			"there has been a problem with contacting Wolfram Alpha")
		return
	}
  if answer == "Wolfram|Alpha did not understand your input" {
    context.Session.ChannelMessageDelete(context.Message.ChannelID, context.Message.ID)
    return
  }

	context.Session.ChannelMessageSend(context.Message.ChannelID, answer)
}

func (h *AnswerCommand) Info() string {
	return `**!answer [searchterm]**
Search Wolfram Alpha for a specific search term`
}

func (h *AnswerCommand) CommandCategory() int {
	return COMMAND_CATEGORY_PRODUCTIVE
}