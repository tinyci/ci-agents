package config

// Websockets is general websocket configuration.
type Websockets struct {
	MaxAttachBuffer    int64 `yaml:"attach_buffer_size"`
	BufSize            int   `yaml:"read_buffer_size"`
	InsecureWebSockets bool  `yaml:"insecure_websockets"`
}
