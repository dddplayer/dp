package application

import (
	"bytes"
	"fmt"
	archFactory "github.com/dddplayer/dp/internal/domain/arch/factory"
	"github.com/dddplayer/dp/internal/domain/arch/repository"
	"github.com/dddplayer/dp/internal/domain/code/entity"
	"github.com/dddplayer/dp/internal/domain/dot/factory"
	"golang.org/x/mod/modfile"
	"os"
	"path"
	"path/filepath"
)

func MessageFlowGraph(mainPkgPath, domain string,
	objRepo repository.ObjectRepository, relRepo repository.RelationRepository) (string, error) {

	goModFilePath, err := findGoModFile(mainPkgPath)
	if err != nil {
		return "", err
	}

	modPath, err := modulePath(goModFilePath)
	if err != nil {
		return "", err
	}

	arch, err := archFactory.NewArch(modPath, objRepo, relRepo)
	if err != nil {
		return "", err
	}

	c, err := entity.NewCode(mainPkgPath, modPath)
	if err != nil {
		return "", err
	}

	if err := c.VisitFast(arch.ObjectHandler()); err != nil {
		return "", err
	}

	g, err := arch.MessageFlowDiagram(c.MainPkgPath(), domain)
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

func modulePath(modFilePath string) (string, error) {
	// 读取go.mod文件内容
	modBytes, err := os.ReadFile(modFilePath)
	if err != nil {
		return "", err
	}

	// 解析go.mod文件内容
	modFile, err := modfile.Parse("go.mod", modBytes, nil)
	if err != nil {
		return "", err
	}

	// 打印模块的名称和版本
	return modFile.Module.Mod.Path, nil
}

func findGoModFile(startDir string) (string, error) {
	// 从指定的目录开始向上逐级查找包含 go.mod 文件的目录
	dir := startDir
	for {
		goModPath := filepath.Join(dir, "go.mod")
		_, err := os.Stat(goModPath)
		if err == nil {
			return path.Join(dir, "go.mod"), nil
		}

		// 如果已经到达根目录，仍未找到 go.mod 文件，则返回错误
		if dir == "/" || dir == "." {
			return "", fmt.Errorf("cannot find go.mod file")
		}

		// 向上级目录移动一级
		dir = filepath.Dir(dir)
	}
}
