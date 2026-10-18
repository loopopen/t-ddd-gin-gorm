package migrate

import (
	"fmt"
	"log/slog"

	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/infra/gorm/po"

	"gorm.io/gorm"
)

type TableCommentter interface {
	TableComment() string
}

func RunMigrate(db *gorm.DB) error {
	if err := migrateWithComment(db,
		&po.User{},
	); err != nil {
		return err
	}

	// if !db.Migrator().HasIndex(&po.User{}, "<idx-name>") {
	// 	if err := db.Exec("CREATE UNIQUE INDEX <idx-name> ON user(<colum-list>)").Error; err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func migrateWithComment(db *gorm.DB, models ...any) error {
	for _, model := range models {
		stmt := &gorm.Statement{DB: db}
		if err := stmt.Parse(model); err != nil {
			return err
		}
		tableName := stmt.Schema.Table

		var tableComment string
		if tc, ok := model.(TableCommentter); ok {
			tableComment = tc.TableComment()
		}

		if err := db.AutoMigrate(model); err != nil {
			return err
		}

		if tableComment != "" {
			var sql string
			switch db.Dialector.Name() {
			case "mysql":
				sql = fmt.Sprintf("ALTER TABLE `%s` COMMENT = '%s'", tableName, tableComment)
			// case "postgres":
			// 	sql = fmt.Sprintf("COMMENT ON TABLE \"%s\" IS '%s'", tableName, tableComment)
			default:
				slog.Warn("unsupported dialector for table comment", "dialector", db.Dialector.Name())
				continue
			}

			if err := db.Exec(sql).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
