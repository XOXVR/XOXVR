package db

import (
	"fmt"
	"log"

	"xox_tokens/db/types"
	"xox_tokens/configs"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

var (
	Instance *pg.DB
)

// Initiate creates DB instance.
func Initiate(conf configs.DB) (err error) {
	log.Println("initiate database...")
	Instance = pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		User:     conf.User,
		Database: conf.Name,
		Password: conf.Password,
	})
	err = createSchema()
	if err != nil { return }

	log.Println("database initialization complete")

	return
}

func createSchema() error {
	for _, model := range []interface{}{
		&types.Buyer{},
		&types.Transaction{},
		&types.NotHandledTransaction{},
	} {
		err := Instance.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// IsNotFoundError return true if no rows found.
func IsNotFoundError(err error) bool {
	return err.Error() == pg.ErrNoRows.Error()
}
