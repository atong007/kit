package sql

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/url"
	"time"
)

func NewMySQL(host, port, user, pass, name string, debug bool) (*gorm.DB, error) {
	con := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, name)
	v := url.Values{}
	v.Add("parseTime", "1")
	v.Add("loc", "Asia/Shanghai")
	dsn := fmt.Sprintf("%s?%s", con, v.Encode())
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql err:%w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("ping mysql err:%w", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	g, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}
	if debug {
		g = g.Debug()
	}
	return g, nil
}
