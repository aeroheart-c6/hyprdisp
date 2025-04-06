package main

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"aeroheart.io/hyprdisp/cli"
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
		logger *slog.Logger
	)

	ctx = setup()

	err = exec(ctx)
	if err != nil {
		logger, _ = sys.GetLogger(ctx)
		logger.Info("encountered an error: %v", err)
	}
}

func setup() context.Context {
	var ctx context.Context

	ctx = context.Background()
	ctx = setupLogger(ctx)
	ctx = setupActions(ctx)

	return ctx
}

func setupLogger(ctx context.Context) context.Context {
	var logger *slog.Logger = slog.Default()

	return sys.SetLogger(ctx, logger)
}

func setupActions(ctx context.Context) context.Context {
	var logger *slog.Logger
	logger, _ = sys.GetLogger(ctx)

	var (
		hyprlandSrv  hyprland.Service  = hyprland.NewDefaultService("./var")
		hyprpanelSrv hyprpanel.Service = hyprpanel.NewDefaultService("./var")
		profilesSrv  profiles.Service  = profiles.NewDefaultService(hyprlandSrv, hyprpanelSrv)
	)

	logger.Info("Configuring Actions")
	var registry cli.ActionRegistry = cli.ActionRegistry{}
	registry.Add(&cli.DetectAction{
		HyprLand: hyprlandSrv,
		Profiles: profilesSrv,
	})
	registry.Add(&cli.ApplyAction{
		HyprLand: hyprlandSrv,
		Profiles: profilesSrv,
	})

	return context.WithValue(ctx, sys.ContextKeyCLIActions, registry)
}

func exec(ctx context.Context) error {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return err
	}

	// Validate arguments
	if len(os.Args) < 2 {
		return errors.New("subcommand is required")
	}

	var (
		registry      cli.ActionRegistry
		actionHandler cli.ActionHandler
		actionConfig  cli.ActionConfigurer
		ok            bool
	)
	registry, ok = ctx.Value(sys.ContextKeyCLIActions).(cli.ActionRegistry)
	if !ok {
		return errors.New("invalid registry in context")
	}

	actionHandler, err = registry.Get(ctx, os.Args[1])
	if err != nil {
		return err
	}
	logger.Info("Action - Found", slog.String("actionID", actionHandler.ID()))

	actionConfig, ok = actionHandler.(cli.ActionConfigurer)
	if ok {
		logger.Info("Action - Configuring", slog.String("actionID", actionHandler.ID()))
		actionConfig.Configure(os.Args[2:])
	} else {
		logger.Info("Action - Skipping Configuration", slog.String("actionID", actionHandler.ID()))
	}

	logger.Info("Action - Running", slog.String("actionID", actionHandler.ID()))
	logger.Info("========================================")
	return actionHandler.Execute(ctx)
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
