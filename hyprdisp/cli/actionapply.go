package cli

import (
	"context"
	"errors"
	"flag"
	"log/slog"

	"aeroheart.io/hyprdisp/hyprland"
	"aeroheart.io/hyprdisp/profiles"
	"aeroheart.io/hyprdisp/sys"
)

type ApplyAction struct {
	HyprLand hyprland.Service
	Profiles profiles.Service
	faked    *bool
}

func (action ApplyAction) ID() string {
	return "apply"
}

func (action ApplyAction) Execute(ctx context.Context) error {
	var (
		logger *slog.Logger
		err    error
	)
	logger, err = sys.GetLogger(ctx)
	if err != nil {
		return err
	}

	var (
		monitors   []hyprland.Monitor
		profileCfg profiles.Config
	)
	monitors, err = action.HyprLand.GetMonitors()
	if err != nil {
		return err
	}

	logger.Info("Detecting current configuration")
	profileCfg, err = action.Profiles.Detect(ctx, monitors)
	if err != nil {
		return err
	}

	if profileCfg.IsZero() {
		return errors.New("empty profile found")
	}

	logger.Info("Applying profile")
	return action.Profiles.Apply(ctx, profileCfg)
}

func (action *ApplyAction) Configure(arguments []string) error {
	var fs *flag.FlagSet = flag.NewFlagSet(action.ID(), flag.ExitOnError)

	action.faked = fs.Bool("faked", false, "whether it will overwrite the actual configuration files or not")

	return fs.Parse(arguments)
}
