package pack

import (
	"bytes"
	"errors"
	"github.com/Kingson4Wu/fast_proxy/common/cerror"
	"github.com/Kingson4Wu/fast_proxy/common/compress"
	"github.com/Kingson4Wu/fast_proxy/common/config"
	"github.com/Kingson4Wu/fast_proxy/common/pool"
	"github.com/Kingson4Wu/fast_proxy/common/proto/protobuf"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/common/sign"
	"github.com/Kingson4Wu/fast_proxy/outproxy/internal/encrypt"
	"github.com/Kingson4Wu/fast_proxy/outproxy/outconfig"
	"io"
	"net/http"
	"sync"

	"google.golang.org/protobuf/proto"
)

var (
	scPool = sync.Pool{
		New: func() interface{} {
			return new(config.ServiceConfig)
		},
	}

	pbPool = sync.Pool{
		New: func() interface{} {
			return new(protobuf.ProxyData)
		},
	}

	pbRespPool = sync.Pool{
		New: func() interface{} {
			return new(protobuf.ProxyRespData)
		},
	}
)

func newSc(serviceConfig *config.ServiceConfig) *config.ServiceConfig {
	sc := scPool.Get().(*config.ServiceConfig)
	sc.CompressEnable = serviceConfig.CompressEnable
	sc.EncryptEnable = serviceConfig.EncryptEnable
	sc.EncryptKeyName = serviceConfig.EncryptKeyName
	sc.SignEnable = serviceConfig.SignEnable
	sc.SignKeyName = serviceConfig.SignKeyName
	return sc
}

func EncodeReq(req *http.Request) ([]byte, *cerror.Err) {
	reqServiceName := server.Center().ClientName(req)

	if reqServiceName == "" {
		return nil, cerror.NewError(http.StatusBadRequest, "illegal args")
	}

	var bodyBytes []byte
	var err error
	b := pool.GetDataBufferChunk(req.ContentLength)
	defer pool.PutDataBufferChunk(b)
	if b == nil {
		bodyBytes, err = io.ReadAll(req.Body)
	} else {
		bf := bytes.NewBuffer(*b)
		bf.Reset()
		_, err = io.Copy(bf, req.Body)
		bodyBytes = bf.Bytes()
	}

	if err != nil {
		return nil, cerror.NewError(http.StatusBadRequest, "illegal body")
	}

	pData, err := Encode(bodyBytes, reqServiceName)
	if err != nil {
		return nil, cerror.NewError(http.StatusInternalServerError, err.Error())
	}

	return pData, nil
}

func DecodeResp(resp *http.Response) ([]byte, *cerror.Err) {

	defer resp.Body.Close()

	var bodyBytes []byte
	var err error
	b := pool.GetDataBufferChunk(resp.ContentLength)
	defer pool.PutDataBufferChunk(b)

	if b == nil {
		bodyBytes, err = io.ReadAll(resp.Body)
	} else {
		bf := bytes.NewBuffer(*b)
		bf.Reset()
		_, err = io.Copy(bf, resp.Body)
		bodyBytes = bf.Bytes()
	}

	if err != nil {
		return nil, cerror.NewError(http.StatusInternalServerError, "illegal body")
	}

	pData, err := Decode(bodyBytes)

	if err != nil {
		return nil, cerror.NewError(http.StatusInternalServerError, "decrypt failure")
	}

	return pData, nil
}

func DecodeFastResp(bodyBytes []byte) ([]byte, *cerror.Err) {
	pData, err := Decode(bodyBytes)

	if err != nil {
		return nil, cerror.NewError(http.StatusInternalServerError, "DecodeFastResp failure")
	}

	return pData, nil
}

//-------

type Middleware func([]byte, *config.ServiceConfig) ([]byte, error)

func ApplyMiddlewares(input []byte, serviceConfig *config.ServiceConfig, middlewares ...Middleware) ([]byte, error) {
	var err error
	output := input
	for _, middleware := range middlewares {
		output, err = middleware(output, serviceConfig)
		if err != nil {
			return nil, err
		}
	}
	return output, nil
}

func Encrypt(input []byte, serviceConfig *config.ServiceConfig) ([]byte, error) {
	if serviceConfig.EncryptEnable {
		return encrypt.EncodeReq(input, serviceConfig)
	}
	return input, nil
}

func Compress(input []byte, serviceConfig *config.ServiceConfig) ([]byte, error) {
	if serviceConfig.CompressEnable {
		return compress.Encode(input, serviceConfig.CompressAlgorithm)
	}
	return input, nil
}

func ProtobufEncode(input []byte, sc *config.ServiceConfig) ([]byte, error) {
	var bodySign string
	var err error
	if sc.SignEnable {
		bodySign, err = sign.GenerateBodySign(input, sc)
		if err != nil {
			server.GetLogger().Error("%s", err)
			return nil, errors.New("sign failure")
		}
	}

	stSend := pbPool.Get().(*protobuf.ProxyData)
	defer pbPool.Put(stSend)
	stSend.Sign = bodySign
	stSend.Compress = sc.CompressEnable
	stSend.Payload = input
	stSend.SignEnable = sc.SignEnable
	stSend.SignKeyName = sc.SignKeyName
	stSend.EncryptEnable = sc.EncryptEnable
	stSend.EncryptKeyName = sc.EncryptKeyName
	stSend.CompressAlgorithm = sc.CompressAlgorithm

	pData, err := proto.Marshal(stSend)
	if err != nil {
		return nil, errors.New("protobuf encode failure")
	}
	return pData, nil
}

type UnpackMiddleware func([]byte, *protobuf.ProxyRespData) ([]byte, error)

func ApplyUnpackMiddlewares(input []byte, proxyData *protobuf.ProxyRespData, middlewares ...UnpackMiddleware) ([]byte, error) {
	var err error
	output := input
	for _, middleware := range middlewares {
		output, err = middleware(output, proxyData)
		if err != nil {
			return nil, err
		}
	}
	return output, nil
}

func Decrypt(input []byte, proxyData *protobuf.ProxyRespData) ([]byte, error) {
	if proxyData.EncryptEnable {
		return encrypt.DecodeResp(input, proxyData.EncryptKeyName)
	}
	return input, nil
}

func Decompress(input []byte, proxyData *protobuf.ProxyRespData) ([]byte, error) {
	if proxyData.Compress {
		input, err := compress.Decode(input, proxyData.CompressAlgorithm)

		if err != nil {
			return nil, errors.New("compress decode failure")
		}
		return input, nil
	}
	return input, nil
}

//----------

func Encode(bodyBytes []byte, serviceName string) ([]byte, error) {
	serviceConfig := outconfig.Get().GetServiceConfig(serviceName)
	if serviceConfig == nil {
		return nil, errors.New("get serviceConfig failure")
	}

	/** Data copy, to ensure that the current data remains unchanged */
	sc := newSc(serviceConfig)
	defer scPool.Put(sc)

	middlewares := []Middleware{Encrypt, Compress, ProtobufEncode}

	result, err := ApplyMiddlewares(bodyBytes, sc, middlewares...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Decode(bodyBytes []byte) ([]byte, error) {

	reData := pbRespPool.Get().(*protobuf.ProxyRespData)
	defer pbRespPool.Put(reData)
	err := proto.Unmarshal(bodyBytes, reData)

	if err != nil {
		return nil, errors.New("protobuf decode failure")
	}

	bodyBytes = reData.Payload

	middlewares := []UnpackMiddleware{Decompress, Decrypt}

	result, err := ApplyUnpackMiddlewares(bodyBytes, reData, middlewares...)
	if err != nil {
		return nil, err
	}
	return result, nil
}
