package pack

import (
	"bytes"
	"errors"
	"github.com/Kingson4Wu/fast_proxy/common/cerror"
	"github.com/Kingson4Wu/fast_proxy/common/compress"
	"github.com/Kingson4Wu/fast_proxy/common/pool"
	"github.com/Kingson4Wu/fast_proxy/common/proto/protobuf"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/Kingson4Wu/fast_proxy/common/sign"
	"github.com/Kingson4Wu/fast_proxy/inproxy/internal/encrypt"
	"io"
	"net/http"
	"sync"

	"google.golang.org/protobuf/proto"
)

var (
	PbPool = sync.Pool{
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

func DecodeReq(req *http.Request) ([]byte, *protobuf.ProxyData, *cerror.Err) {

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
		return nil, nil, cerror.NewError(http.StatusInternalServerError, "illegal body")
	}

	bb, pb, err := Decode(bodyBytes)

	if err != nil {
		return nil, nil, cerror.NewError(http.StatusInternalServerError, "decode failure")
	}

	return bb, pb, nil

}

func EncodeResp(resp *http.Response, reData *protobuf.ProxyData) ([]byte, *cerror.Err) {

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
		return nil, cerror.NewError(http.StatusBadRequest, "illegal body")
	}

	pData, err := Encode(bodyBytes, reData)
	if err != nil {
		return nil, cerror.NewError(http.StatusInternalServerError, err.Error())
	}

	return pData, nil
}

func EncodeFastResp(bodyBytes []byte, reData *protobuf.ProxyData) ([]byte, *cerror.Err) {
	pData, err := Encode(bodyBytes, reData)

	if err != nil {
		return nil, cerror.NewError(http.StatusInternalServerError, "EncodeFastResp failure")
	}

	return pData, nil
}

func Encode(bodyBytes []byte, reData *protobuf.ProxyData) ([]byte, error) {

	var resultBody []byte
	var err error
	if reData.EncryptEnable {

		resultBody, err = encrypt.EncodeResp(bodyBytes, reData.EncryptKeyName)

		if err != nil {
			server.GetLogger().Errorf("%s", err)
			return nil, errors.New("encrypt failure")
		}

	}

	if reData.Compress {

		resultBody, err = compress.Encode(resultBody, reData.CompressAlgorithm)
		if err != nil {
			//logger.GetLogger().Error("")
			//return nil, fmt.Errorf("parsing %s as HTML: %v", url,err)
			return nil, errors.New("compress failure")
		}
	}

	stSend := pbRespPool.Get().(*protobuf.ProxyRespData)
	defer pbRespPool.Put(stSend)
	stSend.Compress = reData.Compress
	stSend.Payload = resultBody
	stSend.EncryptEnable = reData.EncryptEnable
	stSend.EncryptKeyName = reData.EncryptKeyName
	stSend.CompressAlgorithm = reData.CompressAlgorithm

	pData, err := proto.Marshal(stSend)
	if err != nil {
		return nil, errors.New("protobuf encode failure")
	}

	return pData, nil
}

func Decode(bodyBytes []byte) ([]byte, *protobuf.ProxyData, error) {

	reData := PbPool.Get().(*protobuf.ProxyData)
	//defer pbPool.Put(reData)
	err := proto.Unmarshal(bodyBytes, reData)

	if err != nil {
		return nil, nil, errors.New("protobuf decode failure")
	}

	bodyBytes = reData.Payload

	if reData.SignEnable {

		bodySignature, err := sign.GenerateBodySignWithName(bodyBytes, reData.SignKeyName)
		if err != nil {
			return nil, nil, errors.New("generate sign failure")
		}
		reqSignature := reData.Sign

		if reqSignature != bodySignature {
			return nil, nil, errors.New("body sign check failure")
		}

	}

	if reData.Compress {
		bodyBytes, err = compress.Decode(bodyBytes, reData.CompressAlgorithm)

		if err != nil {
			return nil, nil, errors.New("compress decode failure")
		}
	}

	if reData.EncryptEnable {

		bodyBytes, err = encrypt.DecodeReq(bodyBytes, reData.EncryptKeyName)
		if err != nil {
			return nil, nil, errors.New("encrypt decode failure")
		}
	}

	return bodyBytes, reData, nil
}
