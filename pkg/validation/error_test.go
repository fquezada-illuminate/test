package validation

import (
	"fmt"
	"strings"
	"testing"
)

func TestConsolidateValidationErrors(t *testing.T) {
	t.Run("With a list of errors", func(t *testing.T) {

		errs := make([]Error, 0)

		errs = append(
			errs,
			MockValidationError{
				field: "field1",
				rule:  "rule 1",
				value: "bad",
			},
			MockValidationError{
				field: "field2",
				rule:  "rule 2",
				value: "very-bad",
			})

		testError := consolidateValidationErrors(errs)

		splitErrors := strings.Split(testError.Error(), ERROR_DELIMITER)

		if len(splitErrors) != 2 {
			t.Error(fmt.Sprintf("List of errors was expected to be exactly 2 got %d.", len(splitErrors)))
		}

		if !strings.Contains(splitErrors[0], errs[0].Field()) {
			t.Error(fmt.Sprintf("Error message should contain the invalid field."))
		}

		if !strings.Contains(splitErrors[0], errs[0].Rule()) {
			t.Error(fmt.Sprintf("Error message should contain the invalid rule."))
		}
	})
}
