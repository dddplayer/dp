package entity

import "strings"

type Digraph struct {
	Name  string
	Label string

	Nodes []*Node
	Edges []*Edge
}

func (d *Digraph) FirstNodeName() string {
	if len(d.Nodes) > 0 {
		return d.Nodes[0].Name
	}
	return ""
}

type Node struct {
	Name string
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
	From string
	To   string
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
	d.Port = portStr(id)
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

func portStr(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, ".", DotJoiner), "-", DotJoiner)
}
