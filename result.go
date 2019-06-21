package dice

import (
	"context"

	"github.com/pkg/errors"
)

// A Result is a value a die has rolled. By default, CritSuccess and CritFailure
// should be set to true if the maximum or minimum value of a die is rolled
// respectively, but the range in which a critical success/failure must be
// overridable through modifiers.
type Result struct {
	Value       float64 `json:"value"`
	Dropped     bool    `json:"dropped,omitempty"`
	CritSuccess bool    `json:"crit,omitempty"`
	CritFailure bool    `json:"fumble,omitempty"`
}

// NewResult returns a new un-dropped, non-critical Result.
func NewResult(result float64) *Result {
	return &Result{
		Value: result,
	}
}

// Drop marks a Result as dropped.
func (r *Result) Drop(ctx context.Context, drop bool) {
	r.Dropped = drop
}

// Total returns the Result's value or 0 if the result was dropped.
func (r *Result) Total(ctx context.Context) (float64, error) {
	if r == nil {
		return 0.0, errors.New("nil result")
	}
	if r.Dropped {
		return 0.0, nil
	}
	return r.Value, nil
}
