package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "print version info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("version")
		},
	}
}
