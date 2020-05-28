package main

import (
	"os"
	hjr "github.com/ovlad32/hjr"
)

func main()  {
	rc, err := os.Open("/Dev/ge.data/ORCL.CK3.TX")
	s, err := hjr.NewScanner("ORCL.CK3.TX", rc, false)
	_ = err
	_ = s
	rc.Close()
	//fmt.Printf("Hi\r")	
}