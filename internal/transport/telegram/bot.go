package telegram

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Bot struct {
	api     *bot.Bot
	log     *slog.Logger
	players PlayerService
}

func New(token string, log *slog.Logger, players PlayerService) (*Bot, error) {
	tb := &Bot{log: log.With("component", "telegram"), players: players}

	api, err := bot.New(token,
		bot.WithDefaultHandler(tb.handleDefault),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, tb.handleStart),
		bot.WithErrorsHandler(tb.handleError),
	)
	if err != nil {
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}
	tb.api = api

	return tb, nil
}

func (tb *Bot) Start(ctx context.Context) {
	tb.log.Info("telegram bot started")
	tb.api.Start(ctx)
	tb.log.Info("telegram bot stopped")
}

func (tb *Bot) handleDefault(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}
	tb.log.Info("telegram message received",
		"update ID", update.ID,
		"chat ID", update.Message.Chat.ID,
	)
}

func (tb *Bot) handleError(err error) {
	tb.log.Error("telegram error", "err", err)
}
