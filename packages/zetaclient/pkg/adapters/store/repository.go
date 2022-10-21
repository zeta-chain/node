package store

// "github.com/zeta-chain/zetacore/packages/zetaclient/pkg/model"

//go:generate mockery --name Repository

// Repository interface to db
type Repository interface {
	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	NewIterator(prefix []byte) Iterator
}

//go:generate mockery --name Iterator
type Iterator interface {
	Next() bool
	Key() []byte
	Value() []byte
	Release()
	Error() error
}
