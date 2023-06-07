package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vitorsalgado/mocha/v3/mhttp"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "print version info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(mhttp.Version)
		},
	}
}
