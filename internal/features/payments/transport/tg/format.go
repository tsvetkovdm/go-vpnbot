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
			"⏳ Счёт ожидает подтверждения\n\n%s\nКод: %s\n\n"+
				"Администратор проверит перевод и откроет доступ — вам придёт сообщение.\n"+
				"Если что-то пошло не так, напишите %s и назовите код.",
			summary,
			core_tg_html.Escape(code),
			core_tg_html.Escape(h.adminContact),
		)
	}

	return core_tg_html.Sprintf(
		"🧾 Счёт на оплату\n\n%s\n\n"+
			"Переведите сумму любым способом:\n"+
			"💳 Карта: %s\n"+
			"📱 СБП: %s\n"+
			"🏦 Банк: %s\n"+
			"👤 Получатель: %s\n\n"+
			"💬 Комментарий к переводу: %s\n"+
			"Обязательно укажите его - по нему администратор найдёт ваш платёж.\n\n"+
			"Нажмите на номер или комментарий, чтобы скопировать.\n"+
			"После перевода нажмите \n«✅ Я перевёл».",
		summary,
		core_tg_html.Code(h.config.CardNumber),
		core_tg_html.Code(h.config.Phone),
		core_tg_html.Escape(h.config.Bank),
		core_tg_html.Escape(h.config.Recipient),
		core_tg_html.Code(code),
	)
}

func paymentMarkup(p domain.Payment) tgbotapi.InlineKeyboardMarkup {
	homeRow := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("👤 Личный кабинет", "users:home"),
	)

	if p.Status != domain.PaymentStatusCreated {
		return tgbotapi.NewInlineKeyboardMarkup(homeRow)
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Я перевёл", paymentCallback("paid", p.ID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❌ Отменить", paymentCallback("cancel", p.ID)),
		),
		homeRow,
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
