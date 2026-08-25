package replication

import (
	"datasync/internal/codec"
	"datasync/internal/state"
)

type ExportService struct{ Cache *state.BatchCache }

func (s *ExportService) Queue(key string, receiveBuffer []byte, release <-chan struct{}) <-chan codec.Change {
	result := make(chan codec.Change, 1)
	change := codec.Decode(key, receiveBuffer)
	s.Cache.Put(change)
	go func(snapshot codec.Change) {
		<-release
		result <- snapshot
		close(result)
	}(change)
	return result
}
