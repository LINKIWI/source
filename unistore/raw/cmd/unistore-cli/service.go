package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/timestamppb"

	"unistore/internal/config"
	"unistore/schemas/common"
	"unistore/schemas/service"
)

// runInfo prints server metadata information via the Meta service Info RPC.
func runInfo(cmd *cobra.Command, args []string) error {
	rdr := newRenderer(flagFormat, os.Stdout)

	store, rpc, err := newClientConfig(flagConfig, flagStore)
	if err != nil {
		return err
	}

	client, err := newUnistoreClient(store, rpc)
	if err != nil {
		return err
	}

	request := &service.InfoRequest{}

	if flagVerbose {
		rdr.proto(request)
	}

	response, err := client.Meta().Info(context.Background(), request)
	if err != nil {
		return err
	}

	return rdr.proto(response)
}

// runConfig prints parsed client configuration.
func runConfig(cmd *cobra.Command, args []string) error {
	cfg, err := config.New(flagConfig)
	if err != nil {
		return err
	}

	// Proto text format is not supported for the configuration file; use YAML instead.
	if flagFormat == outputFormatText || flagFormat == outputFormatDefault {
		flagFormat = outputFormatYAML
	}

	return newRenderer(flagFormat, os.Stdout).any(cfg.Client)
}

func runHeadBucket(cmd *cobra.Command, args []string) error {
	rdr := newRenderer(flagFormat, os.Stdout)
	bucket, _ := cmd.Flags().GetString("bucket")

	store, rpc, err := newClientConfig(flagConfig, flagStore)
	if err != nil {
		return err
	}

	client, err := newUnistoreClient(store, rpc)
	if err != nil {
		return err
	}

	request := &service.HeadBucketRequest{
		Resource: &common.Resource{
			Backend: common.Backend(common.Backend_value[strings.ToUpper(store.Backend)]),
			Bucket:  bucket,
		},
	}

	if flagVerbose {
		rdr.proto(request)
	}

	response, err := client.HeadBucket(context.Background(), request)
	if err != nil {
		return err
	}

	return rdr.proto(response)
}

// runHeadObject is a command line interface for the HeadObject RPC.
func runHeadObject(cmd *cobra.Command, args []string) error {
	rdr := newRenderer(flagFormat, os.Stdout)
	bucket, _ := cmd.Flags().GetString("bucket")

	store, rpc, err := newClientConfig(flagConfig, flagStore)
	if err != nil {
		return err
	}

	client, err := newUnistoreClient(store, rpc)
	if err != nil {
		return err
	}

	request := &service.HeadObjectRequest{
		Resource: &common.Resource{
			Backend: common.Backend(common.Backend_value[strings.ToUpper(store.Backend)]),
			Bucket:  bucket,
		},
		Key: args[0],
	}

	if flagVerbose {
		rdr.proto(request)
	}

	response, err := client.HeadObject(context.Background(), request)
	if err != nil {
		return err
	}

	return rdr.proto(response)
}

// runDownload is a command line interface for the GetObject and StreamGetObject RPCs.
func runDownload(cmd *cobra.Command, args []string) error {
	rdr := newRenderer(flagFormat, os.Stderr)
	bucket, _ := cmd.Flags().GetString("bucket")
	rangeStart, _ := cmd.Flags().GetUint64("range-start")
	rangeEnd, _ := cmd.Flags().GetUint64("range-end")
	stream, _ := cmd.Flags().GetBool("stream")
	chunkSize, _ := cmd.Flags().GetUint64("chunk-size")
	progress, _ := cmd.Flags().GetBool("progress")

	store, rpc, err := newClientConfig(flagConfig, flagStore)
	if err != nil {
		return err
	}

	client, err := newUnistoreClient(store, rpc)
	if err != nil {
		return err
	}

	get := &service.GetObjectRequest{
		Resource: &common.Resource{
			Backend: common.Backend(common.Backend_value[strings.ToUpper(store.Backend)]),
			Bucket:  bucket,
		},
		Key: args[0],
	}

	if rangeStart > 0 || rangeEnd > 0 {
		get.Range = &common.Range{
			Unit:  "bytes",
			Start: rangeStart,
			End:   rangeEnd,
		}
	}

	if !stream {
		if flagVerbose {
			rdr.proto(get)
		}

		response, err := client.GetObject(context.Background(), get)
		if err != nil {
			return err
		}

		if flagVerbose {
			rdr.proto(response)
		}

		_, err = os.Stdout.Write(response.Data)
		return err
	}

	request := &service.GetObjectStreamRequest{
		Request:   get,
		ChunkSize: chunkSize,
	}

	if flagVerbose {
		rdr.proto(request)
	}

	response, err := client.StreamGetObject(context.Background(), request)
	if err != nil {
		return err
	}

	start := time.Now()
	downloaded := new(uint64)

	if progress {
		interval := 5 * time.Second
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		go func() {
			prev := atomic.LoadUint64(downloaded)

			for range ticker.C {
				current := atomic.LoadUint64(downloaded)
				fmt.Fprintf(
					os.Stderr,
					"%v: downloaded %s (%s/s)\n",
					time.Since(start).Round(interval),
					humanize.Bytes(current),
					humanize.Bytes((current-prev)/5),
				)
				prev = current
			}
		}()
	}

	for {
		message, err := response.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		atomic.AddUint64(downloaded, uint64(len(message.Response.Data)))

		if flagVerbose {
			rdr.proto(message)
		}

		if _, err = os.Stdout.Write(message.Response.Data); err != nil {
			return err
		}
	}

	if progress {
		fmt.Fprintf(
			os.Stderr,
			"%v: completed download of %s\n",
			time.Since(start),
			humanize.Bytes(atomic.LoadUint64(downloaded)),
		)
	}

	return nil
}

// runUpload is a command line interface for the PutObject and StreamPutObject RPCs.
func runUpload(cmd *cobra.Command, args []string) error {
	rdr := newRenderer(flagFormat, os.Stdout)
	bucket, _ := cmd.Flags().GetString("bucket")
	timestamp, _ := cmd.Flags().GetUint64("timestamp")
	stream, _ := cmd.Flags().GetBool("stream")
	chunkSize, _ := cmd.Flags().GetUint64("chunk-size")
	checksum, _ := cmd.Flags().GetString("checksum")
	progress, _ := cmd.Flags().GetBool("progress")

	store, rpc, err := newClientConfig(flagConfig, flagStore)
	if err != nil {
		return err
	}

	client, err := newUnistoreClient(store, rpc)
	if err != nil {
		return err
	}

	attributes := &common.Attributes{}
	if timestamp > 0 {
		attributes.Timestamp = timestamppb.New(time.Unix(0, int64(timestamp)))
	}

	put := &service.PutObjectRequest{
		Resource: &common.Resource{
			Backend: common.Backend(common.Backend_value[strings.ToUpper(store.Backend)]),
			Bucket:  bucket,
		},
		Key:        args[0],
		Checksum:   checksum,
		Attributes: attributes,
	}

	if !stream {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		put.Data = data

		if flagVerbose {
			rdr.proto(put)
		}

		response, err := client.PutObject(context.Background(), put)
		if err != nil {
			return err
		}

		return rdr.proto(response)
	}

	request, err := client.StreamPutObject(context.Background())
	if err != nil {
		return err
	}

	start := time.Now()
	uploaded := new(uint64)

	if progress {
		interval := 5 * time.Second
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		go func() {
			prev := atomic.LoadUint64(uploaded)

			for range ticker.C {
				current := atomic.LoadUint64(uploaded)
				fmt.Printf(
					"%v: uploaded %s (%s/s)\n",
					time.Since(start).Round(interval),
					humanize.Bytes(current),
					humanize.Bytes((current-prev)/5),
				)
				prev = current
			}
		}()
	}

	for {
		chunk := make([]byte, chunkSize)
		n, err := os.Stdin.Read(chunk)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		message := &service.PutObjectStreamRequest{Request: put}
		message.Request.Data = chunk[:n]

		atomic.AddUint64(uploaded, uint64(n))

		if flagVerbose {
			rdr.proto(message)
		}

		if err := request.Send(message); err != nil {
			break
		}
	}

	response, err := request.CloseAndRecv()
	if err != nil {
		return err
	}

	if progress {
		fmt.Printf(
			"%v: completed upload of %s\n",
			time.Since(start),
			humanize.Bytes(atomic.LoadUint64(uploaded)),
		)
	}

	return rdr.proto(response)
}

// runDelete is a command line interface for the DeleteObject RPC.
func runDelete(cmd *cobra.Command, args []string) error {
	rdr := newRenderer(flagFormat, os.Stdout)
	bucket, _ := cmd.Flags().GetString("bucket")
	recursive, _ := cmd.Flags().GetBool("recursive")

	store, rpc, err := newClientConfig(flagConfig, flagStore)
	if err != nil {
		return err
	}

	client, err := newUnistoreClient(store, rpc)
	if err != nil {
		return err
	}

	request := &service.DeleteObjectRequest{
		Resource: &common.Resource{
			Backend: common.Backend(common.Backend_value[strings.ToUpper(store.Backend)]),
			Bucket:  bucket,
		},
		Key:       args[0],
		Recursive: recursive,
	}

	if flagVerbose {
		rdr.proto(request)
	}

	response, err := client.DeleteObject(context.Background(), request)
	if err != nil {
		return err
	}

	return rdr.proto(response)
}

// runListBuckets is a command line interface for the ListObjects RPC.
func runListBuckets(cmd *cobra.Command, args []string) error {
	rdr := newRenderer(flagFormat, os.Stdout)

	store, rpc, err := newClientConfig(flagConfig, flagStore)
	if err != nil {
		return err
	}

	client, err := newUnistoreClient(store, rpc)
	if err != nil {
		return err
	}

	request := &service.ListBucketsRequest{
		Resource: &common.Resource{
			Backend: common.Backend(common.Backend_value[strings.ToUpper(store.Backend)]),
		},
	}

	if flagVerbose {
		rdr.proto(request)
	}

	response, err := client.ListBuckets(context.Background(), request)
	if err != nil {
		return err
	}

	return rdr.proto(response)
}

// runListObjects is a command line interface for the ListObjects RPC.
func runListObjects(cmd *cobra.Command, args []string) error {
	rdr := newRenderer(flagFormat, os.Stdout)
	bucket, _ := cmd.Flags().GetString("bucket")
	recursive, _ := cmd.Flags().GetBool("recursive")

	prefix := ""
	if len(args) > 0 {
		prefix = args[0]
	}

	store, rpc, err := newClientConfig(flagConfig, flagStore)
	if err != nil {
		return err
	}

	client, err := newUnistoreClient(store, rpc)
	if err != nil {
		return err
	}

	request := &service.ListObjectsRequest{
		Resource: &common.Resource{
			Backend: common.Backend(common.Backend_value[strings.ToUpper(store.Backend)]),
			Bucket:  bucket,
		},
		Prefix:    prefix,
		Recursive: recursive,
	}

	if flagVerbose {
		rdr.proto(request)
	}

	response, err := client.ListObjects(context.Background(), request)
	if err != nil {
		return err
	}

	return rdr.proto(response)
}
