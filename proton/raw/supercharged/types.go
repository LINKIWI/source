package supercharged

const (
	// CodeOK describes a successful response.
	CodeOK = 0
	// CodeServerUndefined describes an undefined server-side error.
	CodeServerUndefined = -1
	// CodeClientUndefined describes an undefined client-side error.
	CodeClientUndefined = -2
	// CodeInvalidParameters indicates the client supplied invalid input parameters.
	CodeInvalidParameters = -3
	// CodeNotFound indicates the client attempted to access an unknown route.
	CodeNotFound = -4
)

const (
	// requestHeaderClientID is the outgoing request header indicating the client ID.
	requestHeaderClientID = "X-Supercharged-Client"
	// requestHeaderClientVersion is the outgoing request header indicating the client version.
	requestHeaderClientVersion = "X-Supercharged-Client-Version"
	// requestHeaderData is the outgoing request header used for a request data payload used in
	// lieu of a request body, for HTTP methods where bodies are not supported.
	requestHeaderData = "X-Supercharged-Data"
)

const (
	// clientUserAgent is a format string for the User-Agent header attached to all requests.
	clientUserAgent = "proton/1.0 (client:%s; version:%s)"
)
