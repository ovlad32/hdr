package hjr

import (
	"bufio"
	"fmt"
	"log"
)

type AcceptFunc = func(seq int, data []byte) error

type ByteCleaver struct {
	splitFunc bufio.SplitFunc 
	Accept AcceptFunc
	sep byte
}


func (rl ByteCleaver) LineProviderFunc() bufio.SplitFunc {
	var lineCount int
	log.Printf("Split is gotten")
	return func(data []byte, eof bool) (adv int, token []byte, err error) {
		//adv, token, err = bufio.ScanLines(data, eof)
		adv, token, err = rl.splitFunc(data, eof)

		if err != nil {
			err = fmt.Errorf("scanning line %v: %v", lineCount+1, err)
			return
		}
		lineCount++
		if len(token) > 0 {
			err = rl.Accept(lineCount, token)
			//	ctkn := make([]byte,len(tkn))
			//	copy(ctkn,tkn)
			//	err = s.a.accept(ctx,  s.table, line, ctkn)
			if err != nil {
				err = fmt.Errorf("Accepting line %v: %v", lineCount, err)
				return
			}
		}
		return
	}
}
func(rl ByteCleaver) onSeparator(raw []byte, atEOF bool) (advance int, token []byte, err error){
		for i := 0; i < len(raw); i++ {
			if raw[i] == rl.sep {
				return i + 1, raw[:i], nil
			}
		}
		if !atEOF {
			return 0, nil, nil
		}
		// There is one final token to be delivered, which may be the empty string.
		// Returning bufio.ErrFinalToken here tells Scan there are no more tokens after this
		// but does not trigger an error to be returned from Scan itself.
		return 0, raw, bufio.ErrFinalToken
}

