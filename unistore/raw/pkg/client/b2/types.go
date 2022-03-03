package b2

import (
	"io"
	"net/url"
)

// Base API constants.
const (
	APIV2BasePath = "/b2api/v2/"
	APIV2BaseURL  = "https://api.backblazeb2.com" + APIV2BasePath
)

// API endpoint names, used during URL construction.
const (
	EndpointAuthorizeAccount  = "b2_authorize_account"
	EndpointCancelLargeFile   = "b2_cancel_large_file"
	EndpointCopyFile          = "b2_copy_file"
	EndpointCopyPart          = "b2_copy_part"
	EndpointDeleteFileVersion = "b2_delete_file_version"
	EndpointFinishLargeFile   = "b2_finish_large_file"
	EndpointGetUploadPartURL  = "b2_get_upload_part_url"
	EndpointGetUploadURL      = "b2_get_upload_url"
	EndpointListBuckets       = "b2_list_buckets"
	EndpointListFileNames     = "b2_list_file_names"
	EndpointStartLargeFile    = "b2_start_large_file"
)

const (
	// DefaultClientID is the default client identifier for this B2 client, used as the
	// User-Agent for all outbound HTTP requests to the B2 API.
	DefaultClientID = "unistore/1.0.0"
)

// Range specifies a seek range describing part of an object.
type Range struct {
	Unit  string
	Start int
	End   int
	Total int
}

// AuthorizeAccountRequest is the request schema for b2_authorize_account.
type AuthorizeAccountRequest struct {
	ApplicationKeyID string
	ApplicationKey   string
}

// AuthorizeAccountResponse is the response schema for b2_authorize_account.
type AuthorizeAccountResponse struct {
	AccountID          string `json:"accountId"`
	AuthorizationToken string `json:"authorizationToken"`
	Allowed            struct {
		Capabilities []string `json:"capabilities"`
		BucketID     string   `json:"bucketId"`
		BucketName   string   `json:"bucketName"`
		NamePrefix   string   `json:"namePrefix"`
	} `json:"allowed"`
	APIURL                  *URL `json:"apiUrl"`
	DownloadURL             *URL `json:"downloadUrl"`
	RecommendedPartSize     int  `json:"recommendedPartSize"`
	AbsoluteMinimumPartSize int  `json:"absoluteMinimumPartSize"`
	S3APIURL                *URL `json:"s3ApiUrl"`
}

// CancelLargeFileRequest is the request schema for b2_cancel_large_file.
type CancelLargeFileRequest struct {
	Authorization string   `json:"-" header:"Authorization"`
	BaseURL       *url.URL `json:"-"`
	FileID        string   `json:"fileId"`
}

// CancelLargeFileResponse is the response schema for b2_cancel_large_file.
type CancelLargeFileResponse struct {
	FileID    string `json:"fileId"`
	AccountID string `json:"accountId"`
	BucketID  string `json:"bucketId"`
	FileName  string `json:"fileName"`
}

// CopyFileRequest is the request schema for b2_copy_file.
type CopyFileRequest struct {
	Authorization       string            `json:"-" header:"Authorization"`
	BaseURL             *url.URL          `json:"-"`
	SourceFileID        string            `json:"sourceFileId"`
	DestinationBucketID string            `json:"destinationBucketId"`
	FileName            string            `json:"fileName"`
	MetadataDirective   string            `json:"metadataDirective"`
	ContentType         string            `json:"contentType"`
	FileInfo            map[string]string `json:"fileInfo"`
}

// CopyFileResponse is the response schema for b2_copy_file.
type CopyFileResponse struct {
	AccountID       string            `json:"accountId"`
	Action          string            `json:"action"`
	BucketID        string            `json:"bucketId"`
	ContentLength   int               `json:"contentLength"`
	ContentSHA1     string            `json:"contentSha1"`
	ContentMD5      string            `json:"contentMd5"`
	ContentType     string            `json:"contentType"`
	FileID          string            `json:"fileId"`
	FileInfo        map[string]string `json:"fileInfo"`
	FileName        string            `json:"fileName"`
	UploadTimestamp *UnixTimestamp    `json:"uploadTimestamp"`
}

// CopyPartRequest is the request schema for b2_copy_part.
type CopyPartRequest struct {
	Authorization string   `json:"-" header:"Authorization"`
	BaseURL       *url.URL `json:"-"`
	SourceFileID  string   `json:"sourceFileId"`
	LargeFileID   string   `json:"largeFileId"`
	PartNumber    int      `json:"partNumber"`
}

// CopyPartResponse is the response schema for b2_copy_part.
type CopyPartResponse struct {
	FileID          string         `json:"fileId"`
	PartNumber      int            `json:"partNumber"`
	ContentLength   int            `json:"contentLength"`
	ContentSHA1     string         `json:"contentSha1"`
	ContentMD5      string         `json:"contentMd5"`
	UploadTimestamp *UnixTimestamp `json:"uploadTimestamp"`
}

// DeleteFileVersionRequest is the request schema for b2_delete_file_version.
type DeleteFileVersionRequest struct {
	Authorization string   `json:"-" header:"Authorization"`
	BaseURL       *url.URL `json:"-"`
	FileName      string   `json:"fileName"`
	FileID        string   `json:"fileId"`
}

// DeleteFileVersionResponse is the response schema for b2_delete_file_version.
type DeleteFileVersionResponse struct {
	FileID   string `json:"fileId"`
	FileName string `json:"fileName"`
}

// DownloadFileByNameRequest is the request schema for b2_download_file_by_name.
type DownloadFileByNameRequest struct {
	Authorization string `header:"Authorization"`
	DownloadURL   *url.URL
	Range         Range
	ContentType   string `header:"Content-Type"`
	Bucket        string
	File          string
}

// DownloadFileByNameResponse is the response schema for b2_download_file_by_name.
type DownloadFileByNameResponse struct {
	FileID          string
	FileName        string
	FileInfo        map[string]string
	ContentSHA1     string
	Size            int
	ContentType     string
	Range           Range
	UploadTimestamp *UnixTimestamp
	Data            io.ReadCloser
}

// FinishLargeFileRequest is the request schema for b2_finish_large_file.
type FinishLargeFileRequest struct {
	Authorization string   `json:"-" header:"Authorization"`
	BaseURL       *url.URL `json:"-"`
	FileID        string   `json:"fileId"`
	PartSHA1Array []string `json:"partSha1Array"`
}

// FinishLargeFileResponse is the response schema for b2_finish_large_file.
type FinishLargeFileResponse struct {
	AccountID       string            `json:"accountId"`
	Action          string            `json:"action"`
	BucketID        string            `json:"bucketId"`
	ContentLength   int               `json:"contentLength"`
	ContentSHA1     string            `json:"contentSha1"`
	ContentMD5      string            `json:"contentMd5"`
	ContentType     string            `json:"contentType"`
	FileID          string            `json:"fileId"`
	FileInfo        map[string]string `json:"fileInfo"`
	FileName        string            `json:"fileName"`
	UploadTimestamp *UnixTimestamp    `json:"uploadTimestamp"`
}

// GetUploadPartURLRequest is the request schema for b2_get_upload_part_url.
type GetUploadPartURLRequest struct {
	Authorization string   `json:"-" header:"Authorization"`
	BaseURL       *url.URL `json:"-"`
	FileID        string   `json:"fileId"`
}

// GetUploadPartURLResponse is the response schema for b2_get_upload_part_url.
type GetUploadPartURLResponse struct {
	FileID             string `json:"fileId"`
	UploadURL          *URL   `json:"uploadUrl"`
	AuthorizationToken string `json:"authorizationToken"`
}

// GetUploadURLRequest is the request schema for b2_get_upload_url.
type GetUploadURLRequest struct {
	Authorization string   `json:"-" header:"Authorization"`
	BaseURL       *url.URL `json:"-"`
	BucketID      string   `json:"bucketId"`
}

// GetUploadURLResponse is the response schema for b2_get_upload_url.
type GetUploadURLResponse struct {
	BucketID           string `json:"bucketId"`
	UploadURL          *URL   `json:"uploadUrl"`
	AuthorizationToken string `json:"authorizationToken"`
}

// ListBucketsRequest is the request schema for b2_list_buckets.
type ListBucketsRequest struct {
	Authorization string   `json:"-" header:"Authorization"`
	BaseURL       *url.URL `json:"-"`
	AccountID     string   `json:"accountId"`
	BucketID      string   `json:"bucketId,omitempty"`
	BucketName    string   `json:"bucketName,omitempty"`
	BucketTypes   []string `json:"bucketTypes,omitempty"`
}

// ListBucketsResponse is the response schema for b2_list_buckets.
type ListBucketsResponse struct {
	Buckets []struct {
		AccountID      string            `json:"accountId"`
		BucketID       string            `json:"bucketId"`
		BucketName     string            `json:"bucketName"`
		BucketType     string            `json:"bucketType"`
		BucketInfo     map[string]string `json:"bucketInfo"`
		CORSRules      []string          `json:"corsRules"`
		LifecycleRules []struct {
			FileNamePrefix            string `json:"fileNamePrefix"`
			DaysFromUploadingToHiding int    `json:"daysFromUploadingToHiding"`
			DaysFromHidingToDeleting  int    `json:"daysFromHidingToDeleting"`
		} `json:"lifecycleRules"`
		Revision int      `json:"revision"`
		Options  []string `json:"options"`
	} `json:"buckets"`
}

// ListFileNamesRequest is the request schema for b2_list_file_names.
type ListFileNamesRequest struct {
	Authorization string   `json:"-" header:"Authorization"`
	BaseURL       *url.URL `json:"-"`
	BucketID      string   `json:"bucketId"`
	StartFileName string   `json:"startFileName,omitempty"`
	MaxFileCount  int      `json:"maxFileCount,omitempty"`
	Prefix        string   `json:"prefix,omitempty"`
	Delimiter     string   `json:"delimiter,omitempty"`
}

// ListFileNamesResponse is the response schema for b2_list_file_names.
type ListFileNamesResponse struct {
	Files []struct {
		AccountID       string            `json:"accountId"`
		BucketID        string            `json:"bucketId"`
		ContentLength   int               `json:"contentLength"`
		ContentSHA1     string            `json:"contentSha1"`
		ContentMD5      string            `json:"contentMd5"`
		ContentType     string            `json:"contentType"`
		FileID          string            `json:"fileId"`
		FileInfo        map[string]string `json:"fileInfo"`
		FileName        string            `json:"fileName"`
		UploadTimestamp *UnixTimestamp    `json:"uploadTimestamp"`
	} `json:"files"`
	NextFileName string `json:"nextFileName"`
}

// StartLargeFileRequest is the request schema for b2_start_large_file.
type StartLargeFileRequest struct {
	Authorization string            `json:"-" header:"Authorization"`
	BaseURL       *url.URL          `json:"-"`
	BucketID      string            `json:"bucketId"`
	FileName      string            `json:"fileName"`
	FileInfo      map[string]string `json:"fileInfo"`
	ContentType   string            `json:"contentType"`
}

// StartLargeFileResponse is the response schema for b2_start_large_file.
type StartLargeFileResponse struct {
	AccountID       string            `json:"accountId"`
	Action          string            `json:"action"`
	BucketID        string            `json:"bucketId"`
	ContentLength   int               `json:"contentLength"`
	ContentSHA1     string            `json:"contentSha1"`
	ContentMD5      string            `json:"contentMd5"`
	ContentType     string            `json:"contentType"`
	FileID          string            `json:"fileId"`
	FileInfo        map[string]string `json:"fileInfo"`
	FileName        string            `json:"fileName"`
	UploadTimestamp *UnixTimestamp    `json:"uploadTimestamp"`
}

// UploadFileRequest is the request schema for b2_upload_file.
type UploadFileRequest struct {
	Authorization string `header:"Authorization"`
	UploadURL     *url.URL
	FileName      string `header:"X-Bz-File-Name,encode"`
	FileInfo      map[string]string
	ContentType   string `header:"Content-Type"`
	ContentLength int    `header:"Content-Length"`
	ContentSHA1   string `header:"X-Bz-Content-Sha1"`
	Data          io.ReadCloser
}

// UploadFileResponse is the response schema for b2_upload_file.
type UploadFileResponse struct {
	AccountID       string            `json:"accountId"`
	Action          string            `json:"action"`
	BucketID        string            `json:"bucketId"`
	ContentLength   int               `json:"contentLength"`
	ContentSHA1     string            `json:"contentSha1"`
	ContentMD5      string            `json:"contentMd5"`
	ContentType     string            `json:"contentType"`
	FileID          string            `json:"fileId"`
	FileInfo        map[string]string `json:"fileInfo"`
	FileName        string            `json:"fileName"`
	UploadTimestamp *UnixTimestamp    `json:"uploadTimestamp"`
}

// UploadPartRequest is the request schema for b2_upload_part.
type UploadPartRequest struct {
	Authorization string `header:"Authorization"`
	UploadURL     *url.URL
	PartNumber    int    `header:"X-Bz-Part-Number"`
	ContentLength int    `header:"Content-Length"`
	ContentSHA1   string `header:"X-Bz-Content-Sha1"`
	Data          io.ReadCloser
}

// UploadPartResponse is the response schema for b2_upload_part.
type UploadPartResponse struct {
	FileID          string         `json:"fileId"`
	PartNumber      int            `json:"partNumber"`
	ContentLength   int            `json:"contentLength"`
	ContentSHA1     string         `json:"contentSha1"`
	ContentMD5      string         `json:"contentMd5"`
	UploadTimestamp *UnixTimestamp `json:"uploadTimestamp"`
}
