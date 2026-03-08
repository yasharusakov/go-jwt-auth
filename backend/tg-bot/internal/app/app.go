package app

import (
	"context"
	"fmt"
	"log"
	"tg-bot/internal/client"
	"tg-bot/internal/config"
	"time"

	tele "gopkg.in/telebot.v4"
)

func Run() {
	cfg := config.GetConfig()

	pref := tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	userClient, err := client.NewGRPCUserClient(cfg.GrpcUserServiceInternalAddr)
	if err != nil {
		log.Fatal("Failed to create user client: " + err.Error())
		return
	}

	b.Handle("/users", func(c tele.Context) error {
		if c.Sender().ID != int64(cfg.AdminID) {
			return c.Send("You are not authorized to use this command.")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		count, err := userClient.GetUsersCount(ctx)
		if err != nil {
			return c.Send("Failed to get users count: " + err.Error())
		}

		msg := fmt.Sprintf("You have %d users.", count)

		return c.Send(msg)
	})

	b.Start()
}
