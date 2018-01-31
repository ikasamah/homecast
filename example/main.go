package main

import (
	"context"

	"github.com/ikasamah/homecast"
)

func main() {
	devices := homecast.LookupGoogleHome()

	ctx := context.Background()
	for _, device := range devices {
		_ = device.Speak(ctx, "Hello World", "en")
	}
}
