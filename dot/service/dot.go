package service

import (
	"bytes"
	"github.com/dddplayer/core/dot/entity"
	"github.com/dddplayer/core/dot/factory"
	"github.com/dddplayer/core/dot/valueobject"
	"io"
	"text/template"
)

func WriteDot(g entity.DotGraph, w io.Writer) error {
	gb := factory.NewGraphBuilder(g)
	digraph, err := gb.Build()
	if err != nil {
		return err
	}

	t := template.New("dot")
	for _, s := range []string{valueobject.TmplColumn, valueobject.TmplRow,
		valueobject.TmplNode, valueobject.TmplEdge, valueobject.TmplGraph} {
		if _, err := t.Parse(s); err != nil {
			return err
		}
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, digraph); err != nil {
		return err
	}
	_, err = buf.WriteTo(w)
	return err
}
