package command

import (
	"fmt"

	"github.com/travis-g/dice/math"
	"github.com/urfave/cli"
)

// EvalCommand will evaluate the first argument it is provided as a
// math.DiceExpression and print the result or return any errors during
// evaluation.
func EvalCommand(c *cli.Context) error {
	eval := c.Args().Get(0)
	exp, err := math.Evaluate(eval)
	if err != nil {
		return err
	}
	fmt.Println(exp)
	return nil
}
