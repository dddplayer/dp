package entity

import (
	"errors"
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/pkg/datastructure/directory"
	"path"
	"path/filepath"
	"strings"
)

type Directory struct {
	root     *directory.TreeNode
	walkErrs []error
}

func NewDirectory(paths []string) *Directory {
	return &Directory{
		root:     directory.BuildDirectoryTree(paths),
		walkErrs: []error{},
	}
}

func (d *Directory) ArchDesignPattern() arch.DesignPattern {
	if d.isHexagon() {
		return arch.DesignPatternHexagon
	}

	return arch.DesignPatternPlain
}

func (d *Directory) isHexagon() bool {
	if d.root.Children[string(arch.HexagonDirectoryCmd)] != nil &&
		d.root.Children[string(arch.HexagonDirectoryInternal)] != nil &&
		d.root.Children[string(arch.HexagonDirectoryPkg)] != nil {

		internal := d.root.Children[string(arch.HexagonDirectoryInternal)]
		if internal.Children[string(arch.HexagonDirectoryDomain)] != nil {
			domain := internal.Children[string(arch.HexagonDirectoryDomain)]

			if len(domain.Children) > 0 {
				if domain.Children[string(arch.HexagonDirectoryEntity)] == nil &&
					domain.Children[string(arch.HexagonDirectoryValueObject)] == nil {
					return true
				}
			}
		}
	}

	return false
}

func (d *Directory) DomainDir() (string, error) {
	if d.isHexagon() {
		return path.Join(d.root.Name,
			string(arch.HexagonDirectoryInternal),
			string(arch.HexagonDirectoryDomain)), nil
	}

	return "", errors.New("invalid arch")
}

func (d *Directory) RootDir() string {
	return d.root.Name
}

func (d *Directory) HexagonDirectory(dir string) arch.HexagonDirectory {
	if arch.HexagonDirectoryDomain == arch.HexagonDirectory(dir) {
		return arch.HexagonDirectoryDomain
	}
	switch arch.HexagonDirectory(path.Base(dir)) {
	case arch.HexagonDirectoryEntity:
		return arch.HexagonDirectoryEntity
	case arch.HexagonDirectoryValueObject:
		return arch.HexagonDirectoryValueObject
	case arch.HexagonDirectoryRepository:
		return arch.HexagonDirectoryRepository
	case arch.HexagonDirectoryFactory:
		return arch.HexagonDirectoryFactory
	default:
		if arch.HexagonDirectoryDomain == arch.HexagonDirectory(path.Dir(dir)) {
			return arch.HexagonDirectoryAggregate
		}
	}

	return arch.HexagonDirectoryInvalid
}

func (d *Directory) WalkDir(dir string, cb func(string, []arch.ObjIdentifier) error) {
	targetDir, err := d.getTargetDir(dir)
	if err != nil {
		d.walkErrs = append(d.walkErrs, err)
		return
	}

	if node := d.root.GetNode(targetDir); node != nil {
		directory.Walk(node, func(dir string, val any) {
			if val != nil {
				if err := cb(dir, val.([]arch.ObjIdentifier)); err != nil {
					d.walkErrs = append(d.walkErrs, err)
				}
			} else {
				if err := cb(dir, nil); err != nil {
					d.walkErrs = append(d.walkErrs, err)
				}
			}
		})
	}
}

func (d *Directory) WalkRootDir(cb func(string, []arch.ObjIdentifier) error) {
	directory.Walk(d.root, func(dir string, val any) {
		if val != nil {
			if err := cb(dir, val.([]arch.ObjIdentifier)); err != nil {
				d.walkErrs = append(d.walkErrs, err)
			}
		} else {
			if err := cb(dir, nil); err != nil {
				d.walkErrs = append(d.walkErrs, err)
			}
		}
	})
}

func (d *Directory) WalkErrs() []error {
	return d.walkErrs
}

func (d *Directory) ParentDir(dir string) string {
	return path.Base(path.Dir(dir))
}

func (d *Directory) isValid(dir string) bool {
	return strings.HasPrefix(dir, d.root.Name)
}

func (d *Directory) getTargetDir(dir string) (string, error) {
	if d.isValid(dir) {
		return strings.TrimPrefix(dir, d.root.Name+string(filepath.Separator)), nil
	}

	return "", errors.New("invalid dir")
}

func (d *Directory) AddObjs(dir string, objs []arch.ObjIdentifier) error {
	if d.isRoot(dir) {
		d.root.Value = objs
		return nil
	}

	targetDir, err := d.getTargetDir(dir)
	if err != nil {
		return err
	}

	if err := d.root.AddValue(targetDir, objs); err != nil {
		return err
	}
	return nil
}

func (d *Directory) GetObjs(targetDir string) ([]arch.ObjIdentifier, error) {
	if objs, err := d.root.GetValue(targetDir); err == nil {
		return objs.([]arch.ObjIdentifier), nil
	}
	return nil, errors.New("not found")
}

func (d *Directory) isRoot(dir string) bool {
	return dir == d.root.Name
}
