# pgmig

**Migration management library using [go-pg](https://github.com/go-pg/pg) and [pkger](https://github.com/markbates/pkger)**

Provides a simple interface for creating, running and packaging Postgres schema migrations. Generated migrations are `up` and `down` SQL files, packaged into the project binary using pkger. Put the migration files directory wherever you want in the project.

## Dependencies

```bash
go get github.com/markbates/pkger/cmd/pkger
```

## Example Usage

```go
package main

import (
	"github.com/go-pg/pg/v9"
	"github.com/markbates/pkger"

	"github.com/jacobmoe/pgmig"
)

func main() {
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "dbuser",
		Password: "dbpassword",
		Database: "dbname",
	})

	// path to migrations dir is from project root
	mig := pgmig.New(db, "/path/to/migrations/dir")

	// initialize the migrations table in your db
	// only need to be run this once for a database
	err := mig.Init()
	check(err)

	// create new up and down migration files:
	//  - /path/to/migrations/dir/200405153854_create_users.up.sql
	//  - /path/to/migrations/dir/200405153854_create_users.down.sql
	err = mig.Create("create_users")
	check(err)

	// run new migrations (after updating the migration files)
	err = mig.Migrate()
	check(err)

	// rollback most recent migration
	err = mig.Rollback()
	check(err)
}

// this function never needs to be called. used for pkger static analysis.
func pkgerinclude() {
	// pkger uses static analysis to determine what to include
	// in the packaged file, while pgmig builds migrations
	// dynamically. so, must explicitly include migrations dir.
	pkger.Include("/path/to/migrations/dir")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
```
