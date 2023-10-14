package entity

import (
	"errors"
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/arch/repository"
	"github.com/dddplayer/dp/pkg/datastructure/directed"
	"path"
	"strings"
)

type MessageFlow struct {
	directory       *Directory
	objRepo         repository.ObjectRepository
	relationDigraph *RelationDigraph
	mainPkgPath     string
	endPkgPath      string
}

type DirFilter struct {
	pkgSet []string
	paths  [][]*directed.Node
}

func (sf *DirFilter) IsValid(dir string) bool {
	for _, pkg := range sf.pkgSet {
		if strings.HasPrefix(pkg, dir) {
			return true
		}
	}
	return false
}

func (mf *MessageFlow) newDirFilter() (*DirFilter, error) {
	if n := mf.relationDigraph.FindNodeByKey(mf.mainFuncPath()); n != nil {
		ps := mf.relationDigraph.FindPathsToPrefix(mf.mainFuncPath(), mf.endPkgPath)
		var validPkgs []string
		for _, p := range ps {
			for _, n := range p {
				validPkgs = append(validPkgs, path.Dir(n.Key))
			}
		}
		return &DirFilter{pkgSet: validPkgs, paths: ps}, nil
	}
	return nil, errors.New("main func not found")
}

func (mf *MessageFlow) buildDiagram() (*Diagram, error) {
	dirFilter, err := mf.newDirFilter()
	if err != nil {
		return nil, err
	}

	gm, err := NewGeneralModel(mf.objRepo, mf.directory)
	if err != nil {
		return nil, err
	}
	if err := gm.GroupingWithFilter(dirFilter); err != nil {
		return nil, err
	}

	g, err := NewDiagram(mf.modulePath(), arch.TableDiagram)
	if err != nil {
		return nil, err
	}

	if err := gm.addRootGroupToDiagram(g); err != nil {
		return nil, err
	}

	var preIdentifier arch.ObjIdentifier
	for _, p := range dirFilter.paths {
		for i, n := range p {
			if preIdentifier == nil {
				preIdentifier = n.Value.(arch.ObjIdentifier)
				continue
			}
			current := n.Value.(arch.ObjIdentifier)
			metas, err := mf.relationDigraph.RelationMetas(
				preIdentifier,
				current,
			)
			if err != nil {
				return nil, err
			}
			repeatMetas := repeatElements(metas, i)
			if err := g.AddRelations(preIdentifier.ID(), current.ID(), repeatMetas); err != nil {
				return nil, err
			}

			preIdentifier = current
		}
	}

	return g, nil
}

func repeatElements(input []arch.RelationMeta, N int) []arch.RelationMeta {
	result := make([]arch.RelationMeta, 0)

	for _, element := range input {
		if element.Type() == arch.RelationTypeDependency {
			for i := 0; i < N; i++ {
				result = append(result, element)
			}
		} else {
			result = append(result, element)
		}
	}

	return result
}

func (mf *MessageFlow) mainFuncPath() string {
	mfp := fmt.Sprintf("%s/%s", mf.mainPkgPath, "main")
	return mfp
}

func (mf *MessageFlow) modulePath() string {
	return trimSuffixSlash(findCommonSubstring(mf.mainPkgPath, mf.endPkgPath))
}

func findCommonSubstring(str1, str2 string) string {
	var commonSubstring string

	// 确保 str1 的长度小于等于 str2 的长度
	if len(str1) > len(str2) {
		str1, str2 = str2, str1
	}

	// 遍历较短的字符串，逐个字符比较
	for i := 0; i < len(str1); i++ {
		if str1[i] == str2[i] {
			commonSubstring += string(str1[i])
		} else {
			break // 一旦不匹配，停止比较
		}
	}

	return commonSubstring
}

func trimSuffixSlash(s string) string {
	return strings.TrimSuffix(s, "/")
}
