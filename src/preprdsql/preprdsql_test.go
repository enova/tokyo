package preprdsql

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func failIfError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("err == nil, got '%s'", err)
	}
}

func failIfNoError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("err == nil, got '%s'", err)
	}
}

// loadFromString is a test helper that loads sql statements passed in as a
// string to the Repo.
func loadFromString(sql string) (*Repo, error) {
	buf := bytes.NewBufferString(sql)
	return loadToScanner(buf)
}

func TestLoad(t *testing.T) {
	_, err := loadToScanner(strings.NewReader(""))
	failIfError(t, err)
}

func TestLoadFromFile(t *testing.T) {
	file, err := os.Create("somesqlstatement.sql")
	defer os.Remove(file.Name())

	if err != nil {
		panic(err)
	}

	_, err = LoadFromFile(file.Name())
	failIfError(t, err)
}

func TestLoadFromFile_FileNotFound(t *testing.T) {
	_, err := LoadFromFile("/manbearpig/globalwarming.sql")
	failIfNoError(t, err)
}

func TestLoadFromFile_InvalidFileType(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "notadotsqlfile")
	defer os.Remove(file.Name())

	if err != nil {
		panic(err)
	}

	_, err = LoadFromFile(file.Name())
	failIfNoError(t, err)
}

func TestLookupSQL(t *testing.T) {
	expectedQuery := "SELECT 1+1"

	dot, err := loadFromString("--name: foo-query\n" + expectedQuery)
	failIfError(t, err)

	got, err := dot.LookupSQL("boo-query")
	failIfNoError(t, err)

	got, err = dot.LookupSQL("foo-query")
	failIfError(t, err)

	got = strings.TrimSpace(got)
	if got != expectedQuery {
		t.Errorf("Raw() == '%s', expected '%s'", got, expectedQuery)
	}
}

func TestMergeHaveBothQueries(t *testing.T) {
	expectedSQLStatements := map[string]string{
		"query-a": "SELECT * FROM a",
		"query-b": "SELECT * FROM b",
	}

	a, err := loadFromString("--name: query-a\nSELECT * FROM a")
	failIfError(t, err)

	b, err := loadFromString("--name: query-b\nSELECT * FROM b")
	failIfError(t, err)

	c := Merge(a, b)

	got := c.sqlStatements
	if len(got) != len(expectedSQLStatements) {
		t.Errorf("SQLStatements() len (%d) differ from expected (%d)", len(got), len(expectedSQLStatements))
	}
}

func TestMergeTakesPresecendeFromLastArgument(t *testing.T) {
	expectedQuery := "SELECT * FROM c"

	a, err := loadFromString("--name: query\nSELECT * FROM a")
	failIfError(t, err)

	b, err := loadFromString("--name: query\nSELECT * FROM b")
	failIfError(t, err)

	c, err := loadFromString("--name: query\nSELECT * FROM c")
	failIfError(t, err)

	x := Merge(a, b, c)

	got := x.sqlStatements["query"]
	if expectedQuery != got {
		t.Errorf("Expected query: '%s', got: '%s'", expectedQuery, got)
	}
}
