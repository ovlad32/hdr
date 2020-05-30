package main

import (
	"context"
	"os"

	hjr "github.com/ovlad32/hjr"
)

func main()  {
	ctx := context.Background()
	rc, err := os.Open("/Dev/ge.data/ORCL.CK3.TX")
	s, err := hjr.NewScanner("ORCL.CK3.TX", rc, false)
	_ = err
	_ = s
	s.Read(ctx)
	rc.Close()
	//fmt.Printf("Hi\r")	
}