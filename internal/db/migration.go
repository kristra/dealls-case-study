package db

import (
	"dealls-case-study/internal/models"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "202506100930",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&models.Attendance{}, &models.Overtime{}, &models.Payroll{}, &models.Payslip{}, &models.Reimbursement{}, &models.Role{}, &models.User{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&models.Attendance{}, &models.Overtime{}, &models.Payroll{}, &models.Payslip{}, &models.Reimbursement{}, &models.Role{}, &models.User{})
			},
		},
	})

	return m.Migrate()
}
