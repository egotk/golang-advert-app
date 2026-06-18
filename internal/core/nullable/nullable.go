package nullable

import "encoding/json"

type Nullable[T any] struct {
	Set   bool
	Value *T
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	n.Set = true
	return json.Unmarshal(data, &n.Value)
}
