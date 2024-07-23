package main

import (
	"errors"
	"fmt"
	"os"

	osm_extractor_workflow "github.com/bikehopper/osm-extractor/workflow"
)

func main() {
	var argsWithoutProg []string
	if len(os.Args) != 2 {
		panic(errors.New("only accepts one arguemnt"))
	} else {
		argsWithoutProg = os.Args[1:]
	}
	switch argsWithoutProg[0] {
	case "worker":
		osm_extractor_workflow.Worker()
	default:
		fmt.Printf("Must pass 'worker' \n")
	}
}
