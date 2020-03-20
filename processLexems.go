package main

import (
	"encoding/json"
	"github.com/bhbosman/gocommon"
	"github.com/bhbosman/lexidl"
)

type processor interface {
	run()
}
type processLexemsImpl struct {
	ctx *GoLexIdlAppCtx
}

func (self processLexemsImpl) run() {
	encoder := json.NewEncoder(self.ctx.outputWriter)
	handler, _ := lexidl.NewLexIdlHandler(
		self.ctx.inputFileName,
		self.ctx.idlDefinitions.AssignFlags(),
		gocommon.NewByteReaderNoCloser(self.ctx.inputReader))
	for {
		lexem, _ := handler.ReadLexem()
		if _, ok := self.ctx.suppressTokens.excludedTokens[lexem.TokenName]; !ok {
			_ = encoder.Encode(lexem)
		}

		if lexem.Eof {
			_ = self.ctx.outputWriter.Flush()
			break
		}
	}
}
