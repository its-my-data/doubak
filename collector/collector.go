package collector

import (
	"flag"
	"fmt"
	p "github.com/its-my-data/doubak/proto"
)

// Collect starts the major collection process.
func Collect() {
	user := flag.Lookup(p.Flag_categories.String()).Value.(flag.Getter).Get().(string)
	fmt.Println(user)
}
