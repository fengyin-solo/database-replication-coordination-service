// Package codec owns decoding buffers used by the replication pipeline.
package codec

type Change struct {
	Key     string
	Payload []byte
}

// Decode detaches the result from the caller's reusable receive buffer.
func Decode(key string, receiveBuffer []byte) Change {
	payload := append([]byte(nil), receiveBuffer...)
	return Change{Key: key, Payload: payload}
}
