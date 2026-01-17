package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zeta-chain/node/zetaclient/types"
)

func TestNew(t *testing.T) {
	t.Run("in memory alias", func(t *testing.T) {
		// ARRANGE
		// Given a database
		db, err := NewFromSqliteInMemory(true)
		require.NoError(t, err)
		require.NotNil(t, db)

		// ACT
		runSampleSetGetTest(t, db)

		// Close the database
		assert.NoError(t, db.Close())
	})

	t.Run("in memory", func(t *testing.T) {
		// ARRANGE
		// Given a database
		db, err := NewFromSqlite(SqliteInMemory, "", true)
		require.NoError(t, err)
		require.NotNil(t, db)

		// ACT
		runSampleSetGetTest(t, db)

		// Close the database
		assert.NoError(t, db.Close())
	})

	t.Run("file based", func(t *testing.T) {
		// ARRANGE
		// Given a tmp path
		directory, dbName := t.TempDir(), "test.db"

		// Given a database
		db, err := NewFromSqlite(directory, dbName, true)
		require.NoError(t, err)
		require.NotNil(t, db)

		// Check that the database file exists
		assert.FileExists(t, directory+"/"+dbName)

		// ACT
		runSampleSetGetTest(t, db)

		// Close the database
		assert.NoError(t, db.Close())

		t.Run("close twice", func(t *testing.T) {
			require.NoError(t, db.Close())
		})
	})

	t.Run("invalid file path", func(t *testing.T) {
		// ARRANGE
		// Given a tmp path
		directory, dbName := "///hello", "test.db"

		// Given a database
		db, err := NewFromSqlite(directory, dbName, true)
		require.ErrorContains(t, err, "unable to ensure database path")
		require.Nil(t, db)
	})
}

func runSampleSetGetTest(t *testing.T, db *DB) {
	// Given a dummy sql type
	entity := types.ToLastBlockSQLType(444)

	// ACT #1
	// Create entity
	result := db.Client().Create(&entity)

	// ASSERT
	assert.NoError(t, result.Error)

	// ACT #2
	// Fetch entity
	var entity2 types.LastBlockSQLType

	result = db.Client().First(&entity2)

	// ASSERT
	assert.NoError(t, result.Error)
	assert.Equal(t, entity.Num, entity2.Num)
}
