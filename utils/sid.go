package utils

import "crypto/rand"

const lexicon = "23456789abcdefghjkmnpqrstuvwxyz"

func Sid(size int) string {
	sid := make([]byte, size)

	_, err := rand.Read(sid)
	if err != nil {
		panic(err)
	}

	for i := range size {
		index := sid[i] % byte(len(lexicon))
		sid[i] = lexicon[index]
	}

	return string(sid)
}
