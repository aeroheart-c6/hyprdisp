package profiles

import (
	"context"
	"fmt"

	"aeroheart.io/hyprdisp/hypr"
)

func Listen(ctx context.Context) error {
	var (
		events chan hypr.Event = make(chan hypr.Event)
		errs   chan error      = make(chan error, 1)
	)
	defer func() {
		close(events)
		close(errs)
	}()

	hypr.StreamEvents(ctx, events, errs)

	for event := range events {
		if event.Name != hypr.EventNameMonitorAdded && event.Name != hypr.EventNameMonitorRemoved {
			fmt.Printf("Got irrelevant event: %v\n", event)
			continue
		}

		fmt.Printf("Got monitor event: %v\n", event)
	}

	for err := range errs {
		return err
	}

	return nil
}
