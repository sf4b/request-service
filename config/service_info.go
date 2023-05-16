package config

// ServiceInfo информация о сервисе.
type ServiceInfo struct {
	Name    string `json:"name" env:"SERVICE_INFO_NAME"`
	Version string `json:"version" env:"SERVICE_INFO_VERSION"`
}
