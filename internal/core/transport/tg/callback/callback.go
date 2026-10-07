package core_tg_callback

import "strings"

func Build(segments ...string) string {
	return strings.Join(segments, ":")
}

func Segments(data string) []string {
	if data == "" {
		return nil
	}
	return strings.Split(data, ":")
}
