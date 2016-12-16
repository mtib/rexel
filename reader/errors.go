package reader

import "fmt"

type RexelError string

func (r RexelError) Error() string {
	return fmt.Sprintf("%s", string(r))
}
