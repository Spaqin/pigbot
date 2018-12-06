package main

import "github.com/bwmarrin/discordgo"

type BotRegistry struct {
	ThreadHandler *ThreadHandler
	BotConfig     *BotConfig
}

type MessageHandler struct {
	Receivers   []MessageReceiver
	BotRegistry *BotRegistry
}

type MessageContext struct {
	Session     *discordgo.Session
	Message     *discordgo.MessageCreate
	BotRegistry *BotRegistry
}

type MessageReceiver interface {
	Satisfies(*MessageContext) bool
	Exec(*MessageContext)
}

func CreateMessageHandler() *MessageHandler {
	h := MessageHandler{}
	h.Receivers = make([]MessageReceiver, 0)
	return &h
}

func (m *MessageHandler) RegisterReceiver(handler MessageReceiver) {
	m.Receivers = append(m.Receivers, handler)
}

func (m *MessageHandler) OnMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	context := &MessageContext{Session: session, Message: message, BotRegistry: m.BotRegistry}
	for _, handler := range m.Receivers {
		if handler.Satisfies(context) {
			go handler.Exec(context)
		}
	}

}