package app

import (
	"context"
)

type Instance struct{}

func New(ctx context.Context, conf *Config) (*Instance, error) {
}

func (i *Instance) Start(ctx context.Context) error {
	errCh := make(chan error)
}

func (i *Instance) Stop(ctx context.Context) error {

}
