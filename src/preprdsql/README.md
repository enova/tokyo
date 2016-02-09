======

# Prepared SQL (preprdsql)

This package is not a DB interface like [sqlx] or [sql] nor is it a
query builder. This package was created with the intent of being able to
organize prepared SQL statements by giving the ability to load  the file and
search for the needed prepared statements via tags.

The benefits of using this package over building your queries through your go
code are:

1. Easily maintain your queries that your application uses.
2. Scalability is easier because of the separation of your SQL code from your
   go code.
3. Version control your prepared statements.
4. Syntax highlighting for your SQL code.

## USAGE

Create your prepared statements using the reference below:

### Tag -- name: foo-sql

```SQL
-- name: create-user
INSERT INTO users (name, email) VALUES (?, ?);

-- name: select-all-users-list
SELECT u.name, u.email, a.city, a.state
FROM users u
INNER JOIN addresses a ON u.user_id = a.user_id

-- name: select-a-user
SELECT * FROM users WHERE name = ?
```

```go
sqlRepo, err := preprdsql.LoadFromFile("sqls/users-table.sql")
sql, err := sqlRepo.LookupSQL("create-user")

db.Query(sql)
```

You need to add a tag for each sql file

### Using with [sqlx] or [sql]

Integration with [sqlx] or [sql] is very straight forward. When
executing a `func` from either packages simply pass in the looked up sql
statement from the `repo` with the loaded SQL files.

See below for reference:

```go
  db, err := sqlx.Connect("sqlite3", ":memory:")

  if  err != nil {
      panic(err)
  }

  repo, err = preprdsql.LoadFromFile("manbearpig.sql")

  if  err != nil {
      panic(err)
  }

  sql = repo.LookupSQL("select-manbearpig-in-a-city")

  // optional args if you are passing in named statements or parametarized statements
  result = db.QueryRow(sql, ...args)
```

[sqlx]: https://github.com/jmoiron/sqlx
[sql]: https://github.com/golang/go/wiki/SQLInterface
