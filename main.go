package main

import (
	"github.com/acrucetta/anki-builder/cmd/parser"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	parser.Execute()
}
