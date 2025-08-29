package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/oktaharis/uji-teknis-godigi/internal/config"
	"github.com/oktaharis/uji-teknis-godigi/internal/models"
)

func Connect(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(mysql.Open(cfg.DBDSN), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,
	})
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}

	// Pastikan tabel perusahaan ada (tanpa ALTER)
	if err := EnsureCompanyTables(db); err != nil {
		log.Fatalf("bootstrap company tables error: %v", err)
	}

	// Migrasi tabel internal saja
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").
		AutoMigrate(
			&models.User{},
			&models.PasswordReset{},
			&models.Project{},
		); err != nil {
		log.Fatalf("auto-migrate error: %v", err)
	}

	return db
}
