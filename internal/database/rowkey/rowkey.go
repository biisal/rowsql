// Package rowkey identifies a single table row by its primary key, a physical row id, or all of its columns.
package rowkey

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

type Kind string

const (
	KindPrimary    Kind = "primary"
	KindPhysical   Kind = "physical"
	KindAllColumns Kind = "all"
)

var ErrInvalidKey = errors.New("invalid row key")

type Key struct {
	Kind    Kind     `json:"kind"`
	Columns []string `json:"columns"`
	Values  []any    `json:"values"`
}

func (k Key) Encode() (string, error) {
	encoded, err := json.Marshal(k)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(encoded), nil
}

func Decode(encoded string) (Key, error) {
	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return Key{}, ErrInvalidKey
	}

	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var key Key
	if err := decoder.Decode(&key); err != nil {
		return Key{}, ErrInvalidKey
	}
	if err := key.validate(); err != nil {
		return Key{}, err
	}
	return key, nil
}

func (k Key) validate() error {
	switch k.Kind {
	case KindPrimary, KindPhysical, KindAllColumns:
	default:
		return ErrInvalidKey
	}
	if len(k.Columns) == 0 || len(k.Columns) != len(k.Values) {
		return ErrInvalidKey
	}
	return nil
}

func NormalizeValue(value any) any {
	switch typed := value.(type) {
	case []byte:
		return string(typed)
	case time.Time:
		return typed.Format(time.RFC3339Nano)
	default:
		return value
	}
}
