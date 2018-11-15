package main

import (
	"fmt"
	"os"

	consumer "github.com/harlow/kinesis-consumer"
	"github.com/spf13/pflag"
)

var (
	sourceStream string
)

func main() {
	pflag.Parse()

	if pflag.NFlag() == 0 {
		fmt.Printf("Usage: %s [options]\n", os.Args[0])
		pflag.PrintDefaults()
		os.Exit(1)

	}

	fmt.Printf("Reading from stream %s\n", sourceStream)
}

func init() {
	pflag.StringVarP(&sourceStream, "source", "s", "", "Source Stream")
}
