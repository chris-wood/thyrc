package message

import (
	"testing"
	"reflect"
)

func TestMessage(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		messageString string
		message *Message
	}{
		{"pass 1234\n", New("", "pass", []string{"1234"}) },
	}

	for i, testCase := range testCases {
		got := Parse(testCase.messageString)
		if !reflect.DeepEqual(got, testCase.message) {
			t.Errorf("Case[%v]: got %v, want %v", i, got, testCase.message)
		}
	}
}
