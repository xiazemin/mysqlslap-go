package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func GetClient(ctx context.Context, host string,
	port int64,
	user string,
	password string,
	database string, concurrency int, timeout int64) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", user, password, "tcp", host, port, database)
	DB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("Open mysql failed,err:%v\n", err)

	}
	DB.SetConnMaxLifetime(time.Duration(timeout) * time.Second) //最大连接周期，超过时间的连接就close
	DB.SetMaxOpenConns(concurrency)                             //设置最大连接数
	DB.SetMaxIdleConns(16)                                      //设置闲置连接数
	return DB, nil
}

func ExecQuery(ctx context.Context, db *sql.DB, sql string) (int64, error) {
	result, err := db.Exec(sql)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
