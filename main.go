package main

import (
	"github.com/Clay294/forum/command"
	"github.com/spf13/cobra"
)

func main() {
	err := command.Execute()
	if err != nil {
		cobra.CheckErr(err)
	}
}
