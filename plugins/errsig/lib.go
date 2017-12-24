package errsig

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
)

func dostuff() {
	filename := "foo.go"

	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, filename, nil, parser.AllErrors)
	if err != nil {
		fmt.Printf("failed!, %v", err)
	}

	conf := types.Config{Importer: importer.Default()}
	info := &types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}
	if _, err := conf.Check(filename, fs, []*ast.File{f}, info); err != nil {
		log.Fatal(err) // type error
	}
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncType:
			fmt.Printf("%#v\n", x.Results.List)
			for _, o := range x.Results.List {
				tv, ok := info.Types[o.Type]
				if !ok {
					fmt.Printf("nil...\n")
					return false
				}

				if true == false {
					fmt.Println("blah", tv)
				}

				// fmt.Printf("info: %#v\n", tv)

				fmt.Printf("field: %#v \n", o.Type)
				switch t := o.Type.(type) {
				case *ast.Ident:
					fmt.Printf("object: %#v \n", t.Obj)
					if t.Obj != nil {
						fmt.Printf("typespec: %#v \n", t.Obj.Decl)
					}
				// case *ast.Object:
				// 	fmt.Printf("type: %#v \n", t)
				default:
					fmt.Println("unknown type.(type)")
				}
			}

		case *ast.FuncDecl:
			if x.Name.Name == "test" {

				fmt.Printf("input: \n")
				for _, p := range x.Type.Params.List {
					tv, ok := info.Types[p.Type]
					if !ok {
						fmt.Printf("nil...\n")
						return false
					}
					fmt.Printf("%v %v \n", p.Names, tv.Type)
				}

				fmt.Printf("output: \n")
				for _, o := range x.Type.Results.List {
					tv, ok := info.Types[o.Type]
					if !ok {
						fmt.Printf("nil...\n")
						return false
					}
					fmt.Printf("%v %v \n", o.Names, tv.Type)
				}
			}

		}
		return true
	})
}
