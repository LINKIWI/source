package main

import (
	"os"

	"github.com/spf13/cobra"

	"unistore/internal/meta"
)

var (
	rootCmd = &cobra.Command{
		Use:               "unistore",
		Short:             "Unistore object storage server",
		PersistentPreRunE: validateCmd,
		Version:           meta.Version,
	}
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start the Unistore gRPC server",
		Long: "Start the Unistore gRPC server at the address specified in configuration. " +
			"A valid configuration file is required.",
		Version: rootCmd.Version,
		Args:    cobra.NoArgs,
		RunE:    runServer,
	}
	configCmd = &cobra.Command{
		Use:     "config",
		Short:   "Configuration commands",
		Long:    "Utilities for reading and validating Unistore server configuration files",
		Version: rootCmd.Version,
		Args:    cobra.ExactArgs(1),
	}
	configShowCmd = &cobra.Command{
		Use:   "show",
		Short: "Show parsed configuration",
		Long: "Read and parse the specified configuration file, and print its structured " +
			"contents to standard output.",
		Version: rootCmd.Version,
		Args:    cobra.NoArgs,
		RunE:    runConfigShow,
	}
	configValidateCmd = &cobra.Command{
		Use:     "validate",
		Short:   "Validate specified configuration",
		Long:    "Check the provided configuration file for errors.",
		Version: rootCmd.Version,
		Args:    cobra.NoArgs,
		RunE:    runConfigValidate,
	}
)

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", os.Getenv("UNISTORE_CONFIG"), "path to the config file")

	rootCmd.AddCommand(serverCmd)

	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configValidateCmd)
	rootCmd.AddCommand(configCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
