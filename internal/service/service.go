package service

type Service interface {
	GetHost() string
	GetPort() string
	GetName() string
	GetID() string
	GetMetadata() map[string]string
}
