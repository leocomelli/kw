package main

import (
	"fmt"
	"os"

	"github.com/leocomelli/kube-wide/cmd"
)

func main() {
	cmd := cmd.NewCmdKubeWide(os.Stdin, os.Stdout, os.Stderr)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
