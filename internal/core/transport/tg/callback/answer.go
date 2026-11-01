package core_tg_callback

import "context"

type answerKey struct{}

func WithAnswer(ctx context.Context) (context.Context, func() string) {
	text := new(string)

	return context.WithValue(ctx, answerKey{}, text), func() string {
		return *text
	}
}

func SetAnswer(ctx context.Context, text string) {
	if holder, ok := ctx.Value(answerKey{}).(*string); ok {
		*holder = text
	}
}
