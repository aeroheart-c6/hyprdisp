package profiles

import (
	"context"
	"log"
	"time"

	"aeroheart.io/hyprdisp/hyprland"
	"aeroheart.io/hyprdisp/sys"
)

type state string

// TODO it's actually just 2 - 3 seconds wait
const triggerDuration time.Duration = 5 * time.Second

func (s defaultService) ListenEvents(ctx context.Context, errs chan error, events chan hyprland.Event) {
	var logger *log.Logger = ctx.Value(sys.ContextKeyLogger).(*log.Logger)

	for event := range events {
		if event.Name != hyprland.EventNameMonitorAdded && event.Name != hyprland.EventNameMonitorRemoved {
			logger.Printf("Got irrelevant event: %v\n", event)
			continue
		}

		logger.Printf("Got monitor event: %v\n", event)
		s.timer.Reset(triggerDuration)
	}
}

func (s defaultService) ListenTimer(ctx context.Context, errs chan error) {
	var logger *log.Logger = ctx.Value(sys.ContextKeyLogger).(*log.Logger)

	for moment := range s.timer.C {
		logger.Printf("whoopps I just triggered! at %v\n", moment.Format(time.RFC3339))

		triggerConfigUpdates()
	}
}

func triggerConfigUpdates() {

}
