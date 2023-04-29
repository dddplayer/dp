package entity

import (
	"fmt"
	"golang.org/x/exp/slices"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

func TestMainPackages(t *testing.T) {
	// 创建临时目录
	tempDir, err := ioutil.TempDir(".", "test_main_packages")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 创建真实的Go文件
	goFile := fmt.Sprintf("package main\n\nfunc main() {\n}\n")
	err = ioutil.WriteFile(tempDir+"/main.go", []byte(goFile), 0644)
	if err != nil {
		t.Fatalf("Failed to write Go file: %v", err)
	}

	// 使用packages.Load方法加载包
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
	}
	pkgs, err := packages.Load(cfg, tempDir)
	if err != nil {
		t.Fatalf("Failed to load packages: %v", err)
	}

	// 转换为SSA
	prog, ps := ssautil.AllPackages(pkgs, ssa.SanityCheckFunctions)
	if err != nil {
		t.Fatalf("Failed to create SSA: %v", err)
	}
	prog.Build()

	// 测试函数
	mains, err := mainPackages(ps)

	// 检查返回值
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(mains) != 1 {
		t.Errorf("Expected 1 main package, but got %d", len(mains))
	}

	if mains[0].Pkg.Name() != "main" {
		t.Errorf("Expected package with name 'main', but got %v", mains[0].Pkg.Name())
	}
}

func TestPkg_VisitFile(t *testing.T) {
	// 创建临时文件夹
	tmpDir, err := ioutil.TempDir(".", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 写入临时文件
	sourceCode := `
package main
import "fmt"

type Greeter interface {
	Greet()
}

type SuperMan struct {
	*Person
}

type Person struct {
	Name string
	Age  int
}

func (p *Person) SayHello() {
	fmt.Println("Hello, I'm", p.Name)
}

func (p *Person) Greet() {
	p.SayHello()
}

func main() {
	p := &Person{
		Name: "John",
		Age:  18,
	}
	p.SayHello()
}
`
	tmpFile := filepath.Join(tmpDir, "main.go")
	if err := ioutil.WriteFile(tmpFile, []byte(sourceCode), 0644); err != nil {
		t.Fatal(err)
	}

	// 加载包
	cfg := &packages.Config{
		Mode: packages.LoadAllSyntax,
		Dir:  tmpDir,
	}
	pkgs, err := packages.Load(cfg, "./...")
	if err != nil {
		t.Fatal(err)
	}

	// 构造测试对象
	pkg := &Pkg{
		Path:          "example",
		DomainPkgPath: "example",
		Initial:       pkgs,
		mainPkgPath:   "example",
	}

	// 测试 VisitFile 方法
	nodeList := make([]*Node, 0)
	linkList := make([]*Link, 0)
	pkg.VisitFile(func(node *Node) {
		nodeList = append(nodeList, node)
	}, func(link *Link) {
		linkList = append(linkList, link)
	})

	// 检查结果
	if len(nodeList) != 10 {
		t.Errorf("unexpected number of nodes: %d", len(nodeList))
	}
	expectedNodes := []string{"Greeter", "SuperMan", "Person", "main", "Name", "Age", "SayHello", "Greet"}
	for _, n := range nodeList {
		if !slices.Contains(expectedNodes, n.ID.Name()) {
			t.Errorf("unexpected node: %s", n.ID.Name())
		}
	}
	if len(linkList) != 1 {
		t.Errorf("unexpected number of links: %d", len(linkList))
	}
	expectedLinks := []string{"from Person to Person"}
	for _, l := range linkList {
		if !slices.Contains(expectedLinks, fmt.Sprintf("from %s to %s", l.From.ID.Name(), l.To.ID.Name())) {
			t.Errorf("unexpected link: from %v to %v", l.From, l.To)
		}
	}
}

func TestPkg_CallGraph(t *testing.T) {
	tmpdir, err := ioutil.TempDir(".", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	// create a test package with two functions calling each other
	src := `package main

func foo() {
	bar()
}

func bar() {
	foo()
}

func main() {
	foo()
	bar()
}
`
	err = ioutil.WriteFile(filepath.Join(tmpdir, "test.go"), []byte(src), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// load the test package
	p := &Pkg{Path: tmpdir}
	if err := p.Load(); err != nil {
		t.Fatal(err)
	}

	// initialize the call graph
	var links []*Link
	linkCB := func(link *Link) {
		links = append(links, link)
	}
	if err := p.CallGraph(linkCB); err != nil {
		t.Fatal(err)
	}

	// verify the call graph contains two nodes and one link
	if len(links) != 4 {
		t.Fatalf("expected 1 link, got %d", len(links))
	}

	expectedLinks := []string{"from main to foo", "from main to bar", "from foo to bar", "from bar to foo"}
	for _, l := range links {
		if !slices.Contains(expectedLinks, fmt.Sprintf("from %s to %s", l.From.ID.Name(), l.To.ID.Name())) {
			t.Errorf("unexpected link: from %v to %v", l.From, l.To)
		}
	}
}

func TestPkg_InterfaceImplements(t *testing.T) {
	// create a temporary directory for test files
	tmpDir, err := ioutil.TempDir(".", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// 写入临时文件
	sourceCode := `
package main
import "fmt"

type Greeter interface {
	Greet()
}

type Person struct {
	Name string
	Age  int
}

func (p *Person) Greet() {
	fmt.Println(p.Name, p.Age)
}

type Foo interface {
	Bar() string
}

type fooImpl struct {}

func (f *fooImpl) Bar() string {
	return "Hello, world!"
}

func main() {
	p := &Person{
		Name: "John",
		Age:  18,
	}
	p.Greet()

	fi := &fooImpl{}
	fmt.Println(fi.Bar())
}
`
	tmpFile := filepath.Join(tmpDir, "main.go")
	if err := ioutil.WriteFile(tmpFile, []byte(sourceCode), 0644); err != nil {
		t.Fatal(err)
	}

	// load the test package
	p := &Pkg{
		Path:          tmpDir,
		DomainPkgPath: "test",
	}
	if err := p.Load(); err != nil {
		t.Fatal(err)
	}

	// test the InterfaceImplements method
	var links []*Link
	linkCB := func(l *Link) {
		links = append(links, l)
	}
	p.InterfaceImplements(linkCB)

	// check the result
	if len(links) != 2 {
		t.Fatalf("expected 2 links, got %d", len(links))
	}

	expectedLinks := []string{"from Greeter to Person", "from Foo to fooImpl"}
	for _, l := range links {
		if !slices.Contains(expectedLinks, fmt.Sprintf("from %s to %s", l.From.ID.Name(), l.To.ID.Name())) {
			t.Errorf("unexpected link: from %v to %v", l.From, l.To)
		}
	}
}
