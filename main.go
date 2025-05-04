package main

import (
	"github.com/P001water/P1finger/cmd"
	_ "github.com/P001water/P1finger/cmd/P1finger/finger"
	_ "github.com/P001water/P1finger/cmd/P1finger/fofa"
	_ "github.com/P001water/P1finger/cmd/P1finger/rule"
)

func main() {
	cmd.Execute()
}
