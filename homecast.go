package homecast

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/barnybug/go-cast"
	"github.com/barnybug/go-cast/controllers"
	castnet "github.com/barnybug/go-cast/net"
	"github.com/hashicorp/mdns"
)

const (
	googleCastServiceName = "_googlecast._tcp"
	googleHomeModelInfo   = "md=Google Home"
)

// CastDevice is cast-able device contains cast client
type CastDevice struct {
	*mdns.ServiceEntry
	client *cast.Client
}

// Connect connects required services to cast
func (g *CastDevice) Connect(ctx context.Context) error {
	return g.client.Connect(ctx)
}

// Close calls client's close func
func (g *CastDevice) Close() {
	g.client.Close()
}

// Speak speaks given text on cast device
func (g *CastDevice) Speak(ctx context.Context, text, lang string) error {
	url, err := tts(text, lang)
	if err != nil {
		return err
	}
	return g.Play(ctx, url)
}

// Play plays media contents on cast device
func (g *CastDevice) Play(ctx context.Context, url *url.URL) error {
	conn := castnet.NewConnection()
	if err := conn.Connect(ctx, g.AddrV4, g.Port); err != nil {
		return err
	}
	defer conn.Close()

	status, err := g.client.Receiver().LaunchApp(ctx, cast.AppMedia)
	if err != nil {
		return err
	}
	app := status.GetSessionByAppId(cast.AppMedia)

	cc := controllers.NewConnectionController(conn, g.client.Events, cast.DefaultSender, *app.TransportId)
	if err := cc.Start(ctx); err != nil {
		return err
	}
	media := controllers.NewMediaController(conn, g.client.Events, cast.DefaultSender, *app.TransportId)
	if err := media.Start(ctx); err != nil {
		return err
	}

	mediaItem := controllers.MediaItem{
		ContentId:   url.String(),
		ContentType: "audio/mp3",
		StreamType:  "BUFFERED",
	}

	log.Printf("[INFO] Load media: content_id=%s", mediaItem.ContentId)
	_, err = media.LoadMedia(ctx, mediaItem, 0, true, nil)

	return err
}

// LookupAndConnect retrieves cast-able google home devices
func LookupAndConnect(ctx context.Context) []*CastDevice {
	entriesCh := make(chan *mdns.ServiceEntry, 4)

	results := make([]*CastDevice, 0, 4)
	go func() {
		for entry := range entriesCh {
			log.Printf("[INFO] ServiceEntry detected: [%s:%d]%s", entry.AddrV4, entry.Port, entry.Name)
			for _, field := range entry.InfoFields {
				if field == googleHomeModelInfo {
					client := cast.NewClient(entry.AddrV4, entry.Port)
					if err := client.Connect(ctx); err != nil {
						log.Printf("[ERROR] Failed to connect: %s", err)
						continue
					}
					results = append(results, &CastDevice{entry, client})
				}
			}
		}
	}()

	mdns.Lookup(googleCastServiceName, entriesCh)
	close(entriesCh)

	return results
}

// tts provides text-to-speech sound url.
// NOTE: it seems to be unofficial.
func tts(text, lang string) (*url.URL, error) {
	base := "https://translate.google.com/translate_tts?client=tw-ob&ie=UTF-8&q=%s&tl=%s"
	return url.Parse(fmt.Sprintf(base, url.QueryEscape(text), url.QueryEscape(lang)))
}
