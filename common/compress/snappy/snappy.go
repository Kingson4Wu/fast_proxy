package snappy

import (
	"errors"
	"github.com/golang/snappy"
	"go.uber.org/zap"
)

type Snappy struct {
	Log *zap.Logger
}

func (s *Snappy) Encode(data []byte) (result []byte, erro error) {
	defer func() {
		if err := recover(); err != nil {
			erro = errors.New("snappy encode panic")
			s.Log.Error("", zap.Any("Encode err", err))
		}
	}()

	return snappy.Encode(nil, data), nil
}

func (s *Snappy) Decode(data []byte) (result []byte, erro error) {
	defer func() {
		if err := recover(); err != nil {
			erro = errors.New("snappy decode panic")
			s.Log.Error("", zap.Any("Decode err", err))
		}
	}()

	return snappy.Decode(nil, data)
}

func compressPanic(data []byte, needPanic bool) (result []byte, ok bool) {
	defer func() {
		if err := recover(); err != nil {
			ok = false
			//logger.GetLogger().Error("", zap.Any("err", err))
		}
	}()

	if needPanic {
		panic("miao")
	}

	return []byte("kxw"), true
}
