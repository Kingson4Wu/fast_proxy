package encrypt

import (
	"errors"
	"github.com/Kingson4Wu/fast_proxy/common/encrypt/aes"
	"github.com/Kingson4Wu/fast_proxy/common/logger"
	"go.uber.org/zap"
	"io"
)

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
