package main

import (
	"os"

	"main/cli"
	"main/log"
)

func main() {
	app := cli.NewApp()

	if err := app.Run(os.Args); err != nil {
		log.Log.Fatal().Msg(err.Error())
	}
}
