package request

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

var TestCriterion = &Criterion{
	Field: "level0",
	Conditions: &Condition{
		And: &Criterion{
			Field: "level1",
			Conditions: &Condition{
				Or: &Criterion{
					Field: "level2",
					Conditions: &Condition{
						Or: &Criterion{
							Field: "level0",
						},
					},
				},
			},
		},
	},
}

func TestConditionGetExist(t *testing.T) {
	tests := []struct {
		expect    *Criterion
		condition *Condition
		name      string
	}{
		{name: "getAnd", condition: &Condition{And: TestCriterion}, expect: TestCriterion},
		{name: "getOr", condition: &Condition{Or: TestCriterion}, expect: TestCriterion},
		{name: "getNil", condition: &Condition{}, expect: nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cond := tc.condition.GetExist()
			assert.Equal(t, cond, tc.expect)
		})
	}
}

func TestCriterionGetFields(t *testing.T) {
	fields := TestCriterion.GetFields()
	assert.Equal(t, fields, []string{"level0", "level1", "level2"})
}
