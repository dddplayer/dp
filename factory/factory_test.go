package factory

import (
	"github.com/dddplayer/core/entity"
	"github.com/dddplayer/core/valueobject"
	"testing"
)

type mockRepository struct {
	data map[valueobject.Identifier]entity.DomainObject
}

func (r *mockRepository) Find(id valueobject.Identifier) entity.DomainObject {
	return r.data[id]
}

func (r *mockRepository) Insert(obj entity.DomainObject) error {
	r.data[*obj.Identifier()] = obj
	return nil
}

func (r *mockRepository) Walk(cb func(obj entity.DomainObject) error) {
	for _, obj := range r.data {
		if err := cb(obj); err != nil {
			return
		}
	}
}

func TestNewDomainModel(t *testing.T) {
	repo := &mockRepository{data: make(map[valueobject.Identifier]entity.DomainObject)}
	dm, err := NewDomainModel("test", repo)

	if err != nil {
		t.Errorf("NewDomainModel() error = %v, wantErr %v", err, nil)
		return
	}

	if dm == nil {
		t.Errorf("NewDomainModel() dm = %v, want %v", dm, &entity.DomainModel{Name: "", Repo: repo})
		return
	}

	if dm.Name != "test" {
		t.Errorf("NewDomainModel() dm.Name = %v, want %v", dm.Name, "")
	}

	if dm.Repo != repo {
		t.Errorf("NewDomainModel() dm.Repo = %v, want %v", dm.Repo, repo)
	}
}
