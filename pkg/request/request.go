package request

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	and            string = "AND"
	or             string = "OR"
	not            string = "NOT"
	equal          string = "="
	greater        string = ">"
	less           string = "<"
	greaterOrEqual string = ">="
	lessOrEqual    string = "<="
)

type symbols []string

func (s symbols) List() []string {
	return []string{and, or, not, greaterOrEqual, lessOrEqual, equal, greater, less}
}

// Request is a struct which defines main parameters of the request:
// select, from and where.
type Request struct {
	Where  *Criterion
	From   string
	Select []string
}

// NewRequest parses the given string and returns Request object.
func NewRequest(str string) (*Request, error) {
	r := &Request{}
	preparedStr := removeCharacters(str, " \n\t;")

	fromIndex, whereIndex, err := getIndexes(preparedStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", err, str)
	}

	fromFile := fmt.Sprint(preparedStr[fromIndex+4 : whereIndex])
	r.From = fromFile
	headers, err := getHeaders(fromFile)
	if err != nil {
		return nil, err
	}

	reqSelect := fmt.Sprint(preparedStr[6:fromIndex])
	var selectItems []string
	switch reqSelect {
	case "*":
		selectItems = headers
	default:
		selectItems = strings.Split(reqSelect, ",")
		for _, key := range selectItems {
			if !sliceHasString(key, headers) {
				return nil, fmt.Errorf("cannot find selected option: %s in headers: %v", key, headers)
			}
		}
	}
	r.Select = selectItems

	if whereIndex != len(preparedStr) {
		criterions, err := parseWhere(fmt.Sprint(preparedStr[whereIndex+5:]))
		if err != nil {
			return nil, err
		}
		fields := criterions.GetFields()
		if err := checkWhere(headers, fields); err != nil {
			return nil, err
		}
		r.Where = criterions
	}

	return r, nil
}

func getIndexes(preparedStr string) (int, int, error) {
	if !strings.Contains(preparedStr, "SELECT") {
		return -1, -1, errors.New("cannot find SELECT in your request")
	}

	fromIndex := strings.Index(preparedStr, "FROM")
	if fromIndex == -1 {
		return -1, -1, errors.New("cannot find FROM in your request")
	}

	if !strings.Contains(preparedStr, "WHERE") {
		return fromIndex, len(preparedStr), nil
	}

	whereIndex := strings.Index(preparedStr, "WHERE")
	if whereIndex == -1 {
		return -1, -1, errors.New("cannot find WHERE in your request")
	}

	return fromIndex, whereIndex, nil
}

func checkWhere(headers, fields []string) error {
	var found bool
	for _, key := range fields {
		found = false
		for _, el := range headers {
			if key == el {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("cannot find where condition: %s in headers: %v", key, headers)
		}
	}

	return nil
}

// Do starts the request to a csv file with the request object.
func (r *Request) Do(ctx context.Context, csvSep string) (*Results, error) {
	reqResult := &Results{Request: r}

	resultDataCh := make(chan RowData, 1)
	defer close(resultDataCh)

	doneCh := make(chan error, 1)
	defer close(doneCh)

	headers, err := getHeaders(r.From)
	if err != nil {
		return nil, err
	}

	fieldsInd := make(IndexMap)
	maxLength := make(IndexMap)
	for ind, val := range headers {
		for _, field := range r.Select {
			if val == field {
				fieldsInd[field] = ind
				maxLength[field] = len(field)
			}
		}
	}

	reqResult.Lock()
	reqResult.SelectInd = fieldsInd
	reqResult.MaxLength = maxLength

	reqResult.fillConditionIndexes(headers)
	reqResult.Unlock()
	reqResult.HasData = true

	go reqResult.ParseCSVFile(ctx, csvSep, resultDataCh, doneCh)

	for {
		select {
		case resultData := <-resultDataCh:
			reqResult.Data = append(reqResult.Data, resultData)
		case err := <-doneCh:
			if err != nil {
				return reqResult, err
			}
			return reqResult, nil
		}
	}
}

func removeCharacters(input string, characters string) string {
	filter := func(r rune) rune {
		if !strings.ContainsRune(characters, r) {
			return r
		}
		return -1
	}
	return strings.Map(filter, input)
}

func getHeaders(csvFile string) ([]string, error) {
	f, err := os.Open(csvFile)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	r := csv.NewReader(f)
	row, err := r.Read()
	if err != nil {
		return nil, err
	}

	return row, nil
}

func sliceHasString(key string, arr []string) bool {
	for _, el := range arr {
		if key == el {
			return true
		}
	}
	return false
}

func parseWhere(where string) (*Criterion, error) {
	if kw, ok := hasKeyword(where); ok {
		return define(where, kw), nil
	}

	return nil, fmt.Errorf("cannot find any keyword in WHERE statement: %s", where)
}

func define(where string, sep string) *Criterion {
	sepIndex := strings.Index(where, sep)

	mainWhere := where[:sepIndex]
	restWhere := where[sepIndex+len(sep):]

	if mainSep, ok := hasKeyword(mainWhere); ok {
		mainCriterion := define(mainWhere, mainSep)

		if restSep, ok := hasKeyword(restWhere); ok {
			crit := define(restWhere, restSep)

			cond := &Condition{}
			switch sep {
			case and:
				crit.Strict = true
				cond.And = crit
			case or:
				cond.Or = crit
			}

			criterionWoCondition := mainCriterion.getLastCondition()
			criterionWoCondition.Conditions = cond
		}
		return mainCriterion
	}

	return &Criterion{Field: mainWhere, Symbol: sep, Value: Data(strings.ToLower(restWhere)), Strict: false}
}

func hasKeyword(str string) (string, bool) {
	var s symbols
	for _, kw := range s.List() {
		if strings.Contains(str, kw) {
			return kw, true
		}
	}
	return "", false
}
