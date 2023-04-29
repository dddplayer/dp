package entity

import "github.com/dddplayer/core/codeanalysis/valueobject"

type NodeType string

const (
	TypeGenIdent           NodeType = "general_identifier"
	TypeGenFunc            NodeType = "general_function"
	TypeGenArray           NodeType = "general_array"
	TypeGenStruct          NodeType = "general_struct"
	TypeGenStructField     NodeType = "general_struct_field"
	TypeGenInterface       NodeType = "general_interface"
	TypeGenInterfaceMethod NodeType = "general_interface_method"
	TypeFunc               NodeType = "function"
)

type Node struct {
	ID     valueobject.Identifier
	Pos    valueobject.Position
	Parent *Node
	Type   NodeType
}
