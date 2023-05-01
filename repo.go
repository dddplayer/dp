package main

import (
	"fmt"
	"github.com/dddplayer/core/datastructure/radix"
	"github.com/dddplayer/core/entity"
	"github.com/dddplayer/core/valueobject"
)

type Repo struct {
	Tree *radix.Tree
}

func (r *Repo) Find(id valueobject.Identifier) entity.DomainObject {
	if obj, ok := r.Tree.Get(id.String()); ok {
		if obj != nil {
			return obj.(entity.DomainObject)
		}
	}
	return nil
}

func (r *Repo) Insert(obj entity.DomainObject) error {
	if ok := r.Tree.Insert(obj.Identifier().String(), obj); ok {
		return nil
	}
	return fmt.Errorf("insert failed")
}

func (r *Repo) Walk(cb func(obj entity.DomainObject) error) {
	r.Tree.Walk(func(prefix string, v any, ws radix.WalkState) radix.WalkStatus {
		if ws == radix.WalkIn {
			if obj, ok := v.(entity.DomainObject); ok {
				if err := cb(obj); err != nil {
					return radix.WalkStop
				}
			}
		}
		return radix.WalkContinue
	})
}
