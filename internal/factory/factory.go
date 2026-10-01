package factory

import (
	"bulk-parser/internal/parser"
	"fmt"
)

type ParserFactory struct {
	parsers map[string]parser.Parser
}

func Newparsefactory() *ParserFactory {
	return &ParserFactory{
		parsers: map[string]parser.Parser{
			"go":   parser.GoParser{},
			"java": parser.JavaParser{},
		},
	}
}

func (f *ParserFactory) GetParser(language string) (parser.Parser, error) {
	p, ok := f.parsers[language]
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
	return p, nil
}
