// main.go
package main

import (
	"log"
	"mcp-go-tutorials/internal/app"
)

func main() {
	if err := app.NewAppCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
