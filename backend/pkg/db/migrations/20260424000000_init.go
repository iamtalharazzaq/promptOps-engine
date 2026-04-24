package migrations

import (
	"context"

	"github.com/promptops/backend/pkg/models"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		models := []interface{}{
			(*models.User)(nil),
			(*models.Chat)(nil),
		}

		for _, model := range models {
			_, err := db.NewCreateTable().
				Model(model).
				IfNotExists().
				Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		models := []interface{}{
			(*models.Chat)(nil),
			(*models.User)(nil),
		}

		for _, model := range models {
			_, err := db.NewDropTable().
				Model(model).
				IfExists().
				Exec(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
