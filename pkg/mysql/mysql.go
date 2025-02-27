package mysql

import (
	"fmt"
	"github.com/SaveljevRoman/go-layout-project/internal/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewConnection(cfg config.MySQLConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Установка параметров подключения
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}
