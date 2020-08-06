package args

import (
	"strings"
  "fmt"
)

type Strings []string

func (i *Strings) String() string {
	return fmt.Sprint(*i)
}

func (i *Strings) Set(value string) error {
	*i = append(*i, strings.TrimSpace(value))
	return nil
}
