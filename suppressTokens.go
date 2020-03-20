package main

import "strings"

type SuppressTokens struct {
	excludedTokens map[string]int
}

func NewSuppressTokens() *SuppressTokens {
	return &SuppressTokens{
		excludedTokens: make(map[string]int),
	}
}

func (self SuppressTokens) String() string {
	var s []string
	for k, _ := range self.excludedTokens {
		s = append(s, k)
	}
	return strings.Join(s, ",")
}

func (self *SuppressTokens) Set(value string) error {
	list := strings.Split(value, ",")
	for _, s := range list {
		self.excludedTokens[s] = 1
	}
	return nil
}
