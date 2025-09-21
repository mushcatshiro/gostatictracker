package main

import "github.com/mushcatshiro/gostatictracker/cli"

func main() {
	app := cli.NewCliApp()
	app.Execute()
}
