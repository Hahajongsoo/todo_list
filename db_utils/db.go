package db_utils

import (
	"database/sql"
	"log"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
)

func UseDB() *sql.DB {
	db, err := sql.Open("mysql", "todo_admin:1234@tcp(127.0.0.1:3306)/todo")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func RowsToStructs(rows *sql.Rows, dest interface{}) error {
	destv := reflect.ValueOf(dest).Elem()

	args := make([]interface{}, destv.Type().Elem().NumField())

	for rows.Next() {
		rowp := reflect.New(destv.Type().Elem())
		rowv := rowp.Elem()

		for i := 0; i < rowv.NumField(); i++ {
			args[i] = rowv.Field(i).Addr().Interface()
		}

		if err := rows.Scan(args...); err != nil {
			return err
		}

		destv.Set(reflect.Append(destv, rowv))
	}
	return nil
}
