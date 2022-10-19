package cmd

import (
	"fmt"
	"os"

	"github.com/frederik-jatzkowski/hermes/admin"
	"github.com/frederik-jatzkowski/hermes/logs"
	"github.com/frederik-jatzkowski/hermes/params"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.Flags().StringVarP(
		&params.AdminHost,
		"admin",
		"a",
		"",
		"Sets the host name of the admin panel is running on",
	)
	rootCmd.Flags().StringVarP(
		&params.LogLevel,
		"logLevel",
		"l",
		"info",
		"Sets the minimum level for log output",
	)
	rootCmd.Flags().StringVarP(
		&params.EmailAdress,
		"email",
		"e",
		"",
		"Sets the initial email address for obtaining certificates",
	)
	rootCmd.MarkFlagRequired("email")
	rootCmd.Flags().StringVarP(
		&params.User,
		"user",
		"u",
		"",
		"Sets the hermes admin user",
	)
	rootCmd.MarkFlagRequired("user")
	rootCmd.Flags().StringVarP(
		&params.Password,
		"password",
		"p",
		"",
		"Sets the hermes admin password",
	)
	rootCmd.MarkFlagRequired("password")
}

var rootCmd = &cobra.Command{
	Use:     "hermes",
	Short:   "Hermes is a level 4 reverse proxy and load balancer",
	Version: params.Version,
	RunE: func(cmd *cobra.Command, args []string) error {
		var (
			err error
		)

		// setup logger
		logs.PrepareLogger(params.LogLevel)

		// setup admin panel
		admin.Start()

		return err
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
