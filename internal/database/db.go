package database

import (
	"log"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

// InitDB initializes the database connection and runs migrations.
// Pass the DSN from your config (e.g., from LoadConfig().DBDsn)
func InitDB(dsn string) *gorm.DB {
	dbOnce.Do(func() {
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("[error] Database connection failed: %v", err)
		}

		if err := runMigrations(db); err != nil {
			log.Fatalf("[error] Migrations failed: %v", err)
		}

		log.Println("✅ Database connected and migrated successfully")
	})
	return db
}

func runMigrations(gormDB *gorm.DB) error {
	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}

	driver, err := migratePostgres.WithInstance(sqlDB, &migratePostgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations", // make sure your migrations folder is correctly mounted
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("✅ Migrations applied successfully")
	return nil
}

// GetDB returns the singleton DB instance
func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database not initialized. Call InitDB(dsn) first.")
	}
	return db
}

// CloseDB closes the underlying SQL connection
func CloseDB() {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("[warn] Error getting SQL DB: %v", err)
			return
		}
		sqlDB.Close()
		log.Println("✅ Database connection closed")
	}
}
