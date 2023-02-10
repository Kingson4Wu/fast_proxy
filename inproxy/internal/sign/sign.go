package sign

import (
	"errors"
	"github.com/Kingson4Wu/fast_proxy/common/sign"
	"github.com/Kingson4Wu/fast_proxy/inproxy/config"
)

func GenerateBodySign(body []byte, signName string) (string, error) {

	signKey := config.Get().GetSignKeyByName(signName)

	if signKey == "" {
		return "", errors.New("sign key not exist")
	}
	return sign.GenerateSign(body, signKey)
}
