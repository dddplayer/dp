package persistence

import (
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/pkg/datastructure/radix"
)

type RadixTree struct {
	Tree   *radix.Tree
	objIds []arch.ObjIdentifier
}

func (r *RadixTree) Find(id arch.ObjIdentifier) arch.Object {
	if obj, ok := r.Tree.Get(id.ID()); ok {
		if obj != nil {
			return obj.(arch.Object)
		}
	}
	return nil
}

func (r *RadixTree) Insert(obj arch.Object) error {
	if ok := r.Tree.Insert(obj.Identifier().ID(), obj); ok {
		r.objIds = append(r.objIds, obj.Identifier())
		return nil
	}
	return fmt.Errorf("insert failed")
}

func (r *RadixTree) All() []arch.ObjIdentifier {
	return r.objIds
}

func (r *RadixTree) Walk(cb func(obj arch.Object) error) {
	r.Tree.Walk(func(prefix string, v any, ws radix.WalkState) radix.WalkStatus {
		if ws == radix.WalkIn {
			if obj, ok := v.(arch.Object); ok {
				if err := cb(obj); err != nil {
					return radix.WalkStop
				}
			}
		}
		return radix.WalkContinue
	})
}

func (r *RadixTree) GetObjects(ids []arch.ObjIdentifier) ([]arch.Object, error) {
	objs := make([]arch.Object, 0)
	for _, id := range ids {
		obj := r.Find(id)
		if obj == nil {
			return nil, fmt.Errorf("object %s not found", id)
		}
		objs = append(objs, obj)
	}
	return objs, nil
}
