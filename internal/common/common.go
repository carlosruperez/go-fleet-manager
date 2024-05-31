package common

type Version string

type RepositoryTypes string

const (
	Microservice RepositoryTypes = "ms"
	SDK          RepositoryTypes = "sdk"
	Others       RepositoryTypes = "others"
)

type Environment string

const (
	Development Environment = "development"
	Staging     Environment = "staging"
	Production  Environment = "production"
)
