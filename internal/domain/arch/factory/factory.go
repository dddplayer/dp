package factory

import (
	"github.com/dddplayer/dp/internal/domain/arch/entity"
	"github.com/dddplayer/dp/internal/domain/arch/repository"
	"github.com/dddplayer/dp/internal/domain/arch/valueobject"
)

func NewArch(scope string, objRepo repository.ObjectRepository, relRepo repository.RelationRepository) (*entity.Arch, error) {
	return &entity.Arch{
		CodeHandler: &valueobject.CodeHandler{
			ObjRepo: objRepo,
			RelRepo: relRepo,
			Scope:   scope,
		},
	}, nil
}
