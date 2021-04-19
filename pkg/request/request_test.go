package request

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRemoveCharacters(t *testing.T) {
	output := removeCharacters("SELECT test FROM some/place WHERE something;", " ")
	assert.Equal(t, output, "SELECTtestFROMsome/placeWHEREsomething;")
}

func TestSliceHasKey(t *testing.T) {
	tests := []struct {
		name   string
		kw     string
		array  []string
		result bool
	}{
		{name: "has", array: []string{"test1", "test2", "test3"}, kw: "test2", result: true},
		{name: "not", array: []string{"test1", "test2", "test3"}, kw: "test", result: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := sliceHasString(tc.kw, tc.array)
			assert.Equal(t, result, tc.result)
		})
	}
}

func TestHasKeyword(t *testing.T) {
	tests := []struct {
		name   string
		str    string
		kw     string
		result bool
	}{
		{name: "hasAnd", str: "something AND something", kw: and, result: true},
		{name: "hasOr", str: "something OR something", kw: or, result: true},
		{name: "hasNot", str: "something NOT something", kw: not, result: true},
		{name: "hasEqual", str: "something = something", kw: equal, result: true},
		{name: "hasGreater", str: "something > something", kw: greater, result: true},
		{name: "hasLess", str: "something < something", kw: less, result: true},
		{name: "hasGreaterOrEqual", str: "something >= something", kw: greaterOrEqual, result: true},
		{name: "hasLessOrEqual", str: "something <= something", kw: lessOrEqual, result: true},
		{name: "no", str: "something something", kw: "", result: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			kw, result := hasKeyword(tc.str)
			assert.Equal(t, kw, tc.kw)
			assert.Equal(t, result, tc.result)
		})
	}
}

func TestNewRequest(t *testing.T) {
	requestString := `SELECT location, new_cases, date 
	FROM ./test/owid-covid-data.csv
	WHERE location = Ukraine OR location = Russia OR location = United Arab Emirates AND new_cases > 0;
	`
	req, err := NewRequest(requestString)
	assert.Nil(t, err)
	want := &Request{
		Select: []string{"location", "new_cases", "date"},
		From:   "./test/owid-covid-data.csv",
		Where: &Criterion{
			Field: "location",
			Conditions: &Condition{
				Or: &Criterion{
					Field: "location",
					Conditions: &Condition{
						Or: &Criterion{
							Field: "location",
							Conditions: &Condition{
								And: &Criterion{
									Field:      "new_cases",
									Conditions: nil,
									Symbol:     ">",
									Value:      Data("0"),
									Strict:     true,
								},
							},
							Symbol: "=",
							Value:  Data("unitedarabemirates"),
							Strict: false,
						},
					},
					Symbol: "=",
					Value:  Data("russia"),
					Strict: false,
				},
			},
			Symbol: "=",
			Value:  Data("ukraine"),
			Strict: false,
		},
	}
	if !reflect.DeepEqual(req, want) {
		t.Errorf("incorrect new requets: expected: %#v; got: %#v", want, req)
	}
	assert.Equal(t, req, want)

	headers, err := getHeaders("./test/owid-covid-data.csv")
	if err != nil {
		t.Errorf("cannot get headers: %s", err)
	}

	requestString = `SELECT * 
	FROM ./test/owid-covid-data.csv
	WHERE location = Ukraine OR location = Russia AND new_cases > 0;`
	req, err = NewRequest(requestString)
	assert.Nil(t, err)
	want = &Request{
		Select: headers,
		From:   "./test/owid-covid-data.csv",
		Where: &Criterion{
			Field: "location",
			Conditions: &Condition{
				Or: &Criterion{
					Field: "location",
					Conditions: &Condition{
						And: &Criterion{
							Field:      "new_cases",
							Conditions: nil,
							Symbol:     ">",
							Value:      Data("0"),
							Strict:     true,
						},
					},
					Symbol: "=",
					Value:  Data("russia"),
					Strict: false,
				},
			},
			Symbol: "=",
			Value:  Data("ukraine"),
			Strict: false,
		},
	}
	if !reflect.DeepEqual(req, want) {
		t.Errorf("incorrect new requets: expected: %#v; got: %#v", want, req)
	}
	assert.Equal(t, req, want)
}

type TestError struct {
	name      string
	reqString string
	err       string
}

func TestGetIndexes(t *testing.T) {
	tests := []TestError{
		{
			name:      "noSelect",
			reqString: "FROM somewhere WHERE something",
			err:       "cannot find SELECT in your request",
		},
		{
			name:      "noFrom",
			reqString: "SELECT something WHERE something",
			err:       "cannot find FROM in your request",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := NewRequest(tc.reqString)
			assert.Nil(t, req)
			assert.Equal(t, err.Error(), fmt.Sprintf("%s: %s", tc.err, tc.reqString))
		})
	}
}

func TestNewRequestError(t *testing.T) {
	req, err := NewRequest("SELECT something FROM somewhere WHERE something")
	assert.Nil(t, req)
	assert.Equal(t, err.Error(), "open somewhere: no such file or directory")
}

func TestGetHeadersError(t *testing.T) {
	_, err := getHeaders("somewhere")

	assert.Equal(t, err.Error(), "open somewhere: no such file or directory")
}

func TestNotInHeaders(t *testing.T) {
	headers, err := getHeaders("./test/owid-covid-data.csv")
	if err != nil {
		t.Errorf("cannot get headers: %s", err)
	}

	tests := []TestError{
		{
			name:      "noSelect",
			reqString: "SELECT something FROM ./test/owid-covid-data.csv WHERE location = russia;",
			err:       fmt.Sprintf("cannot find selected option: something in headers: %v", headers),
		},
		{
			name:      "noWhere",
			reqString: "SELECT location FROM ./test/owid-covid-data.csv WHERE something = russia;",
			err:       fmt.Sprintf("cannot find where condition: something in headers: %v", headers),
		},
		{
			name:      "noKeywordInWhere",
			reqString: "SELECT location FROM ./test/owid-covid-data.csv WHERE location IS russia;",
			err:       "cannot find any keyword in WHERE statement: locationISrussia",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewRequest(tc.reqString)
			assert.Equal(t, err.Error(), tc.err)
		})
	}
}

func TestRequestDo(t *testing.T) {
	requestString := `
	SELECT location, new_cases, date 
	FROM ./test/owid-covid-data.csv
	WHERE location = Russia OR location = Ukraine AND date >= 2020-04-20 AND date <= 2020-04-30 
	AND date NOT 2020-04-23 AND new_cases > 500 AND new_cases < 5500;
	`
	req, err := NewRequest(requestString)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result, err := req.Do(ctx, ",")
	if errors.Is(err, context.DeadlineExceeded) {
		t.Error("deadline exceeded, try to increase the processing time in config file or specify the request")
	} else if err != nil {
		t.Errorf("error: %s", err.Error())
	}

	want := &Results{
		Request:      req,
		SelectInd:    IndexMap{"date": 3, "location": 2, "new_cases": 5},
		ConditionInd: IndexMap{"date": 3, "location": 2, "new_cases": 5},
		MaxLength:    IndexMap{"date": 10, "location": 8, "new_cases": 9},
		Data: []RowData{
			{"date": "2020-04-20", "location": "Russia", "new_cases": "4268.0"},
			{"date": "2020-04-22", "location": "Russia", "new_cases": "5236.0"},
			{"date": "2020-04-30", "location": "Ukraine", "new_cases": "540.0"},
		},
		HasData: true,
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("incorrect result: expected: %#v; got: %#v", want, result)
	}
	assert.Equal(t, result, want)
}

func TestDoCtxTimeout(t *testing.T) {
	requestString := `
	SELECT location 
	FROM ./test/owid-covid-data.csv
	WHERE new_cases > 0;
	`
	req, err := NewRequest(requestString)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond)
	defer cancel()

	_, err = req.Do(ctx, ",")
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}
