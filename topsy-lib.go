package topsy

import "os"
import "io"
import "bufio"

type Kind int
const (
	KCons = iota
	KText
)
type Datum interface {
	kind() Kind
}
type Cons struct {
	car Datum
	cdr Datum
}
type Text string
func (f Cons) kind() Kind { return KCons }
func (f Text) kind() Kind { return KText }

func Average(xs []float64) float64 {
  total := float64(0)
  for _, x := range xs {
    total += x
  }
  return total / float64(len(xs))
}

//Reads a token, and the remaining unmatched input
func ReadToken(input *string) (Datum, string, error) {
	return Text(*input), "", nil
}

func Lex(f *os.File) (Datum, error) {
	var tree Datum=nil;
	reader:=bufio.NewReader(f)
	unconsumed:=""
	for {
		s, e := reader.ReadString('\n')
		s=unconsumed+s
		var token Datum
		var tokenErr error
		token, unconsumed, tokenErr=ReadToken(&s)
		if tokenErr!=nil { return nil, tokenErr }
		tree=Cons{token, tree}
		if e == io.EOF { break }
	}
	return tree, nil
}

