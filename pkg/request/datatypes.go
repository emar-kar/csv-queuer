package request

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	typeInteger = "INT"
	typeFloat   = "FLOAT"
	typeString  = "STRING"
	typeDate    = "DATE"
)

// Variable is an interface which determine methods of the Condition Value.
type Variable interface {
	defineType() string
	isInteger() bool
	toInteger() int
	isFloat() bool
	toFloat() float64
	isDate() bool
	toDate() *Date
}

// Date is a wrapper around the date.
// It defines methods and fields to work with dates easily.
type Date struct {
	Year  int
	Month int
	Day   int
}

// Not defines if dates are not equal.
func (d *Date) Not(ad *Date) bool {
	if d.Year != ad.Year || d.Month != ad.Month || d.Day != ad.Day {
		return true
	}
	return false
}

// Equal defines if dates are Equal.
func (d *Date) Equal(ad *Date) bool {
	if d.Year == ad.Year && d.Month == ad.Month && d.Day == ad.Day {
		return true
	}
	return false
}

// Greater defines if the given date is less than the main one.
func (d *Date) Greater(ad *Date) bool {
	if d.Year > ad.Year {
		return true
	}
	if d.Year == ad.Year && d.Month > ad.Month {
		return true
	}
	if d.Year == ad.Year && d.Month == ad.Month && d.Day > ad.Day {
		return true
	}
	return false
}

// Less defines if the given date is greater than the main one.
func (d *Date) Less(ad *Date) bool {
	if d.Year < ad.Year {
		return true
	}
	if d.Year == ad.Year && d.Month < ad.Month {
		return true
	}
	if d.Year == ad.Year && d.Month == ad.Month && d.Day < ad.Day {
		return true
	}
	return false
}

// GreaterOrEqual defines if the given date is less or equals the main one.
func (d *Date) GreaterOrEqual(ad *Date) bool {
	if d.Greater(ad) || d.Equal(ad) {
		return true
	}

	return false
}

// LessOrEqual defines if the given date is greater or equals the main one.
func (d *Date) LessOrEqual(ad *Date) bool {
	if d.Less(ad) || d.Equal(ad) {
		return true
	}

	return false
}

// Data is a redefined custom type from string.
// Implements Variable interface.
type Data string

func (d Data) defineType() string {
	if d.isInteger() {
		return typeInteger
	}
	if d.isFloat() {
		return typeFloat
	}
	if d.isDate() {
		return typeDate
	}
	return typeString
}

func (d Data) isInteger() bool {
	if _, err := strconv.Atoi(string(d)); err != nil {
		return false
	}

	return true
}

func (d Data) toInteger() int {
	num, err := strconv.Atoi(string(d))
	if err != nil {
		return 0
	}

	return num
}

func (d Data) isFloat() bool {
	if _, err := strconv.ParseFloat(string(d), 64); err != nil {
		return false
	}

	return true
}

func (d Data) toFloat() float64 {
	num, _ := strconv.ParseFloat(string(d), 64)
	return num
}

func (d Data) isDate() bool {
	match, err := regexp.MatchString("[1-2][0-9][0-9][0-9]-[0-1][0-9]-[0-3][0-9]", string(d))
	if err != nil || !match {
		return false
	}
	return true
}

func (d Data) toDate() *Date {
	dateStringSlice := strings.Split(string(d), "-")
	dateIntSlice := stringSliceToInt(dateStringSlice)
	return &Date{
		Year:  dateIntSlice[0],
		Month: dateIntSlice[1],
		Day:   dateIntSlice[2],
	}
}

func stringSliceToInt(s []string) []int {
	newSlice := make([]int, len(s))
	for ind, el := range s {
		num, _ := strconv.Atoi(el)
		newSlice[ind] = num
	}

	return newSlice
}
