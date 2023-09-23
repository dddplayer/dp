package persistence

import (
	"github.com/dddplayer/dp/internal/domain/arch"
)

type Relations struct {
	relations []arch.Relation
}

func (kv *Relations) Insert(rel arch.Relation) error {
	kv.relations = append(kv.relations, rel)

	return nil
}

func (kv *Relations) Walk(walker func(rel arch.Relation) error) {
	for _, rel := range kv.relations {
		if err := walker(rel); err != nil {
			break
		}
	}
}
