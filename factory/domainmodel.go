package factory

import (
	"github.com/dddplayer/core/entity"
)

func NewDomainModel() (*entity.DomainModel, error) {
	return &entity.DomainModel{}, nil
}
