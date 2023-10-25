package entity

import (
	"errors"
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/arch/repository"
	"github.com/dddplayer/dp/internal/domain/arch/valueobject"
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
	modulePath      string
}

func (mf *MessageFlow) newDirFilter() (*DirFilter, error) {
	if n := mf.relationDigraph.FindNodeByKey(mf.mainFuncPath()); n != nil {
		ps := mf.relationDigraph.FindPathsToPrefix(mf.mainFuncPath(), mf.endPkgPath)

		validPkgs := make(map[string]bool)
		objs := make([]arch.ObjIdentifier, 0, len(ps))
		for _, p := range ps {
			for _, n := range p {
				dir := path.Dir(n.Key)
				if ok := validPkgs[dir]; !ok {
					validPkgs[dir] = true
				}

				objs = append(objs, n.Value.(arch.ObjIdentifier))
			}
		}

		keys := make([]string, 0, len(validPkgs))
		for key := range validPkgs {
			keys = append(keys, key)
		}
		return &DirFilter{pkgSet: keys, paths: ps, objs: objs}, nil
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
	gm.GroupingWithFilter(dirFilter)

	g, err := NewDiagram(mf.modulePath, arch.TableDiagram)
	if err != nil {
		return nil, err
	}

	if err := gm.addRootGroupToDiagram(g); err != nil {
		return nil, err
	}

	var preIdentifier arch.ObjIdentifier

	for _, p := range dirFilter.paths {
		for _, n := range p {
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

			fromId := preIdentifier.ID()
			toId := current.ID()
			if err := g.AddRelations(fromId, toId, metas); err != nil {
				return nil, err
			}
			preIdentifier = current
		}
	}

	return g, nil
}

func (mf *MessageFlow) mainFuncPath() string {
	mfp := fmt.Sprintf("%s/%s", mf.mainPkgPath, "main")
	return mfp
}

type DirFilter struct {
	pkgSet []string
	paths  [][]*directed.Node
	objs   []arch.ObjIdentifier
}

func (sf *DirFilter) IsValid(dir string) bool {
	for _, pkg := range sf.pkgSet {
		if strings.HasPrefix(pkg, dir) {
			return true
		}
	}
	return false
}

func (sf *DirFilter) isExist(id arch.ObjIdentifier) bool {
	for _, obj := range sf.objs {
		if obj.ID() == id.ID() {
			return true
		}
	}
	return false
}

func (sf *DirFilter) FilterObjs(sourceData []arch.Object) []arch.Object {
	visitedObjects := make(map[string]struct{})
	result := make([]arch.Object, 0)

	for _, sourceObject := range sourceData {
		if function, ok := sourceObject.(*valueobject.Function); ok {
			for _, targetIdentifier := range sf.objs {
				if targetIdentifier.ID() == function.Identifier().ID() {
					if _, visited := visitedObjects[function.Identifier().ID()]; !visited {
						result = append(result, function)
						visitedObjects[function.Identifier().ID()] = struct{}{}
					}
					if function.Receiver != nil {
						if rcv := getReceiverFromSourceData(function.Receiver, sourceData); rcv != nil {
							if _, visited := visitedObjects[function.Receiver.ID()]; !visited {
								result = append(result, rcv)
								visitedObjects[function.Receiver.ID()] = struct{}{}
							}
						}
					}
				}
			}
		}
	}

	return result
}

func getReceiverFromSourceData(receiver arch.ObjIdentifier, sourceData []arch.Object) arch.Object {
	for _, sourceObject := range sourceData {
		if sourceObject.Identifier().ID() == receiver.ID() {
			return sourceObject
		}
	}
	return nil
}
