package access_transport_tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	core_tg_callback "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/callback"
	core_tg_html "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/html"
)

func connectText(link string, profileTitle string) core_tg_html.HTML {
	return core_tg_html.Sprintf(
		"🔑 Подключение\n\n"+
			"Ваша ссылка для приложения (нажмите, чтобы скопировать):\n"+
			"%s\n\n"+
			"Как подключиться:\n"+
			"1. Установите Happ - выберите своё устройство ниже\n"+
			"2. Скопируйте ссылку и добавьте её в приложении: «+» → «Вставить\u00a0из\u00a0буфера\u00a0обмена»\n"+
			"3. Выберите сервер и подключитесь\n\n"+
			"Если что-то перестало работать - обновите список серверов в самом приложении: "+
			"значок с двумя круговыми стрелками рядом с «%s». Заново ничего вводить не нужно.",
		core_tg_html.Code(link),
		core_tg_html.Escape(profileTitle),
	)
}

func connectMarkup() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🍏 iPhone и Mac", core_tg_callback.Build("access", "apps", platformApple)),
			tgbotapi.NewInlineKeyboardButtonData("🤖 Android", core_tg_callback.Build("access", "apps", platformAndroid)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💻 Windows и Linux", core_tg_callback.Build("access", "apps", platformDesktop)),
		),
		accountRow(),
	)
}

func appsText(p platform) string {
	return p.title + "\n\n" + p.note
}

func appsMarkup(p platform) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(p.apps)+2)
	for _, app := range p.apps {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(app.title, app.url),
		))
	}

	rows = append(rows, backRow(), accountRow())

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func backRow() []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Назад", "access:connect"),
	)
}

func accountRow() []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("👤 Личный кабинет", "users:home"),
	)
}
