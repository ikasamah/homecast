# homecast

> Make your speaker speak.

`homecast` is a Go package to enable text-to-speech on Google Home in local network.

This is Go version of [noelportugal/google-home-notifier](https://github.com/noelportugal/google-home-notifier)

## Install
```bash
$ go get -u github.com/ikasamah/homecast
```

## Usage
```golang
ctx := context.Background()
devices := homecast.LookupAndConnect(ctx)

for _, device := range devices {
    err := device.Speak(ctx, "Hello World", "en")
}
```

## Run example
```bash
$ go run $GOPATH/src/github.com/ikasamah/homecast/example/main.go
```


## Server erxample
```bash
$ go run $GOPATH/src/github.com/ikasamah/homecast/example/server.go
```
Then, access following URL in your browser.

http://localhost:8080/?text=Ciao&lang=it 


## Author
[Masayuki Hamasaki](https://github.com/ikasamah)
