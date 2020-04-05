# pgmig

*Migration management library for [go-pg](https://github.com/go-pg/pg) using [pkger](https://github.com/markbates/pkger)*

Simple interface for creating and running Postgres schema migrations. Generated migrations are up and down SQL files, packaged into the project binary using [pkger](https://github.com/markbates/pkger). Put the migration files directory wherever you want in the project.

## Dependencies

- [pkger](https://github.com/markbates/pkger)

```bash
go get github.com/markbates/pkger/cmd/pkger
```

## Example Usage

```go
import (
	"github.com/go-pg/pg"
    "github.com/markbates/pkger"

    "github.com/jacobmoe/pgmig"
)

func main() {
	db := pg.Connect(&pg.Options{
		Addr:     addr,
		User:     user,
		Password: password,
		Database: database,
	})
    
    // path to migrations dir is from project root
    mig := pgmig.New(db, "/path/to/migrations/dir")
    
    // initialize the migrations table in your db
    // only need to be run this once for a database
    mig.Init()
    
    // create new up and down migration files:
    //  - /path/to/migrations/dir/200405153854_create_users.up.sql
    //  - /path/to/migrations/dir/200405153854_create_users.down.sql
    mig.Create("create_users")
    
    // run new migrations (after updating the migration files)
    mig.Migrate()
    
    // rollback most recent migration
    mig.Rollback()
}

// this function never need to be called. used for pkger static analysis.
func pkgerinclude() {
	// pkger uses static analysis to determine what to include
	// in the packaged file, while pgmig builds migrations
	// dynamically. so, must explicitly include migrations dir.
	pkger.Include("/path/to/migrations/dir")
}
```
