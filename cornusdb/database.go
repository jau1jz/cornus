package cornusdb

import (
	"errors"
	"github.com/jau1jz/cornus/commons"
	slog "github.com/jau1jz/cornus/commons/log"
	"github.com/jau1jz/cornus/config"
	serveries "github.com/jau1jz/cornus/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
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
			SingularTable: true, //单表名
		}, Logger: logger.New(
			&slog.Slog,
			logger.Config{
				SlowThreshold:             time.Second,                                         // 慢 SQL 阈值
				LogLevel:                  commons.GormLogLevel[config.SC.SConfigure.LogLevel], // 日志级别
				IgnoreRecordNotFoundError: true,                                                // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,                                               // 禁用彩色打印
			},
		)})
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
