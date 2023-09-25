package application

import (
	"bytes"
	archFactory "github.com/dddplayer/dp/internal/domain/arch/factory"
	"github.com/dddplayer/dp/internal/domain/arch/repository"
	"github.com/dddplayer/dp/internal/domain/code/entity"
	"github.com/dddplayer/dp/internal/domain/dot/factory"
)

func TacticGraph(mainPkgPath, domain string,
	objRepo repository.ObjectRepository,
	relRepo repository.RelationRepository) (string, error) {

	return generateTacticGraph(mainPkgPath, domain, objRepo, relRepo, false, false)
}

func DetailTacticGraph(mainPkgPath, domain string,
	objRepo repository.ObjectRepository,
	relRepo repository.RelationRepository) (string, error) {

	return generateTacticGraph(mainPkgPath, domain, objRepo, relRepo, true, false)
}

func generateTacticGraph(mainPkgPath, domain string,
	objRepo repository.ObjectRepository,
	relRepo repository.RelationRepository,
	all, composition bool) (string, error) {

	arch, err := archFactory.NewArch(domain, objRepo, relRepo)
	if err != nil {
		return "", err
	}

	c, err := entity.NewCode(mainPkgPath, domain)
	if err != nil {
		return "", err
	}

	if err := c.VisitFast(arch.ObjectHandler()); err != nil {
		return "", err
	}

	g, err := arch.TacticGraph(&options{
		all:         all,
		composition: composition,
	})
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
