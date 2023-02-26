package encrypt

import (
	"errors"
	"github.com/Kingson4Wu/fast_proxy/common/encrypt/aes"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"io"
)

func Encode(data []byte, key string) (result []byte, erro error) {

	defer func() {
		if err := recover(); err != nil {
			erro = errors.New("aes Encode panic")
			server.GetLogger().Errorf("Encode err:%s", err)
		}
	}()

	result, erro = aes.AesEncrypt(data, []byte(key))
	return
}

func Decode(data []byte, key string) (result []byte, erro error) {

	defer func() {
		if err := recover(); err != nil {
			erro = errors.New("aes Decode panic")
			server.GetLogger().Errorf("Decode err:%s", err)
		}
	}()

	result, erro = aes.AesDecrypt(data, []byte(key))
	return
}

func encrypt(dst io.Writer, src io.Reader) {
	//src.Read()
}
