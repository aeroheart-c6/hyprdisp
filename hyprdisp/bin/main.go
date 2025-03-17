package main

import (
	"context"
	"log"

	"aeroheart.io/hyprdisp/hyprland"
	"aeroheart.io/hyprdisp/hyprpanel"
	"aeroheart.io/hyprdisp/profiles"
	"aeroheart.io/hyprdisp/sys"
)

/*

Here's the plan:

1. SHA hash of all monitor's:
	ID# + ID + description
	Order is important

2. Have a config file for each profile

3. Have a default profile where we auto add into the left

read data: {name:monitorremoved data:DP-2}
read data: {name:monitoraddedv2 data:1,DP-2,Beihai Century Joint Innovation Technology Co.Ltd F240v 0000000000001}
read data: {name:monitorremoved data:DP-2}
read data: {name:monitoraddedv2 data:1,DP-2,Beihai Century Joint Innovation Technology Co.Ltd F240v 0000000000001}

*/

func main() {
	var (
		ctx    context.Context
		err    error
		logger *log.Logger
	)

	ctx = setup()

	err = exec(ctx)
	if err != nil {
		logger = ctx.Value(sys.ContextKeyLogger).(*log.Logger)
		logger.Printf("encountered an error: %v", err)
	}
}

func setup() context.Context {
	var ctx context.Context

	ctx = context.Background()
	ctx = setupLogger(ctx)
	return ctx
}

func setupLogger(ctx context.Context) context.Context {
	var logger *log.Logger = log.Default()

	return context.WithValue(ctx, sys.ContextKeyLogger, logger)
}

func exec(ctx context.Context) error {
	var (
		logger       *log.Logger       = ctx.Value(sys.ContextKeyLogger).(*log.Logger)
		hyprlandSrv  hyprland.Service  = hyprland.NewDefaultService()
		hyprpanelSrv hyprpanel.Service = hyprpanel.NewDefaultService()
		profilesSrv  profiles.Service  = profiles.NewDefaultService(hyprlandSrv, hyprpanelSrv)
		monitors     []hyprland.Monitor
		profileCfg   profiles.Config
		err          error
	)

	monitors, err = hyprlandSrv.GetMonitors()
	if err != nil {
		return err
	}

	profileCfg, err = profilesSrv.Detect(ctx, monitors)
	if err != nil {
		return err
	}

	if profileCfg.IsZero() {
		logger.Printf("Configuration for monitors not found. Creating...")

		profileCfg, err = profilesSrv.Init(ctx, monitors)
		if err != nil {
			return err
		}
	} else {
		logger.Printf("Found configuration for monitors doing nothing")
	}

	err = profilesSrv.Apply(ctx, profileCfg)
	if err != nil {
		logger.Printf("oh no: %v", err)
	}

	return nil
}

// func exec(ctx context.Context) error {
// 	var (
// 		hyprCtx          context.Context
// 		hyprCancelFn     context.CancelFunc
// 		profilesCtx      context.Context
// 		profilesCancelFn context.CancelFunc
// 		err              error
// 	)

// 	var (
// 		hyprEvents chan hypr.Event
// 		hyprErrs   chan error
// 		profErrs   chan error = make(chan error, 1)
// 	)

// 	hyprCtx, hyprCancelFn = context.WithCancel(ctx)
// 	defer hyprCancelFn()

// 	profilesCtx, profilesCancelFn = context.WithCancel(ctx)
// 	defer profilesCancelFn()

// 	hyprEvents, hyprErrs, err = hypr.StreamEvents(hyprCtx)
// 	if err != nil {
// 		return err
// 	}

// 	profiles.Init(profilesCtx)
// 	go profiles.ListenEvents(profilesCtx, profErrs, hyprEvents)
// 	go profiles.ListenTimer(profilesCtx, profErrs)

// 	// Wait for SIGTERM / SIGINT
// 	var signals chan os.Signal

// 	signals = make(chan os.Signal, 1)
// 	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
// 	defer close(signals)

// 	select {
// 	case <-signals:
// 		hyprCancelFn()
// 		profilesCancelFn()
// 	case err = <-hyprErrs:
// 		profilesCancelFn()
// 	case err = <-profErrs:
// 		hyprCancelFn()
// 		profilesCancelFn()
// 	}

// 	return err
// }
