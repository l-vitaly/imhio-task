package postgres

import (
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/l-vitaly/imhio-task/pkg/catalog"
	"github.com/pkg/errors"
)

var catalogMigrations = []migrations.Migration{
	{
		Version: 1,
		Up: func(db orm.DB) error {
			_, err := db.Exec(`
				CREATE TABLE catalogs
				(
					id SERIAL PRIMARY KEY,
					type VARCHAR(255) NOT NULL
				);
				CREATE UNIQUE INDEX catalogs_type_uindex ON catalogs (type);
				
				INSERT INTO catalogs (id, type) VALUES (1, 'Develop.mr_robot');
			`)
			return err
		},
		Down: func(db orm.DB) error {
			_, err := db.Exec(`
				DROP TABLE catalogs;
			`)
			return err
		},
	},
}

type catalogRepository struct {
	db *pg.DB
}

func (r *catalogRepository) Find(t string) (*catalog.Catalog, error) {
	result := new(catalog.Catalog)
	err := r.db.Model(result).Where("type = ?", t).Select()
	if err != nil {
		if err == pg.ErrNoRows {
			return nil, catalog.ErrUnknown
		}
		return nil, err
	}
	return result, nil
}

func (r *catalogRepository) CreateSchemas() (int64, int64, error) {
	migrations.SetTableName("catalogs_migration")
	oldVersion, newVersion, err := migrations.RunMigrations(r.db, catalogMigrations)
	if err != nil {
		return 0, 0, errors.Wrap(err, "Could not migrate sql schema catalogs")
	}
	return oldVersion, newVersion, nil
}

// NewCatalogRepository returns a new instance of a Postgres catalog repository.
func NewCatalogRepository(db *pg.DB) catalog.Repository {
	return &catalogRepository{
		db: db,
	}
}
