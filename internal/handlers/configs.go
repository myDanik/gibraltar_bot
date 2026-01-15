package handlers

import (
	"bot/internal/services"
	"bot/internal/shared"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type CfgHandler struct {
	timerService  *services.TimerService
	configService *services.ConfigService
}

func NewConfigHandler(configService *services.ConfigService, timerService *services.TimerService) *CfgHandler {
	return &CfgHandler{
		configService: configService,
		timerService:  timerService,
	}
}

func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   shared.StartMessage,
	})
	if err != nil {
		log.Println(err)
	}
}

func (h *CfgHandler) GetConfigsHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.timerService.AddNewChatToTimer(update.Message.Chat.ID)
	configs, err := h.configService.GetConfigs()
	if err != nil {
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   shared.DataRecievingError,
		})
		if err != nil {
			log.Println(err)
		}
		return
	}
	if len(configs) == 0 {
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   shared.EmptyListError,
		})
		if err != nil {
			log.Println(err)
		}
		return
	}
	messages := make([]string, 1)
	if len(configs) > 4080 {
		splittedString := strings.Split(configs, "\n")
		mesIdx := 0
		for _, v := range splittedString {
			if len(messages[mesIdx]+v) < 4080 {
				messages[mesIdx] += v + "\n"
				continue
			}
			mesIdx++
			messages = append(messages, v+"\n")
		}

	} else {
		messages = append(messages, configs)
	}
	for _, msg := range messages {
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:    update.Message.Chat.ID,
			Text:      "```\n" + msg + "```",
			ParseMode: models.ParseModeMarkdown,
		})
		if err != nil {
			log.Println(err)
		}
		time.Sleep(100 * time.Millisecond)
	}

}

func (h *CfgHandler) UpdateConfigs(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.timerService.AddNewChatToTimer(update.Message.Chat.ID)
	err := h.configService.UpdateConfigs()
	if err != nil {
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   shared.DataRecievingError,
		})
		if err != nil {
			log.Println(err)
		}
		return
	}
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   shared.UpdateStarted,
	})
	if err != nil {
		log.Println(err)
	}
}

func (h *CfgHandler) GetHelp(ctx context.Context, b *bot.Bot, update *models.Update) {
	h.timerService.AddNewChatToTimer(update.Message.Chat.ID)
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      shared.HelpMessage,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		log.Println(err)
	}
}

func (h *CfgHandler) SendConfigByTimer(ctx context.Context, b *bot.Bot) {
	fmt.Println("timer is working")
	chatIDList := h.timerService.Cache.GetKeys()
	configs, err := h.configService.GetConfigs()
	if err != nil {
		log.Println(err)
		return
	}
	if len(configs) == 0 {
		return
	}
	messages := make([]string, 1)
	if len(configs) > 4080 {
		splittedString := strings.Split(configs, "\n")
		mesIdx := 0
		for _, v := range splittedString {
			if len(messages[mesIdx]+v) < 4080 {
				messages[mesIdx] += v + "\n"
				continue
			}
			mesIdx++
			messages = append(messages, v+"\n")
		}

	} else {
		messages = append(messages, configs)
	}
	for _, chatID := range chatIDList {
		for _, msg := range messages {

			_, err = b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:    chatID,
				Text:      "```\n" + msg + "```",
				ParseMode: models.ParseModeMarkdown,
			})
			if err != nil {
				log.Println(err)
			}
			time.Sleep(100 * time.Millisecond)
		}

	}

}
