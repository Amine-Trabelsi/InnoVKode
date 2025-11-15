package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func main() {
	api, err := maxbot.New(os.Getenv("TOKEN"))
	if err != nil {
		panic(err)
	}
	// Some methods demo:
	info, err := api.Bots.GetBot(context.Background())
	fmt.Printf("Get me: %#v %#v", info, err)

	ctx, cancel := context.WithCancel(context.Background()) // создам
	go func() {
		exit := make(chan os.Signal)
		signal.Notify(exit, os.Kill, os.Interrupt)
		<-exit
		cancel()
	}()

	for upd := range api.GetUpdates(ctx) { // Чтение из канала с обновлениями
		switch upd := upd.(type) { // Определение типа пришедшего обновления
		case *schemes.MessageCreatedUpdate:
			// Отправка сообщения
			_, err := api.Messages.Send(ctx, maxbot.NewMessage().SetChat(upd.Message.Recipient.ChatId).SetText("Привет! ✨"))
			if err != nil && err.Error() != "" {
				panic(fmt.Sprintf("HERE2: %+v", err))
			}
		}
	}
}
