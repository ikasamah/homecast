# homecast

> Make your speaker speak.

`homecast` is a Go package to enable text-to-speech on Google Home in local network.

This is Go version of [noelportugal/google-home-notifier](https://github.com/noelportugal/google-home-notifier)

## Install
```bash
$ go get github.com/ikasamah/homecast
```

## Usage
```golang
devices := homecast.LookupGoogleHome()

ctx := context.Background()
for _, device := range devices {
    err := device.Speak(ctx, "Hello World", "en")
}
```

## Run example
```bash
$ dep ensure
$ go run example/main.go
```

## Author
[Masayuki Hamasaki](https://github.com/ikasamah)
