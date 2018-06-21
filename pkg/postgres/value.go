package postgres

import (
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/l-vitaly/imhio-task/pkg/value"
	"github.com/pkg/errors"
)

var valueMigrations = []migrations.Migration{
	{
		Version: 1,
		Up: func(db orm.DB) error {
			_, err := db.Exec(`
				CREATE TABLE values
				(
				id        SERIAL                      NOT NULL
					      CONSTRAINT values_pkey
					      PRIMARY KEY,
				catalog_id INTEGER                     NOT NULL,
				name      VARCHAR(255)                NOT NULL,
				data      JSONB DEFAULT '{}' :: JSONB NOT NULL
				);
				CREATE UNIQUE INDEX values_name_uindex ON values (name);

				INSERT INTO values (id, catalog_id, name, data) VALUES (1, 1, 'Database.processing', '{"host": "localhost", "port": "5432", "user": "mr_robot", "schema": "public", "database": "devdb", "password": "secret"}');
			`)
			return err
		},
		Down: func(db orm.DB) error {
			_, err := db.Exec(`
				DROP TABLE values;
			`)
			return err
		},
	},
}

type valueRepository struct {
	db *pg.DB
}

func (r *valueRepository) Find(catalogID int, name string) (*value.Value, error) {
	result := new(value.Value)
	err := r.db.Model(result).Where("catalog_id = ? AND name = ?", catalogID, name).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, value.ErrUnknown
		}
		return nil, err
	}
	return result, nil
}

func (r *valueRepository) CreateSchemas() (int64, int64, error) {
	migrations.SetTableName("values_migration")
	oldVersion, newVersion, err := migrations.RunMigrations(r.db, valueMigrations)
	if err != nil {
		return 0, 0, errors.Wrap(err, "Could not migrate sql schema values")
	}
	return oldVersion, newVersion, nil
}

// NewValueRepository returns a new instance of a PostgreSQL value repository.
func NewValueRepository(db *pg.DB) value.Repository {
	return &valueRepository{
		db: db,
	}
}
