package lib

import "math/rand"

func GetRandomAlias(i int) string {
	runes := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	alias := make([]rune, i)
	for k := range alias {
		alias[k] = runes[rand.Intn(len(runes))]
	}
	return string(alias)
}
