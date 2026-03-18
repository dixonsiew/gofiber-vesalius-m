package shared

import (
    "github.com/nleeper/goment"
)

func LittleKidsEligibleAge(dob string) bool {
    // only 12 years old and below
    // subtract another 1 day because same date of 12 years back is invalid because of hour/min differences
    now, _ := goment.New()
    minDate := now.Subtract(12, "years")
    minDate = minDate.Subtract(1, "days")

    dobDate, err := goment.New(dob, "DD/MM/YYYY")
    if err != nil {
        return false
    }

    return dobDate.IsAfter(minDate)
}

func GoldenPearlEligibleAge(dob string) bool {
    // only 60 years old and above
    now, _ := goment.New()
    minDate := now.Subtract(60, "years")

    dobDate, err := goment.New(dob, "DD/MM/YYYY")
    if err != nil {
        return false
    }

    return dobDate.IsBefore(minDate)
}