package entity

import (
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/arch/valueobject"
	"testing"
)

func TestObjColor(t *testing.T) {
	mockAggregate := &valueobject.Aggregate{}
	mockEntity := &valueobject.Entity{}
	mockValueObject := &valueobject.ValueObject{}
	mockDomainAttr := &valueobject.DomainAttr{}
	mockDomainFunction := &valueobject.DomainFunction{}
	mockStringObj := &valueobject.StringObj{}
	mockObject := &MockObject{}

	t.Run("Test objColor with Aggregate", func(t *testing.T) {
		color := objColor(mockAggregate)
		if color != arch.ColorAggregate {
			t.Errorf("Expected color %s, but got %s", arch.ColorAggregate, color)
		}
	})

	t.Run("Test objColor with Entity", func(t *testing.T) {
		color := objColor(mockEntity)
		if color != arch.ColorEntity {
			t.Errorf("Expected color %s, but got %s", arch.ColorEntity, color)
		}
	})

	t.Run("Test objColor with ValueObject", func(t *testing.T) {
		color := objColor(mockValueObject)
		if color != arch.ColorValueObject {
			t.Errorf("Expected color %s, but got %s", arch.ColorValueObject, color)
		}
	})

	t.Run("Test objColor with DomainAttr", func(t *testing.T) {
		color := objColor(mockDomainAttr)
		if color != arch.ColorAttribute {
			t.Errorf("Expected color %s, but got %s", arch.ColorAttribute, color)
		}
	})

	t.Run("Test objColor with DomainFunction", func(t *testing.T) {
		color := objColor(mockDomainFunction)
		if color != arch.ColorFunc {
			t.Errorf("Expected color %s, but got %s", arch.ColorFunc, color)
		}
	})

	t.Run("Test objColor with StringObj", func(t *testing.T) {
		color := objColor(mockStringObj)
		if color != arch.ColorWhite {
			t.Errorf("Expected color %s, but got %s", arch.ColorWhite, color)
		}
	})

	t.Run("Test objColor with default", func(t *testing.T) {
		color := objColor(mockObject)
		if color != arch.ColorGeneral {
			t.Errorf("Expected color %s, but got %s", arch.ColorGeneral, color)
		}
	})
}

func TestObjColorWithParent(t *testing.T) {
	mockAggregate := &valueobject.Aggregate{}
	mockEntity := &valueobject.Entity{}
	mockValueObject := &valueobject.ValueObject{}
	mockDomainFunction := &valueobject.DomainFunction{}
	mockObject := &MockObject{}

	t.Run("Test objColorWithParent with Aggregate parent and DomainFunction object", func(t *testing.T) {
		color := objColorWithParent(mockDomainFunction, mockAggregate)
		if color != arch.ColorMethod {
			t.Errorf("Expected color %s, but got %s", arch.ColorMethod, color)
		}
	})

	t.Run("Test objColorWithParent with Entity parent and DomainFunction object", func(t *testing.T) {
		color := objColorWithParent(mockDomainFunction, mockEntity)
		if color != arch.ColorMethod {
			t.Errorf("Expected color %s, but got %s", arch.ColorMethod, color)
		}
	})

	t.Run("Test objColorWithParent with ValueObject parent and DomainFunction object", func(t *testing.T) {
		color := objColorWithParent(mockDomainFunction, mockValueObject)
		if color != arch.ColorMethod {
			t.Errorf("Expected color %s, but got %s", arch.ColorMethod, color)
		}
	})

	t.Run("Test objColorWithParent with other parent and DomainFunction object", func(t *testing.T) {
		color := objColorWithParent(mockDomainFunction, mockObject)
		if color != arch.ColorFunc {
			t.Errorf("Expected color %s, but got %s", arch.ColorFunc, color)
		}
	})
}
