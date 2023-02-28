package encrypt

import (
	"errors"
	"github.com/Kingson4Wu/fast_proxy/common/encrypt/aes"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"io"
)

func Encode(data []byte, key string) (result []byte, error error) {

	defer func() {
		if err := recover(); err != nil {
			error = errors.New("aes Encode panic")
			server.GetLogger().Errorf("Encode err:%s", err)
		}
	}()

	result, error = aes.Encrypt(data, []byte(key))
	return
}

func Decode(data []byte, key string) (result []byte, error error) {

	defer func() {
		if err := recover(); err != nil {
			error = errors.New("aes Decode panic")
			server.GetLogger().Errorf("Decode err:%s", err)
		}
	}()

	result, error = aes.Decrypt(data, []byte(key))
	return
}

func encrypt(dst io.Writer, src io.Reader) {
	//src.Read()
}
