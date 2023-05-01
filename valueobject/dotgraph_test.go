package valueobject

import (
	"github.com/dddplayer/core/dot/entity"
	"reflect"
	"testing"
)

func TestNewDotGraph(t *testing.T) {
	g := NewDotGraph("test")
	if g == nil {
		t.Error("NewDotGraph returned nil")
	}
	if g.Name() != "test" {
		t.Errorf("Expected name %q, got %q", "test", g.Name())
	}
	if len(g.Nodes()) != 0 {
		t.Errorf("Expected empty nodes slice, got %v", g.Nodes())
	}
	if len(g.Edges()) != 0 {
		t.Errorf("Expected empty edges slice, got %v", g.Edges())
	}
}

func TestNewDotEdge(t *testing.T) {
	e := NewDotEdge("from", "to")
	if e == nil {
		t.Error("NewDotEdge returned nil")
	}
	if e.From() != "from" {
		t.Errorf("Expected From %q, got %q", "from", e.From())
	}
	if e.To() != "to" {
		t.Errorf("Expected To %q, got %q", "to", e.To())
	}
}

func TestNewNode(t *testing.T) {
	n := NewNode("test")
	if n == nil {
		t.Error("NewNode returned nil")
	}
	if n.Name() != "test" {
		t.Errorf("Expected name %q, got %q", "test", n.Name())
	}
	if len(n.Elements()) != 0 {
		t.Errorf("Expected empty elements slice, got %v", n.Elements())
	}
}

func TestNewElement(t *testing.T) {
	e := NewElement("id", "name")
	if e == nil {
		t.Error("NewElement returned nil")
	}
	if e.Port() != "id" {
		t.Errorf("Expected Port %q, got %q", "id", e.Port())
	}
	if e.Name() != "name" {
		t.Errorf("Expected Name %q, got %q", "name", e.Name())
	}
	if e.Color() != "" {
		t.Errorf("Expected empty Color, got %q", e.Color())
	}
	if len(e.Attributes()) != 0 {
		t.Errorf("Expected empty Attributes slice, got %v", e.Attributes())
	}
}

func TestNewAttribute(t *testing.T) {
	a := NewAttribute("id", "name", "color")
	if a == nil {
		t.Error("NewAttribute returned nil")
	}
	if a.Port() != "id" {
		t.Errorf("Expected Port %q, got %q", "id", a.Port())
	}
	if a.Name() != "name" {
		t.Errorf("Expected Name %q, got %q", "name", a.Name())
	}
	if a.Color() != "color" {
		t.Errorf("Expected Color %q, got %q", "color", a.Color())
	}
}

func TestDotGraph_AppendNode(t *testing.T) {
	g := NewDotGraph("test")
	n := NewNode("node1")
	g.AppendNode(n)
	if !reflect.DeepEqual(g.Nodes(), []entity.DotNode{n}) {
		t.Errorf("Expected nodes %v, got %v", []entity.DotNode{n}, g.Nodes())
	}
}

func TestDotNode_Append(t *testing.T) {
	n := NewNode("node1")
	e1 := NewElement("id1", "element1")
	e2 := NewElement("id2", "element2")
	n.Append(e1, e2)
	if !reflect.DeepEqual(n.Elements(), []entity.DotElement{e1, e2}) {
		t.Errorf("Expected elements %v, got %v", []entity.DotElement{e1, e2}, n.Elements())
	}
}

func TestDotElement_Append(t *testing.T) {
	e := NewElement("id1", "element1")
	a1 := NewAttribute("id1", "attribute1", "red")
	a2 := NewAttribute("id2", "attribute2", "green")
	e.Append(a1)
	e.Append(a2)

	if !reflect.DeepEqual(e.Attributes(), []interface{}{a1, a2}) {
		t.Errorf("Expected attributes %v, got %v", []interface{}{a1, a2}, e.Attributes())
	}
}
