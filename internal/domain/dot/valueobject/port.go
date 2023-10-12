package valueobject

import (
	"github.com/dddplayer/dp/internal/domain/dot"
	"hash/fnv"
	"strings"
)

func PortStr(name string) string {
	ps := name
	for _, old := range []string{".", "-", "/"} {
		ps = strings.ReplaceAll(ps, old, dot.Joiner)
	}
	return ps
	//return GenerateShortURL(ps)
}

func GenerateShortURL(originalURL string) string {
	h := fnv.New32a()
	_, _ = h.Write([]byte(originalURL))
	hash := h.Sum32()

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortURL := ""
	for hash > 0 {
		index := int(hash % 62)
		shortURL = charset[index:index+1] + shortURL
		hash /= 62
	}

	return "d" + shortURL
}
