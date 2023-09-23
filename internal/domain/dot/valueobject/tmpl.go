package valueobject

const TmplEdge = `{{define "edge" -}}
    {{printf "%s -> %s  [style=%s arrowhead=%s label=%q tooltip=%q]" .From .To .T .A .L .Tooltip}}
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
    {{printf "%s [label=<" .Scope}}
        <table border="0" cellpadding="10">
			{{range .Table.Rows}}
				{{template "row" .}}
			{{- end}}
        </table>
        > ]
{{- end}}`

const TmplSubGraph = `{{define "subgraph" -}}
{{printf "subgraph cluster_%s {" .Scope}}
	node [style=dotted shape=rect]

	{{range .Nodes}}
		{{template "node" .}}
	{{- end}}

    }
{{- end}}`

const TmplGraph = `digraph {
    node [style=dotted shape=rect]

    subgraph cluster_ddd_concept{
        ddd_concept [label=<
        <table border="0" cellpadding="10">
		<tr>
			<td bgcolor="#ffffff00" rowspan="1" colspan="1">BoundedContext</td>
			<td bgcolor="#ffd966ff" rowspan="1" colspan="1">AggregateRoot</td>
			<td bgcolor="#ffe599ff" rowspan="1" colspan="1">Entity</td>
			<td bgcolor="#a2c4c9ff" rowspan="1" colspan="1">ValueObject</td>
			<td bgcolor="#e69138ff" rowspan="1" colspan="1">Service</td>
		</tr>
		<tr>
			<td bgcolor="white" rowspan="1" colspan="1"></td>
			<td bgcolor="#a4c2f4ff" rowspan="1" colspan="1">Command</td>
			<td bgcolor="#f6b26bff" rowspan="1" colspan="1">Event</td>
			<td bgcolor="#cfe2f3ff" rowspan="1" colspan="1">Factory</td>
			<td bgcolor="#b4a7d6ff" rowspan="1" colspan="1">Class</td>
			
		</tr>
		<tr>
			<td bgcolor="white" rowspan="1" colspan="1"></td>
			<td bgcolor="#f4ccccff" rowspan="1" colspan="1">General</td>
			<td bgcolor="#ead1dcff" rowspan="1" colspan="1">Function</td>
			<td bgcolor="#9fc5e8ff" rowspan="1" colspan="1">Interface</td>
			<td bgcolor="#f3f3f3ff" rowspan="1" colspan="1">Attribute</td>
		</tr>
        </table>
        > ]
    }

    {{range .SubGraphs}}
		{{template "subgraph" .}}
	{{- end}}

	{{range .Edges}}
		{{template "edge" .}}
	{{- end}}

	{{printf "label = %q;" .Label}}
    fontsize=20;
}
`

const TmplSimpleNode = `{{define "simple_node" -}}
	{{if .Table}}
    	{{printf "%s [label=<" .ID}}
        <table border="0" cellpadding="10">
			{{range .Table.Rows}}
				{{template "row" .}}
			{{- end}}
        </table>
        > ]
	{{else}}
		{{printf "%s [label=%q style=filled fillcolor=%q]" .ID .Name .BgColor}}
	{{end}}
{{- end}}`

const TmplSimpleSubGraph = `{{define "simple_subgraph" -}}
{{printf "subgraph cluster_%s {" .Name}}
	{{range .Nodes}}
		{{template "simple_node" .}}
	{{- end}}

	{{printf "label = %q" .Label}}

	{{range .SubGraphs}}
		{{template "simple_subgraph" .}}
	{{- end}}
    }
{{- end}}`

const TmplSimpleGraph = `digraph {
	node [style=dotted shape=rect]

    subgraph cluster_ddd_concept{
		node [color=white]

        ddd_concept [label=<
        <table border="0" cellpadding="10">
		<tr>
			<td bgcolor="#ffffff00" rowspan="1" colspan="1">BoundedContext</td>
			<td bgcolor="#ffd966ff" rowspan="1" colspan="1">AggregateRoot</td>
			<td bgcolor="#ffe599ff" rowspan="1" colspan="1">Entity</td>
			<td bgcolor="#a2c4c9ff" rowspan="1" colspan="1">ValueObject</td>
			<td bgcolor="#e69138ff" rowspan="1" colspan="1">Service</td>
		</tr>
		<tr>
			<td bgcolor="white" rowspan="1" colspan="1"></td>
			<td bgcolor="#a4c2f4ff" rowspan="1" colspan="1">Command</td>
			<td bgcolor="#f6b26bff" rowspan="1" colspan="1">Event</td>
			<td bgcolor="#cfe2f3ff" rowspan="1" colspan="1">Factory</td>
			<td bgcolor="#b4a7d6ff" rowspan="1" colspan="1">Class</td>
			
		</tr>
		<tr>
			<td bgcolor="white" rowspan="1" colspan="1"></td>
			<td bgcolor="#f4ccccff" rowspan="1" colspan="1">General</td>
			<td bgcolor="#ead1dcff" rowspan="1" colspan="1">Function</td>
			<td bgcolor="#9fc5e8ff" rowspan="1" colspan="1">Interface</td>
			<td bgcolor="#f3f3f3ff" rowspan="1" colspan="1">Attribute</td>
		</tr>
        </table>
        > ]
	}

    {{range .SubGraphs}}
		{{template "simple_subgraph" .}}
	{{- end}}

	{{range .Edges}}
		{{template "edge" .}}
	{{- end}}

	{{printf "label = %q;" .Label}}
    fontsize=20;
}
`
