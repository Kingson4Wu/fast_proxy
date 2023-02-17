package sign

import (
	"errors"
	"github.com/Kingson4Wu/fast_proxy/common/config"
	"github.com/Kingson4Wu/fast_proxy/common/sign"
	"github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
)

func GenerateBodySign(body []byte, serviceConfig *config.ServiceConfig) (string, error) {

	signKey := outconfig.Get().GetSignKey(serviceConfig)

	if signKey == "" {
		return "", errors.New("sign key not exist")
	}
	return sign.GenerateSign(body, signKey)
}
