package arch

type DesignPattern string

const (
	DesignPatternPlain   DesignPattern = "plain"
	DesignPatternHexagon DesignPattern = "hexagon"
)

type HexagonDirectory string

const (
	HexagonDirectoryCmd            HexagonDirectory = "cmd"
	HexagonDirectoryPkg            HexagonDirectory = "pkg"
	HexagonDirectoryInternal       HexagonDirectory = "internal"
	HexagonDirectoryDomain         HexagonDirectory = "domain"
	HexagonDirectoryAggregate      HexagonDirectory = "aggregate"
	HexagonDirectoryApplication    HexagonDirectory = "application"
	HexagonDirectoryInfrastructure HexagonDirectory = "infrastructure"
	HexagonDirectoryInterfaces     HexagonDirectory = "interfaces"
	HexagonDirectoryRepository     HexagonDirectory = "repository"
	HexagonDirectoryFactory        HexagonDirectory = "factory"
	HexagonDirectoryEntity         HexagonDirectory = "entity"
	HexagonDirectoryValueObject    HexagonDirectory = "valueobject"
	HexagonDirectoryInvalid        HexagonDirectory = "invalid"
)

type ObjectWalker func(Object) error

type Object interface {
	Identifier() ObjIdentifier
	Position() Position
}

type ObjIdentifier interface {
	Identifier
	Name() string
	NameSeparatorLength() int
	Dir() string
}

type Identifier interface {
	ID() string
}

type Position interface {
	Filename() string
	Offset() int
	Line() int
	Column() int
}

type RelationType uint8

const (
	RelationTypeAssociationOneOne RelationType = iota + 1
	RelationTypeAssociationOneMany
	RelationTypeAssociation
	RelationTypeComposition
	RelationTypeEmbedding
	RelationTypeAggregation
	RelationTypeAggregationRoot
	RelationTypeDependency
	RelationTypeImplementation
	RelationTypeAbstraction
	RelationTypeAttribution
	RelationTypeBehavior
	RelationTypeNone
)

type RelationPos interface {
	From() Position
	To() Position
}

type RelationMeta interface {
	Type() RelationType
	Position() RelationPos
}

type Relation interface {
	Type() RelationType
	From() Object
}

type DependenceRelation interface {
	Relation
	DependsOn() Object
}

type CompositionRelation interface {
	Relation
	Child() Object
}

type EmbeddingRelation interface {
	Relation
	Embedded() Object
}

type ImplementationRelation interface {
	Relation
	Implements() []Object
	Implemented(object Object)
}

type AssociationRelation interface {
	Relation
	Refer() Object
	AssociationType() RelationType
}

type DiagramType uint8

const (
	PlainDiagram DiagramType = iota + 1
	TableDiagram
)

type Diagram interface {
	Name() string
	SubDiagrams() []SubDiagram
	Type() DiagramType
	Edges() []Edge
}

type SubDiagram interface {
	Name() string
	Nodes() []Node
	Summary() []Element
	SubGraphs() []SubDiagram
	Print()
}

type Node interface {
	ID() string
	Name() string
	Color() string
}

type Nodes []Node

type Element interface {
	Node
	Children() []Nodes
}

type Edge interface {
	From() string
	To() string
	Count() int
	Type() RelationType
	Pos() []RelationPos
}

type ObjColor string

const (
	ColorAggregate   ObjColor = "#ffd966ff"
	ColorEntity      ObjColor = "#ffe599ff"
	ColorValueObject ObjColor = "#a2c4c9ff"
	ColorFactory     ObjColor = "#cfe2f3ff"
	ColorService     ObjColor = "#e69138ff"
	ColorWhite       ObjColor = "#ffffffff"
	ColorMethod      ObjColor = "#a4c2f4ff"
	ColorAttribute   ObjColor = "#f3f3f3ff"
	ColorInterface   ObjColor = "#9fc5e8ff"
	ColorClass       ObjColor = "#b4a7d6ff"
	ColorGeneral     ObjColor = "#f4ccccff"
	ColorFunc        ObjColor = "#ead1dcff"
)

type Domain interface {
	Domain() string
}

type DomainObj interface {
	Object
	Domain
	OriginIdentifier() ObjIdentifier
}

type DomainObjs []DomainObj

func (dos DomainObjs) Objects() []Object {
	var objs []Object
	for _, obj := range dos {
		objs = append(objs, obj)
	}
	return objs
}
