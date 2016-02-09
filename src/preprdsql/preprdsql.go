// Package preprdsql provides a way to separate your SQL statements from your
// code.
package preprdsql

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Repo stores the sql statements scanned from a loaded file
type Repo struct {
	sqlStatements map[string]string
}

// LookupSQL searches the Repo using the tag and returns the sql statement if
// found. It returns an error if the tag is not located.
func (d Repo) LookupSQL(tag string) (sql string, err error) {
	sqlStatement, ok := d.sqlStatements[tag]
	if !ok {
		err = fmt.Errorf("preprdsql: '%s' could not be found", tag)
	}

	return sqlStatement, err
}

// loadToScanner loads the opened file to the Scanner
func loadToScanner(r io.Reader) (*Repo, error) {
	scanner := &Scanner{}
	statements := scanner.Scan(bufio.NewScanner(r))

	repo := &Repo{
		sqlStatements: statements,
	}

	return repo, nil
}

// LoadFromFile opens a file and passes it load.
// It returns an error if a file is not located.
func LoadFromFile(file string) (*Repo, error) {
	f, err := os.Open(file)
	defer f.Close()

	if err != nil {
		return nil, err
	}

	base := filepath.Base(file)

	if ext := filepath.Ext(base); ext != ".sql" {
		return nil, errors.New(
			"not a recognized file type, only .sql files",
		)
	}

	return loadToScanner(f)
}

// Merge takes one or more *Repo and merges its sqlStatements.
// The merge is done in sequence so the last prepared statement will override
// any existing entry that matches it's given tag or name
func Merge(repos ...*Repo) *Repo {
	sqlStatements := make(map[string]string)

	for _, repo := range repos {
		for k, v := range repo.sqlStatements {
			sqlStatements[k] = v
		}
	}

	return &Repo{
		sqlStatements: sqlStatements,
	}
}
