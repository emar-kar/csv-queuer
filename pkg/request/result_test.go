package request

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func ExamplePrint() {
	// Create your request string.
	requestString := `
	SELECT location, new_cases, new_deaths, date 
	FROM /Users/lemarkar/Documents/Projects/course_project/test_data/owid-covid-data.csv 
	WHERE location = Russia AND new_cases >= 0 OR new_deaths > 60 AND new_cases <= 6500 AND date > 2020-04-20 AND 
	date < 2020-04-30 AND date NOT 2020-04-25;
	`

	// Create context for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := NewRequest(requestString)
	if err != nil {
		log.Printf("error: %s\n", err)
	}

	result, err := req.Do(ctx, ",")
	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Println("deadline exceeded, try to increase the processing time in config file or specify the request")
	} else if err != nil {
		log.Printf("error: %s\n", err)
	} else if result != nil {
		result.Print()
	}

	// Output
	// ==================================================
	// | location | new_cases | new_deaths |    date    |
	// ==================================================
	// |  Russia  |  6361.0   |    66.0    | 2020-04-26 |
	// |  Russia  |  6411.0   |    73.0    | 2020-04-28 |
	// |  Russia  |  5841.0   |   105.0    | 2020-04-29 |
	// ==================================================
}

type TestVariable struct {
	name      string
	value     Data
	lineValue Data
	result    bool
}

func TestCheckNot(t *testing.T) {
	tests := []TestVariable{
		{name: "notIntegerTrue", value: Data("45"), lineValue: Data("50"), result: true},
		{name: "notIntegerFalse", value: Data("45"), lineValue: Data("45"), result: false},
		{name: "notFloatTrue", value: Data("4.5"), lineValue: Data("5.0"), result: true},
		{name: "notFloatFalse", value: Data("4.5"), lineValue: Data("4.5"), result: false},
		{name: "notStringTrue", value: Data("something"), lineValue: Data("anotherthing"), result: true},
		{name: "notStringFalse", value: Data("something"), lineValue: Data("something"), result: false},
		{name: "notDateTrue", value: Data("2020-11-18"), lineValue: Data("2020-11-19"), result: true},
		{name: "notDateFalse", value: Data("2020-11-18"), lineValue: Data("2020-11-18"), result: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := checkNot(tc.value, tc.lineValue)
			assert.Equal(t, result, tc.result)
		})
	}
}

func TestCheckEqual(t *testing.T) {
	tests := []TestVariable{
		{name: "equalIntegerFalse", value: Data("45"), lineValue: Data("50"), result: false},
		{name: "equalIntegerTrue", value: Data("45"), lineValue: Data("45"), result: true},
		{name: "equalFloatFalse", value: Data("4.5"), lineValue: Data("5.0"), result: false},
		{name: "equalFloatTrue", value: Data("4.5"), lineValue: Data("4.5"), result: true},
		{name: "equalStringFalse", value: Data("something"), lineValue: Data("aequalherthing"), result: false},
		{name: "equalStringTrue", value: Data("something"), lineValue: Data("something"), result: true},
		{name: "equalDateFalse", value: Data("2020-11-18"), lineValue: Data("2020-11-19"), result: false},
		{name: "equalDateTrue", value: Data("2020-11-18"), lineValue: Data("2020-11-18"), result: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := checkEqual(tc.value, tc.lineValue)
			assert.Equal(t, result, tc.result)
		})
	}
}

func TestCheckGreater(t *testing.T) {
	tests := []TestVariable{
		{name: "greaterIntegerTrue", value: Data("45"), lineValue: Data("50"), result: true},
		{name: "greaterIntegerFalse", value: Data("45"), lineValue: Data("40"), result: false},
		{name: "greaterFloatTrue", value: Data("4.0"), lineValue: Data("4.5"), result: true},
		{name: "greaterFloatFalse", value: Data("4.5"), lineValue: Data("3.5"), result: false},
		{name: "greaterDateTrue", value: Data("2020-11-18"), lineValue: Data("2020-11-20"), result: true},
		{name: "greaterDateFalse", value: Data("2020-11-18"), lineValue: Data("2020-11-15"), result: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := checkGreater(tc.value, tc.lineValue)
			assert.Equal(t, result, tc.result)
		})
	}
}

func TestCheckLess(t *testing.T) {
	tests := []TestVariable{
		{name: "lessIntegerTrue", value: Data("50"), lineValue: Data("45"), result: true},
		{name: "lessIntegerFalse", value: Data("40"), lineValue: Data("50"), result: false},
		{name: "lessFloatTrue", value: Data("5.0"), lineValue: Data("4.5"), result: true},
		{name: "lessFloatFalse", value: Data("4.5"), lineValue: Data("5.5"), result: false},
		{name: "lessDateTrue", value: Data("2020-11-18"), lineValue: Data("2020-11-15"), result: true},
		{name: "lessDateFalse", value: Data("2020-11-18"), lineValue: Data("2020-11-20"), result: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := checkLess(tc.value, tc.lineValue)
			assert.Equal(t, result, tc.result)
		})
	}
}

func TestCheckGreaterOrEqual(t *testing.T) {
	tests := []TestVariable{
		{name: "greaterIntegerTrue", value: Data("45"), lineValue: Data("50"), result: true},
		{name: "equalIntegerTrue", value: Data("45"), lineValue: Data("45"), result: true},
		{name: "greaterIntegerFalse", value: Data("45"), lineValue: Data("40"), result: false},
		{name: "greaterFloatTrue", value: Data("4.0"), lineValue: Data("4.5"), result: true},
		{name: "equalFloatTrue", value: Data("4.0"), lineValue: Data("4.0"), result: true},
		{name: "greaterFloatFalse", value: Data("4.5"), lineValue: Data("3.5"), result: false},
		{name: "greaterDateTrue", value: Data("2020-11-18"), lineValue: Data("2020-11-20"), result: true},
		{name: "equalDateTrue", value: Data("2020-11-18"), lineValue: Data("2020-11-18"), result: true},
		{name: "greaterDateFalse", value: Data("2020-11-18"), lineValue: Data("2020-11-15"), result: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := checkGreaterOrEqual(tc.value, tc.lineValue)
			assert.Equal(t, result, tc.result)
		})
	}
}

func TestCheckLessOrEqual(t *testing.T) {
	tests := []TestVariable{
		{name: "lessIntegerTrue", value: Data("50"), lineValue: Data("45"), result: true},
		{name: "equalIntegerTrue", value: Data("50"), lineValue: Data("50"), result: true},
		{name: "lessIntegerFalse", value: Data("40"), lineValue: Data("50"), result: false},
		{name: "lessFloatTrue", value: Data("5.0"), lineValue: Data("4.5"), result: true},
		{name: "equalFloatTrue", value: Data("5.0"), lineValue: Data("5.0"), result: true},
		{name: "lessFloatFalse", value: Data("4.5"), lineValue: Data("5.5"), result: false},
		{name: "lessDateTrue", value: Data("2020-11-18"), lineValue: Data("2020-11-15"), result: true},
		{name: "equalDateTrue", value: Data("2020-11-18"), lineValue: Data("2020-11-18"), result: true},
		{name: "lessDateFalse", value: Data("2020-11-18"), lineValue: Data("2020-11-20"), result: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := checkLessOrEqual(tc.value, tc.lineValue)
			assert.Equal(t, result, tc.result)
		})
	}
}
