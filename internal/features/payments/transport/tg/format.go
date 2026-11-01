package payments_transport_tg

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shopspring/decimal"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_text "github.com/tsvetkovdm/go-vpnbot/internal/core/text"
	core_tg_html "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/html"
)

func (h *PaymentsTGHandler) formatPayment(p domain.Payment) core_tg_html.HTML {
	code := strconv.Itoa(p.CommentCode)
	summary := core_tg_html.Sprintf(
		"Тариф: %s\nСрок: %s\nСумма: %s",
		core_tg_html.Escape(p.PlanTitle),
		core_tg_html.Escape(formatDuration(p.Duration)),
		core_tg_html.Bold(formatAmount(p.Price)),
	)

	if p.Status == domain.PaymentStatusPending {
		return core_tg_html.Sprintf(
			"⏳ Проверяем ваш перевод\n\n%s\nКод: %s\n\n"+
				"Администратор проверит перевод и откроет доступ — вам придёт сообщение.\n"+
				"Если что-то пошло не так, напишите %s и назовите код.",
			summary,
			core_tg_html.Escape(code),
			core_tg_html.Escape(h.adminContact),
		)
	}

	return core_tg_html.Sprintf(
		"💳 Почти готово - осталось оплатить\n\n%s\n\n"+
			"Переведите сумму любым способом:\n"+
			"💳 Карта: %s\n"+
			"📱 По номеру телефона: %s\n"+
			"🏦 Банк: %s\n"+
			"👤 Получатель: %s\n\n"+
			"💬 Комментарий к переводу: %s\n"+
			"Обязательно укажите его - по нему администратор найдёт ваш перевод.\n\n"+
			"Нажмите на номер или комментарий, чтобы скопировать.\n"+
			"После перевода нажмите «✅\u00a0Я\u00a0перевёл».",
		summary,
		core_tg_html.Code(h.config.CardNumber),
		core_tg_html.Code(h.config.Phone),
		core_tg_html.Escape(h.config.Bank),
		core_tg_html.Escape(h.config.Recipient),
		core_tg_html.Code(code),
	)
}

func paymentMarkup(p domain.Payment) tgbotapi.InlineKeyboardMarkup {
	if p.Status != domain.PaymentStatusCreated {
		return tgbotapi.NewInlineKeyboardMarkup(accountRow())
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Я перевёл", paymentCallback("paid", p.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Отменить", paymentCallback("cancel", p.ID)),
		),
		accountRow(),
	)
}

func formatDuration(d domain.Duration) string {
	if d.Months > 0 {
		n := int(d.Months)
		return fmt.Sprintf("%d %s", n, core_text.PluralRu(n, "месяц", "месяца", "месяцев"))
	}

	n := int(d.Days)
	return fmt.Sprintf("%d %s", n, core_text.PluralRu(n, "день", "дня", "дней"))
}

func formatAmount(price decimal.Decimal) string {
	if price.IsInteger() {
		return price.String() + " ₽"
	}
	return price.StringFixed(2) + " ₽"
}

func formatUnresolvedNotice(p domain.Payment) string {
	choice := fmt.Sprintf("%s, %s — %s", p.PlanTitle, formatDuration(p.Duration), formatAmount(p.Price))

	if p.Status == domain.PaymentStatusPending {
		return "⏳ Проверяем ваш перевод:\n" + choice + "\n\n" +
			"Выбрать новый тариф можно будет после проверки."
	}

	return "⚠️ Вы уже выбрали тариф, но ещё не оплатили:\n" + choice + "\n\n" +
		"Продолжите или отмените, чтобы выбрать другой тариф."
}

func unresolvedNoticeMarkup(p domain.Payment) tgbotapi.InlineKeyboardMarkup {
	rows := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 Продолжить оплату", "pay:current"),
		),
	}

	if p.Status == domain.PaymentStatusCreated {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Отменить", paymentCallback("cancel", p.ID)),
		))
	}

	rows = append(rows, accountRow())

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func noPaymentMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💳 Тарифы", "plans:public"),
		),
		accountRow(),
	)
}

func accountRow() []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("👤 Личный кабинет", "users:home"),
	)
}
