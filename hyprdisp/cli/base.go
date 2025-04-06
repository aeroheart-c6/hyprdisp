package cli

import "context"

type ActionHandler interface {
	ID() string
	Execute(context.Context) error
}

type ActionConfigurer interface {
	Configure([]string) error
}
