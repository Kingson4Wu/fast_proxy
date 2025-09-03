package encrypt_test

import (
	"github.com/Kingson4Wu/fast_proxy/common/encrypt"
	"github.com/Kingson4Wu/fast_proxy/common/logger"
	"github.com/Kingson4Wu/fast_proxy/common/logger/zap"
	"github.com/Kingson4Wu/fast_proxy/common/server"
	"github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEncode(t *testing.T) {

	patchLogger := gomonkey.ApplyFunc(server.GetLogger, func() logger.Logger {
		return zap.DefaultLogger()
	})
	defer patchLogger.Reset()

	text := "hello world"
	key := "ABCDABCDABCD1234"
	result, err := encrypt.Encode([]byte(text), key)

	Convey("data not null", t, func() {
		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
		So(string(result), ShouldNotEqual, text)
	})

	decryptText, err := encrypt.Decode(result, key)

	Convey("data equals", t, func() {
		So(err, ShouldBeNil)
		So(decryptText, ShouldNotBeNil)
		So(string(decryptText), ShouldEqual, text)
	})

}
