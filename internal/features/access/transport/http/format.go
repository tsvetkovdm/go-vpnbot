package access_transport_http

import "net/url"

const maxTokenLength = 64

func isValidToken(token string) bool {
	if token == "" || len(token) > maxTokenLength {
		return false
	}

	for _, r := range token {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
		default:
			return false
		}
	}

	return true
}

func maskToken(token string) string {
	const visible = 4
	if len(token) <= visible {
		return "***"
	}

	return token[:visible] + "***"
}

func noAccessTitle(adminContact string) string {
	return "⛔ Нет доступа — продлите в боте или напишите " + adminContact
}

func placeholderLink(title string) string {
	return "vless://00000000-0000-0000-0000-000000000000@127.0.0.1:1?encryption=none&type=tcp#" +
		url.PathEscape(title)
}
