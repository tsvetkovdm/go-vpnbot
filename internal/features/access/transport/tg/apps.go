package access_transport_tg

type appLink struct {
	title string
	url   string
}

type platform struct {
	title string
	note  string
	apps  []appLink
}

const (
	platformApple   = "apple"
	platformAndroid = "android"
	platformDesktop = "desktop"
)

func findPlatform(key string) (platform, bool) {
	switch key {
	case platformApple:
		return platform{
			title: "🍏 iPhone и Mac",
			note: "Рекомендуем Happ. v2RayTun - запасной вариант, если с Happ что-то не так.\n\n" +
				"Если основная ссылка на Happ не открывается, у вас, скорее всего, российский Apple ID - выберите кнопку с 🇷🇺.",
			apps: []appLink{
				{"🍏 Happ - iPhone и Mac", "https://apps.apple.com/app/id6504287215"},
				{"🇷🇺 Happ - iPhone (российский App Store)", "https://apps.apple.com/ru/app/happ-lite/id6799917773"},
				{"🇷🇺 Happ - Mac (российский App Store)", "https://apps.apple.com/ru/app/happ-proxy-utility/id6783623643"},
				{"v2RayTun - iPhone", "https://apps.apple.com/app/id6476628951"},
			},
		}, true

	case platformAndroid:
		return platform{
			title: "🤖 Android",
			note:  "Рекомендуем Happ. v2RayTun - запасной вариант, если с Happ что-то не так.",
			apps: []appLink{
				{"🤖 Happ - Google Play", "https://play.google.com/store/apps/details?id=com.happproxy"},
				{"v2RayTun - Google Play", "https://play.google.com/store/apps/details?id=com.v2raytun.android"},
			},
		}, true

	case platformDesktop:
		return platform{
			title: "💻 Windows и Linux",
			note: "Рекомендуем Happ: на странице загрузки выберите файл для своей системы " +
				"(для Windows - setup-Happ.x64.exe). v2rayN - запасной вариант.",
			apps: []appLink{
				{"💻 Happ - Windows и Linux", "https://github.com/Happ-proxy/happ-desktop/releases/latest"},
				{"v2rayN - Windows и Linux", "https://github.com/2dust/v2rayN/releases/latest"},
			},
		}, true
	}

	return platform{}, false
}
