package core_tg_html

import (
	"fmt"
	"html"
)

type HTML string

func Escape(s string) HTML {
	return HTML(html.EscapeString(s))
}

func Code(s string) HTML {
	return "<code>" + Escape(s) + "</code>"
}

func Bold(s string) HTML {
	return "<b>" + Escape(s) + "</b>"
}

func Sprintf(format string, args ...HTML) HTML {
	values := make([]any, len(args))
	for i, arg := range args {
		values[i] = string(arg)
	}

	return HTML(fmt.Sprintf(format, values...))
}
