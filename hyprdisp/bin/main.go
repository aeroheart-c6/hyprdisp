package main

import (
	"fmt"
	"strings"
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
	// var (
	// 	// ctx     context.Context = context.Background()
	// 	signals chan os.Signal = make(chan os.Signal, 1)
	// 	// hyprCtx context.Context
	// 	hyprCFn context.CancelFunc
	// 	// err     error
	// )

	// hyprCtx, hyprCFn = context.WithCancel(ctx)
	// err = hypr.Listen(hyprCtx)
	// if err != nil {
	// 	fmt.Printf("error encountered: %v\n", err)
	// 	os.Exit(1)
	// }

	// Wait for SIGTERM / SIGINT
	// signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// for range signals {
	// 	hyprCFn()
	// 	close(signals)
	// }

	paragraph := "" +
		"this is sparta\n" +
		"this is mondstadt\n" +
		"this is liyue\n" +
		"this is"

	fmt.Println(paragraph)

	lines := strings.Split(paragraph, "\n")

	if len(lines[len(lines)-1]) == 0 {
		fmt.Println("we ended with a new line!")
	}
}
