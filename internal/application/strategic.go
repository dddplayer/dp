package application

import (
	"bytes"
	archFactory "github.com/dddplayer/dp/internal/domain/arch/factory"
	"github.com/dddplayer/dp/internal/domain/arch/repository"
	"github.com/dddplayer/dp/internal/domain/code/entity"
	"github.com/dddplayer/dp/internal/domain/dot/factory"
)

func StrategicGraph(mainPkgPath, domain string, deep bool,
	objRepo repository.ObjectRepository,
	relRepo repository.RelationRepository) (string, error) {

	arch, err := archFactory.NewArch(domain, objRepo, relRepo)
	if err != nil {
		return "", err
	}

	c, err := entity.NewCode(mainPkgPath, domain)
	if err != nil {
		return "", err
	}

	if deep {
		if err := c.VisitDeep(arch.ObjectHandler()); err != nil {
			return "", err
		}
	} else {
		if err := c.VisitFast(arch.ObjectHandler()); err != nil {
			return "", err
		}
	}

	g, err := arch.StrategicGraph()
	if err != nil {
		return "", err
	}

	dot, err := factory.NewDotBuilder(g).Build()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := dot.Write(&buf); err != nil {
		return "", err
	}

	return string(buf.Bytes()), nil
}
