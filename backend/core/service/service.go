package service

// ServiceEnvironment is the environment of the service
type ServiceEnvironment string

const (
	Testing    ServiceEnvironment = "test"
	Staging    ServiceEnvironment = "stag"
	Production ServiceEnvironment = "prod"
)
