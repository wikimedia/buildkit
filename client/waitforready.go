package client

import (
	"context"

	controlapi "github.com/moby/buildkit/api/services/control"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// WaitForReady performs a preflight request (currently ListWorkers) with
// grpc.WaitForReady(true) to ensure that grpc.ClientConn has established the
// underlying connection and that it can be considered available. Performing
// this request prior to solves makes the client more robust in environments
// where the server is behind a proxy or part of a service mesh (e.g.
// Istio/Envoy). In these environments, connections may be prematurely closed
// prior to any client requests due circuit breaking on max connections.
//
// TODO it might be better to make this a grpc.health.v1.Health/Check request.
// However the main server does not currently implement that API.
func (c *Client) WaitForReady(ctx context.Context) error {
	if c.readyTimeout != nil {
		newCtx, cancel := context.WithTimeout(ctx, *c.readyTimeout)
		ctx = newCtx
		defer cancel()
	}

	_, err := c.ControlClient().ListWorkers(ctx, &controlapi.ListWorkersRequest{}, grpc.WaitForReady(true))
	if err != nil {
		return errors.Wrap(err, "failed to ensure an available connection")
	}
	return nil
}
