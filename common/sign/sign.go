package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

func GenerateSign(body []byte, signKey string) (string, error) {

	if signKey == "" {
		return "", errors.New("sign key not exist")
	}

	key := []byte(signKey)
	m := hmac.New(sha256.New, key)
	m.Write(body)
	bodySignature := hex.EncodeToString(m.Sum(nil))

	return bodySignature, nil
}
