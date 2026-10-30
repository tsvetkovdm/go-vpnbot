package core_text

func PluralRu(n int, one, few, many string) string {
	if n < 0 {
		n = -n
	}
	n = n % 100
	if n >= 11 && n <= 14 { // 11–14 — исключение: «11 месяцев», не «11 месяц»
		return many
	}
	switch n % 10 {
	case 1:
		return one
	case 2, 3, 4:
		return few
	default:
		return many
	}
}
