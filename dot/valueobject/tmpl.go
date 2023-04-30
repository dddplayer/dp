package valueobject

const TmplEdge = `{{define "edge" -}}
    {{printf "%s -> %s  [style=dotted]" .From .To}}
{{- end}}`

const TmplColumn = `{{define "column" -}}
    {{printf "<td port=%q bgcolor=%q rowspan=\"%d\" colspan=\"%d\">%s</td>" .Port .BgColor .RowSpan .ColSpan .Text}}
{{- end}}`

const TmplRow = `{{define "row" -}}
	<tr>
		{{range .Data}}
			{{template "column" .}}
		{{- end}}
	</tr>
{{- end}}`

const TmplNode = `{{define "node" -}}
    {{printf "%s [label=<" .Name}}
        <table border="0" cellpadding="10">
			{{range .Rows}}
				{{template "row" .}}
			{{- end}}
        </table>
        > ]
{{- end}}`

const TmplGraph = `digraph {
    fontname="Helvetica, Arial, sans-serif"
    node [fontname="Helvetica, Arial, sans-serif"]
    edge [fontname="Helvetica, Arial, sans-serif"]
    graph [style=invis]

    subgraph cluster_ddd_concept{
        node [shape=rect]
        edge [style=invis]

        {rank=same node0 node1 node2 node3 node4 node5 node6 node7 node8 node9 node10}

		node10 [label="Function" style=filled fillcolor="#ead1dcff"]
        node9 [label="General" style=filled fillcolor="#f4ccccff"]
        node8 [label="Class" style=filled fillcolor="#b4a7d6ff"]
        node7 [label="Factory" style=filled fillcolor="#cfe2f3ff"]
        node6 [label="Event" style=filled fillcolor="#f6b26bff"]
        node5 [label="Command" style=filled fillcolor="#a4c2f4ff"]
        node4 [label="Service" style=filled fillcolor="#e69138ff"]
        node3 [label="ValueObject" style=filled fillcolor="#a2c4c9ff"]
        node2 [label="Entity" style=filled fillcolor="#ffe599ff"]
        node1 [label="AggregateRoot" style=filled fillcolor="#ffd966ff"]
        node0 [label="BoundedContext" style=dashed fillcolor="#ffffff00"]

        node0 -> node1 -> node2 -> node3 -> node4 -> node5 -> node6 -> node7 -> node8 -> node9 -> node10

    }

    subgraph cluster_entity {
        node [style=dotted shape=rect]

		{{range .Nodes}}
			{{template "node" .}}
		{{- end}}

    }
	
	{{printf "node0 -> %s:first_blank_row [style=invis]" .FirstNodeName}}

	{{range .Edges}}
		{{template "edge" .}}
	{{- end}}

	{{printf "label = %q;" .Label}}
    fontsize=20;
}
`
