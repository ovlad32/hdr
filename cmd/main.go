package main

import (
	"context"
	"fmt"
	"hash/fnv"
	"log"
	"os"

	hjr "github.com/ovlad32/hjr"
)


func main()  {
	in := "/Dev/ge.data/ORCL.CK3.TX"
	ctx := context.Background()
	rc, err := os.Open(in)
	if err != nil {
		err = fmt.Errorf("Opening file %v: %v", in, err)
		log.Fatal(err)
	}
	log.Printf("reading data")
	defer rc.Close()
	
	// if zip {
	// 	r, err = gzip.NewReader(r);
	// 	if err != nil {
	// 		return
	// 	}
	// }
	opts := hjr.NewIndexOptions().
		SetHashing(fnv.New32a()).
		SetSeparator([]byte{0x7}).
		SetStorage(hjr.NewMemStorage())
	
	idx := hjr.NewIndexProc(opts)

	bc := &hjr.ByteCleaver{
		Accept :idx.ParseFunc(),
	}
	s := hjr.NewScanner(rc, bc.LineProviderFunc())
	err = s.Serve(ctx) 
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("done")
	//"CK3.TX"
	//fmt.Printf("Hi\r")	
}