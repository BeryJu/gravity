package storage

import (
	"encoding/json"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Parse Attempt to parse data as protobuf first, otherwise try json
func Parse[T protoreflect.ProtoMessage](raw []byte, v T) (T, error) {
	// Try loading protobuf first
	err := proto.Unmarshal(raw, v)
	if err == nil {
		return v, nil
	}
	// Otherwise try json
	err = json.Unmarshal(raw, &v)
	if err == nil {
		return v, nil
	}
	return v, err
}
