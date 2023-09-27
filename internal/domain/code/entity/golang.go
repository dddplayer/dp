package entity

import (
	"fmt"
	"github.com/dddplayer/dp/internal/domain/code"
	"github.com/dddplayer/dp/internal/domain/code/valueobject"
	"go/ast"
	"go/build"
	"go/types"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/cha"
	"golang.org/x/tools/go/callgraph/rta"
	"golang.org/x/tools/go/callgraph/static"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"strings"
)

type Go struct {
	Path          string
	DomainPkgPath string
	Initial       []*packages.Package
	mainPkgPath   string
	prog          *ssa.Program
}

func (golang *Go) Load() error {
	cfg := &packages.Config{
		Mode:       packages.LoadAllSyntax,
		Tests:      false,
		Dir:        "",
		BuildFlags: build.Default.BuildTags,
	}

	initial, err := packages.Load(cfg, golang.Path)
	if err != nil {
		return err
	}
	if packages.PrintErrors(initial) > 0 {
		return fmt.Errorf("packages contain errors")
	}
	if len(initial) == 0 {
		return fmt.Errorf("package empty error")
	}

	golang.Initial = initial
	golang.mainPkgPath = golang.Initial[0].PkgPath

	golang.buildProg()

	return nil
}

func (golang *Go) buildProg() {
	prog, _ := ssautil.AllPackages(golang.Initial, 0)
	prog.Build()
	golang.prog = prog
}

func (golang *Go) VisitFile(nodeCB code.NodeCB, linkCB code.LinkCB) {
	packages.Visit(golang.Initial, nil, func(pkg *packages.Package) {
		if strings.Contains(pkg.String(), golang.DomainPkgPath) {

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
								node := &code.Node{
									Meta:   valueobject.NewMeta(pkg.ID, typeSpec.Name.Name),
									Pos:    declPos,
									Parent: nil,
									Type:   code.TypeGenIdent,
								}

								switch typeSpec.Type.(type) {
								case *ast.Ident:
									node.Type = code.TypeGenIdent
									nodeCB(node)
								case *ast.FuncType:
									node.Type = code.TypeGenFunc
									nodeCB(node)
								case *ast.ArrayType:
									node.Type = code.TypeGenArray
									nodeCB(node)

								case *ast.StructType:
									node.Type = code.TypeGenStruct
									nodeCB(node)

									var params valueobject.Params
									if typeSpec.TypeParams != nil {
										for _, f := range typeSpec.TypeParams.List {
											for _, name := range f.Names {
												params = append(params, valueobject.NewParam(name.Name))
											}
										}
									}

									structType := typeSpec.Type.(*ast.StructType)
									for _, field := range structType.Fields.List {
										fieldPos := valueobject.AstPosition(pkg, field)
										var fieldNode *code.Node
										if len(field.Names) > 0 {
											fieldNode = &code.Node{
												Meta:   valueobject.NewMetaWithParent(node.Meta.Pkg(), field.Names[0].Name, node.Meta.Name()),
												Pos:    fieldPos,
												Parent: node,
												Type:   code.TypeGenStructField,
											}
											nodeCB(fieldNode)
											linkCB(&code.Link{
												From:     node,
												To:       fieldNode,
												Relation: code.OneOne,
											})
										}

										exp := &expression{
											expr: field.Type,
											pkg:  pkg,
										}
										exp.visit(func(path, name string, ship code.RelationShip) {
											if params.Contains(name) {
												return
											}

											if fieldNode == nil {
												fieldNode = &code.Node{
													Meta:   valueobject.NewMetaWithParent(node.Meta.Pkg(), name, node.Meta.Name()),
													Pos:    fieldPos,
													Parent: node,
													Type:   code.TypeGenStructEmbeddedField,
												}
												nodeCB(fieldNode)
												linkCB(&code.Link{
													From:     node,
													To:       fieldNode,
													Relation: ship,
												})
											}
											linkCB(&code.Link{
												From: fieldNode,
												To: &code.Node{
													Meta: valueobject.NewMeta(path, name),
													Pos:  nil,
													Type: code.TypeAny,
												},
												Relation: ship,
											})
										})
									}

								case *ast.InterfaceType:
									node.Type = code.TypeGenInterface
									nodeCB(node)

									interfaceType := typeSpec.Type.(*ast.InterfaceType)
									for _, method := range interfaceType.Methods.List {
										for _, name := range method.Names {
											methodPos := valueobject.AstPosition(pkg, method)
											methodNode := &code.Node{
												Meta:   valueobject.NewMetaWithParent(node.Meta.Pkg(), name.Name, node.Meta.Name()),
												Pos:    methodPos,
												Parent: node,
												Type:   code.TypeGenInterfaceMethod,
											}
											nodeCB(methodNode)
											linkCB(&code.Link{
												From:     node,
												To:       methodNode,
												Relation: code.OneOne,
											})
										}
									}
								}
							}
						}

					case *ast.FuncDecl:
						funcDecl := decl.(*ast.FuncDecl)
						if funcDecl.Name.Name == "init" { // ignore init function
							continue
						}

						funcNode := &code.Node{
							Meta:   valueobject.NewMeta(pkg.ID, funcDecl.Name.Name),
							Pos:    declPos,
							Parent: nil,
							Type:   code.TypeFunc,
						}

						if funcDecl.Recv != nil {
							for _, rcv := range funcDecl.Recv.List {
								switch rcv.Type.(type) {
								case *ast.StarExpr:
									star := rcv.Type.(*ast.StarExpr)
									switch star.X.(type) {
									case *ast.Ident:
										rcvName := star.X.(*ast.Ident).Name
										rcvIdent := valueobject.NewMeta(pkg.ID, rcvName)
										funcNode.Parent = &code.Node{
											Meta: rcvIdent,
											Type: code.TypeAny,
											Pos:  valueobject.AstPosition(pkg, rcv),
										}
										funcNode.Meta = valueobject.NewMetaWithParent(rcvIdent.Pkg(), funcDecl.Name.Name, rcvIdent.Name())
									case *ast.IndexListExpr:
										//if funcDecl.Name.Name == "Signature" {
										//	l := star.X.(*ast.IndexListExpr)
										//
										//	fmt.Println("Signature recv: ")
										//	fmt.Printf("%#v\n", l.X)
										//	for _, index := range l.Indices {
										//		fmt.Printf("%#v\n", index)
										//	}
										//}
										rcvName := star.X.(*ast.IndexListExpr).X.(*ast.Ident).Name
										rcvIdent := valueobject.NewMeta(pkg.ID, rcvName)
										funcNode.Parent = &code.Node{
											Meta: rcvIdent,
											Type: code.TypeAny,
											Pos:  valueobject.AstPosition(pkg, rcv),
										}
										funcNode.Meta = valueobject.NewMetaWithParent(rcvIdent.Pkg(), funcDecl.Name.Name, rcvIdent.Name())
									}
								case *ast.Ident:
									rcvName := rcv.Type.(*ast.Ident).Name
									rcvIdent := valueobject.NewMeta(pkg.ID, rcvName)
									funcNode.Parent = &code.Node{
										Meta: rcvIdent,
										Type: code.TypeAny,
										Pos:  valueobject.AstPosition(pkg, rcv),
									}
									funcNode.Meta = valueobject.NewMetaWithParent(rcvIdent.Pkg(), funcDecl.Name.Name, rcvIdent.Name())
								}
							}
						}

						nodeCB(funcNode)
						if funcNode.Parent != nil {
							linkCB(&code.Link{
								From:     funcNode.Parent,
								To:       funcNode,
								Relation: code.OneOne,
							})
						}
					}
				}
			}
		}
	})
}

func (golang *Go) InterfaceImplements(linkCB code.LinkCB) {
	pkgs := golang.prog.AllPackages()

	namedMap := map[*types.Named]*ssa.Package{}
	var namedInterface []*types.Named
	var namedObj []*types.Named

	for _, pkg := range pkgs {
		if strings.Contains(pkg.String(), golang.DomainPkgPath) {
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
		iNode := &code.Node{
			Meta: valueobject.NewMeta(i.Obj().Pkg().Path(), i.Obj().Name()),
			Pos:  valueobject.SsaPosition(namedMap[i], i.Obj()),
			Type: code.TypeGenInterface,
		}
		impls := implMap[i.Obj().Name()]
		for _, impl := range impls {
			linkCB(&code.Link{
				From: &code.Node{
					Meta: valueobject.NewMeta(impl.Obj().Pkg().Path(), impl.Obj().Name()),
					Pos:  valueobject.SsaPosition(namedMap[impl], impl.Obj()),
					Type: code.TypeAny,
				},
				To:       iNode,
				Relation: code.OneOne,
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

func (golang *Go) callGraph(algo code.CallGraphType) (*callgraph.Graph, error) {
	var graph *callgraph.Graph

	switch algo {
	case code.CallGraphTypeStatic: // RA, reachability analysis
		graph = static.CallGraph(golang.prog)
	case code.CallGraphTypeCha: // Cha, class hierarchy analysis
		graph = cha.CallGraph(golang.prog)
	case code.CallGraphTypeRta: // Rta,rapid type analysis
		mains, err := mainPackages(golang.prog.AllPackages())
		if err != nil {
			return nil, err
		}
		var roots []*ssa.Function
		for _, main := range mains {
			roots = append(roots, main.Func("main"))
		}
		graph = rta.Analyze(roots, true).CallGraph
	case code.CallGraphTypePointer: // Pointer
		mains, err := mainPackages(golang.prog.AllPackages())
		if err != nil {
			return nil, err
		}
		config := &pointer.Config{
			Mains:          mains,
			BuildCallGraph: true,
		}
		res, err := pointer.Analyze(config)
		if err != nil {
			return nil, err
		}
		graph = res.CallGraph
	default:
		return nil, fmt.Errorf("invalid call graph type: %s", string(algo))
	}

	return graph, nil
}

func (golang *Go) CallGraph(linkCB code.LinkCB, mode code.CallGraphMode) error {
	var algo code.CallGraphType
	switch mode {
	case code.CallGraphFastMode:
		algo = code.CallGraphTypeStatic
	case code.CallGraphDeepMode:
		algo = code.CallGraphTypePointer
	default:
		algo = code.CallGraphTypeStatic
	}

	callGraph, err := golang.callGraph(algo)
	if err != nil {
		return err
	}

	callGraph.DeleteSyntheticNodes()
	err = callgraph.GraphVisitEdges(callGraph, func(edge *callgraph.Edge) error {
		caller := edge.Caller
		callee := edge.Callee

		if strings.Contains(caller.String(), golang.DomainPkgPath) &&
			strings.Contains(callee.String(), golang.DomainPkgPath) {

			if ignore(caller.Func.Name()) || ignore(callee.Func.Name()) {
				return nil
			}

			if caller.Func.Pkg != nil && callee.Func.Pkg != nil {
				callerNode := golang.functionNode(caller)
				calleeNode := golang.functionNode(callee)

				callerNode.Pos = valueobject.SsaInstructionPosition(caller.Func.Pkg, edge.Site)

				linkCB(&code.Link{
					From:     callerNode,
					To:       calleeNode,
					Relation: code.OneOne,
				})
			}
		}

		return nil
	})

	return err
}

func ignore(funcName string) bool {
	if funcName == "init" || strings.HasPrefix(funcName, "init#") || strings.HasPrefix(funcName, "init$") {
		return true
	}
	return false
}

func (golang *Go) functionNode(node *callgraph.Node) *code.Node {
	pkgPath := node.Func.Pkg.Pkg.Path()
	funcName := node.Func.Name()
	objName := funcName

	n := &code.Node{
		Meta:   valueobject.NewMeta(pkgPath, funcName),
		Pos:    valueobject.SsaFuncPosition(node.Func.Pkg, node.Func),
		Type:   code.TypeFunc,
		Parent: nil,
	}

	recv := node.Func.Signature.Recv()
	if strings.Contains(funcName, "$1") && node.Func.Parent() != nil {
		count := strings.Count(funcName, "$1")
		p := node.Func.Parent()
		for i := 1; i < count; i++ {
			if p != nil {
				p = p.Parent()
			}
		}
		if p != nil {
			recv = p.Signature.Recv()
		}
	}
	if recv != nil {
		objName = recv.Name()
		if obj, ok := recv.Type().(*types.Pointer); ok {
			switch o := obj.Elem().(type) {
			case *types.Named:
				objName = o.Obj().Name()
				n.Parent = &code.Node{
					Meta: valueobject.NewMeta(pkgPath, objName),
				}
				n.Meta = valueobject.NewMetaWithParent(n.Parent.Meta.Pkg(), funcName, n.Parent.Meta.Name())
			}
		} else if obj, ok := recv.Type().(*types.Named); ok {
			objName = obj.Obj().Name()
			n.Parent = &code.Node{
				Meta: valueobject.NewMeta(pkgPath, objName),
			}
			n.Meta = valueobject.NewMetaWithParent(n.Parent.Meta.Pkg(), funcName, n.Parent.Meta.Name())
		}
	}

	return n
}
