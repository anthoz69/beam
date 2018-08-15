package main

import (
	"log"

	"github.com/acoshift/hime"
	"github.com/anthoz69/beam/app"
)

func main() {
	appHime := hime.New()
	appFactory := app.New(appHime)

	err := appHime.
		Handler(appFactory).
		GracefulShutdown().
		Address(":8080").
		ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
