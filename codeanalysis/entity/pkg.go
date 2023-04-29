package entity

import (
	"fmt"
	"github.com/dddplayer/core/codeanalysis/valueobject"
	"go/ast"
	"go/build"
	"go/types"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"strings"
)

type NodeCB func(node *Node)
type LinkCB func(link *Link)

type Pkg struct {
	Path          string
	DomainPkgPath string
	Initial       []*packages.Package
	mainPkgPath   string
	prog          *ssa.Program
}

func (p *Pkg) Load() error {
	cfg := &packages.Config{
		Mode:       packages.LoadAllSyntax,
		Tests:      false,
		Dir:        "",
		BuildFlags: build.Default.BuildTags,
	}

	initial, err := packages.Load(cfg, p.Path)
	if err != nil {
		return err
	}
	if packages.PrintErrors(initial) > 0 {
		return fmt.Errorf("packages contain errors")
	}
	if len(initial) == 0 {
		return fmt.Errorf("package empty error")
	}

	p.Initial = initial
	p.mainPkgPath = p.Initial[0].PkgPath

	p.buildProg()

	return nil
}

func (p *Pkg) buildProg() {
	prog, _ := ssautil.AllPackages(p.Initial, 0)
	prog.Build()
	p.prog = prog
}

func (p *Pkg) VisitFile(nodeCB NodeCB, linkCB LinkCB) {
	packages.Visit(p.Initial, nil, func(pkg *packages.Package) {
		if strings.Contains(pkg.String(), p.DomainPkgPath) {

			for _, f := range pkg.Syntax {
				for _, decl := range f.Decls {
					declPos := valueobject.AstPosition(pkg, decl)

					switch decl.(type) {
					case *ast.GenDecl:
						genDecl := decl.(*ast.GenDecl)
						for _, spec := range genDecl.Specs {
							switch spec.(type) {
							case *ast.TypeSpec:
								typeSpec := spec.(*ast.TypeSpec)
								node := &Node{
									ID:     valueobject.NewIdentifier(pkg.ID, typeSpec.Name.Name),
									Pos:    declPos,
									Parent: nil,
									Type:   TypeGenIdent,
								}

								switch typeSpec.Type.(type) {
								case *ast.Ident:
									node.Type = TypeGenIdent
									nodeCB(node)
								case *ast.FuncType:
									node.Type = TypeGenFunc
									nodeCB(node)
								case *ast.ArrayType:
									node.Type = TypeGenArray
									nodeCB(node)

								case *ast.StructType:
									node.Type = TypeGenStruct
									nodeCB(node)

									structType := typeSpec.Type.(*ast.StructType)
									for _, field := range structType.Fields.List {
										fieldPos := valueobject.AstPosition(pkg, field)
										var fieldNode *Node
										if len(field.Names) > 0 {
											fieldNode = &Node{
												ID:     valueobject.NewIdentifier(node.ID.String(), field.Names[0].Name),
												Pos:    fieldPos,
												Parent: node,
												Type:   TypeGenStructField,
											}
											nodeCB(fieldNode)
										}

										exp := &expression{
											expr: field.Type,
											pkg:  pkg,
										}
										exp.visit(func(path, name string, ship RelationShip) {
											if fieldNode == nil {
												fieldNode = &Node{
													ID:     valueobject.NewIdentifier(node.ID.String(), name),
													Pos:    fieldPos,
													Parent: node,
													Type:   TypeGenStructField,
												}
												nodeCB(fieldNode)
											}
											linkCB(&Link{
												From: fieldNode,
												To: &Node{
													ID:  valueobject.NewIdentifier(path, name),
													Pos: nil,
												},
												Relation: ship,
											})
										})
									}

								case *ast.InterfaceType:
									node.Type = TypeGenInterface
									nodeCB(node)

									interfaceType := typeSpec.Type.(*ast.InterfaceType)
									for _, method := range interfaceType.Methods.List {
										for _, name := range method.Names {
											methodPos := valueobject.AstPosition(pkg, method)
											nodeCB(&Node{
												ID:     valueobject.NewIdentifier(node.ID.String(), name.Name),
												Pos:    methodPos,
												Parent: node,
												Type:   TypeGenInterfaceMethod,
											})
										}
									}
								}
							}
						}

					case *ast.FuncDecl:
						funcDecl := decl.(*ast.FuncDecl)
						funcNode := &Node{
							ID:     valueobject.NewIdentifier(pkg.ID, funcDecl.Name.Name),
							Pos:    declPos,
							Parent: nil,
							Type:   TypeFunc,
						}

						if funcDecl.Recv != nil {
							for _, rcv := range funcDecl.Recv.List {
								switch rcv.Type.(type) {
								case *ast.StarExpr:
									star := rcv.Type.(*ast.StarExpr)
									switch star.X.(type) {
									case *ast.Ident:
										rcvName := star.X.(*ast.Ident).Name
										rcvIdent := valueobject.NewIdentifier(pkg.ID, rcvName)
										funcNode.Parent = &Node{
											ID: rcvIdent,
										}
									}
								}
							}
						}

						nodeCB(funcNode)
					}
				}
			}
		}
	})
}

func (p *Pkg) InterfaceImplements(linkCB LinkCB) {
	pkgs := p.prog.AllPackages()

	namedMap := map[*types.Named]*ssa.Package{}
	var namedInterface []*types.Named
	var namedObj []*types.Named

	for _, pkg := range pkgs {
		if strings.Contains(pkg.String(), p.DomainPkgPath) {
			for _, m := range pkg.Members {
				o := m.Object()
				if obj, ok := o.(*types.TypeName); ok {
					if named, ok := obj.Type().(*types.Named); ok {
						if types.IsInterface(obj.Type()) {
							namedInterface = append(namedInterface, named)
						} else {
							namedObj = append(namedObj, named)
						}
						namedMap[named] = pkg
					}
				}
			}
		}
	}

	implMap := map[string][]*types.Named{}
	for _, i := range namedInterface {
		implMap[i.Obj().Name()] = []*types.Named{}
	}

	addImpl := func(iName string, obj *types.Named) {
		impls := implMap[iName]
		impls = append(impls, obj)
		implMap[iName] = impls
	}

	for _, o := range namedObj {
		for _, i := range namedInterface {
			if types.AssignableTo(o, i) {
				addImpl(i.Obj().Name(), o)
			} else if pU := types.NewPointer(o); types.AssignableTo(pU, i) {
				addImpl(i.Obj().Name(), o)
			}
		}
	}

	for _, i := range namedInterface {
		iNode := &Node{
			ID:   valueobject.NewIdentifier(i.Obj().Pkg().Path(), i.Obj().Name()),
			Pos:  valueobject.SsaPosition(namedMap[i], i.Obj()),
			Type: TypeGenInterface,
		}
		impls := implMap[i.Obj().Name()]
		for _, impl := range impls {
			linkCB(&Link{
				From: iNode,
				To: &Node{
					ID:  valueobject.NewIdentifier(impl.Obj().Pkg().Path(), impl.Obj().Name()),
					Pos: valueobject.SsaPosition(namedMap[impl], impl.Obj()),
				},
				Relation: OneOne,
			})
		}
	}
}

func mainPackages(pkgs []*ssa.Package) ([]*ssa.Package, error) {
	var mains []*ssa.Package
	for _, p := range pkgs {
		if p != nil && p.Pkg.Name() == "main" && p.Func("main") != nil {
			mains = append(mains, p)
		}
	}
	if len(mains) == 0 {
		return nil, fmt.Errorf("no main packages")
	}
	return mains, nil
}

func (p *Pkg) CallGraph(linkCB LinkCB) error {
	mains, err := mainPackages(p.prog.AllPackages())
	if err != nil {
		return err
	}
	res, err := pointer.Analyze(&pointer.Config{
		Mains:          mains,
		BuildCallGraph: true,
	})
	if err != nil {
		return err
	}

	res.CallGraph.DeleteSyntheticNodes()
	err = callgraph.GraphVisitEdges(res.CallGraph, func(edge *callgraph.Edge) error {
		caller := edge.Caller
		callee := edge.Callee

		if strings.Contains(caller.String(), p.DomainPkgPath) &&
			strings.Contains(callee.String(), p.DomainPkgPath) {

			if caller.Func.Name() == "init" || callee.Func.Name() == "init" {
				return nil
			}

			if caller.Func.Pkg != nil && callee.Func.Pkg != nil {
				pkgPath := caller.Func.Pkg.Pkg.Path()
				callerFuncName := caller.Func.Name()
				objName := callerFuncName

				recv := caller.Func.Signature.Recv()
				if strings.Contains(callerFuncName, "$") && caller.Func.Parent() != nil {
					recv = caller.Func.Parent().Signature.Recv()
					split := strings.Split(callerFuncName, "$")
					callerFuncName = split[0]
				}
				if recv != nil {
					objName = recv.Name()
					if obj, ok := recv.Type().(*types.Pointer); ok {
						switch o := obj.Elem().(type) {
						case *types.Named:
							objName = o.Obj().Name()
						}
					}
				}

				callerNode := &Node{
					ID:   valueobject.NewIdentifier(pkgPath, objName),
					Pos:  valueobject.SsaFuncPosition(caller.Func.Pkg, caller.Func),
					Type: TypeFunc,
				}

				pkgPath = callee.Func.Pkg.Pkg.Path()
				funcName := callee.Func.Name()
				objName = funcName

				recv = callee.Func.Signature.Recv()
				if recv != nil {
					objName = recv.Name()
					if obj, ok := recv.Type().(*types.Pointer); ok {
						switch o := obj.Elem().(type) {
						case *types.Named:
							objName = o.Obj().Name()
						}
					}
				}

				calleeNode := &Node{
					ID:  valueobject.NewIdentifier(pkgPath, objName),
					Pos: valueobject.SsaFuncPosition(callee.Func.Pkg, callee.Func),
				}

				callerNode.Pos = valueobject.SsaInstructionPosition(caller.Func.Pkg, edge.Site)

				linkCB(&Link{
					From:     callerNode,
					To:       calleeNode,
					Relation: OneOne,
				})
			}
		}

		return nil
	})

	return err
}
