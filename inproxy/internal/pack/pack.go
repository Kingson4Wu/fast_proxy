package pack

import (
	"bytes"
	"errors"
	"github.com/Kingson4Wu/fast_proxy/common/cerror"
	"github.com/Kingson4Wu/fast_proxy/common/logger"
	"github.com/Kingson4Wu/fast_proxy/common/pool"
	"github.com/Kingson4Wu/fast_proxy/common/proto/protobuf"
	"github.com/Kingson4Wu/fast_proxy/inproxy/internal/compress"
	"github.com/Kingson4Wu/fast_proxy/inproxy/internal/encrypt"
	"github.com/Kingson4Wu/fast_proxy/inproxy/internal/sign"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	"go.uber.org/zap"
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
		bodyBytes, err = ioutil.ReadAll(req.Body)
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
		bodyBytes, err = ioutil.ReadAll(resp.Body)
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

func Encode(bodyBytes []byte, reData *protobuf.ProxyData) ([]byte, error) {

	var resultBody []byte
	var err error
	var ok bool
	/** 加密 */
	if reData.EncryptEnable {

		resultBody, err = encrypt.EncodeResp(bodyBytes, reData.EncryptKeyName)

		if err != nil {
			logger.GetLogger().Error("", zap.Error(err))
			return nil, errors.New("encrypt failure")
		}

	}

	/** 压缩 */
	if reData.Compress {

		resultBody, ok = compress.SnappyCompress.Encode(resultBody)
		if !ok {
			//logger.GetLogger().Error("")
			//return nil, fmt.Errorf("parsing %s as HTML: %v", url,err)
			return nil, errors.New("compress failure")
		}
	}

	/** protobuf编码 */
	stSend := pbRespPool.Get().(*protobuf.ProxyRespData)
	defer pbRespPool.Put(stSend)
	stSend.Compress = reData.Compress
	stSend.Payload = resultBody
	stSend.EncryptEnable = reData.EncryptEnable
	stSend.EncryptKeyName = reData.EncryptKeyName

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

	/** 签名验证 */
	if reData.SignEnable {

		bodySignature, err := sign.GenerateBodySign(bodyBytes, reData.SignKeyName)
		if err != nil {
			return nil, nil, errors.New("generate sign failure")
		}
		reqSignature := reData.Sign

		if reqSignature != bodySignature {
			return nil, nil, errors.New("body sign check failure")
		}

	}

	if reData.Compress {
		//解压
		bodyBytes, err = compress.SnappyCompress.Decode(bodyBytes)

		if err != nil {
			return nil, nil, errors.New("compress decode failure")
		}
	}

	if reData.EncryptEnable {

		/** 解密*/
		bodyBytes, err = encrypt.DecodeReq(bodyBytes, reData.EncryptKeyName)
		if err != nil {
			return nil, nil, errors.New("encrypt decode failure")
		}
	}

	return bodyBytes, reData, nil
}
