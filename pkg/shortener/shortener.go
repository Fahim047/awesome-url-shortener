package shortener

import (
	"crypto/rand"
	"math/big"
)

const (
	ALPHABETS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	KEY_LENGTH = 6
)

func GenerateShortKey() (string, error) {
    b := make([]byte, KEY_LENGTH)
    for i := range b {
        n, err := rand.Int(rand.Reader, big.NewInt(int64(len(ALPHABETS))))
        if err != nil {
            return "", err
        }
        b[i] = ALPHABETS[n.Int64()]
    }
    return string(b), nil
}