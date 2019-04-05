package main

import (
	"github.com/didil/volusnap/pkg/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func buildRootCmd() *cobra.Command {
	var p int

	rootCmd := &cobra.Command{
		Use:   "volusnapd",
		Short: "Cloud Volume Auto Snapshot Server",
		Run: func(cmd *cobra.Command, args []string) {
			logrus.Fatal(api.StartServer(p))
		},
	}

	rootCmd.Flags().IntVarP(&p, "port", "p", 8080, "Server Port")

	return rootCmd
}
