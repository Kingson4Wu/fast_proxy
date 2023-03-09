package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/Kingson4Wu/fast_proxy/common/config"
	"github.com/Kingson4Wu/fast_proxy/common/server"
)

func getSignKey(serviceConfig *config.ServiceConfig, signName string) (string, error) {
	if signName != "" {
		signKey := server.Config().GetSignKeyByName(signName)
		if signKey == "" {
			return "", errors.New("sign key not exist")
		}
		return signKey, nil
	}
	signKey := server.Config().GetSignKey(serviceConfig)
	if signKey == "" {
		return "", errors.New("sign key not exist")
	}
	return signKey, nil
}

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

func GenerateBodySignWithName(body []byte, signName string) (string, error) {

	signKey, err := getSignKey(nil, signName)
	if err != nil {
		return "", err
	}
	return GenerateSign(body, signKey)
}

func GenerateBodySign(body []byte, serviceConfig *config.ServiceConfig) (string, error) {

	signKey, err := getSignKey(serviceConfig, "")
	if err != nil {
		return "", err
	}
	return GenerateSign(body, signKey)
}
