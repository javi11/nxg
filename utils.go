package nxg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	mathRand "math/rand"
	"strconv"
	"time"
)

func encrypt(plaintext string, key string) (string, error) {
	key32Byte := make([]byte, 32)
	copy(key32Byte[:], []byte(key))

	c, err := aes.NewCipher(key32Byte)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	return hex.EncodeToString(gcm.Seal(nonce, nonce, []byte(plaintext), nil)), nil
}

func randomString(length uint32, salt uint64, num bool) string {
	var charset string

	if num {
		charset = "0123456789abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	} else {
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	source, _ := strconv.ParseInt((fmt.Sprintf("%v%v", salt, time.Now().UnixNano())), 10, 64)

	seededRand := mathRand.New(mathRand.NewSource(source))
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset)-1)]
	}

	return string(b)
}

func getSHA256Hash(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))

	return hex.EncodeToString(hasher.Sum(nil))
}
