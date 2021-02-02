package sql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var RowsNotAffectedErr = fmt.Errorf("rows not affected")

// Config sql config.
type Config struct {
	DSN    string // write data source name.
	BizTag string
}

type MySQL struct {
	db   *sql.DB
	conf *Config
	step int
}

// NewMySQL new db and retry connection when has error.
func NewMySQL(c *Config) *MySQL {
	db, err := sql.Open("mysql", c.DSN)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Hour) //最大连接周期，超过时间的连接就close
	db.SetMaxOpenConns(10)           //设置最大连接数
	db.SetMaxIdleConns(2)            //设置闲置连接数

	return &MySQL{
		db:   db,
		conf: c,
	}
}

func (m *MySQL) InitBizTag(ctx context.Context, bizTag string, maxID uint64, step int, description string) error {
	db := m.db

	_, err := db.Exec(_insertSql, bizTag, maxID, step, description)
	if err != nil {
		return err
	}

	return nil
}

func (m *MySQL) GetEndID(ctx context.Context) (startID uint64, endID uint64, step int, err error) {
	db := m.db
	tx, err := db.Begin()
	if err != nil {
		return
	}

	bizTag := m.conf.BizTag
	res, err := tx.Exec(_updateSql, bizTag)
	if err != nil {
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAffected == 0 {
		err = RowsNotAffectedErr
		return
	}

	row := tx.QueryRow(_selectSql, bizTag)

	var maxID uint64
	err = row.Scan(&maxID, &step)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	m.step = step
	startID = maxID - uint64(step) + 1
	endID = maxID
	return
}
