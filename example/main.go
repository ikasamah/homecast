package main

import (
	"context"
	"fmt"

	"github.com/ikasamah/homecast"
)

func main() {
	ctx := context.Background()
	devices := homecast.LookupAndConnect(ctx)

	for _, device := range devices {
		fmt.Printf("Device: [%s:%d]%s", device.AddrV4, device.Port, device.Name)

		if err := device.Speak(ctx, "Hello World", "en"); err != nil {
			fmt.Printf("Failed to speak: %v", err)
		}
	}
}
