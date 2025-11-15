package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	maxbot "github.com/max-messenger/max-bot-api-client-go"

	"github.com/escalopa/inno-vkode/internal/adapters/backend/httpclient"
	maxadapter "github.com/escalopa/inno-vkode/internal/adapters/messenger/max"
	"github.com/escalopa/inno-vkode/internal/adapters/notifier/email"
	"github.com/escalopa/inno-vkode/internal/app/bot"
	"github.com/escalopa/inno-vkode/internal/config"
	"github.com/escalopa/inno-vkode/internal/logger"
	"github.com/escalopa/inno-vkode/internal/state"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	log := logger.New(cfg.LogLevel)

	api, err := maxbot.New(cfg.MaxBotToken)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to init MAX bot API")
	}

	info, err := api.Bots.GetBot(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", info)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		osSig := make(chan os.Signal, 1)
		signal.Notify(osSig, syscall.SIGINT, syscall.SIGTERM)
		<-osSig
		log.Info().Msg("shutdown signal received")
		cancel()
	}()

	backend := httpclient.New(cfg.BackendBaseURL, cfg.HTTPTimeout, log)
	messenger := maxadapter.New(api, log)
	emailSender := email.NewLogSender(log)
	store := state.NewMemoryStore(time.Now)

	service := bot.New(cfg, log, backend, messenger, emailSender, store)

	if err := service.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Error().Err(err).Msg("service stopped with error")
	}
	log.Info().Msg("bot service stopped")
}
