package seed

import (
	"dealls-case-study/internal/models"
	"dealls-case-study/internal/utils"
	"fmt"
	"log"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	log.Println("Seeding data to database...")
	roles := []models.Role{
		{Name: "Admin", CreatedBy: 999},
		{Name: "Employee", CreatedBy: 999},
	}

	// seed roles
	for _, role := range roles {
		err := db.FirstOrCreate(&role, models.Role{Name: role.Name}).Error
		if err != nil {
			log.Printf("Failed to seed role %s: %v", role.Name, err)
			return err
		}
	}

	// seed admin
	password, err := utils.HashPassword("password")
	if err != nil {
		log.Printf("Failed creating admin password: %v", err)
		return err
	}
	admin := models.User{
		ID:       999,
		Username: "admin",
		Password: password,
		RoleID:   1,
	}
	err = db.FirstOrCreate(&admin, models.User{ID: admin.ID}).Error
	if err != nil {
		log.Printf("Failed to seed admin: %v", err)
		return err
	}

	// seed users
	var users = []models.User{}
	for i := uint(1); i <= 100; i++ {
		username := fmt.Sprintf("%s%d", "user", i)

		if err != nil {
			log.Printf("Failed to seed user %s: %v", username, err)
			return err
		}

		user := models.User{
			ID:        i,
			Username:  username,
			Password:  password,
			Salary:    float64(i * 10),
			RoleID:    uint(2),
			CreatedBy: 999,
		}
		users = append(users, user)
	}

	err = db.Create(&users).Error
	if err != nil {
		log.Println("Failed to seed users: ", err)
		return err
	}
	log.Println("Seeding completed!")

	return nil
}
