package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAttendance_DateOnlyString(t *testing.T) {
	date := time.Date(2025, 6, 11, 15, 4, 5, 0, time.UTC)
	a := Attendance{Date: date}
	assert.Equal(t, "2025-06-11", a.DateOnlyString())
}

func TestAttendance_BeforeSave_HoursWorkedCalculation(t *testing.T) {
	checkIn := time.Date(2025, 6, 11, 9, 0, 0, 0, time.UTC)
	checkOut := time.Date(2025, 6, 11, 17, 0, 0, 0, time.UTC)

	a := Attendance{
		CheckInAt:  &checkIn,
		CheckOutAt: &checkOut,
	}

	err := a.BeforeSave(&gorm.DB{})
	assert.Nil(t, err)
	assert.Equal(t, 8.0, a.HoursWorked)
}

func TestAttendance_BeforeSave_InvalidDuration(t *testing.T) {
	checkIn := time.Date(2025, 6, 11, 17, 0, 0, 0, time.UTC)
	checkOut := time.Date(2025, 6, 11, 9, 0, 0, 0, time.UTC)

	a := Attendance{
		CheckInAt:  &checkIn,
		CheckOutAt: &checkOut,
	}

	err := a.BeforeSave(&gorm.DB{})
	assert.Nil(t, err)
	assert.Equal(t, 0.0, a.HoursWorked)
}

func TestAttendance_BeforeSave_NoCheckOut(t *testing.T) {
	checkIn := time.Date(2025, 6, 11, 9, 0, 0, 0, time.UTC)

	a := Attendance{
		CheckInAt: &checkIn,
	}

	err := a.BeforeSave(&gorm.DB{})
	assert.Nil(t, err)
	assert.Equal(t, 0.0, a.HoursWorked)
}
