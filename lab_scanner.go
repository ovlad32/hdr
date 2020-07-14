package hjr

import (
	"bufio"
	"context"
	"io"
)

//bufio.SplitFunc
type scanner struct {
	reader    io.Reader
	splitFunc bufio.SplitFunc
}

func NewScanner(r io.Reader, s bufio.SplitFunc) scanner {
	return scanner{reader: r, splitFunc: s}
}

func (s scanner) Serve(ctx context.Context) (err error) {
	sc := bufio.NewScanner(s.reader)
	sc.Split(s.splitFunc)
	for sc.Scan() {
		select {
		case <-ctx.Done():
			break
		default:
			if err = sc.Err(); err != nil {
				break
			}
		}
	}
	err = sc.Err()
	return err
}
