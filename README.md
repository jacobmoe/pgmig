# pgmig

**Migration management library using [go-pg](https://github.com/go-pg/pg)**

Provides a simple interface for creating, running and packaging Postgres schema migrations. Generated migrations are `up` and `down` SQL files.

## Dependencies

go 1.16+

## Example Usage

```go
package main

import (
	"embed"
	"github.com/go-pg/pg/v10"
	"github.com/jacobmoe/pgmig"
)

//go:embed path/to/migrations/dir
var migrationsFS embed.FS

func main() {
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "dbuser",
		Password: "dbpassword",
		Database: "dbname",
	})

	// initialize new migrator
	mig := pgmig.New(db, migrationsFS)

	// initialize the migrations table in your db
	// only need to be run this once for a database
	err := mig.Init()
	check(err)

	// create new up and down migration files:
	//  - /path/to/migrations/dir/200405153854_create_users.up.sql
	//  - /path/to/migrations/dir/200405153854_create_users.down.sql
	err = mig.Create("/full/path/to/migrations/dir", "create_users")
	check(err)

	// run new migrations (after updating the migration files)
	err = mig.Migrate()
	check(err)

	// rollback most recent migration
	err = mig.Rollback()
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
```
