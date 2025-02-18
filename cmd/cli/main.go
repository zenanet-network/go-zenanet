package main

import (
	"os"

	"github.com/zenanet-network/go-zenanet/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
