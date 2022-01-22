package main

import (
	"bank_ledger_api/app"
	"bank_ledger_api/config"
)

func main() {
	config := config.GetConfig()
	app := &app.App{}
	app.Initialize(config)
	app.Run(":8000")
}
