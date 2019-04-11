package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/temprory/util"
)

// func clearTransaction(tx *sql.Tx) {
// 	err := tx.Rollback()
// 	if err != nil && err != sql.ErrTxDone {
// 		fmt.Printf("clearTransaction failed: %v\n", err)
// 	}
// }

func one(i int, db *sql.DB, rollback bool) error {
	defer fmt.Printf("--- one done %v\n", i)
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("one Begin() failed: %v, %v", i, err)
		return err
	}

	if rollback {
		defer util.ClearTx(tx)
	}

	sqlstr := "select 1 from dual"
	_, err = tx.Exec(sqlstr)
	if err != nil {
		fmt.Printf("one Exec() failed: %v, %v", i, err)
		return err
	}

	panic(fmt.Errorf("test err"))

	err = tx.Commit()
	return err
}

func main() {
	db, err := sql.Open("mysql", "root:123qwe@tcp(localhost:3306)/test")
	if err != nil {
		panic(fmt.Errorf("db.NewMysql sql.Open Failed: %v\n", err))
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(1)

	fmt.Println("********* with rollback")
	rollback := true
	for i := 0; i < 10; i++ {
		one(i, db, rollback)
	}

	fmt.Println("********* without rollback")
	rollback = false
	for i := 0; i < 10; i++ {
		one(i, db, rollback)
	}

}
