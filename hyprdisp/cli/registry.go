package cli

import (
	"context"
	"fmt"
)

type ActionRegistry map[string]ActionHandler

func (registry ActionRegistry) Add(action ActionHandler) {
	var found bool

	_, found = registry[action.ID()]
	if found {
		return
	}

	registry[action.ID()] = action
}

func (registry ActionRegistry) Get(ctx context.Context, request string) (ActionHandler, error) {
	var (
		action ActionHandler
		ok     bool
	)
	action, ok = registry[request]
	if !ok {
		return nil, fmt.Errorf("unsupported action \"%v\"", request)
	}

	return action, nil
}
