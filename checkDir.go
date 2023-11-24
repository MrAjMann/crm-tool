package main

import (
	"fmt"
	"log"
	"os"
)

func checkDir() {
	files, err := os.ReadDir("/web/templates")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}

}
