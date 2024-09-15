package gormwrapper

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

type sampleType struct {
	ID       int64  `gorm:"column:id;primaryKey"`
	Name     string `gorm:"column:name"`
	Price    int64  `gorm:"column:price"`
	IsPublic bool   `gorm:"column:is_public"`
}

func (t sampleType) TableName() string {
	return "sample_type_table"
}

func TestDb_Create(t *testing.T) {
	t.Run("Specified primary value", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer sqlDB.Close()

		row1 := &sampleType{
			ID:       1,
			Name:     "Super Phone",
			Price:    3000,
			IsPublic: true,
		}

		mock.ExpectBegin()
		mock.
			ExpectExec(regexp.QuoteMeta("INSERT INTO `sample_type_table` (`name`,`price`,`is_public`,`id`) VALUES (?,?,?,?)")).
			WithArgs(row1.Name, row1.Price, row1.IsPublic, row1.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		conn, err := newMySQLConn(sqlDB)
		assert.NoError(t, err)

		err = conn.Create(row1)
		assert.NoError(t, err)
	})

	t.Run("Not specified primary value", func(t *testing.T) {
		sqlDB, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer sqlDB.Close()

		row1 := &sampleType{
			// ID:       1, unspecified
			Name:     "Super Phone",
			Price:    3000,
			IsPublic: true,
		}
		newID := int64(3)

		mock.ExpectBegin()
		mock.
			ExpectExec(regexp.QuoteMeta("INSERT INTO `sample_type_table` (`name`,`price`,`is_public`) VALUES (?,?,?)")).
			WithArgs(row1.Name, row1.Price, row1.IsPublic).
			WillReturnResult(sqlmock.NewResult(newID, 1))
		mock.ExpectCommit()

		conn, err := newMySQLConn(sqlDB)
		assert.NoError(t, err)

		err = conn.Create(row1)
		assert.NoError(t, err)
		assert.Equal(t, newID, row1.ID)
	})
}

func newMySQLConn(sqlDB *sql.DB) (DB, error) {
	conn, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}))
	if err != nil {
		return nil, err
	}
	return NewDB(conn), nil
}
