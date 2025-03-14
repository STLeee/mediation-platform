package service

type Environment string

const (
	Testing    Environment = "test"
	Staging    Environment = "stag"
	Production Environment = "prod"
)
