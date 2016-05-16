package main

import (
	"golang.org/x/net/context"
)

type Split struct {
	inputs []<-chan context.Context
}

func (s *Split) Left(stage Stage) *Split {
	return s
}

func (s *Split) Right(stage Stage) *Split {
	return s
}
