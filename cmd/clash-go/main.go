// Command clash-go is the Clash Go desktop app.
package main

import (
	"log"

	"github.com/charlzyx/clash-go/internal/app"
)

func main() {
	if err := app.New(app.Options{}).Run(); err != nil {
		log.Fatal(err)
	}
}
