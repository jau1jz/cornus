package cornusdb

import (
	"errors"
	slog "github.com/jau1jz/cornus/commons/log"
	serveries "github.com/jau1jz/cornus/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

func (slf *CornusDB) StartSqlite(config serveries.DataBaseConfig) error {
	if slf.db != nil {
		return errors.New("db already open")
	}
	slf.name = config.Name

	var err error
	slf.db, err = gorm.Open(sqlite.Open(config.DBFilePath), &gorm.Config{PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //单表名
		}, Logger: &slog.Slog})
	slf.db.Logger.LogMode(logger.Info)

	if err != nil {
		slog.Slog.InfoF("conn database error %s", err)
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
