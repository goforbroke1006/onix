package prometheus

import "fmt"

var (
	ErrZeroStep     = fmt.Errorf("step should not be zero")
	ErrNegativeStep = fmt.Errorf("step should not be negative")
)
