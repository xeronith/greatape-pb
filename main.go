package main

import (
	"greateape-pb/hooks"
	"log"

	"github.com/pocketbase/pocketbase"
)

func main() {
	app := pocketbase.New()

	hooks.Initialize(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
