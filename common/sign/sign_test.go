package sign_test

import (
	"github.com/Kingson4Wu/fast_proxy/common/sign"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGenerateSign(t *testing.T) {

	text := "hello"
	salt := "123456"
	signature, err := sign.GenerateSign([]byte(text), salt)
	Convey("signature not null", t, func() {
		So(err, ShouldBeNil)
		So(signature, ShouldNotBeEmpty)
	})
}
