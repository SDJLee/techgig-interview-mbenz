package util

const (
	EnvDev             = "dev"
	EnvProd            = "prod"
	AppEnv             = "APP_ENV"
	BasePath           = "BASE_PATH"
	DefaultBasePath    = "/app"
	ConfigFileFormat   = "app-%s"
	ConfigFileType     = "env"
	ApiAddress         = "API_ADDRESS"
	LogPath            = "LOG_PATH"
	Port               = "PORT"
	ServerReadTimeout  = "SERVER_READ_TIMEOUT"
	ServerWriteTimeout = "SERVER_WRITE_TIMEOUT"
	ApiHealthCheck     = "health"
	ApiComputeRoute    = "/compute-route"
	ApiBasePath        = "/api"
	ApiV1              = "/v1"
	ErrUnreachableId   = 8888
	ErrUnreachableMsg  = "Unable to reach the destination with the current charge level"
	ErrTechExpId       = 9999
	ErrTechExpMsg      = "Technical Exception"
)
