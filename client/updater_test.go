package updater

import (
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	apipb "github.com/cesarghali/lit/proto/api"
)

type Environment struct {
	ctx context.Context
	u   *Updater
}

func NewEnvironment(t *testing.T) *Environment {
	ctx := context.Background()
	u := New()

	return &Environment{ctx, u}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		update *apipb.UpdateRequest
		outErr codes.Code
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
