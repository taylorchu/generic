package main

import (
	"io/ioutil"
	"log"

	"github.com/taylorchu/generic/rewrite"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	b, err := ioutil.ReadFile("GoRewrite.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	var c rewrite.Config
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		log.Fatalln(err)
	}

	err = c.RewritePackage()
	if err != nil {
		log.Fatalln(err)
	}
}
