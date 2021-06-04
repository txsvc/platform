package validate

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/johngb/langreg"
)

const (
	// AssertionWarning indicates a potential issue
	AssertionWarning = 0
	// AssertionError indicates an error in the validation
	AssertionError = 1

	MsgStringMismatch      = "expected '%s', found '%s'"
	MsgNonEmptyString      = "expected non empty string '%s'"
	MsgNotNil              = "expected not-nil attribute, got '%s'"
	MsgNonZero             = "expected non-zero value '%s'"
	MsgInvalidLanguageCode = "invalid language code '%s'"
	MsgInvalidTimestamp    = "invalid timestamp '%d'"
	MsgNonEmptyMap         = "expected none empty map '%s'"
	MsgExpectedKey         = "expected key '%s' in map '%s'"
	MsgEmptyReport         = "validation '%s' has zero errors/warnings"
)

type (
	// Assertion is used to collect validation information
	Assertion struct {
		Type int    // 0 == warning, 1 == error
		Txt  string // description of the problem
		Err  error
	}

	// Validator collects assertions
	Validator struct {
		Name     string
		Issues   []*Assertion
		Errors   int
		Warnings int
	}

	// Validatable is the interface that maust be implemented to support recursive validations of strucs
	Validatable interface {
		Validate(*Validator) *Validator
	}
)

// New initializes and returns a new Validator
func New(name string) *Validator {
	v := Validator{
		Name:   name,
		Issues: make([]*Assertion, 0),
	}
	return &v
}

// Validate starts the chain of validations
func (v *Validator) Validate(src Validatable) *Validator {
	return src.Validate(v)
}

// AddError adds an error assertion
func (v *Validator) AddError(txt string) {
	v.Issues = append(v.Issues, &Assertion{Type: AssertionError, Txt: txt})
	v.Errors++
}

// AddWarning adds an warning assertion
func (v *Validator) AddWarning(txt string) {
	v.Issues = append(v.Issues, &Assertion{Type: AssertionWarning, Txt: txt})
	v.Errors++
}

// StringEquals verifies a string
func (v *Validator) StringEquals(src, expected string) {
	if len(src) != len(expected) {
		v.AddError(fmt.Sprintf(MsgStringMismatch, expected, src))
		return
	}

	if src != expected {
		v.AddError(fmt.Sprintf(MsgStringMismatch, expected, src))
	}
}

// StringNotEmpty verifies a string is not empty
func (v *Validator) StringNotEmpty(src, name string) {
	if len(src) == 0 {
		v.AddError(fmt.Sprintf(MsgNonEmptyString, name))
	}
}

// NotNil verifies that an attribute is not nil
func (v *Validator) NotNil(src interface{}) {
	if src == nil || (reflect.ValueOf(src).Kind() == reflect.Ptr && reflect.ValueOf(src).IsNil()) {
		v.AddError(fmt.Sprintf(MsgNotNil, reflect.TypeOf(src).String()))
	}
}

// NonZero verifies that a map is not empty
func (v *Validator) NonZero(src int, name string) {
	if src == 0 {
		v.AddError(fmt.Sprintf(MsgNonZero, name))
	}
}

// ISO639 verifies that src complies with ISO 639-1
func (v *Validator) ISO639(src string) {
	lang := src
	if !strings.Contains(src, "_") {
		lang = src + "_" + strings.ToUpper(src)
	}
	if !langreg.IsValidLangRegCode(lang) {
		v.AddError(fmt.Sprintf(MsgInvalidLanguageCode, src))
	}
}

// RFC1123Z verifies that src complies with RFC 1123-Z
func (v *Validator) RFC1123Z(src string) {
	_, err := time.Parse(time.RFC1123Z, src)
	if err != nil {
		v.AddError(err.Error())
	}
}

// Timestamp validates that src is a valid UNIX timestamp
func (v *Validator) Timestamp(src string) {
	ts, err := strconv.Atoi(src)
	if err != nil {
		v.AddError(err.Error())
	}
	if ts < 0 {
		v.AddError(fmt.Sprintf(MsgInvalidTimestamp, ts))
	}
}

// MapNotEmpty verifies that a map is not empty
func (v *Validator) MapNotEmpty(src map[string]string, name string) {
	if len(src) == 0 {
		v.AddError(fmt.Sprintf(MsgNonEmptyMap, name))
	}
}

// MapContains verifies that a map contains key
func (v *Validator) MapContains(src map[string]string, key, name string) {
	if len(src) == 0 {
		v.AddError(fmt.Sprintf(MsgNonEmptyMap, name))
		return
	}
	if _, ok := src[key]; !ok {
		v.AddError(fmt.Sprintf(MsgExpectedKey, key, name))
	}
}

// IsValid returns true if NError == 0. Warnings are ignored
func (v *Validator) IsValid() bool {
	return v.Errors == 0
}

// IsClean returns true if NError == 0 AND NWarnings == 0
func (v *Validator) IsClean() bool {
	return v.Errors == 0 && v.Warnings == 0
}

// NErrors returns the number of erros
func (v *Validator) NErrors() int {
	return v.Errors
}

// NWarnings returns the number of warnings
func (v *Validator) NWarnings() int {
	return v.Warnings
}

// AsError returns an error if NError > 0, nil otherwise
func (v *Validator) AsError() error {
	if v.Errors == 0 {
		return nil
	}
	return fmt.Errorf(v.Error())
}

// Error returns an error text
func (v *Validator) Error() string {
	return v.Report()
}

// Report returns a description of all issues
func (v *Validator) Report() string {
	if v.Errors == 0 {
		return fmt.Sprintf(MsgEmptyReport, v.Name)
	}
	r := "\n"
	for i, issue := range v.Issues {
		r = r + fmt.Sprintf("issue-%d: %s\n", i+1, issue.Txt)
	}
	return r
}
