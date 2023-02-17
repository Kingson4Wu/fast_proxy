package encrypt

import (
	"errors"
	"github.com/Kingson4Wu/fast_proxy/common/aes"
	"github.com/Kingson4Wu/fast_proxy/common/logger"
	"github.com/Kingson4Wu/fast_proxy/inproxy/inconfig"
	"go.uber.org/zap"
	"io"
)

func EncodeResp(data []byte, encryptKeyName string) (result []byte, erro error) {
	encryptKey := inconfig.Get().GetEncryptKeyByName(encryptKeyName)

	if encryptKey == "" {
		return nil, errors.New("get encryptKey failure")
	}
	result, erro = Encode(data, encryptKey)
	return
}

func DecodeReq(data []byte, encryptKeyName string) (result []byte, erro error) {

	encryptKey := inconfig.Get().GetEncryptKeyByName(encryptKeyName)
	if encryptKey == "" {
		return nil, errors.New("encrypt key not exist")
	}
	result, erro = Decode(data, encryptKey)
	return
}

func Encode(data []byte, key string) (result []byte, erro error) {

	defer func() {
		if err := recover(); err != nil {
			erro = errors.New("aes Encode panic")
			logger.GetLogger().Error("", zap.Any("Encode err", err))
		}
	}()

	result, erro = aes.AesEncrypt(data, []byte(key))
	return
}

func Decode(data []byte, key string) (result []byte, erro error) {

	defer func() {
		if err := recover(); err != nil {
			erro = errors.New("aes Decode panic")
			logger.GetLogger().Error("", zap.Any("Decode err", err))
		}
	}()

	result, erro = aes.AesDecrypt(data, []byte(key))
	return
}

func encrypt(dst io.Writer, src io.Reader) {
	//src.Read()
}
