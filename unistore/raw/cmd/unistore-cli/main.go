package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"unistore/internal/meta"
)

var (
	flagFormat  string
	flagConfig  string
	flagStore   string
	flagTimeout time.Duration
	flagVerbose bool
)

var (
	rootCmd = &cobra.Command{
		Use:               "uni",
		Short:             "Command line client for Unistore",
		Long:              "Command line interfaces for create, read, update, and deletion of objects in a remote Unistore gRPC server",
		Version:           meta.Version,
		PersistentPreRunE: preRunGlobalTimeout,
	}
	infoCmd = &cobra.Command{
		Use:   "info",
		Short: "Show server metadata information",
		Long:  "Show server metadata information, including the build version and configured backend composition tree.",
		Example: strings.Join([]string{
			"Show server metadata information for a configured store:",
			"  $ uni info --store default",
			"Show server metadata information in a machine-readable JSON format:",
			"  $ uni info --store default --format json",
		}, "\n"),
		Version: rootCmd.Version,
		Args:    cobra.NoArgs,
		RunE:    runInfo,
	}
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Show parsed client configuration",
		Long: "Show structured client configuration parsed from the file indicated by the --config flag\n" +
			"or UNI_CONFIG environment variable.",
		Example: strings.Join([]string{
			"Show parsed client configuration from default paths:",
			"  $ uni config",
			"Show parsed client configuration for a specific file:",
			"  $ uni config --config config.yaml",
			"Display parsed configuration as JSON:",
			"  $ uni config --format json",
		}, "\n"),
		Version: rootCmd.Version,
		Args:    cobra.NoArgs,
		RunE:    runConfig,
	}
	headBucketCmd = &cobra.Command{
		Use:     "head-bucket",
		Aliases: []string{"headb", "hb"},
		Short:   "Retrieve metadata for a bucket",
		Long:    "Retrieve metadata for a bucket by name.\nStructured data is printed to standard output.",
		Example: strings.Join([]string{
			"Head a bucket by name:",
			"  $ uni head-bucket --store default --bucket bucket",
		}, "\n"),
		Version: rootCmd.Version,
		Args:    cobra.NoArgs,
		RunE:    runHeadBucket,
	}
	headObjectCmd = &cobra.Command{
		Use:     "head-object <key>",
		Aliases: []string{"head", "ho", "stat"},
		Short:   "Retrieve metadata for an object",
		Long:    "Retrieve metadata for an object by key.\nStructured data is printed to standard output.",
		Example: strings.Join([]string{
			"Head an object by key:",
			"  $ uni head --store default --bucket bucket key",
		}, "\n"),
		Version: rootCmd.Version,
		Args:    cobra.ExactArgs(1),
		RunE:    runHeadObject,
	}
	downloadCmd = &cobra.Command{
		Use:     "download <key>",
		Aliases: []string{"cat", "dl"},
		Short:   "Download an object",
		Long: "Download an object by key, either in its entirety via a unary call or in chunks via a streaming response.\n" +
			"Object data is written to standard output, while verbose logs (when enabled) are written to standard error.",
		Example: strings.Join([]string{
			"Print the contents of an object to stdout:",
			"  $ uni download --store default --bucket bucket key",
			"Download an object to a file on disk:",
			"  $ uni download --store default --bucket bucket key > key",
			"Download an object using streaming with the default chunk size:",
			"  $ uni download --store default --bucket bucket --stream key",
			"Download a partial byte range chunk of an object:",
			"  $ uni download --store default --bucket bucket --range-start 10 --range-end 100 key",
		}, "\n"),
		Version: rootCmd.Version,
		Args:    cobra.ExactArgs(1),
		RunE:    runDownload,
	}
	uploadCmd = &cobra.Command{
		Use:     "upload <key>",
		Aliases: []string{"up"},
		Short:   "Upload an object",
		Long: "Upload an object with the specified key, either in its entirety via a unary call or in chunks via a streaming request.\n" +
			"Object data is read from standard input. Verbose logs (when enabled) are written to standard output.",
		Example: strings.Join([]string{
			"Upload an object from stdin:",
			"  $ cat file | uni upload --store default --bucket bucket key",
			"Upload an object using shell redirection:",
			"  $ uni upload --store default --bucket bucket key < file",
			"Upload an object from stdin using streaming with the default chunk size:",
			"  $ cat file | uni upload --store default --bucket bucket --stream key",
			"Upload an object, preserving its original modification time in associated metadata:",
			"  $ uni upload --store default --bucket bucket --timestamp $(stat --format=\"%.9Y\" file | sed 's/\\.//g') key < file",
		}, "\n"),
		Version: rootCmd.Version,
		Args:    cobra.ExactArgs(1),
		RunE:    runUpload,
	}
	deleteCmd = &cobra.Command{
		Use:     "delete <key>",
		Aliases: []string{"del", "rm"},
		Short:   "Delete an object",
		Long:    "Delete an object by key.\nThe response status is printed to standard output.",
		Example: strings.Join([]string{
			"Delete an object by key:",
			"  $ uni delete --store default --bucket bucket key",
		}, "\n"),
		Version: rootCmd.Version,
		Args:    cobra.ExactArgs(1),
		RunE:    runDelete,
	}
	listBucketsCmd = &cobra.Command{
		Use:     "list-buckets",
		Aliases: []string{"listb", "lb", "ls-buckets"},
		Short:   "List all buckets",
		Long:    "List all buckets. Structured data is printed to standard output.",
		Example: strings.Join([]string{
			"List all buckets:",
			"  $ uni list-buckets --store default",
		}, "\n"),
		Version: rootCmd.Version,
		Args:    cobra.NoArgs,
		RunE:    runListBuckets,
	}
	listObjectsCmd = &cobra.Command{
		Use:     "list-objects <prefix>",
		Aliases: []string{"list", "lo", "ls", "ls-objects"},
		Short:   "List objects by prefix",
		Long: "List objects by prefix, recursively through any hierarchical delimiters (e.g. directories).\n" +
			"Structured data is printed to standard output.",
		Example: strings.Join([]string{
			"List all objects for a bucket:",
			"  $ uni list --store default --bucket bucket",
			"List objects with a key prefix literal:",
			"  $ uni list --store default --bucket bucket prefix",
			"List objects in a machine-readable JSON format:",
			"  $ uni list --store default --bucket bucket --format json | jq .",
		}, "\n"),
		Version: rootCmd.Version,
		Args:    cobra.MaximumNArgs(1),
		RunE:    runListObjects,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&flagFormat,
		"format",
		"f",
		os.Getenv("UNI_FORMAT"),
		fmt.Sprintf(
			"output render format; choose from %v",
			[]string{outputFormatJSON, outputFormatText, outputFormatHuman},
		),
	)
	rootCmd.PersistentFlags().StringVarP(
		&flagConfig,
		"config",
		"c",
		os.Getenv("UNI_CONFIG"),
		"path to the configuration file",
	)
	rootCmd.PersistentFlags().StringVarP(
		&flagStore,
		"store",
		"s",
		os.Getenv("UNI_STORE"),
		"server store alias in configuration",
	)
	rootCmd.PersistentFlags().DurationVarP(
		&flagTimeout,
		"timeout",
		"w",
		0,
		"global timeout for all operations",
	)
	rootCmd.PersistentFlags().BoolVarP(
		&flagVerbose,
		"verbose",
		"v",
		false,
		"enable verbose output",
	)

	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(configCmd)

	headBucketCmd.Flags().StringP("bucket", "b", "", "bucket name")
	headBucketCmd.MarkFlagRequired("bucket")
	rootCmd.AddCommand(headBucketCmd)

	headObjectCmd.Flags().StringP("bucket", "b", "", "bucket name")
	headObjectCmd.MarkFlagRequired("bucket")
	rootCmd.AddCommand(headObjectCmd)

	downloadCmd.Flags().StringP("bucket", "b", "", "bucket name")
	downloadCmd.Flags().Uint64P("range-start", "m", 0, "inclusive byte range start index for partial downloads")
	downloadCmd.Flags().Uint64P("range-end", "n", 0, "exclusive byte range end index for partial downloads")
	downloadCmd.Flags().BoolP("stream", "t", false, "stream object contents")
	downloadCmd.Flags().Uint64P("chunk-size", "r", 10*1024, "stream response chunk size")
	downloadCmd.Flags().BoolP("progress", "p", false, "periodically print downloaded bytes and transfer rate (only for stream RPC)")
	downloadCmd.MarkFlagRequired("bucket")
	rootCmd.AddCommand(downloadCmd)

	uploadCmd.Flags().StringP("bucket", "b", "", "bucket name")
	uploadCmd.Flags().Uint64P("timestamp", "m", 0, "object timestamp override (Unix timestamp, nanoseconds)")
	uploadCmd.Flags().BoolP("stream", "t", false, "stream object contents")
	uploadCmd.Flags().Uint64P("chunk-size", "r", 10*1024, "stream request chunk size")
	uploadCmd.Flags().StringP("checksum", "k", "", "object contents checksum (only for unary RPC)")
	uploadCmd.Flags().BoolP("progress", "p", false, "periodically print uploaded bytes and transfer rate (only for stream RPC)")
	uploadCmd.MarkFlagRequired("bucket")
	rootCmd.AddCommand(uploadCmd)

	deleteCmd.Flags().StringP("bucket", "b", "", "bucket name")
	deleteCmd.Flags().BoolP("recursive", "r", false, "delete the object tree recursively")
	deleteCmd.MarkFlagRequired("bucket")
	rootCmd.AddCommand(deleteCmd)

	rootCmd.AddCommand(listBucketsCmd)

	listObjectsCmd.Flags().StringP("bucket", "b", "", "bucket name")
	listObjectsCmd.Flags().BoolP("recursive", "r", false, "traverse key delimiters recursively")
	listObjectsCmd.MarkFlagRequired("bucket")
	rootCmd.AddCommand(listObjectsCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
