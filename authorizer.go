// ********************************************************************
// * authorizer.go                                                    *
// *                                                                  *
// * 2020-03-13 First Version, JR                                     *
// * 2020-03-16 Adds Print output, JR                                 *
// *                                                                  *
// * Go application able to read stdin line by line and retrieve,     *
// * the messages associated to operations read .                     *
// *                                                                  *
// * Usage:                                                           *
// * $ authorizer < $FILE                                             *
// ********************************************************************

package main

import (
	"authorizer/executer"
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Init our operation's execter.
	e := executer.Init()
	// Read stdin line by line, while not empty line.
	stdin := bufio.NewReader(os.Stdin)
	for {
		op, _ := stdin.ReadString('\n')
		if op == "" {
			return
		}
		// Sent input line, and get a json output string.
		out := e.Exec(op)
		fmt.Println(out)
	}
}
