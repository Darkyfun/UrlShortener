// Package aliasname предоставляет функции для генерации псевдонимов.
package aliasname

import "math/rand"

// GetRandomAlias возвращает набор случайных символов, который будет являться псевдонимом для URL.
func GetRandomAlias(i int) string {
	runes := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	alias := make([]rune, i)
	for k := range alias {
		alias[k] = runes[rand.Intn(len(runes))]
	}

	return string(alias)
}
