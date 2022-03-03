package filters

// requestContextKey is used for arbitrary state attached to a request retained throughout its
// entire lifecycle, for use by both request and response processing.
type requestContextKey int

// Request context keys used by filters, centralized in the filters package as an enum to maintain
// isolation of context values across multiple filters.
const (
	ctxInstrumentationStopwatch requestContextKey = iota
	ctxInstrumentationRequestBody

	ctxMatchStatus

	ctxMetadataStopwatch

	ctxLogRequestID
	ctxLogStopwatch
	ctxLogOriginalURL
)
