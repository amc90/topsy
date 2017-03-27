package topsy

import "os"
import "io"
import "bufio"
import "regexp"
import "errors"
import "fmt"

//TODO::
var _ = fmt.Printf

type Kind int
const (
	KCons = iota
	KText
	KSymbol
)
type Datum interface {
	kind() Kind
}
type Cons struct {
	car Datum
	cdr Datum
}
type Text string
type Symbol string
func (f Cons) kind() Kind { return KCons }
func (f Text) kind() Kind { return KText }
func (f Symbol) kind() Kind { return KSymbol }

//Line oriented. Supports partial reads
type sourceFile struct {
	//pending is data which has already been read by reader, but pushed back.
	pending *string
	reader *bufio.Reader
	//lineNo and charNo refer to the start of "pending"
	lineNo int
	charNo int
}

//Reads as much of the line as matches the regular expression
//Returns nil, nil if there is input but it doesn't match
//Returns nil, io.EOF at EOF
func (src *sourceFile) ReadMatch(re *regexp.Regexp) (*string, error) {
	if(src.pending==nil) {
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
		src.pending=&s
	}
	loc:=re.FindStringIndex(*src.pending)
	if loc==nil {
		//Didn't find
		return nil, nil
	}
	ret:=(*src.pending)[loc[0]:loc[1]]
	if loc[1]<len(*src.pending) {
		*src.pending=(*src.pending)[loc[1]:]
	} else {
		src.pending=nil
	}
	src.charNo+=loc[1];
	return &ret, nil
}

var reLeadingToken=regexp.MustCompile(`^([[:space:]]+|"[^"]*"|[()]|[^()[:space:]]+)`)
var reSpace=regexp.MustCompile(`^[[:space:]]+$`)
var reString=regexp.MustCompile(`^"[^"]*"$`)
//parse a file
func Read(src *sourceFile) (Datum, error) {
	d, e:=ReadChild(src);
	switch e {
		case nil:
fmt.Println(d);
			return nil, errors.New("Unexpected termination of Read. Extraneous closing bracket?");
		case io.EOF:
			return d, e;
		default:
			return nil, e;
	}
}

func ReadChild(src *sourceFile) (Datum, error) {
	var tree Datum=nil;
	for {
		s, e:=src.ReadMatch(reLeadingToken);
		switch e {
			case nil:
			case io.EOF:
				return tree, io.EOF
			default:
				return nil, e
		}
		switch {
		case reSpace.MatchString(*s)://Discard
		case reString.MatchString(*s):
			tree=Cons{"TEXT:"+Text((*s)[1:len(*s)-1]), tree}
		case *s=="(":
			children, e:=ReadChild(src);
			switch e {
			case nil:
				tree=Cons{children,tree};
			case io.EOF:
				return nil, errors.New("Unexpected EOF in Read. Missing closing bracket?");
			default:
				return nil, e;
			}
		case *s==")":
			return tree, nil
		default:
			tree=Cons{"SYMBOL:"+Symbol(*s), tree}
		}
	}
}

func Lex(f *os.File) (Datum, error) {
	defer f.Close()
	reader:=bufio.NewReader(f)
	src:=&sourceFile{nil,reader,1,1}
	return Read(src)
}

