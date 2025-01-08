package database

import (
	"fmt"
	"log"
	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

var db *gorm.DB

func InitDB(config configs.PostgreSQL) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.Username,
		config.Password,
		config.Database,
		config.Port,
		config.SSLMode,
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	err = db.AutoMigrate(
        &entities.NursingHouse{},
		&entities.Role{},
		&entities.User{},
		&entities.Image{},
		&entities.News{},
		&entities.Dialog{},
		&entities.Favorite{},
    )
	insertRoles()
	log.Println("Database connection established successfully!")
}

func GetDB() *gorm.DB {
	if db == nil {
		log.Fatal("Database is not initialized")
	}
	
	return db
}

func insertRoles() {
	var adminRole entities.Role
	var userRole entities.Role

	if err := db.First(&adminRole, "role_name = ?", "Admin").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			adminRole = entities.Role{RoleName: "Admin"}
			if err := db.Create(&adminRole).Error; err != nil {
				log.Fatalf("Failed to insert Admin role: %v", err)
			}
			log.Println("Admin role created successfully!")
		} else {
			log.Fatalf("Error checking Admin role: %v", err)
		}
	}

	if err := db.First(&userRole, "role_name = ?", "User").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			userRole = entities.Role{RoleName: "User"}
			if err := db.Create(&userRole).Error; err != nil {
				log.Fatalf("Failed to insert User role: %v", err)
			}
			log.Println("User role created successfully!")
		} else {
			log.Fatalf("Error checking User role: %v", err)
		}
	}
}