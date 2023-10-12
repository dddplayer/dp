package entity

import (
	"fmt"
	"github.com/dddplayer/dp/internal/domain/arch"
	"github.com/dddplayer/dp/internal/domain/arch/repository"
	"strings"
)

type MessageFlow struct {
	directory       *Directory
	objRepo         repository.ObjectRepository
	relationDigraph *RelationDigraph
	mainPkgPath     string
	endPkgPath      string
}

func (mf *MessageFlow) buildDiagram() (*Diagram, error) {
	if n := mf.relationDigraph.FindNodeByKey(mf.mainFuncPath()); n != nil {
		ps := mf.relationDigraph.FindPathsToPrefix(mf.mainFuncPath(), mf.endPkgPath)
		for _, p := range ps {
			fmt.Println("------")
			for _, n := range p {
				fmt.Println(n.Key)
			}
		}
	} else {
		fmt.Println("????", mf.relationDigraph.FindNodeByKey("github.com/dddplayer/markdown/cmd.main"))
	}

	gm, err := NewGeneralModel(mf.objRepo, mf.directory)
	if err != nil {
		return nil, err
	}
	if err := gm.Grouping(); err != nil {
		return nil, err
	}

	g, err := NewDiagram(mf.modulePath(), arch.TableDiagram)
	if err != nil {
		return nil, err
	}

	if err := gm.addRootGroupToDiagram(g); err != nil {
		return nil, err
	}

	return g, nil
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
