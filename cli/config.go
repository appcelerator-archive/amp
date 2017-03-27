package cli

// Configuration is for all configurable client settings
type Configuration struct {
	Address string
	Verbose bool
	Theme   string
}

const (
	// DefaultConnection default address for amp api using grpc protocol
	DefaultConnection = "m1:50101"
)
