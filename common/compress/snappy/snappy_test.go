package snappy_test

import (
	"github.com/Kingson4Wu/fast_proxy/common/compress/snappy"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEncode(t *testing.T) {

	compress := snappy.Snappy{}
	text := "hello world"
	data := []byte(text)
	encodeText, err := compress.Encode(data)
	if err != nil {
		t.Fatal(err)
	}
	result, err := compress.Decode(encodeText)
	if err != nil {
		t.Fatal(err)
	}
	Convey("data equals", t, func() {
		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
		So(string(result), ShouldEqual, text)
	})
}
