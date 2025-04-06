package cli

import "context"

type ListenHandler struct{}

func (action ListenHandler) ID() string {
	return "listen"
}

func (action ListenHandler) Execute(ctx context.Context) error {
	return nil
}
