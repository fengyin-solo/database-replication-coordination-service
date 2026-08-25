package replication

import "datasync/internal/state"

type SchemaService struct{ Registry *state.SchemaRegistry }

func (s *SchemaService) Refresh(name string, builder state.SchemaBuilder) (*state.Schema, error) {
	schema, err := s.Registry.Build(name, builder)
	if err != nil {
		if partial, ok := s.Registry.Get(name); ok {
			return partial, err
		}
	}
	return schema, err
}

func (s *SchemaService) Lookup(name string) (*state.Schema, bool) {
	return s.Registry.Get(name)
}
