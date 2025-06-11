package seed

import (
	"dealls-case-study/internal/models"
	"dealls-case-study/internal/utils"
	"fmt"
	"log"
	"time"

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
			Salary:    float64(i * 1000),
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

func Seed44(db *gorm.DB) {
	log.Println("Seeding data to database...")

	const workingHoursPerDay = 8
	const overtimeHours = 2
	userID := uint(44)
	year := 2024
	month := 9

	startDate := time.Date(2025, 9, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC)

	attendanceCount := 0
	weekendOTCount := 0

	var overtimes = []models.Overtime{}
	var attendances = []models.Attendance{}
	var reimbursements = []models.Reimbursement{
		{UserID: userID, Amount: 400, Date: startDate},
		{UserID: userID, Amount: 40, Date: startDate},
		{UserID: userID, Amount: 4, Date: startDate},
	}

	for d := 0; d < 31; d++ {
		current := startDate.AddDate(0, 0, d)
		if current.Month() != time.Month(month) {
			break
		}

		switch current.Weekday() {
		case time.Saturday, time.Sunday:
			// Add 4 weekend overtime entries
			if weekendOTCount < 4 {
				overtimes = append(overtimes, models.Overtime{
					UserID:      userID,
					Date:        current,
					HoursWorked: overtimeHours,
					CreatedBy:   userID,
				})

				weekendOTCount++
			}
		default:
			// Add up to 20 working day attendances
			if attendanceCount < 20 {
				checkIn := time.Date(current.Year(), current.Month(), current.Day(), 9, 0, 0, 0, time.UTC)
				checkOut := checkIn.Add(time.Hour * workingHoursPerDay)

				attendances = append(attendances, models.Attendance{
					UserID:     userID,
					Date:       current,
					CheckInAt:  &checkIn,
					CheckOutAt: &checkOut,
					CreatedBy:  userID,
				})
				attendanceCount++
			}
		}
	}

	payroll := models.Payroll{
		Month:       month,
		Year:        year,
		CreatedBy:   999,
		PeriodStart: startDate,
		PeriodEnd:   endDate,
	}
	db.Create(&payroll)

	db.Create(overtimes)
	db.Create(attendances)
	db.Create(reimbursements)

	return
}
