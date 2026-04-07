package utils

import "strings"

func EncodeBase62(num int) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	if num == 0 {
		return string(charset[0])
	}

	var result strings.Builder
	for num > 0 {
		result.WriteByte(charset[num%62])
		num = num / 62
	}
	return reverse(result.String())
}

func reverse(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
