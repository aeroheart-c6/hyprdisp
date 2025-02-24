package profiles

import (
	"context"
	"fmt"

	"aeroheart.io/hyprdisp/hypr"
)


func ListenEvents(ctx context.Context, errs chan error, events chan hypr.Event) {
	var logger *log.Logger = ctx.Value(sys.ContextKeyLogger).(*log.Logger)

	for event := range events {
		if event.Name != hypr.EventNameMonitorAdded && event.Name != hypr.EventNameMonitorRemoved {
			fmt.Printf("Got irrelevant event: %v\n", event)
			continue
		}

		fmt.Printf("Got monitor event: %v\n", event)
	}
}

	}

}
