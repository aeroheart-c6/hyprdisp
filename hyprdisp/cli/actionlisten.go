package cli

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"aeroheart.io/hyprdisp/hyprland"
	"aeroheart.io/hyprdisp/profiles"
)

type ListenAction struct {
	Hyprland hyprland.Service
	Profiles profiles.ListenerService
}

func (action ListenAction) ID() string {
	return "listen"
}

func (action ListenAction) Execute(ctx context.Context) error {
	var (
		ctxCancellable context.Context
		ctxCancelFn    context.CancelFunc
		hyprEvents     chan hyprland.Event
		hyprErrs       chan error
		profErrs       chan error = make(chan error, 1)
		err            error
	)

	ctxCancellable, ctxCancelFn = context.WithCancel(ctx)
	defer ctxCancelFn()

	hyprEvents, hyprErrs, err = action.Hyprland.StreamEvents(ctxCancellable)
	if err != nil {
		return err
	}

	go action.Profiles.ListenTimer(ctxCancellable, profErrs)
	go action.Profiles.ListenHyprland(ctxCancellable, hyprEvents, profErrs)

	// Wait for SIGTERM / SIGINT
	var signals chan os.Signal = make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer close(signals)

	select {
	case <-signals:
		ctxCancelFn()
	case <-hyprErrs:
		ctxCancelFn()
	case <-profErrs:
		ctxCancelFn()
	}

	return nil
}
