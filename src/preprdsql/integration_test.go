package preprdsql

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mxk/go-sqlite/sqlite3"
)

func initSQLFiles() *os.File {
	file, err := os.Create("test_sql.sql")
	defer file.Close()

	if err != nil {
		panic(err)
	}

	writer := bufio.NewWriter(file)

	fmt.Fprintln(writer,
		`
		    -- name: create-users-table
		    CREATE TABLE users (
			    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			    name VARCHAR(255),
			    email VARCHAR(255)
		    );

		    -- name: create-user
		    INSERT INTO users (name, email) VALUES(?, ?)

		    -- name: find-one-user-by-email
		    SELECT id,name,email FROM users WHERE email = ? LIMIT 1

		    -- name: drop-users-table
		    DROP TABLE users
		`,
	)
	writer.Flush()
	return file
}

func initArrangement() (*sqlx.DB, *Repo) {
	db, err := sqlx.Open("sqlite3", ":memory:")

	if err != nil {
		panic(err)
	}

	file := initSQLFiles()

	defer os.Remove(file.Name())

	repo, err := LoadFromFile(file.Name())

	if err != nil {
		panic(err)
	}

	sql, err := repo.LookupSQL("create-users-table")

	if err != nil {
		panic(err)
	}

	db.MustExec(sql)

	return db, repo
}

func TestIntegrationWithSqlx_Exec(t *testing.T) {
	db, repo := initArrangement()
	defer db.Close()

	sql, err := repo.LookupSQL("create-user")

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(sql, "Man Bear Pig", "itsreal@mbpig.com")

	if err != nil {
		panic(err)
	}

	row := db.QueryRow(
		"SELECT email FROM users LIMIT 1",
	)

	var email string
	row.Scan(&email)

	if email != "itsreal@mbpig.com" {
		t.Errorf("Expect to find user with email == %s, got %s", "itsreal@mbpig.com", email)
	}
}

func TestIntegrationWithSqlx_Query(t *testing.T) {
	db, repo := initArrangement()
	defer db.Close()

	sql, err := repo.LookupSQL("create-user")

	if err != nil {
		panic(err)
	}

	_, err = db.Exec(sql, "Man Bear Pig", "itsreal@mbpig.com")
	if err != nil {
		panic(err)
	}

	sql, err = repo.LookupSQL("find-one-user-by-email")

	rows, err := db.Query(sql, "itsreal@mbpig.com")

	if err != nil {
		panic(err)
	}

	ok := rows.Next()

	if ok == false {
		t.Errorf("User with email 'itsreal@mbpig.com' cound not be found")
	}

	var id, name interface{}
	var email string

	err = rows.Scan(&id, &name, &email)

	if err != nil {
		panic(err)
	}

	if email != "itsreal@mbpig.com" {
		t.Errorf("Expect to find user with email == %s, got %s", "itsreal@mbpig.com", email)
	}
}

func TestIntegrationWithSqlx_QueryRow(t *testing.T) {
	db, repo := initArrangement()
	defer db.Close()

	sql, err := repo.LookupSQL("find-one-user-by-email")

	db.MustExec("INSERT INTO users(email) VALUES('itsreal@mbpig.com')")

	if err != nil {
		panic(err)
	}

	row := db.QueryRow(sql, "itsreal@mbpig.com")

	if err != nil {
		panic(err)
	}

	var id, name interface{}
	var email string

	err = row.Scan(&id, &name, &email)

	if err != nil {
		t.Errorf("User with email 'itsreal@mbpig.com' cound not be found")
	}

	if email != "itsreal@mbpig.com" {
		t.Errorf("Expect to find user with email == %s, got %s", "itsreal@mbpig.com", email)
	}
}
