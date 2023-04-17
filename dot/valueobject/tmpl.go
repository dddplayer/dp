package valueobject

const TmplNode = `{{define "node" -}}
    {{printf "%s [label=<" .Name}}
        <table border="0" cellpadding="10">
			<tr>
				<td>Hello DDD Player</td>
			</tr>
        </table>
        > ]
{{- end}}`

const TmplGraph = `digraph {
    fontname="Helvetica, Arial, sans-serif"
    node [fontname="Helvetica, Arial, sans-serif"]
    edge [fontname="Helvetica, Arial, sans-serif"]
    graph [style=invis]

    subgraph cluster_entity {
        node [style=dotted shape=rect]

		{{range .Nodes}}
			{{template "node" .}}
		{{- end}}

    }

	{{printf "label = %q;" .Name}}
    fontsize=20;
}
`
