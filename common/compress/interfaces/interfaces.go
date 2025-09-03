package interfaces

type Compress interface {
	Encode(data []byte) ([]byte, error)
	Decode(data []byte) ([]byte, error)
}