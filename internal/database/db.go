package database

import (
	"log"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres" // Use alias
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres" // Add this import for GORM
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

func InitDB() *gorm.DB {
	dbOnce.Do(func() {
		dsn := "host=localhost user=postgres password=secret dbname=mydb port=5432 sslmode=disable"
		
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("Database connection failed:", err)
		}
		
		// Run migrations
		if err := runMigrations(db); err != nil {
			log.Fatal("Migrations failed:", err)
		}
		
		log.Println("Database connected and migrated successfully")
	})
	return db
}

func runMigrations(gormDB *gorm.DB) error {
	// Get underlying sql.DB
	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}

	// Setup migrations
	driver, err := migratePostgres.WithInstance(sqlDB, &migratePostgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", 
		driver,
	)
	if err != nil {
		return err
	}

	// Run migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}

func GetDB() *gorm.DB {
	if db == nil {
		return InitDB()
	}
	return db
}

func CloseDB() {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Error getting SQL DB: %v", err)
			return
		}
		sqlDB.Close()
		log.Println("Database connection closed")
	}
}