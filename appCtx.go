package main

import (
	"bufio"
	"flag"
	"fmt"
	yaccToken "github.com/bhbosman/yaccidl"
	"os"
	"strings"
)

type IdlDefinitions struct {
	definitions map[string]int
}

func (i IdlDefinitions) String() string {
	var result []string
	for k, _ := range i.definitions {
		result = append(result, k)
	}
	return strings.Join(result, ",")
}

func (i *IdlDefinitions) Set(v string) error {
	ss := strings.Split(v, ",")
	for _, s := range ss {
		i.definitions[s] = 1
	}
	return nil
}

func (i IdlDefinitions) AssignFlags() []string {
	var result []string
	for k, _ := range i.definitions {
		result = append(result, k)
	}
	return result
}

func NewIdlDefinitions() IdlDefinitions {
	return IdlDefinitions{
		definitions: make(map[string]int),
	}
}

type GoLexIdlAppCtx struct {
	suppressTokens               SuppressTokens
	oflag                        string
	verbose                      bool
	hflag, tflag, dumpTokenNames bool
	inputFileName                string
	inputReader                  *bufio.Reader
	outputWriter                 *bufio.Writer
	idlDefinitions               IdlDefinitions
}

func NewAppCtx() *GoLexIdlAppCtx {
	return &GoLexIdlAppCtx{
		suppressTokens: *NewSuppressTokens(),
		oflag:          "",
		verbose:        false,
		hflag:          false,
		tflag:          false,
		dumpTokenNames: false,
		inputFileName:  "",
		idlDefinitions: NewIdlDefinitions(),
	}
}

func (self *GoLexIdlAppCtx) run() error {
	flag.BoolVar(&self.verbose, "v", self.verbose, "verbose")
	flag.BoolVar(&self.hflag, "h", false, "show help and exit")
	flag.StringVar(&self.oflag, "o", "", "lexer output")
	flag.BoolVar(&self.tflag, "t", false, "write scanner on stdout")
	flag.BoolVar(&self.dumpTokenNames, "d", false, "dump token names")
	flag.Var(&self.suppressTokens, "remove_token", "remove token")
	flag.Parse()
	self.inputFileName = flag.Arg(0)
	if self.hflag || flag.NArg() > 1 {
		flag.Usage()
		_, _ = fmt.Fprintf(stderr, "\n%s [-o out_name] [other_options] [in_name]\n", os.Args[0])
		_, _ = fmt.Fprintln(stderr, "  If no in_name is given then read from stdin.")
		stderr.Flush()
		return fmt.Errorf("exit as help was required")
	}

	if self.dumpTokenNames {
		for i := 0; i < len(yaccToken.YaccIdlTok2)-yaccToken.YaccIdlErrCode; i++ {
			fmt.Println(yaccToken.YaccIdlTokname(yaccToken.YaccIdlTok2[i] + yaccToken.YaccIdlErrCode))
		}
		return fmt.Errorf("exit as token dump")
	}
	return nil
}
