package profiles

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"aeroheart.io/hyprdisp/hyprland"
	"aeroheart.io/hyprdisp/sys"
)

type state string

// TODO it's actually just 2 - 3 seconds wait
const triggerDuration time.Duration = 5 * time.Second

func (s defaultService) ListenEvents(ctx context.Context, errs chan error, events chan hyprland.Event) {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return
	}

	for event := range events {
		if event.Name != hyprland.EventNameMonitorAdded && event.Name != hyprland.EventNameMonitorRemoved {
			logger.Info("Got irrelevant event: %v\n", slog.Any("event", event))
			continue
		}

		logger.Info("Got monitor event: %v\n", slog.Any("event", event))
		s.timer.Reset(triggerDuration)
	}
}

func (s defaultService) ListenTimer(ctx context.Context, errs chan error) {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return
	}

	for moment := range s.timer.C {
		logger.Info(fmt.Sprintf("whoopps I just triggered! at %v\n", moment.Format(time.RFC3339)))
		triggerConfigUpdates()
	}
}

func triggerConfigUpdates() {

}
