package application

type options struct {
	all         bool
	composition bool
}

func (o *options) ShowAllRelations() bool {
	return o.all
}
func (o *options) ShowStructEmbeddedRelations() bool {
	return o.composition
}
