package sntp

import (
	"reflect"
	"testing"
)

func TestInt2Bytes(t *testing.T) {
	tests := []struct {
		input    int64
		expected []byte
	}{
		{0, []byte{0, 0, 0, 0}},
		{255, []byte{0, 0, 0, 255}},
		{65535, []byte{0, 0, 255, 255}},
		{16777215, []byte{0, 255, 255, 255}},
		{4294967295, []byte{255, 255, 255, 255}},
	}

	for _, test := range tests {
		result := int2bytes(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("For input %d, expected %v, but got %v", test.input, test.expected, result)
		}
	}
}
