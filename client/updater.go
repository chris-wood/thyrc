package updater

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	apipb "github.com/cesarghali/lit/proto/api"
)

type Updater struct {
}

// New creates a new instance of the Publisher object.
func New() *Updater {
	return &Updater{}
}

// Upload an update to a published file so that it can be broadcast
// to all subscribed users.
func (s *Updater) Update(ctx context.Context, req *apipb.UpdateRequest) (*apipb.UpdateResponse, error) {
	return nil, grpc.Errorf(codes.Unimplemented, "")
}
