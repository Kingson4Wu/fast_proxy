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

	reqServiceName := server.Center().ClientName(resp.Request)

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

	pData, err := Decode(bodyBytes, reqServiceName)

	if err != nil {
		return nil, cerror.NewError(http.StatusInternalServerError, "decrypt failure")
	}

	return pData, nil
}

func Encode(bodyBytes []byte, serviceName string) ([]byte, error) {

	serviceConfig := outconfig.Get().GetServiceConfig(serviceName)

	if serviceConfig == nil {
		return nil, errors.New("get serviceConfig failure")
	}

	/** 数据拷贝, 保证当前数据不变 */
	sc := newSc(serviceConfig)
	defer scPool.Put(sc)

	var resultBody []byte
	var err error

	/** 加密 */
	if sc.EncryptEnable {

		resultBody, err = encrypt.EncodeReq(bodyBytes, sc)

		if err != nil {
			server.GetLogger().Errorf("%s", err)
			return nil, errors.New("encrypt failure")
		}

	}

	/** 压缩 */
	if sc.CompressEnable {

		resultBody, err = compress.Encode(resultBody, sc.CompressAlgorithm)
		if err != nil {
			//logger.GetLogger().Error("")
			//return nil, fmt.Errorf("parsing %s as HTML: %v", url,err)
			return nil, errors.New("compress failure")
		}
	}

	/** 生成签名 */
	var bodySign string
	if sc.SignEnable {
		bodySign, err = sign.GenerateBodySign(resultBody, sc)

		if err != nil {
			server.GetLogger().Error("%s", err)
			return nil, errors.New("sign failure")
		}

	}

	//protobuf编码
	stSend := pbPool.Get().(*protobuf.ProxyData)
	defer pbPool.Put(stSend)
	stSend.Sign = bodySign
	stSend.Compress = sc.CompressEnable
	stSend.Payload = resultBody
	stSend.SignEnable = sc.SignEnable
	stSend.SignKeyName = sc.SignKeyName
	stSend.EncryptEnable = sc.EncryptEnable
	stSend.EncryptKeyName = sc.EncryptKeyName
	stSend.CompressAlgorithm = sc.CompressAlgorithm

	/* stSend := &protobuf.ProxyData{
		Sign:           bodySign,
		Compress:       sc.CompressEnable,
		Payload:        resultBody,
		SignEnable:     sc.SignEnable,
		SignKeyName:    sc.SignKeyName,
		EncryptEnable:  sc.EncryptEnable,
		EncryptKeyName: sc.EncryptKeyName,
	} */
	pData, err := proto.Marshal(stSend)
	if err != nil {
		return nil, errors.New("protobuf encode failure")
	}

	return pData, nil
}

func Decode(bodyBytes []byte, serviceName string) ([]byte, error) {

	//reData := &protobuf.ProxyRespData{}
	reData := pbRespPool.Get().(*protobuf.ProxyRespData)
	defer pbRespPool.Put(reData)
	err := proto.Unmarshal(bodyBytes, reData)

	if err != nil {
		return nil, errors.New("protobuf decode failure")
	}

	bodyBytes = reData.Payload

	//解压
	if reData.Compress {
		bodyBytes, err = compress.Decode(bodyBytes, reData.CompressAlgorithm)

		if err != nil {
			return nil, errors.New("compress decode failure")
		}
	}

	/** 解密*/
	if reData.EncryptEnable {
		bodyBytes, _ = encrypt.DecodeResp(bodyBytes, reData.EncryptKeyName)

	}

	return bodyBytes, nil
}
