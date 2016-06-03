package message

import (
	"testing"

	"golang.org/x/net/context"
)

type Environment struct {
	ctx context.Context
	m   *Message
}

func NewEnvironment(t *testing.T) *Environment {
	ctx := context.Background()
	m := New()

	return &Environment{ctx, m}
}

func TestClient(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		messageString string
		message *Message
	}{
		
		{nil, codes.Unimplemented},
	}

	for i, testCase := range testCases {
		env := NewEnvironment(t)
		_, err := env.u.Update(env.ctx, testCase.update)
		if got, want := grpc.Code(err), testCase.outErr; got != want {
			t.Errorf("Case[%v]: Update(%v)=%v, want %v", i, testCase.update, got, want)
		}
	}
}
