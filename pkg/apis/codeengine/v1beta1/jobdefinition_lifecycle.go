/*******************************************************************************
 * Licensed Materials - Property of IBM
 * IBM Cloud Code Engine, 5900-AB0
 * © Copyright IBM Corp. 2020
 * US Government Users Restricted Rights - Use, duplication or
 * disclosure restricted by GSA ADP Schedule Contract with IBM Corp.
 ******************************************************************************/

package v1beta1

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// CalculateArrayIndices returns indices based on arraySpec
func (jds *JobDefinitionSpec) CalculateArrayIndices() (map[int64]interface{}, error) {
	result := make(map[int64]interface{})

	if jds.ArraySpec == nil {
		return result, NewInvalidFieldError("missing field in JobDefinitionSpec: arraySpec")
	}

	indexRanges := strings.Split(*jds.ArraySpec, ",")
	for _, indexRange := range indexRanges {
		indices := strings.Split(indexRange, "-")
		var start, end int64
		var err error
		if len(indices) > 2 {
			return nil, errors.Errorf("error parsing arraySpec range: '%s'. Expect 2, got %d", indexRange, len(indices))
		}
		startString := strings.Replace(indices[0], " ", "", -1)
		start, err = strconv.ParseInt(startString, 10, 64)
		if err != nil {
			return nil, errors.Errorf("error getting start index of range '%s': %s", indexRange, err.Error())
		}

		if len(indices) > 1 {
			endString := strings.Replace(indices[1], " ", "", -1)
			end, err = strconv.ParseInt(endString, 10, 64)
			if err != nil {
				return nil, errors.Errorf("error getting end index of range '%s': %s", indexRange, err.Error())
			}
		} else {
			end = start
		}
		if start > end {
			start, end = end, start
		}
		if start < 0 || end > MaxIndexValue {
			return nil, errors.Errorf("exceeded index range, must between 0 and %d, got: %d and %d", MaxIndexValue, start, end)
		}
		for i := start; i <= end; i++ {
			if _, exist := result[i]; !exist {
				result[i] = nil
			}
		}
		if len(result) > maxArraySize {
			return nil, errors.Errorf("exceeded maximum array size: %d", maxArraySize)
		}
	}

	return result, nil
}

// InvalidFieldError is a custom error which captures invalid fields in jobDefinitionSpec
type InvalidFieldError struct {
	message string
}

// NewInvalidFieldError returns a new InvalidFieldError
func NewInvalidFieldError(message string) *InvalidFieldError {
	return &InvalidFieldError{message: message}
}

// Error returns the error message
func (e *InvalidFieldError) Error() string {
	return e.message
}

// IsInvalidFieldError returns if the error is an InvalidFieldError
func IsInvalidFieldError(o interface{}) bool {
	err := o.(error)
	err = errors.Cause(err)
	_, ok := err.(*InvalidFieldError)
	return ok
}
