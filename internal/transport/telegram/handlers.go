package telegram

import (
	"context"

	"github.com/ARKremlin/go_dnd/internal/domain"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type PlayerService interface {
	GetOrCreatePlayer(ctx context.Context, telegramID int64, telegramUsername string) (*domain.User, error)
}

func (tb *Bot) handleStart(ctx context.Context, _ *bot.Bot, update *models.Update) {
	msg := update.Message
	if msg == nil || msg.From == nil {
		return
	}

	u, err := tb.players.GetOrCreatePlayer(ctx, msg.From.ID, msg.From.Username)
	if err != nil {
		tb.log.Error("failed to register player", "err", err, "telegram id", msg.From.ID)
		tb.reply(ctx, msg.Chat.ID, "Что-то пошло не так. Попробуйте повторить запрос позже.")
		return
	}

	text := "Вы успешно зарегистрировались. Сообщите своё Имя пользователя telegram мастеру, чтобы он смог пригласить Вас за стол."
	if u.TelegramUsername == nil {
		text += "\n\nЧтобы иметь возможность играть, создайте его в настройках и нажмите /start ещё раз."
	}
	tb.reply(ctx, msg.Chat.ID, text)
}

func (tb *Bot) reply(ctx context.Context, chatID int64, text string) {
	_, err := tb.api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text})
	if err != nil {
		tb.log.Error("failed to send message", "err", err, "chat id", chatID)
	}

}
