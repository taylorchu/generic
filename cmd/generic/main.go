package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/taylorchu/generic"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) < 3 {
		log.Fatalln("generic [SRCPATH] [DEST] [TypeXXX->OtherType]...")
	}

	if os.Args[1] == "" {
		log.Fatalln("SRCPATH cannot be empty")
	}

	if os.Args[2] == "" {
		log.Fatalln("DEST cannot be empty")
	}

	typeMap, err := generic.ParseTypeMap(os.Args[3:])
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := os.Stat(fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), os.Args[1])); err != nil {
		err := exec.Command("go", "get", "-u", os.Args[1]).Run()
		if err != nil {
			log.Fatalln(err)
		}
	}

	err = generic.RewritePackage(os.Args[1], os.Args[2], typeMap)
	if err != nil {
		log.Fatalln(err)
	}
}
