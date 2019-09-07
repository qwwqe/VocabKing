package storage

type Storage interface {
	Open() error
	Close() error

	GetHomepage(*responses.Homepage) error
}
