package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	hyprdispIPC "aeroheart.io/hyprdisp/ipc"
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
		ctx      context.Context = context.Background()
		sigChan  chan os.Signal  = make(chan os.Signal, 1)
		hyprChan chan hyprdispIPC.HyprEvent
		hyprCtx  context.Context
		hyprCFn  context.CancelFunc
		err      error
	)

	// Wait for SIGTERM / SIGINT

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	hyprCtx, hyprCFn = context.WithCancel(ctx)
	hyprChan, err = hyprdispIPC.ListenHyprEvents(hyprCtx)

	if err != nil {
		fmt.Printf("error encountered: %w\n", err)
		os.Exit(1)
	}

	// envRuntimePath, found = os.LookupEnv("XDG_RUNTIME_DIR")
	// if !found {
	// 	fmt.Println("Could not find variable: XDG_RUNTIME_DIR")
	// 	os.Exit(1)
	// }

	// envHyprlandSig, found = os.LookupEnv("HYPRLAND_INSTANCE_SIGNATURE")
	// if !found {
	// 	fmt.Println("Could not find variable: HYPRLAND_INSTANCE_SIGNATURE")
	// 	os.Exit(1)
	// }

	// var (
	// 	conn net.Conn
	// 	err  error
	// )
	// conn, err = net.Dial("unix", fmt.Sprintf("%s/hypr/%s/.socket2.sock",
	// 	envRuntimePath,
	// 	envHyprlandSig,
	// ))
	// if err != nil {
	// 	fmt.Println("oh no: %w", err)
	// 	os.Exit(1)
	// }

	// var (
	// 	sigChan  chan os.Signal = make(chan os.Signal, 1)
	// 	hyprChan chan hyprEvent = make(chan hyprEvent)
	// 	hyprCtx  context.Context
	// 	hyprCFn  context.CancelFunc
	// )

	// hyprdisp.Something()

	// // Wait for SIGTERM / SIGINT
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// // Wait for input...
	// hyprCtx, hyprCFn = context.WithCancel(ctx)
	// go listen(hyprCtx, conn, hyprChan)

listener:
	for {
		select {
		case data := <-hyprChan:
			fmt.Printf("read data: %+v\n", data)
		case <-sigChan:
			hyprCFn()
			break listener
		}
	}
}
