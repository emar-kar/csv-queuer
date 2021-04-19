package request

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
)

// IndexMap is a custom type that defines several options of the Result object.
type IndexMap map[string]int

// RowData is a custom type that defines lines of results from the csv file.
type RowData map[string]string

// Results is an object containing result information.
type Results struct {
	Request      *Request
	SelectInd    IndexMap
	ConditionInd IndexMap
	MaxLength    IndexMap
	Data         []RowData
	sync.Mutex
	HasData bool
}

func (r *Results) fillConditionIndexes(headers []string) {
	conditionInd := make(IndexMap)
	fields := r.Request.Where.GetFields()

	for ind, val := range headers {
		for _, key := range fields {
			if key == val {
				conditionInd[key] = ind
			}
		}
	}

	r.ConditionInd = conditionInd
}

// ParseCSVFile contains main logic of the file parsing procedure.
func (r *Results) ParseCSVFile(ctx context.Context, csvSep string, resultDataCh chan<- RowData, doneCh chan<- error) {
	f, err := os.Open(r.Request.From)
	if err != nil {
		doneCh <- err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	firstRow := true
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			doneCh <- ctx.Err()
			return
		default:
			if firstRow {
				firstRow = false
				continue
			}

			line := strings.Split(scanner.Text(), csvSep)
			if r.checkConditions(line) {
				data := r.createData(line)
				resultDataCh <- data
			}
		}
	}
	doneCh <- nil
}

func (r *Results) checkConditions(line []string) bool {
	for key := range r.ConditionInd {
		lineData := line[r.ConditionInd[key]]

		if r.Request.Where == nil {
			return true
		}

		if !r.Request.Where.Condition(key, strings.ToLower(lineData)) {
			return false
		}
	}
	return true
}

func (r *Results) createData(line []string) RowData {
	data := make(RowData)
	for _, field := range r.Request.Select {
		r.Lock()
		ind := r.SelectInd[field]
		data[field] = line[ind]
		if r.MaxLength[field] < len(line[ind]) {
			r.MaxLength[field] = len(line[ind])
		}
		r.Unlock()
	}
	return data
}

// Print prints results in table.
func (r *Results) Print() {
	if !r.HasData {
		fmt.Println("nothing to print")
		return
	}
	var line string = "|"
	for _, key := range r.Request.Select {
		length := r.MaxLength[key] + 2
		leftSide := (length - len(key)) / 2
		rightSide := length - len(key) - leftSide
		line = fmt.Sprintf("%s%s", line, strings.Repeat(" ", leftSide))
		line += key
		line = fmt.Sprintf("%s%s|", line, strings.Repeat(" ", rightSide))
	}
	fmt.Println(strings.Repeat("=", len(line)))
	fmt.Println(line)
	fmt.Println(strings.Repeat("=", len(line)))

	for _, data := range r.Data {
		line = "|"
		for _, key := range r.Request.Select {
			length := r.MaxLength[key] + 2
			leftSide := (length - len(data[key])) / 2
			rightSide := length - len(data[key]) - leftSide
			line = fmt.Sprintf("%s%s", line, strings.Repeat(" ", leftSide))
			line += data[key]
			line = fmt.Sprintf("%s%s|", line, strings.Repeat(" ", rightSide))
		}
		fmt.Println(line)
	}
	fmt.Println(strings.Repeat("=", len(line)))
}
