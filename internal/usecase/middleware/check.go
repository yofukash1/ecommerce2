package middleware

import (
	"errors"
	"strconv"
	"time"
)

func CheckExpirationDate(expirationDate string) error {
	currentTime := time.Now()
	f := currentTime.Format("2006-01-02")
	year := f[2:4]
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		return err
	}

	month := f[5:7]
	monthInt, err := strconv.Atoi(month)
	if err != nil {
		return err
	}

	expYear := expirationDate[3:]
	expYearInt, err := strconv.Atoi(expYear)
	if err != nil {
		return err
	}

	expMonth := expirationDate[:2]
	expMonthInt, err := strconv.Atoi(expMonth)
	if err != nil {
		return err
	}

	if expYearInt < yearInt {
		return errors.New("invalid expiration date")
	}

	if expYearInt == yearInt && expMonthInt < monthInt {
		return errors.New("invalid expiration date")
	}
	return nil

}
