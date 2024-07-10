package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	var email, password, storagePath string

	flag.StringVar(&email, "email", "", "new admin's email")
	flag.StringVar(&password, "password", "", "new password's password")
	flag.StringVar(&storagePath, "storage-path", "", "storage path")

	flag.Parse()

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Printf("failed to generate password hash: %s\n", err.Error())
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", storagePath)

	if err != nil {
		fmt.Printf("failed to connect db: %s\n", err.Error())
		os.Exit(1)
	}

	stmt, _ := db.Prepare("INSERT INTO users(email, pass_hash, is_admin) VALUES(?, ?, ?)")

	res, err := stmt.ExecContext(context.Background(), email, passHash, true)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			fmt.Printf("email %s is already busy\n", email)
			os.Exit(1)
		}

		fmt.Printf("failed to save admin: %s\n", err.Error())
		os.Exit(1)
	}

	userId, _ := res.LastInsertId()

	fmt.Printf("successfily create admin with id: %d\n", userId)
}
