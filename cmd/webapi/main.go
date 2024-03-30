package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ffigari/stored-strings/internal/auth"
	"github.com/ffigari/stored-strings/internal/config"
	"github.com/ffigari/stored-strings/internal/personalwebsite"
)

func main() {
	config, err := config.Get()
	if err != nil {
		log.Fatal("[webapi] could not read config: ", err)
	}

	authenticator, err := auth.NewFromConfig()
	if err != nil {
		log.Fatal("[webapi] could not instantiate authenticator: ", err)
	}

	if m, err := personalwebsite.NewMux(
		context.Background(), "storedstrings", authenticator,
		config.WebPassword,
	); err != nil {
		log.Fatal(err)
	} else {
		log.Println("[webapi] about to start http server")
		if err := http.ListenAndServe(":3000", m); err != nil {
			log.Fatal(err)
		}
	}
}
