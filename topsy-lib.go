package topsy

import "os"
import "io"
import "bufio"
import "regexp"

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

//Line oriented. Supports partial reads
type sourceFile struct {
	//pending is data which has already been read by reader, but pushed back.
	pending string
	reader *bufio.Reader
	//lineNo and charNo refer to the start of "pending"
	lineNo int
	charNo int
}

//Reads as much of the line as matches the regular expression
//Returns nil, nil if there is input but it doesn't match
//Returns nil, io.EOF at EOF
func (src *sourceFile) ReadMatch(re *regexp.Regexp) (*string, error) {
	if(src.pending=="") {
		s, e := src.reader.ReadString('\n')
		if e!=nil {
			if e==io.EOF {
				return nil, e
			}
			src.reader=nil
			return nil, e
		}
		src.lineNo++
		src.charNo=1
		src.pending=s
	}
	loc:=re.FindStringIndex(src.pending)
	if loc==nil {
		//Didn't find
		return nil, nil
	}
	ret:=src.pending[loc[0]:loc[1]]
	src.pending=src.pending[loc[1]+1:]
	src.charNo+=loc[1];
	return &ret, nil
}

//Returns a datum, and the remaining unmatched input
func Read(src *sourceFile) (Datum, error) {
	var tree Datum=nil;
	for {
		var reAll=regexp.MustCompile(`^.*`)
		s, e:=src.ReadMatch(reAll);
		if e != nil {
			if e == io.EOF {
				return tree, nil
			}
			return nil, e
		}
		tree=Cons{Text(*s), tree}
	}
}

func Lex(f *os.File) (Datum, error) {
	defer f.Close()
	reader:=bufio.NewReader(f)
	src:=&sourceFile{"",reader,1,1}
	return Read(src)
}

