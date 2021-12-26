package cornusdb

import (
	"context"
	"errors"
	slog "github.com/jau1jz/cornus/commons/log"
	serveries "github.com/jau1jz/cornus/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type CornusDB struct {
	db     *gorm.DB
	name   string //db name
	dbType string //db 类型 mysql sqlite
}

func (slf *CornusDB) GormDB() *gorm.DB {
	return slf.db
}

func (slf *CornusDB) Name() string {
	return slf.name
}

func (slf *CornusDB) StartSqlite(dbConfig serveries.DataBaseConfig) error {
	if slf.db != nil {
		return errors.New("db already open")
	}
	slf.name = dbConfig.Name
	var err error
	slf.db, err = gorm.Open(sqlite.Open(dbConfig.DBFilePath), &gorm.Config{PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		}, Logger: &slog.Gorm})
	if err != nil {
		slog.Slog.InfoF(context.Background(), "conn database error %s", err)
		return err
	}
	return nil
}

func (slf *CornusDB) StopDb() error {
	if slf.db != nil {
		db, err := slf.db.DB()
		if err != nil {
			slf.db = nil
		} else {
			_ = db.Close()
		}
		return err
	} else {
		return errors.New("db is nil")
	}
}
