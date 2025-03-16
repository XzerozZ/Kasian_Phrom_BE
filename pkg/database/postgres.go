package database

import (
	"fmt"
	"log"

	"github.com/XzerozZ/Kasian_Phrom_BE/configs"
	"github.com/XzerozZ/Kasian_Phrom_BE/modules/entities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	_ = db.AutoMigrate(
		&entities.NursingHouse{},
		&entities.Role{},
		&entities.User{},
		&entities.Image{},
		&entities.News{},
		&entities.Dialog{},
		&entities.Favorite{},
		&entities.Asset{},
		&entities.RetirementPlan{},
		&entities.SelectedHouse{},
		&entities.OTP{},
		&entities.Loan{},
		&entities.History{},
		&entities.Risk{},
		&entities.Quiz{},
		&entities.Transaction{},
		&entities.Notification{},
		&entities.NursingHouseHistory{},
	)

	insertRoles()
	insertRisk()
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

func insertRisk() {
	var Role1 entities.Risk
	var Role2 entities.Risk
	var Role3 entities.Risk
	var Role4 entities.Risk
	var Role5 entities.Risk

	if err := db.First(&Role1, "risk_name = ?", "ความเสี่ยงต่ำ").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			Role1 = entities.Risk{RiskName: "ความเสี่ยงต่ำ"}
			if err := db.Create(&Role1).Error; err != nil {
				log.Fatalf("Failed to insert risk name: %v", err)
			}

			log.Println("Risk name created successfully!")
		} else {
			log.Fatalf("Error checking risk name: %v", err)
		}
	}

	if err := db.First(&Role2, "risk_name = ?", "ความเสี่ยงปานกลางค่อนข้างต่ำ").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			Role2 = entities.Risk{RiskName: "ความเสี่ยงปานกลางค่อนข้างต่ำ"}
			if err := db.Create(&Role2).Error; err != nil {
				log.Fatalf("Failed to insert risk name: %v", err)
			}

			log.Println("Risk name created successfully!")
		} else {
			log.Fatalf("Error checking risk name: %v", err)
		}
	}

	if err := db.First(&Role3, "risk_name = ?", "ความเสี่ยงปานกลางค่อนข้างสูง").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			Role3 = entities.Risk{RiskName: "ความเสี่ยงปานกลางค่อนข้างสูง"}
			if err := db.Create(&Role3).Error; err != nil {
				log.Fatalf("Failed to insert risk name: %v", err)
			}

			log.Println("Risk name created successfully!")
		} else {
			log.Fatalf("Error checking risk name: %v", err)
		}
	}

	if err := db.First(&Role4, "risk_name = ?", "ความเสี่ยงสูง").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			Role4 = entities.Risk{RiskName: "ความเสี่ยงสูง"}
			if err := db.Create(&Role4).Error; err != nil {
				log.Fatalf("Failed to insert risk name: %v", err)
			}

			log.Println("Risk name created successfully!")
		} else {
			log.Fatalf("Error checking risk name: %v", err)
		}
	}

	if err := db.First(&Role5, "risk_name = ?", "ความเสี่ยงสูงมาก").Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			Role5 = entities.Risk{RiskName: "ความเสี่ยงสูงมาก"}
			if err := db.Create(&Role5).Error; err != nil {
				log.Fatalf("Failed to insert risk name: %v", err)
			}

			log.Println("Risk name created successfully!")
		} else {
			log.Fatalf("Error checking risk name: %v", err)
		}
	}
}
