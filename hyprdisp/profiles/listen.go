package profiles

import (
	"context"
	"io/fs"
	"log/slog"
	"time"

	"aeroheart.io/hyprdisp/hyprland"
	"aeroheart.io/hyprdisp/sys"
)

type state int

func (s state) Valid() bool {
	switch s {
	case
		SetupState,
		WatchState,
		ApplyState:
		return true
	default:
		return false
	}
}

const (
	SetupState state = 0
	WatchState state = 1
	ApplyState state = 2
)

// TODO it's actually just 2 - 3 seconds wait
const triggerDuration time.Duration = 3 * time.Second

// ListenHyprland listens for hyprland events coming from the IPC channel and responds when it's a relevant window event
func (s *defaultService) ListenHyprland(ctx context.Context, events chan hyprland.Event, errs chan error) {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return
	}

	for event := range events {
		var relevant bool = event.Name == hyprland.EventNameMonitorAdded ||
			event.Name == hyprland.EventNameMonitorRemoved

		if !relevant {
			continue
		}

		logger.Info("Received hyprland event",
			slog.Any("event", event),
			slog.Bool("relevant", relevant),
		)

		// if s.state == WatchState {
		// 	s.timer.Reset(triggerDuration)
		// }
	}
}

// ListenTimer listens to the internal timer which is the actual action after receiving a relevant Hyprland event
func (s *defaultService) ListenTimer(ctx context.Context, errs chan error) {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return
	}

	s.state = SetupState
	s.timer = time.NewTimer(0)

	for range s.timer.C {
		if s.state == SetupState {
			s.state = WatchState
			continue
		}

		s.state = ApplyState
		logger.Info("Changing state", slog.String("state", "Apply"), slog.Any("stateValue", s.state))

		s.triggerConfigUpdates(ctx)

		s.state = WatchState
		logger.Info("Changing state", slog.String("state", "Watch"), slog.Any("stateValue", s.state))
	}
}

func (s *defaultService) triggerConfigUpdates(ctx context.Context) error {
	var (
		logger   *slog.Logger
		monitors []hyprland.Monitor
		config   Config
		err      error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return err
	}

	logger.Info("profile: querying monitors...")
	monitors, err = s.hyprland.GetMonitors()
	if err != nil {
		return err
	}

	logger.Info("profile: detecting...")
	config, err = s.Detect(ctx, monitors)
	if _, ok := err.(*fs.PathError); ok {
		logger.Debug("Path error received from profile detection", slog.Any("error", err))
	} else if err != nil {
		return err
	}

	if config.IsZero() {
		logger.Info("profile: no current config. creating...")
		config, err = s.Init(ctx, monitors)
		if err != nil {
			return err
		}
	}

	logger.Info("profile: applying...")
	return s.Apply(ctx, config)
}
