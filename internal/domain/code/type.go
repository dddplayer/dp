package code

type RelationShip int

const (
	None RelationShip = iota + 1
	OneOne
	OneMany
)

type Link struct {
	From     *Node
	To       *Node
	Relation RelationShip
}

type NodeType int

const (
	TypeGenIdent NodeType = 1 << iota
	TypeGenFunc
	TypeGenArray
	TypeGenStruct
	TypeGenStructField
	TypeGenStructEmbeddedField
	TypeGenInterface
	TypeGenInterfaceMethod
	TypeFunc
	TypeAny
	TypeNone
)

type Node struct {
	Meta   MetaInfo
	Pos    Position
	Parent *Node
	Type   NodeType
}

type NodeCB func(node *Node)
type LinkCB func(link *Link)

type Handler interface {
	NodeHandler(node *Node)
	LinkHandler(link *Link)
}

type CallGraphType string

const (
	CallGraphTypeStatic  CallGraphType = "static"
	CallGraphTypeCha                   = "cha"
	CallGraphTypeRta                   = "rta"
	CallGraphTypePointer               = "pointer"
)

type CallGraphMode uint8

const (
	CallGraphFastMode CallGraphMode = iota + 1
	CallGraphDeepMode
)

type Language interface {
	VisitFile(nodeCB NodeCB, linkCB LinkCB)
	InterfaceImplements(linkCB LinkCB)
	CallGraph(linkCB LinkCB, mode CallGraphMode) error
}

type MetaInfo interface {
	Pkg() string
	Name() string
	Parent() string
	HasParent() bool
}

type Position interface {
	Filename() string
	Offset() int
	Line() int
	Column() int
}
