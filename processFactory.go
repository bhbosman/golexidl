package main

import (
	"log"
)

type processFactory interface {
	create() processor
}

type processFactoryImpl struct {
	ctx    *GoLexIdlAppCtx
	logger *log.Logger
}

func (self processFactoryImpl) create() processor {
	return &processLexemsImpl{
		ctx: self.ctx,
	}
}
