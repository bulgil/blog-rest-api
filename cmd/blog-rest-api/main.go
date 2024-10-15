package main

import (
	"github.com/bulgil/blog-rest-api/internal/app"
)

func main() {
	application := app.NewApp()
	if err := application.Run(); err != nil {
		application.Logger.Fatal(err)
	}
}
