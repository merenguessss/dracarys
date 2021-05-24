package server

import "context"

type dispatcher struct {
	serviceMap map[string]Service
}

func (d *dispatcher) Handle(ctx context.Context, b []byte) ([]byte, error) {

}
