package config

// Tag is an identifier for static metadata tag keys.
type Tag string

// Serializer is an identifier for message serializer implementations.
type Serializer string

// SeekPosition is an identifier for file seek behaviors.
type SeekPosition string

const (
	// TagHostname is the tag identifier for the machine hostname.
	TagHostname Tag = "hostname"
	// TagVersion is the tag identifier for the application version.
	TagVersion Tag = "version"

	// SerializerLine is the serializer identifier for the line message serializer.
	SerializerLine Serializer = "line"
	// SerializerJSON is the serializer identifier for the JSON message serializer.
	SerializerJSON Serializer = "json"

	// SeekPositionStart is the seek position identifier for the beginning of the file.
	SeekPositionStart SeekPosition = "start"
	// SeekPositionEnd is the seek position identifier for the end of the file.
	SeekPositionEnd SeekPosition = "end"
)

// Config describes the top-level sections recognized in the configuration.
type Config struct {
	Piper  Piper   `toml:"piper"`
	Relays []Relay `toml:"relay"`
}

// Piper provides configuration of the application itself.
type Piper struct {
	MetricsAddress string `toml:"metrics-address"`
	MetricsProxy   string `toml:"metrics-proxy"`
	MetricsPrefix  string `toml:"metrics-prefix"`
	SentryDSN      string `toml:"sentry-dsn"`
}

// Relay provides configuration of a single log relay stream.
type Relay struct {
	Name           string       `toml:"name"`
	LogFile        glob         `toml:"log-file"`
	Delimiter      string       `toml:"delimiter"`
	ResetDelay     duration     `toml:"reset-delay"`
	Filter         regex        `toml:"filter"`
	Serializer     Serializer   `toml:"serializer"`
	SeekPosition   SeekPosition `toml:"seek-position"`
	TagIdentifiers []Tag        `toml:"tags"`
	BufferLength   int          `toml:"buffer-length"`
	ProxyAddress   address      `toml:"proxy-address"`
	KafkaAddress   address      `toml:"kafka-address"`
	KafkaTopic     string       `toml:"kafka-topic"`
	KafkaAcks      int          `toml:"kafka-acks"`
	KafkaRetries   int          `toml:"kafka-retries"`
	KafkaTimeout   duration     `toml:"kafka-timeout"`
}
