package infra

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zeta-chain/zetacore/packages/zetaclient/pkg/adapters/store"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDBRepository struct {
	db *leveldb.DB
}

var _ store.Repository = (*LevelDBRepository)(nil)

func NewLevelDBRepository(path string) (*LevelDBRepository, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dbpath := fmt.Sprintf("%s/%s", filepath.Join(userDir, ".zetaclient/chainobserver"), path)
	db, err := leveldb.OpenFile(dbpath, nil)
	if err != nil {
		return nil, err
	}
	return &LevelDBRepository{
		db: db,
	}, nil
}

func (rep *LevelDBRepository) Get(key []byte) ([]byte, error) {
	return rep.db.Get(key, nil)
}

func (rep *LevelDBRepository) Put(key, value []byte) error {
	return rep.db.Put(key, value, nil)
}

func (rep *LevelDBRepository) NewIterator(prefix []byte) store.Iterator {
	return rep.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
}
