package request

// Condition represents a wrapper on two main possible conditions.
// Only one condition option can be used at a time.
// Thus, it represents next criterion in the chain of options.
type Condition struct {
	And *Criterion
	Or  *Criterion
}

// GetExist returns Criterion of the Condition.
func (c *Condition) GetExist() *Criterion {
	if c != nil {
		if c.And != nil {
			return c.And
		}
		if c.Or != nil {
			return c.Or
		}
	}
	return nil
}

// Criterion combines all information about WHERE option.
// It contains its value and definition symbol and name of the field.
// Also it has a hook to the next criterion.
type Criterion struct {
	Value      Variable
	Conditions *Condition
	Field      string
	Symbol     string
	Strict     bool
}

// GetFields returns map, representing SELECT fields.
func (c *Criterion) GetFields() []string {
	crit := c
	var fields []string
	fields = append(fields, crit.Field)

	for {
		crit = crit.Conditions.GetExist()
		if crit == nil {
			break
		}
		if !sliceHasString(crit.Field, fields) {
			fields = append(fields, crit.Field)
		}
	}

	return fields
}

func (c *Criterion) getLastCondition() *Criterion {
	if c.Conditions == nil {
		return c
	}
	crit := c.Conditions.GetExist()
	for {
		if crit.Conditions == nil {
			return crit
		}
		crit = crit.Conditions.GetExist()
	}
}

func (c *Criterion) Condition(key, lineData string) bool {
	var result bool
	lineValue := Data(removeCharacters(lineData, " "))

	crit := c
	for {
		if crit.Field == key {
			if !crit.Strict && !result || crit.Strict {
				result = analyze(crit.Symbol, crit.Value, lineValue)
			}
			// if !crit.Strict && result {
			// 	continue
			// }
			if crit.Strict && !result {
				return result
			}
		}
		crit = crit.Conditions.GetExist()
		if crit != nil {
			continue
		}
		break
	}

	return result
}

func analyze(symbol string, data, lineData Variable) bool {
	var result bool
	switch symbol {
	case "NOT":
		result = checkNot(data, lineData)
	case "=":
		result = checkEqual(data, lineData)
	case ">":
		result = checkGreater(data, lineData)
	case "<":
		result = checkLess(data, lineData)
	case ">=":
		result = checkGreaterOrEqual(data, lineData)
	case "<=":
		result = checkLessOrEqual(data, lineData)
	}
	return result
}

func checkNot(value, lineData Variable) bool {
	switch lineData.defineType() {
	case typeInteger:
		intValue := value.toInteger()
		intLineData := lineData.toInteger()
		if intValue != intLineData {
			return true
		}
	case typeFloat:
		floatValue := value.toFloat()
		floatLineData := lineData.toFloat()
		if floatValue != floatLineData {
			return true
		}
	case typeString:
		if value != lineData {
			return true
		}
	case typeDate:
		dateValue := value.toDate()
		dateLineData := lineData.toDate()
		return dateValue.Not(dateLineData)
	}
	return false
}

func checkEqual(value, lineData Variable) bool {
	switch lineData.defineType() {
	case typeInteger:
		intValue := value.toInteger()
		intLineData := lineData.toInteger()
		if intValue == intLineData {
			return true
		}
	case typeFloat:
		floatValue := value.toFloat()
		floatLineData := lineData.toFloat()
		if floatValue == floatLineData {
			return true
		}
	case typeString:
		if value == lineData {
			return true
		}
	case typeDate:
		dateValue := value.toDate()
		dateLineData := lineData.toDate()
		return dateValue.Equal(dateLineData)
	}
	return false
}

func checkGreater(value, lineData Variable) bool {
	switch lineData.defineType() {
	case typeInteger:
		intValue := value.toInteger()
		intLineData := lineData.toInteger()
		if intValue < intLineData {
			return true
		}
	case typeFloat:
		floatValue := value.toFloat()
		floatLineData := lineData.toFloat()
		if floatValue < floatLineData {
			return true
		}
	case typeDate:
		dateValue := value.toDate()
		dateLineData := lineData.toDate()
		return dateLineData.Greater(dateValue)
	}
	return false
}

func checkLess(value, lineData Variable) bool {
	switch lineData.defineType() {
	case typeInteger:
		intValue := value.toInteger()
		intLineData := lineData.toInteger()
		if intValue > intLineData {
			return true
		}
	case typeFloat:
		floatValue := value.toFloat()
		floatLineData := lineData.toFloat()
		if floatValue > floatLineData {
			return true
		}
	case typeDate:
		dateValue := value.toDate()
		dateLineData := lineData.toDate()
		return dateLineData.Less(dateValue)
	}
	return false
}

func checkGreaterOrEqual(value, lineData Variable) bool {
	switch lineData.defineType() {
	case typeInteger:
		intValue := value.toInteger()
		intLineData := lineData.toInteger()
		if intValue <= intLineData {
			return true
		}
	case typeFloat:
		floatValue := value.toFloat()
		floatLineData := lineData.toFloat()
		if floatValue <= floatLineData {
			return true
		}
	case typeDate:
		dateValue := value.toDate()
		dateLineData := lineData.toDate()
		return dateLineData.GreaterOrEqual(dateValue)
	}
	return false
}

func checkLessOrEqual(value, lineData Variable) bool {
	switch lineData.defineType() {
	case typeInteger:
		intValue := value.toInteger()
		intLineData := lineData.toInteger()
		if intValue >= intLineData {
			return true
		}
	case typeFloat:
		floatValue := value.toFloat()
		floatLineData := lineData.toFloat()
		if floatValue >= floatLineData {
			return true
		}
	case typeDate:
		dateValue := value.toDate()
		dateLineData := lineData.toDate()
		return dateLineData.LessOrEqual(dateValue)
	}
	return false
}
