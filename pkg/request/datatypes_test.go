package request

import (
	"reflect"
	"testing"

	"github.com/magiconair/properties/assert"
)

var (
	testString = Data("test")
	testInt    = Data("4")
	testFloat  = Data("4.5")
	testDate   = Data("2020-11-18")
)

func TestDataDefineType(t *testing.T) {
	tests := []struct {
		name   string
		data   Variable
		expect string
	}{
		{name: "string", data: testString, expect: typeString},
		{name: "int", data: testInt, expect: typeInteger},
		{name: "float", data: testFloat, expect: typeFloat},
		{name: "date", data: testDate, expect: typeDate},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dataType := tc.data.defineType()
			assert.Equal(t, dataType, tc.expect)
		})
	}
}

func TestDataToInteger(t *testing.T) {
	num := testInt.toInteger()
	assert.Equal(t, num, 4)
}

func TestDataToIntegerError(t *testing.T) {
	data := Data("45g")
	num := data.toInteger()
	assert.Equal(t, num, 0)
}

func TestDataToFloat(t *testing.T) {
	num := testFloat.toFloat()
	assert.Equal(t, num, 4.5)
}

func TestDataToDate(t *testing.T) {
	date := testDate.toDate()
	want := &Date{Year: 2020, Month: 11, Day: 18}
	assert.Equal(t, date, want)
	if !reflect.DeepEqual(date, want) {
		t.Errorf("incorrect date: expected: %#v; got: %#v", want, date)
	}
}

type TestDate struct {
	date        *Date
	anotherDate *Date
	name        string
	result      bool
}

func TestNotDate(t *testing.T) {
	tests := []TestDate{
		{
			name: "notYear",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2019,
				Month: 11,
				Day:   18,
			},
			result: true},
		{
			name: "notMonth",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 10,
				Day:   18,
			},
			result: true},
		{
			name: "notDay",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 11,
				Day:   17,
			},
			result: true},
		{
			name: "notNot",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			result: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.date.Not(tc.anotherDate)
			assert.Equal(t, result, tc.result)
		})
	}
}

func TestEqualDate(t *testing.T) {
	tests := []TestDate{
		{
			name: "equal",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			result: true},
		{
			name: "notEqual",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 11,
				Day:   17,
			},
			result: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.date.Equal(tc.anotherDate)
			assert.Equal(t, result, tc.result)
		})
	}
}

func TestGreaterDate(t *testing.T) {
	tests := []TestDate{
		{
			name: "year",
			date: &Date{
				Year:  2021,
				Month: 11,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			result: true},
		{
			name: "month",
			date: &Date{
				Year:  2020,
				Month: 12,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			result: true},
		{
			name: "day",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   20,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			result: true},
		{
			name: "not",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   15,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			result: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.date.Greater(tc.anotherDate)
			assert.Equal(t, result, tc.result)
		})
	}
}

func TestLessDate(t *testing.T) {
	tests := []TestDate{
		{
			name: "year",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2021,
				Month: 11,
				Day:   18,
			},
			result: true},
		{
			name: "month",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 12,
				Day:   18,
			},
			result: true},
		{
			name: "day",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 11,
				Day:   20,
			},
			result: true},
		{
			name: "not",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 11,
				Day:   15,
			},
			result: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.date.Less(tc.anotherDate)
			assert.Equal(t, result, tc.result)
		})
	}
}

func TestGreaterOrEqualDate(t *testing.T) {
	tests := []TestDate{
		{
			name: "year",
			date: &Date{
				Year:  2021,
				Month: 11,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			result: true},
		{
			name: "month",
			date: &Date{
				Year:  2020,
				Month: 12,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			result: true},
		{
			name: "day",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   20,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			result: true},
		{
			name: "not",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   15,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			result: false},
		{
			name: "equal",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			result: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.date.GreaterOrEqual(tc.anotherDate)
			assert.Equal(t, result, tc.result)
		})
	}
}

func TestLessOrEqualDate(t *testing.T) {
	tests := []TestDate{
		{
			name: "year",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2021,
				Month: 11,
				Day:   18,
			},
			result: true},
		{
			name: "month",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 12,
				Day:   18,
			},
			result: true},
		{
			name: "day",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 11,
				Day:   20,
			},
			result: true},
		{
			name: "not",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 11,
				Day:   15,
			},
			result: false},
		{
			name: "equal",
			date: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			anotherDate: &Date{
				Year:  2020,
				Month: 11,
				Day:   18,
			},
			result: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.date.LessOrEqual(tc.anotherDate)
			assert.Equal(t, result, tc.result)
		})
	}
}
