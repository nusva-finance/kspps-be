package services

import (
	"fmt"
	"time"

	"nusvakspps/config"
	"nusvakspps/models"

	"gorm.io/gorm"
)

// GenerateMemberNo generates member number in format: YYMMXXZZZZZZ
// YY = 2 digit year, MM = 2 digit month, XX = 01 (Laki-laki) / 02 (Perempuan), ZZZZZZ = 6 digit sequence
func GenerateMemberNo(joinDate time.Time, gender string) (string, error) {
	db := config.GetDB()

	yy := joinDate.Format("06")
	mm := joinDate.Format("01")

	var genderCode string
	switch gender {
	case models.GenderMale:
		genderCode = "01"
	case models.GenderFemale:
		genderCode = "02"
	default:
		genderCode = "01"
	}

	// Use transaction with SELECT FOR UPDATE
	var counter models.MemberSequenceCounter
	err := db.Transaction(func(tx *gorm.DB) error {
		// Lock the counter row
		err := tx.Clauses(gorm.Locking{Strength: "UPDATE"}).
			Where("year_code = ? AND month_code = ? AND gender_code = ?", yy, mm, genderCode).
			First(&counter).Error

		if err == gorm.ErrRecordNotFound {
			// Create new counter
			counter = models.MemberSequenceCounter{
				YearCode:   yy,
				MonthCode:  mm,
				GenderCode: genderCode,
				LastSeq:    0,
			}
			if err := tx.Create(&counter).Error; err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		// Increment sequence
		counter.LastSeq++
		if err := tx.Save(&counter).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate member number: %w", err)
	}

	return fmt.Sprintf("%s%s%s%06d", yy, mm, genderCode, counter.LastSeq), nil
}
