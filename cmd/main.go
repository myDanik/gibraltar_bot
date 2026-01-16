package main

import (
	"bot/internal/handlers"
	"bot/internal/services"
	"context"
	"log"
	"os"
	"time"

	"github.com/go-telegram/bot"
	"github.com/joho/godotenv"
)

const timerDuration = 12*time.Hour + 5*time.Minute

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.DefaultHandler),
	}
	b, err := bot.New(os.Getenv("TELEGRAM_API_KEY"), opts...)
	if err != nil {
		panic(err)
	}
	chatIDCache := services.NewCache()
	timerService := services.NewTimerService(chatIDCache, "chatIDList")
	cfgS := &services.ConfigService{APIUrl: "http://localhost:8080/configs"}
	cfgH := handlers.NewConfigHandler(cfgS, timerService)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/configs", bot.MatchTypeExact, cfgH.GetConfigsHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/update", bot.MatchTypeExact, cfgH.UpdateConfigs)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, cfgH.GetHelp)

	var Timer *time.Timer = time.NewTimer(timerDuration)

	go func(b *bot.Bot) {
		for {
			<-Timer.C
			cfgH.SendConfigByTimer(context.TODO(), b)
			Timer.Reset(timerDuration)
		}
	}(b)

	b.Start(context.TODO())

}
