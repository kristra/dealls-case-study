package seed

import (
	"dealls-case-study/internal/models"
	"log"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	roles := []models.Role{
		{Name: "Admin"},
		{Name: "Employee"},
	}

	// seed roles
	for _, role := range roles {
		err := db.FirstOrCreate(&role, models.Role{Name: role.Name}).Error
		if err != nil {
			log.Printf("Failed to seed role %s: %v", role.Name, err)
			return err
		}
	}

	

	return nil
}
