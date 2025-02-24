package profiles

import (
	"context"
	"fmt"
	"time"

	"aeroheart.io/hyprdisp/hypr"
)

// TODO it's actually just 2 - 3 seconds wait
const triggerDuration time.Duration = 5 * time.Second

var (
	timer        *time.Timer
)

func Init(ctx context.Context) {
	timer = time.NewTimer(0)
}

func ListenEvents(ctx context.Context, errs chan error, events chan hypr.Event) {
	var logger *log.Logger = ctx.Value(sys.ContextKeyLogger).(*log.Logger)

	for event := range events {
		if event.Name != hypr.EventNameMonitorAdded && event.Name != hypr.EventNameMonitorRemoved {
			fmt.Printf("Got irrelevant event: %v\n", event)
			continue
		}

		fmt.Printf("Got monitor event: %v\n", event)
		timer.Reset(triggerDuration)
	}
}

func ListenTimer(ctx context.Context, errs chan error) {
	var logger *log.Logger = ctx.Value(sys.ContextKeyLogger).(*log.Logger)

	for moment := range timer.C {

		triggerConfigUpdates()
	}
}

func triggerConfigUpdates() {

}
