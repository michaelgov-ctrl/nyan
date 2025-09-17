package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/x/fyne/widget"
)

var errInvalidSpanFormat = errors.New("failed to parse lifespan: invalid span")

// sudo ~/go/bin/fyne-cross windows -arch=amd64 -app-id=nyan.cat -icon=rainbow.png
func main() {
	var lifespanFlag string
	flag.StringVar(&lifespanFlag, "lifespan", "30s", "Lifespan of pop up window, ex: 1s, 1m, 1h, 1d, 1w")

	flag.Parse()

	lifespan, err := parseLifespan(lifespanFlag)
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	go func() {
		time.Sleep(lifespan)
		os.Exit(0)
	}()

	a := app.New()
	w := a.NewWindow("nyan")

	animated, err := widget.NewAnimatedGifFromResource(resourceNyanGif)
	if err != nil {
		log.Fatal(err)
	}

	animated.Start()

	w.SetContent(container.NewStack(animated))
	w.Resize(fyne.NewSize(400, 400))
	w.ShowAndRun()
}

func parseLifespan(span string) (time.Duration, error) {
	nilDuration := time.Duration(0)

	n, t := splitOnLastChar(span)
	if t == "" {
		return nilDuration, errInvalidSpanFormat
	}

	num, err := strconv.Atoi(n)
	if err != nil || num <= 0 {
		return nilDuration, errInvalidSpanFormat
	}

	var shift time.Duration
	switch t {
	case "s":
		shift = time.Second
	case "m":
		shift = time.Minute
	case "h":
		shift = time.Hour
	case "d":
		shift = time.Hour * 24
	case "w":
		shift = time.Hour * 24 * 7
	default:
		return nilDuration, errInvalidSpanFormat
	}

	return time.Duration(num) * shift, nil
}

func splitOnLastChar(s string) (string, string) {
	if len(s) > 0 {
		return s[:len(s)-1], s[len(s)-1:]
	}

	return s, ""
}
