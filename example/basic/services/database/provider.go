package database

type DatabaseProvider interface {
	GetData() string
	GetPrefix() string
}
