package config

// ServiceConfig конфигурация сервиса.
type ServiceConfig struct {
	GrpcEndpoint string       `json:"grpc_endpoint" env:"SERVICE_GRPC_ENDPOINT"`
	Info         *ServiceInfo `json:"info"`
}
