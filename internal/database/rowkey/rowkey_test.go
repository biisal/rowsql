package rowkey

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestEncodeDecodeRoundTrip(t *testing.T) {
	original := Key{
		Kind:    KindPrimary,
		Columns: []string{"user_id", "note"},
		Values:  []any{int64(9007199254740993), nil},
	}

	encoded, err := original.Encode()
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}

	want := Key{
		Kind:    KindPrimary,
		Columns: []string{"user_id", "note"},
		Values:  []any{json.Number("9007199254740993"), nil},
	}
	if !reflect.DeepEqual(decoded, want) {
		t.Fatalf("decoded = %#v, want %#v", decoded, want)
	}
}

func TestDecodeRejectsInvalidKeys(t *testing.T) {
	validJSON := func(key Key) string {
		encoded, _ := key.Encode()
		return encoded
	}
	tests := map[string]string{
		"not base64":        "!!!",
		"not json":          "bm90LWpzb24",
		"unknown kind":      validJSON(Key{Kind: "other", Columns: []string{"id"}, Values: []any{1}}),
		"no columns":        validJSON(Key{Kind: KindPrimary}),
		"length mismatched": validJSON(Key{Kind: KindPrimary, Columns: []string{"id"}, Values: []any{1, 2}}),
	}
	for name, encoded := range tests {
		if _, err := Decode(encoded); err == nil {
			t.Errorf("%s: expected error", name)
		}
	}
}

func TestNormalizeValue(t *testing.T) {
	moment := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	if got := NormalizeValue([]byte("abc")); got != "abc" {
		t.Errorf("bytes: %v", got)
	}
	if got := NormalizeValue(moment); got != "2024-01-02T03:04:05Z" {
		t.Errorf("time: %v", got)
	}
	if got := NormalizeValue(int64(5)); got != int64(5) {
		t.Errorf("int: %v", got)
	}
}
