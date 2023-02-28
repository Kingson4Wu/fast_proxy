package pack

import "github.com/Kingson4Wu/fast_proxy/common/config"

type PackHandler interface {
	Fmap(fn func([]byte, *config.ServiceConfig) ([]byte, error)) PackHandler
}
type packHandlerImpl struct {
	b  []byte
	sc *config.ServiceConfig
}

func (phi packHandlerImpl) Fmap(fn func([]byte, *config.ServiceConfig) ([]byte, error)) PackHandler {
	nb, err := fn(phi.b, phi.sc)
	if err != nil {
	}

	return packHandlerImpl{b: nb}
}

func NewPackHandler(b []byte, sc *config.ServiceConfig) PackHandler {
	return packHandlerImpl{b: b, sc: sc}
}

/* func main() {
	b := []byte{1, 2, 3, 4}
	f := NewPackHandler(b, nil)
	mapperFunc1 := func(b []byte, sc *apollo.ServiceConfig) ([]byte, error) {
		return b, nil
	}
	mapped1 := f.Fmap(mapperFunc1)
	fmt.Printf("mapped functor1: %+v\n", mapped1)
} */
