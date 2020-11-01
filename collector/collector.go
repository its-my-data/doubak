package collector

import (
	"flag"
	"fmt"
)

// Collect starts the major collection process.
func Collect() {
	fmt.Println(flag.Lookup("tasks").Value.(flag.Getter).Get().(string))
}
