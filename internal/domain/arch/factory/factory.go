package factory

import (
	"errors"
	"github.com/dddplayer/dp/internal/domain/arch/entity"
	"github.com/dddplayer/dp/internal/domain/arch/repository"
	"github.com/dddplayer/dp/internal/domain/arch/valueobject"
)

func NewArch(scope string, objRepo repository.ObjectRepository, relRepo repository.RelationRepository) (*entity.Arch, error) {
	if objRepo == nil {
		return nil, errors.New("objRepo cannot be nil")
	}
	if relRepo == nil {
		return nil, errors.New("relRepo cannot be nil")
	}
	return &entity.Arch{
		CodeHandler: &valueobject.CodeHandler{
			ObjRepo: objRepo,
			RelRepo: relRepo,
			Scope:   scope,
		},
	}, nil
}
