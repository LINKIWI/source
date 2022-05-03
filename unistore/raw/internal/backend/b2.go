package backend

import (
	"bytes"
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"unistore/internal/config"
	"unistore/internal/meta"
	"unistore/internal/schemas"
	"unistore/internal/util"
	"unistore/pkg/client/b2"
	"unistore/schemas/common"
	"unistore/schemas/service"
)

const (
	// authenticationRefreshInterval is the duration of time to consider an issued API
	// authorization token to be valid; e.g., after this amount of time has elapsed since the
	// most recently issued authorization token, the authorization API should be invoked again
	// to obtain a new set of credentials.
	authenticationRefreshInterval = 3 * time.Hour
	// uploadRetryBudget is the number of upload attempts that will be made for a single upload
	// request, for both full uploads and chunked (part) uploads.
	uploadRetryBudget = 5
	// uploadRetryInitialBackoff is the initial delay used by the client before retrying a
	// failed upload request. Retries follow an exponential backoff policy for each attempt.
	uploadRetryInitialBackoff = 10 * time.Second
)

// authenticationState is a container for current B2 API authentication state.
type authenticationState struct {
	accountID          string
	authorizationToken string
	apiURL             *url.URL
	downloadURL        *url.URL
	refresh            time.Time
}

// B2 is a Backend that implements an object storage server on top of B2.
type B2 struct {
	client           *b2.Client
	applicationKeyID string
	applicationKey   string

	auth      *authenticationState
	authMutex sync.Mutex
}

// NewB2 creates a new B2 backend.
func NewB2(cfg *config.B2) (Backend, error) {
	var opts []b2.Option

	dialer := util.NewRetryDialer(
		net.Dialer{Timeout: cfg.Connection.ConnectTimeout},
		cfg.Connection.ConnectAttempts,
		cfg.Connection.ConnectTimeout/2,
	)
	transport := &http.Transport{
		DialContext:           dialer.DialContext,
		TLSHandshakeTimeout:   cfg.Connection.ConnectTimeout,
		ResponseHeaderTimeout: cfg.Connection.RequestTimeout,
		ForceAttemptHTTP2:     true,
	}

	opts = append(opts, b2.WithTransport(transport))
	opts = append(opts, b2.WithLogger(zap.L()))
	if cfg.Connection.Identity != "" {
		clientID := fmt.Sprintf(
			"unistore/%s (instance:%s)",
			meta.Version,
			cfg.Connection.Identity,
		)

		opts = append(opts, b2.WithClientID(clientID))
	}

	client, err := b2.New(opts...)
	if err != nil {
		return nil, err
	}

	b := &B2{
		client:           client,
		applicationKeyID: cfg.ApplicationKeyID,
		applicationKey:   cfg.ApplicationKey,
		auth:             new(authenticationState),
	}

	return newInstrumentation("b2", b), nil
}

// HeadBucket locates the requested bucket by name followed by listing its metadata.
func (b *B2) HeadBucket(ctx context.Context, request *service.HeadBucketRequest) (*service.HeadBucketResponse, error) {
	var visibility common.Visibility

	auth, err := b.authorizeAccount(ctx)
	if err != nil {
		return nil, b.createError(err)
	}

	lbReq := &b2.ListBucketsRequest{
		Authorization: auth.authorizationToken,
		BaseURL:       auth.apiURL,
		AccountID:     auth.accountID,
		BucketName:    request.Resource.Bucket,
	}

	lbResp, err := b.client.ListBuckets(ctx, lbReq)
	if err != nil {
		return nil, b.createError(err)
	}

	if len(lbResp.Buckets) != 1 {
		return nil, status.Errorf(
			codes.NotFound,
			"backend: unexpected number of buckets matched by name: name=%s buckets=%d",
			request.Resource.Bucket,
			len(lbResp.Buckets),
		)
	}

	switch lbResp.Buckets[0].BucketType {
	case "allPublic":
		visibility = common.Visibility_PUBLIC
	case "allPrivate", "snapshot":
		visibility = common.Visibility_PRIVATE
	default:
		visibility = common.Visibility_DEFAULT
	}

	lifecyclePolicies := make([]*common.LifecyclePolicy, len(lbResp.Buckets[0].LifecycleRules))

	for i, rule := range lbResp.Buckets[0].LifecycleRules {
		// B2 only supports lifecycle policies based on key prefix; adapt these semantics
		// into a regular expression pattern.
		pattern := fmt.Sprintf("^%s", regexp.QuoteMeta(rule.FileNamePrefix))
		ttl := time.Duration(rule.DaysFromUploadingToHiding+rule.DaysFromHidingToDeleting) * 24 * time.Hour

		lifecyclePolicies[i] = &common.LifecyclePolicy{
			Pattern:    pattern,
			Expiration: durationpb.New(ttl),
		}
	}

	return &service.HeadBucketResponse{
		Bucket: &common.Bucket{
			Id:                lbResp.Buckets[0].BucketID,
			Name:              lbResp.Buckets[0].BucketName,
			Visibility:        visibility,
			LifecyclePolicies: lifecyclePolicies,
			Tags:              lbResp.Buckets[0].BucketInfo,
		},
	}, nil
}

// HeadObject locates the requested file ID followed by listing its metadata.
func (b *B2) HeadObject(ctx context.Context, request *service.HeadObjectRequest) (*service.HeadObjectResponse, error) {
	auth, err := b.authorizeAccount(ctx)
	if err != nil {
		return nil, b.createError(err)
	}

	lbReq := &b2.ListBucketsRequest{
		Authorization: auth.authorizationToken,
		BaseURL:       auth.apiURL,
		AccountID:     auth.accountID,
		BucketName:    request.Resource.Bucket,
	}

	lbResp, err := b.client.ListBuckets(ctx, lbReq)
	if err != nil {
		return nil, b.createError(err)
	}

	if len(lbResp.Buckets) != 1 {
		return nil, status.Errorf(
			codes.NotFound,
			"backend: unexpected number of buckets matched by name: name=%s buckets=%d",
			request.Resource.Bucket,
			len(lbResp.Buckets),
		)
	}

	lfReq := &b2.ListFileNamesRequest{
		Authorization: auth.authorizationToken,
		BaseURL:       auth.apiURL,
		BucketID:      lbResp.Buckets[0].BucketID,
		StartFileName: request.Key,
		MaxFileCount:  1,
		Prefix:        request.Key,
	}

	lfResp, err := b.client.ListFileNames(ctx, lfReq)
	if err != nil {
		return nil, b.createError(err)
	}

	if len(lfResp.Files) != 1 || lfResp.Files[0].FileName != request.Key {
		return nil, status.Errorf(
			codes.NotFound,
			"backend: no files match requested name: bucket=%s name=%s",
			request.Resource.Bucket,
			request.Key,
		)
	}

	attributes, err := schemas.UnmarshalAttributes(lfResp.Files[0].FileInfo, nil)
	if err != nil {
		return nil, b.createError(err)
	}

	return &service.HeadObjectResponse{
		Metadata: &common.Metadata{
			Id:           lfResp.Files[0].FileID,
			Key:          lfResp.Files[0].FileName,
			Size:         uint64(lfResp.Files[0].ContentLength),
			ObjectType:   common.Object_REGULAR,
			MimeType:     lfResp.Files[0].ContentType,
			ModifiedTime: timestamppb.New(lfResp.Files[0].UploadTimestamp.Time),
			Checksum:     lfResp.Files[0].ContentSHA1,
			Attributes:   attributes,
		},
	}, nil
}

// GetObject downloads the requested file and buffers its full contents in memory for response.
func (b *B2) GetObject(ctx context.Context, request *service.GetObjectRequest) (*service.GetObjectResponse, error) {
	auth, err := b.authorizeAccount(ctx)
	if err != nil {
		return nil, b.createError(err)
	}

	dfReq := &b2.DownloadFileByNameRequest{
		Authorization: auth.authorizationToken,
		DownloadURL:   auth.downloadURL,
		Bucket:        request.Resource.Bucket,
		File:          request.Key,
	}

	if request.Range != nil {
		dfReq.Range = b2.Range{
			Unit:  request.Range.Unit,
			Start: int(request.Range.Start),
		}

		if request.Range.End > 0 {
			// Note that the B2 API considers the range end index to be inclusive, while
			// the Unistore API considers the end index to be exclusive. Normalize this
			// behavior by subtracting 1 from the desired end index for compatibility.
			dfReq.Range.End = int(request.Range.End) - 1
		}
	}

	dfResp, err := b.client.DownloadFileByName(ctx, dfReq)
	if err != nil {
		return nil, b.createError(err)
	}

	defer dfResp.Data.Close()

	data, err := io.ReadAll(dfResp.Data)
	if err != nil {
		return nil, b.createError(err)
	}

	attributes, err := schemas.UnmarshalAttributes(dfResp.FileInfo, nil)
	if err != nil {
		return nil, b.createError(err)
	}

	return &service.GetObjectResponse{
		Metadata: &common.Metadata{
			Id:           dfResp.FileID,
			Key:          dfResp.FileName,
			Size:         uint64(dfResp.Size),
			ObjectType:   common.Object_REGULAR,
			MimeType:     dfResp.ContentType,
			ModifiedTime: timestamppb.New(dfResp.UploadTimestamp.Time),
			Checksum:     dfResp.ContentSHA1,
			Attributes:   attributes,
		},
		Data: data,
	}, nil
}

// StreamGetObject is a streaming implementation of GetObject. It segments the requests to B2 in
// chunks of the requested size using the Range request directive, and pipes chunked responses into
// a stream channel.
func (b *B2) StreamGetObject(ctx context.Context, request *service.GetObjectStreamRequest) (<-chan *service.GetObjectStreamResponse, chan error) {
	stream := make(chan *service.GetObjectStreamResponse)
	errs := make(chan error, 1)

	if request.ChunkSize == 0 {
		errs <- status.Error(
			codes.InvalidArgument,
			"backend: B2 streaming get chunk size must be greater than zero",
		)
		return nil, errs
	}

	auth, err := b.authorizeAccount(ctx)
	if err != nil {
		errs <- b.createError(err)
		return nil, errs
	}

	go func() {
		var offset uint64

		if request.Request.Range != nil {
			offset = request.Request.Range.Start
		}

		for {
			// Read the requested chunk size, but only up to the end index specified by
			// the range, when available.
			end := offset + request.ChunkSize
			if request.Request.Range != nil &&
				request.Request.Range.End > 0 &&
				end > request.Request.Range.End {
				end = request.Request.Range.End
			}

			dfReq := &b2.DownloadFileByNameRequest{
				Authorization: auth.authorizationToken,
				DownloadURL:   auth.downloadURL,
				Bucket:        request.Request.Resource.Bucket,
				File:          request.Request.Key,
				// Note that the B2 API considers the range end index to be
				// inclusive. To normalize this behavior against the Unistore API,
				// add/subtract 1 from the end index as needed for non-inclusive
				// end index semantics.
				Range: b2.Range{
					Unit:  "bytes",
					Start: int(offset),
					End:   int(end) - 1,
				},
			}

			dfResp, err := b.client.DownloadFileByName(ctx, dfReq)
			if err != nil {
				errs <- b.createError(err)
				break
			}

			chunk, err := io.ReadAll(dfResp.Data)
			if err != nil {
				dfResp.Data.Close()
				errs <- b.createError(err)
				break
			}

			dfResp.Data.Close()

			attributes, err := schemas.UnmarshalAttributes(dfResp.FileInfo, nil)
			if err != nil {
				errs <- b.createError(err)
				break
			}

			stream <- &service.GetObjectStreamResponse{
				Response: &service.GetObjectResponse{
					Metadata: &common.Metadata{
						Id:           dfResp.FileID,
						Key:          dfResp.FileName,
						Size:         uint64(dfResp.Size),
						ObjectType:   common.Object_REGULAR,
						MimeType:     dfResp.ContentType,
						ModifiedTime: timestamppb.New(dfResp.UploadTimestamp.Time),
						Checksum:     dfResp.ContentSHA1,
						Attributes:   attributes,
					},
					Data: chunk,
				},
				Range: &common.Range{
					Unit:  "bytes",
					Start: uint64(dfResp.Range.Start),
					End:   uint64(dfResp.Range.End) + 1,
					Total: uint64(dfResp.Range.Total),
				},
			}

			offset += uint64(dfResp.Range.End - dfResp.Range.Start + 1)

			if offset >= uint64(dfResp.Range.Total) {
				break
			}

			if request.Request.Range != nil &&
				request.Request.Range.End > 0 &&
				offset >= request.Request.Range.End {
				break
			}
		}

		close(stream)
	}()

	return stream, errs
}

// PutObject uploads the request file in full.
func (b *B2) PutObject(ctx context.Context, request *service.PutObjectRequest) (*service.PutObjectResponse, error) {
	retryAttempt := 0
	retryDelay := uploadRetryInitialBackoff

	auth, err := b.authorizeAccount(ctx)
	if err != nil {
		return nil, b.createError(err)
	}

	lbReq := &b2.ListBucketsRequest{
		Authorization: auth.authorizationToken,
		BaseURL:       auth.apiURL,
		AccountID:     auth.accountID,
		BucketName:    request.Resource.Bucket,
	}

	lbResp, err := b.client.ListBuckets(ctx, lbReq)
	if err != nil {
		return nil, b.createError(err)
	}

	if len(lbResp.Buckets) != 1 {
		return nil, status.Errorf(
			codes.NotFound,
			"backend: unexpected number of buckets matched by name: name=%s buckets=%d",
			request.Resource.Bucket,
			len(lbResp.Buckets),
		)
	}

	for {
		uuReq := &b2.GetUploadURLRequest{
			Authorization: auth.authorizationToken,
			BaseURL:       auth.apiURL,
			BucketID:      lbResp.Buckets[0].BucketID,
		}

		uuResp, err := b.client.GetUploadURL(ctx, uuReq)
		if err != nil {
			return nil, b.createError(err)
		}

		attributes, err := schemas.MarshalAttributes(request.Attributes, nil)
		if err != nil {
			return nil, b.createError(err)
		}

		hash := sha1.New()
		hash.Write(request.Data)
		checksum := fmt.Sprintf("%x", hash.Sum(nil))

		ufReq := &b2.UploadFileRequest{
			Authorization: uuResp.AuthorizationToken,
			UploadURL:     uuResp.UploadURL.URL,
			FileName:      request.Key,
			FileInfo:      attributes,
			ContentType:   "b2/x-auto",
			ContentLength: len(request.Data),
			ContentSHA1:   checksum,
			Data:          io.NopCloser(bytes.NewReader(request.Data)),
		}

		ufResp, err := b.client.UploadFile(ctx, ufReq)
		if err != nil {
			if e, ok := err.(*b2.Status); ok && e.Status >= http.StatusInternalServerError && retryAttempt < uploadRetryBudget {
				zap.L().Warn(
					"full file upload failed with server error; retrying after delay",
					zap.Int("attempt", retryAttempt),
					zap.Duration("delay", retryDelay),
					zap.Error(err),
				)

				select {
				case <-ctx.Done():
					return nil, status.FromContextError(ctx.Err()).Err()
				case <-time.After(retryDelay):
				}

				retryAttempt++
				retryDelay *= 2 // Exponential backoff
				continue
			}

			return nil, b.createError(err)
		}

		return &service.PutObjectResponse{
			Metadata: &common.Metadata{
				Id:           ufResp.FileID,
				Key:          ufResp.FileName,
				Size:         uint64(ufResp.ContentLength),
				ObjectType:   common.Object_REGULAR,
				MimeType:     ufResp.ContentType,
				ModifiedTime: timestamppb.New(ufResp.UploadTimestamp.Time),
				Checksum:     ufResp.ContentSHA1,
				Attributes:   request.Attributes,
			},
		}, nil
	}
}

// StreamPutObject is a streaming implementation of PutObject. It identifies the target file, starts
// a large file, authorizes a chunked upload URL, uploads each incoming chunk in a single request,
// and finishes the file when all chunks have been uploaded.
func (b *B2) StreamPutObject(ctx context.Context, stream chan *service.PutObjectStreamRequest) (*service.PutObjectStreamResponse, error) {
	var authorization string
	var uploadURL *url.URL
	var fileID string
	var checksums []string

	part := 1

	auth, err := b.authorizeAccount(ctx)
	if err != nil {
		return nil, b.createError(err)
	}

	for request := range stream {
		retryAttempt := 0
		retryDelay := uploadRetryInitialBackoff

		if request.Range != nil {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"backend: B2 does not support range specifications for streaming puts: range=%v",
				request.Range,
			)
		}

		if fileID == "" {
			lbReq := &b2.ListBucketsRequest{
				Authorization: auth.authorizationToken,
				BaseURL:       auth.apiURL,
				AccountID:     auth.accountID,
				BucketName:    request.Request.Resource.Bucket,
			}

			lbResp, err := b.client.ListBuckets(ctx, lbReq)
			if err != nil {
				return nil, b.createError(err)
			}

			if len(lbResp.Buckets) != 1 {
				return nil, status.Errorf(
					codes.NotFound,
					"backend: unexpected number of buckets matched by name: name=%s buckets=%d",
					request.Request.Resource.Bucket,
					len(lbResp.Buckets),
				)
			}

			// The B2 API only supports writing file info when the large file is
			// started, but at this point, the total size of the file is not yet known.
			// Explicitly override its value to 0 to indicate that the original size is
			// not known.
			request.Request.Attributes.Size = 0
			attributes, err := schemas.MarshalAttributes(request.Request.Attributes, nil)
			if err != nil {
				return nil, b.createError(err)
			}

			sfReq := &b2.StartLargeFileRequest{
				Authorization: auth.authorizationToken,
				BaseURL:       auth.apiURL,
				BucketID:      lbResp.Buckets[0].BucketID,
				FileName:      request.Request.Key,
				FileInfo:      attributes,
				ContentType:   "b2/x-auto",
			}

			sfResp, err := b.client.StartLargeFile(ctx, sfReq)
			if err != nil {
				return nil, b.createError(err)
			}

			fileID = sfResp.FileID
		}

		for {
			if uploadURL == nil {
				uuReq := &b2.GetUploadPartURLRequest{
					Authorization: auth.authorizationToken,
					BaseURL:       auth.apiURL,
					FileID:        fileID,
				}

				uuResp, err := b.client.GetUploadPartURL(ctx, uuReq)
				if err != nil {
					return nil, b.createError(err)
				}

				authorization = uuResp.AuthorizationToken
				uploadURL = uuResp.UploadURL.URL
			}

			hash := sha1.New()
			hash.Write(request.Request.Data)
			checksum := fmt.Sprintf("%x", hash.Sum(nil))

			upReq := &b2.UploadPartRequest{
				Authorization: authorization,
				UploadURL:     uploadURL,
				PartNumber:    part,
				ContentLength: len(request.Request.Data),
				ContentSHA1:   checksum,
				Data:          io.NopCloser(bytes.NewReader(request.Request.Data)),
			}

			if _, err := b.client.UploadPart(ctx, upReq); err != nil {
				if e, ok := err.(*b2.Status); ok && e.Status >= http.StatusInternalServerError && retryAttempt < uploadRetryBudget {
					zap.L().Warn(
						"chunked file upload failed with server error; retrying after delay",
						zap.Int("attempt", retryAttempt),
						zap.Duration("delay", retryDelay),
						zap.Error(err),
					)

					select {
					case <-ctx.Done():
						return nil, status.FromContextError(ctx.Err()).Err()
					case <-time.After(retryDelay):
					}

					// Explicitly reset the upload URL so that the retry also
					// obtains a new upload URL before retrying the actual
					// upload. Such client behavior is mandated by the B2 API.
					// Reference: https://www.backblaze.com/b2/docs/uploading.html
					uploadURL = nil

					retryAttempt++
					retryDelay *= 2 // Exponential backoff
					continue
				}

				return nil, b.createError(err)
			}

			part++
			checksums = append(checksums, checksum)
			break
		}
	}

	ffReq := &b2.FinishLargeFileRequest{
		Authorization: auth.authorizationToken,
		BaseURL:       auth.apiURL,
		FileID:        fileID,
		PartSHA1Array: checksums,
	}

	ffResp, err := b.client.FinishLargeFile(ctx, ffReq)
	if err != nil {
		return nil, b.createError(err)
	}

	attributes, err := schemas.UnmarshalAttributes(ffResp.FileInfo, nil)
	if err != nil {
		return nil, b.createError(err)
	}

	return &service.PutObjectStreamResponse{
		Response: &service.PutObjectResponse{
			Metadata: &common.Metadata{
				Id:           ffResp.FileID,
				Key:          ffResp.FileName,
				Size:         uint64(ffResp.ContentLength),
				ObjectType:   common.Object_REGULAR,
				MimeType:     ffResp.ContentType,
				ModifiedTime: timestamppb.New(ffResp.UploadTimestamp.Time),
				Checksum:     ffResp.ContentSHA1,
				Attributes:   attributes,
			},
		},
	}, nil
}

// DeleteObject deletes the requested file.
func (b *B2) DeleteObject(ctx context.Context, request *service.DeleteObjectRequest) (*service.DeleteObjectResponse, error) {
	auth, err := b.authorizeAccount(ctx)
	if err != nil {
		return nil, b.createError(err)
	}

	lbReq := &b2.ListBucketsRequest{
		Authorization: auth.authorizationToken,
		BaseURL:       auth.apiURL,
		AccountID:     auth.accountID,
		BucketName:    request.Resource.Bucket,
	}

	lbResp, err := b.client.ListBuckets(ctx, lbReq)
	if err != nil {
		return nil, b.createError(err)
	}

	if len(lbResp.Buckets) != 1 {
		return nil, status.Errorf(
			codes.NotFound,
			"backend: unexpected number of buckets matched by name: name=%s buckets=%d",
			request.Resource.Bucket,
			len(lbResp.Buckets),
		)
	}

	lfReq := &b2.ListFileNamesRequest{
		Authorization: auth.authorizationToken,
		BaseURL:       auth.apiURL,
		BucketID:      lbResp.Buckets[0].BucketID,
		StartFileName: request.Key,
		MaxFileCount:  1,
		Prefix:        request.Key,
		Delimiter:     "/",
	}

	if request.Recursive {
		lfReq.Delimiter = ""
		lfReq.MaxFileCount = 10000
	}

	lfResp, err := b.client.ListFileNames(ctx, lfReq)
	if err != nil {
		return nil, b.createError(err)
	}

	if len(lfResp.Files) == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			"backend: no files match requested name: bucket=%s name=%s",
			request.Resource.Bucket,
			request.Key,
		)
	} else if !request.Recursive && len(lfResp.Files) == 1 && lfResp.Files[0].FileID == "" {
		// To emulate Unistore API semantics, disallow non-recursive deletions on non-empty
		// directories. The B2 API represents this a single matched file with no file ID.
		return nil, status.Errorf(
			codes.FailedPrecondition,
			"backend: cannot delete directory without recursive flag: bucket=%s name=%s",
			request.Resource.Bucket,
			request.Key,
		)
	} else if !request.Recursive && len(lfResp.Files) > 0 && strings.HasSuffix(request.Key, "/") {
		// Directories can only exist in B2 if there are nonzero children (potentially
		// nested arbitrarily deeply). Require such deletions to be recursive.
		return nil, status.Errorf(
			codes.FailedPrecondition,
			"backend: cannot delete non-empty directory without recursive flag: bucket=%s name=%s children=%d",
			request.Resource.Bucket,
			request.Key,
			len(lfResp.Files),
		)
	}

	for _, file := range lfResp.Files {
		dfReq := &b2.DeleteFileVersionRequest{
			Authorization: auth.authorizationToken,
			BaseURL:       auth.apiURL,
			FileName:      file.FileName,
			FileID:        file.FileID,
		}

		if _, err := b.client.DeleteFileVersion(ctx, dfReq); err != nil {
			return nil, b.createError(err)
		}
	}

	return &service.DeleteObjectResponse{Deleted: int32(len(lfResp.Files))}, nil
}

// ListBuckets retrieves metadata for all buckets.
func (b *B2) ListBuckets(ctx context.Context, request *service.ListBucketsRequest) (*service.ListBucketsResponse, error) {
	auth, err := b.authorizeAccount(ctx)
	if err != nil {
		return nil, b.createError(err)
	}

	lbReq := &b2.ListBucketsRequest{
		Authorization: auth.authorizationToken,
		BaseURL:       auth.apiURL,
		AccountID:     auth.accountID,
	}

	lbResp, err := b.client.ListBuckets(ctx, lbReq)
	if err != nil {
		return nil, b.createError(err)
	}

	buckets := make([]*service.HeadBucketResponse, len(lbResp.Buckets))

	for i, bucket := range lbResp.Buckets {
		var visibility common.Visibility

		switch bucket.BucketType {
		case "allPublic":
			visibility = common.Visibility_PUBLIC
		case "allPrivate", "snapshot":
			visibility = common.Visibility_PRIVATE
		default:
			visibility = common.Visibility_DEFAULT
		}

		lifecyclePolicies := make([]*common.LifecyclePolicy, len(bucket.LifecycleRules))

		for j, rule := range bucket.LifecycleRules {
			pattern := fmt.Sprintf("^%s", regexp.QuoteMeta(rule.FileNamePrefix))
			ttl := time.Duration(rule.DaysFromUploadingToHiding+rule.DaysFromHidingToDeleting) * 24 * time.Hour

			lifecyclePolicies[j] = &common.LifecyclePolicy{
				Pattern:    pattern,
				Expiration: durationpb.New(ttl),
			}
		}

		buckets[i] = &service.HeadBucketResponse{
			Bucket: &common.Bucket{
				Id:                bucket.BucketID,
				Name:              bucket.BucketName,
				Visibility:        visibility,
				LifecyclePolicies: lifecyclePolicies,
				Tags:              bucket.BucketInfo,
			},
		}
	}

	return &service.ListBucketsResponse{Buckets: buckets}, nil
}

// ListObjects retrieves metadata for all files that start with the specified key prefix.
func (b *B2) ListObjects(ctx context.Context, request *service.ListObjectsRequest) (*service.ListObjectsResponse, error) {
	auth, err := b.authorizeAccount(ctx)
	if err != nil {
		return nil, b.createError(err)
	}

	lbReq := &b2.ListBucketsRequest{
		Authorization: auth.authorizationToken,
		BaseURL:       auth.apiURL,
		AccountID:     auth.accountID,
		BucketName:    request.Resource.Bucket,
	}

	lbResp, err := b.client.ListBuckets(ctx, lbReq)
	if err != nil {
		return nil, b.createError(err)
	}

	if len(lbResp.Buckets) != 1 {
		return nil, status.Errorf(
			codes.NotFound,
			"backend: unexpected number of buckets matched by name: name=%s buckets=%d",
			request.Resource.Bucket,
			len(lbResp.Buckets),
		)
	}

	lfReq := &b2.ListFileNamesRequest{
		Authorization: auth.authorizationToken,
		BaseURL:       auth.apiURL,
		BucketID:      lbResp.Buckets[0].BucketID,
		Prefix:        request.Prefix,
		MaxFileCount:  10000,
	}

	// Specifying a delimiter will request B2 to only return paths whose component ends with
	// the specified delimiter. Unistore uses the forward slash character as the standard
	// delimiter for consistency with Unix filesystem semantics.
	if !request.Recursive {
		lfReq.Delimiter = "/"
	}

	lfResp, err := b.client.ListFileNames(ctx, lfReq)
	if err != nil {
		return nil, b.createError(err)
	}

	objects := make([]*service.HeadObjectResponse, len(lfResp.Files))

	for i, file := range lfResp.Files {
		objectType := common.Object_REGULAR
		if file.ContentType == "" && file.ContentLength == 0 {
			objectType = common.Object_DIRECTORY
		}

		attributes, err := schemas.UnmarshalAttributes(file.FileInfo, nil)
		if err != nil {
			return nil, b.createError(err)
		}

		objects[i] = &service.HeadObjectResponse{
			Metadata: &common.Metadata{
				Id:           file.FileID,
				Key:          file.FileName,
				Size:         uint64(file.ContentLength),
				ObjectType:   objectType,
				MimeType:     file.ContentType,
				ModifiedTime: timestamppb.New(file.UploadTimestamp.Time),
				Checksum:     file.ContentSHA1,
				Attributes:   attributes,
			},
		}
	}

	return &service.ListObjectsResponse{Objects: objects}, nil
}

// Descriptor returns a structured Protobuf-defined descriptor of this backend.
func (b *B2) Descriptor() *common.Node {
	return &common.Node{
		Name: "b2",
		Children: map[string]*common.Node_Value{
			"url": {Child: &common.Node_Value_Value{Value: b2.APIV2BaseURL}},
		},
	}
}

// Close is a noop.
func (b *B2) Close() error {
	return nil
}

// String returns a human-consumable representation of this backend.
func (b *B2) String() string {
	return schemas.MarshalDescriptor(b.Descriptor())
}

// authorizeAccount invokes the B2 account authorization API to obtain a short-lived authorization
// token for subsequent API calls. Authentication state is cached local to the backend instance and
// refreshed at a regular interval.
func (b *B2) authorizeAccount(ctx context.Context) (*authenticationState, error) {
	b.authMutex.Lock()
	defer b.authMutex.Unlock()

	// Authorization tokens are good for at most 24 hours after they have been issued by a call
	// to b2_authorize_account. Somewhat conservatively, this opts to refresh the token if more
	// than 3 hours have elapsed since the last successful refresh.
	if time.Since(b.auth.refresh) < authenticationRefreshInterval {
		return b.auth, nil
	}

	req := &b2.AuthorizeAccountRequest{
		ApplicationKeyID: b.applicationKeyID,
		ApplicationKey:   b.applicationKey,
	}

	resp, err := b.client.AuthorizeAccount(ctx, req)
	if err != nil {
		return nil, err
	}

	b.auth = &authenticationState{
		accountID:          resp.AccountID,
		authorizationToken: resp.AuthorizationToken,
		apiURL:             resp.APIURL.URL,
		downloadURL:        resp.DownloadURL.URL,
		refresh:            time.Now(),
	}

	return b.auth, nil
}

// createError creates a semantically accurate gRPC status for an error status from the B2 API.
func (b *B2) createError(err error) error {
	b2Errors := map[int]codes.Code{
		http.StatusBadRequest:                   codes.InvalidArgument,
		http.StatusUnauthorized:                 codes.PermissionDenied,
		http.StatusForbidden:                    codes.PermissionDenied,
		http.StatusNotFound:                     codes.NotFound,
		http.StatusMethodNotAllowed:             codes.InvalidArgument,
		http.StatusRequestTimeout:               codes.DeadlineExceeded,
		http.StatusRequestedRangeNotSatisfiable: codes.OutOfRange,
		http.StatusTooManyRequests:              codes.ResourceExhausted,
		http.StatusInternalServerError:          codes.Unknown,
		http.StatusServiceUnavailable:           codes.Unavailable,
	}

	switch e := err.(type) {
	case *b2.Status:
		code, ok := b2Errors[e.Status]
		if !ok {
			code = codes.Unknown
		}

		return status.Error(code, fmt.Sprintf("%s (code %s)", e.Message, e.Code))
	default:
		return status.Error(codes.Unknown, e.Error())
	}
}
