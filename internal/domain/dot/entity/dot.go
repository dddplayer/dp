package entity

import (
	"bytes"
	"github.com/dddplayer/dp/internal/domain/dot/valueobject"
	"io"
	"text/template"
)

type Dot struct {
	Name  string
	Label string

	SubGraphs []*SubGraph
	Edges     []*Edge

	Templates []string
}

type SubGraph struct {
	Name  string
	Label string
	Nodes []*Node

	SubGraphs []*SubGraph
}

type Node struct {
	ID      string
	Name    string
	BgColor string
	Table   *Table
}

type Table struct {
	Rows []*Row
}

type Row struct {
	Data []*Data
}

type Data struct {
	Text    string
	Port    string
	BgColor string
	RowSpan int
	ColSpan int
}

type Edge struct {
	From    string
	To      string
	Tooltip string
	L       string
	T       string
	A       string
}

func (d *Dot) Write(w io.Writer) error {
	t := template.New("dot")
	for _, s := range d.Templates {
		if _, err := t.Parse(s); err != nil {
			return err
		}
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, d); err != nil {
		return err
	}
	_, err := buf.WriteTo(w)
	return err
}

func blankRow() *Row {
	r := &Row{Data: []*Data{}}
	for i := 0; i < rowStartEmptyBlockNum; i++ {
		r.Data = append(r.Data, blankColumn())
	}
	return r
}

func nameRow(name string, colSpan int) *Row {
	r := &Row{Data: []*Data{}}
	d := blankColumn()
	d.Text = name
	d.ColSpan = colSpan
	r.Data = append(r.Data, d)
	return r
}

func column(name, id, bgColor string) *Data {
	d := blankColumn()
	d.Text = name
	d.Port = valueobject.PortStr(id)
	d.BgColor = bgColor

	return d
}

func blankColumn() *Data {
	return &Data{
		Text:    "",
		Port:    "",
		BgColor: "white",
		RowSpan: 1,
		ColSpan: 1,
	}
}
