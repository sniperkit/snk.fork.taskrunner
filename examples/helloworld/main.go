/*
Sniperkit-Bot
- Status: analyzed
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"runtime"

	"github.com/sniperkit/snk.fork.taskrunner"
)

// HelloWorldTask wraps a number that prefixes the hello world message
// generated by the task.
type HelloWorldTask struct {
	num int
}

// Task implements the taskrunner.Task interface.
// Prints Hello World! to stdout prefixed with the Task number.
func (p *HelloWorldTask) Task(context.Context) (interface{}, error) {
	log.Printf("%d: Hello World!\n", p.num)
	return nil, nil
}

func main() {
	num := flag.Int("num", 100, "The number of hello worlds to print.")
	workers := flag.Int("workers", runtime.NumCPU()+1, "The number of concurrent workers to start.")

	flag.Parse()

	if *num < 0 {
		log.Fatal("num must be a positive integer")
	}

	// Create the runner.
	runner, err := taskrunner.NewTaskRunner(taskrunner.OptionMaxGoroutines(*workers))
	if err != nil {
		panic(err.Error())
	}

	// Start the runner.
	if err := runner.Start(); err != nil {
		panic(err.Error())
	}

	// Collect the promises after running the HelloWorld tasks.
	promises := make([]taskrunner.Promise, *num)
	for i := 0; i < *num; i++ {
		promises[i] = runner.Run(context.TODO(), &HelloWorldTask{i})
	}

	// Handle the promises.
	for i := range promises {
		res, err := promises[i]()
		if err != nil {
			fmt.Printf("promise failure: res=%+v - err=%+v", res, err)
		}
	}
}
