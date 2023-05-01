package factory

import (
	"github.com/dddplayer/core/entity"
)

func NewDomainModel(name string, repo entity.Repository) (*entity.DomainModel, error) {
	return &entity.DomainModel{
		Name: name,
		Repo: repo,
	}, nil
}
