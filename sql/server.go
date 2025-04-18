package sql

import (
	"strings"

	ksql "github.com/kovey/db-go/v3"
)

type serverOption struct {
	key ksql.ServOptKey
	val string
}

type Server struct {
	*base
	server      string
	wrapperName string
	options     []*serverOption
	isCreate    bool
}

func NewServer() *Server {
	s := &Server{base: newBase(), isCreate: true}
	s.opChain.Append(s._keyword, s._name, s._wrapper, s._options)
	return s
}

func (s *Server) _keyword(builder *strings.Builder) {
	if s.isCreate {
		keywordCreate(builder)
		return
	}

	keywordAlter(builder)
}

func (s *Server) _name(builder *strings.Builder) {
	builder.WriteString(" SERVER ")
	Backtick(s.server, builder)
}

func (s *Server) _wrapper(builder *strings.Builder) {
	if !s.isCreate {
		return
	}

	builder.WriteString(" FOREIGN DATA WRAPPER ")
	builder.WriteString(s.wrapperName)
}

func (s *Server) _options(builder *strings.Builder) {
	builder.WriteString(" OPTIONS (")
	for index, option := range s.options {
		if index > 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(string(option.key))
		builder.WriteString(" ")
		Quote(option.val, builder)
	}
	builder.WriteString(")")
}

func (s *Server) Server(name string) ksql.ServerInterface {
	s.server = name
	return s
}

func (s *Server) WrapperName(wrapperName string) ksql.ServerInterface {
	if !s.isCreate {
		return s
	}

	s.wrapperName = wrapperName
	return s
}

func (s *Server) Option(key ksql.ServOptKey, val string) ksql.ServerInterface {
	s.options = append(s.options, &serverOption{key: key, val: val})
	return s
}

func (s *Server) Alter() ksql.ServerInterface {
	s.isCreate = false
	return s
}
