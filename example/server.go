package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/ikasamah/homecast"
)

func main() {
	port := flag.Int("port", 8080, "Listen port")
	defaultLang := flag.String("lang", "en", "Default language to speak")
	flag.Parse()

	ctx := context.Background()
	devices := homecast.LookupAndConnect(ctx)
	defer func() {
		for _, device := range devices {
			device.Close()
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		text := r.FormValue("text")
		lang := r.FormValue("lang")

		if text == "" {
			log.Printf("[INFO] Skip request due to no text given")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if lang == "" {
			lang = *defaultLang
		}

		for _, device := range devices {
			if err := device.Speak(ctx, text, lang); err != nil {
				log.Printf("[ERROR] Failed to speak: %v", err)
			}
		}
	})

	addr := fmt.Sprintf(":%d", *port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
