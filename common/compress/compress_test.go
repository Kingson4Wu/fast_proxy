package compress

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEncode(t *testing.T) {

	for algorithm := range algorithms {
		encode(t, int32(algorithm))
	}
}

func encode(t *testing.T, algorithm int32) {
	text := "hello world"
	data := []byte(text)
	encodeText, err := Encode(data, algorithm)
	if err != nil {
		t.Fatal(err)
	}
	result, err := Decode(encodeText, algorithm)
	if err != nil {
		t.Fatal(err)
	}
	Convey("data equals", t, func() {
		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
		So(string(result), ShouldEqual, text)
	})
}
