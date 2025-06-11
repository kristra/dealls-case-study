package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReimbursement_DateOnlyString(t *testing.T) {
	date := time.Date(2025, 6, 11, 15, 4, 5, 0, time.UTC)
	a := Reimbursement{Date: date}
	assert.Equal(t, "2025-06-11", a.DateOnlyString())
}
