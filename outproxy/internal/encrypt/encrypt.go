package encrypt

import (
	"errors"
	"github.com/Kingson4Wu/fast_proxy/common/config"
	"github.com/Kingson4Wu/fast_proxy/common/encrypt"
	"github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
)

func EncodeReq(data []byte, serviceConfig *config.ServiceConfig) (result []byte, erro error) {
	encryptKey := outconfig.Get().GetEncryptKey(serviceConfig)

	if encryptKey == "" {
		return nil, errors.New("get encryptKey failure")
	}
	result, erro = encrypt.Encode(data, encryptKey)
	return
}

func DecodeResp(data []byte, encryptKeyName string) (result []byte, erro error) {

	encryptKey := outconfig.Get().GetEncryptKeyByName(encryptKeyName)
	if encryptKey == "" {
		return nil, errors.New("encrypt key not exist")
	}
	result, erro = encrypt.Decode(data, encryptKey)
	return
}
