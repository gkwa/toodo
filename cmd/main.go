package main

import (
	"os"

	"github.com/taylormonacelli/toodo"
)

func main() {
	code := toodo.Execute()
	os.Exit(code)
}
