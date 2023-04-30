package entity

type DomainModel struct {
	Name string
}

func (dm *DomainModel) NameHandler(name string) {
	dm.Name = name
}

func (dm *DomainModel) Output() (string, error) {
	return "", nil
}
