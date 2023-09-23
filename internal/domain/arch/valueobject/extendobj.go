package valueobject

type StringObj struct {
	*obj
}

func NewStringObj(v string) *StringObj {
	return &StringObj{
		obj: &obj{
			id: &ident{
				name: v,
				pkg:  "",
			},
			pos: emptyPosition(),
		},
	}
}
