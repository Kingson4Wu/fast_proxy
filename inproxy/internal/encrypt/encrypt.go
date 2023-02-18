package encrypt

import (
	"errors"
	"github.com/Kingson4Wu/fast_proxy/common/encrypt"
	"github.com/Kingson4Wu/fast_proxy/inproxy/inconfig"
)

func EncodeResp(data []byte, encryptKeyName string) (result []byte, erro error) {
	encryptKey := inconfig.Get().GetEncryptKeyByName(encryptKeyName)

	if encryptKey == "" {
		return nil, errors.New("get encryptKey failure")
	}
	result, erro = encrypt.Encode(data, encryptKey)
	return
}

func DecodeReq(data []byte, encryptKeyName string) (result []byte, erro error) {

	encryptKey := inconfig.Get().GetEncryptKeyByName(encryptKeyName)
	if encryptKey == "" {
		return nil, errors.New("encrypt key not exist")
	}
	result, erro = encrypt.Decode(data, encryptKey)
	return
}
