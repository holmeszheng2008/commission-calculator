package main

import (
	"commission-calculator/internal/service"
	"flag"
	"fmt"
	"log"
)

func main() {
	log.Println("hello")
	arUrlPtr := flag.String("ar", "", "ar.csv location")
	glUrlPtr := flag.String("gl", "", "gl.csv location")
	piUrlPtr := flag.String("pi", "", "pi.csv localtion")
	recUrlPtr := flag.String("rec", "", "rec.csv localtion")
	opPtr := flag.String("op", "", "operation")
	flag.Parse()
	if *opPtr == "reconciliation" {
		if *arUrlPtr == "" || *glUrlPtr == "" {
			log.Fatalln("parameters ar, gl must be specified for reconciliation")
			return
		}
		err := service.LoadCsvs(arUrlPtr, glUrlPtr, nil, nil)
		if err != nil {
			return
		}
		err = service.Reconciliate()
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
	} else if *opPtr == "commission" {
		if *arUrlPtr == "" || *glUrlPtr == "" || *piUrlPtr == "" || *recUrlPtr == "" {
			log.Fatalln("parameters ar, gl, pi, rec must be specified for commission generation")
			return
		}
		err := service.LoadCsvs(arUrlPtr, glUrlPtr, piUrlPtr, recUrlPtr)
		if err != nil {
			return
		}

		err = service.GenCommisionRpt()
		if err != nil {
			log.Fatalln(err.Error())
			return
		}

	} else {
		log.Fatalln(fmt.Sprintf("Operation %s invalid", *opPtr))
		return
	}

}
