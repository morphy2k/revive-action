package main

import (
	"fmt"
	"go/token"
	"strings"
)

type failure struct {
	Position struct {
		Start token.Position
		End   token.Position
	}
	Failure  string
	Severity string
}

func (f *failure) Format() string {
	var sb strings.Builder

	if f.Severity == "warning" {
		sb.WriteString("::warning ")
	} else {
		sb.WriteString("::error ")
	}

	fmt.Fprintf(&sb, "file=%s,line=%d,endLine=%d,col=%d,endColumn=%d::%s",
		f.Position.Start.Filename, f.Position.Start.Line,
		f.Position.End.Line, f.Position.Start.Column,
		f.Position.End.Column, f.Failure)

	return sb.String()
}

type statistics struct {
	Total, Warnings, Errors int
}

func (s statistics) String() string {
	return fmt.Sprintf("%d failures (%d warnings, %d errors)",
		s.Total, s.Warnings, s.Errors)
}
