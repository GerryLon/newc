package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"
)

var fset = token.NewFileSet()

// GetPackageInfo get the Go package information in the dir
func GetPackageInfo(dir string) (*packages.Package, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles,
		Tests: false,
	}, dir)
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("cannot find any package in %v", dir)
	}
	return pkgs[0], nil
}

// IncludeMakeMark check whether a code file contains "newc" comment
func IncludeMakeMark(filepath string) (bool, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return false, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if isMakeComment(line) {
			return true, nil
		}
	}
	return false, nil
}

// BuildAST build an AST from the code file
func BuildAST(filename string) (*ast.File, error) {
	astFile, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to build AST from file(%v): %w", filename, err)
	}
	return astFile, nil
}

// ImportInfo the information of an import
type ImportInfo struct {
	Name string
	Path string
}

// StructInfo the information of a struct
type StructInfo struct {
	StructName string
	InitFlag   bool
	ValueFlag  bool
	InitError  bool
	Fields     []StructField
}

// StructField the information of a struct field
type StructField struct {
	Name    string
	Type    string
	Skipped bool
}

// ParseCodeFile parse structs and imports in a code file
func ParseCodeFile(filename string) ([]StructInfo, []ImportInfo, error) {
	structs := []StructInfo{}
	imports := []ImportInfo{}
	astFile, err := BuildAST(filename)
	if err != nil {
		return structs, imports, err
	}

	// 先收集所有的 init 方法信息
	initMethods := make(map[string]bool) // structName -> hasErrorReturn
	for _, decl := range astFile.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Recv == nil || funcDecl.Name.Name != "init" {
			continue
		}

		// 检查接收者类型
		if len(funcDecl.Recv.List) > 0 {
			recvType := funcDecl.Recv.List[0].Type
			var structName string
			if starExpr, ok := recvType.(*ast.StarExpr); ok {
				if ident, ok := starExpr.X.(*ast.Ident); ok {
					structName = ident.Name
				}
			} else if ident, ok := recvType.(*ast.Ident); ok {
				structName = ident.Name
			}

			if structName != "" {
				// 检查返回值是否只包含 error
				onlyReturnError := false
				if funcDecl.Type.Results != nil {
					if len(funcDecl.Type.Results.List) > 1 {
						// 只支持返回一个参数(error)或者无返回值的情况, 更多参数不支持
						return nil, nil, fmt.Errorf("init方法只支持返回error或者不返回任何值, 实际上结构体\"%v\"有%d个返回值", structName, len(funcDecl.Type.Results.List))
					}
					for _, result := range funcDecl.Type.Results.List {
						if ident, ok := result.Type.(*ast.Ident); ok && ident.Name == "error" {
							onlyReturnError = true
							break
						}
					}
				}
				initMethods[structName] = onlyReturnError
			}
		}
	}

	for _, decl := range astFile.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		var initMode bool
		var valueMode bool
		if genDecl.Tok == token.TYPE {
			if genDecl.Doc == nil {
				continue
			}
			needGen := false
			for _, doc := range genDecl.Doc.List {
				if isMakeComment(doc.Text) {
					needGen = true
					initMode = isInitModeEnable(doc.Text)
					valueMode = isValueModeEnable(doc.Text)
					break
				}
			}
			if !needGen {
				continue
			}
		}

		for _, spec := range genDecl.Specs {
			importSpec, ok := spec.(*ast.ImportSpec)
			if ok {
				var name string
				if importSpec.Name != nil {
					name = importSpec.Name.Name
				}
				imports = append(imports, ImportInfo{
					Name: name,
					Path: importSpec.Path.Value,
				})
				continue
			}

			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}
			structFields := []StructField{}
			for _, field := range structType.Fields.List {
				fieldType := types.ExprString(field.Type)
				var fieldName string
				if len(field.Names) > 0 {
					fieldName = field.Names[0].Name
				} else {
					// handle embeded struct cases just like this:
					// 		type Foo struct {
					//  		pkg.Struct,
					// 		}
					items := strings.Split(fieldType, ".")
					fieldName = items[len(items)-1]
					// handle pointer cases just like this:
					// 		type Foo struct {
					//  		*pkg.Struct,
					// 		}
					fieldName = strings.TrimPrefix(fieldName, "*")
				}
				structFields = append(structFields, StructField{
					Type:    fieldType,
					Name:    fieldName,
					Skipped: isSkippedField(field),
				})
			}

			// 检查是否有 init 方法返回 error
			initError := initMethods[typeSpec.Name.Name]

			structs = append(structs, StructInfo{
				StructName: typeSpec.Name.Name,
				Fields:     structFields,
				InitFlag:   initMode,
				ValueFlag:  valueMode,
				InitError:  initError,
			})
		}
	}
	return structs, imports, nil
}

// isMakeComment ...
func isMakeComment(s string) bool {
	s = strings.TrimSpace(s)
	return strings.HasPrefix(s, "//go:generate") && strings.Contains(s, "newc")
}

// isInitModeEnable check if this struct enable the init mode
func isInitModeEnable(s string) bool {
	return strings.Contains(s, "init")
}

// isValueModeEnable check if this struct enable the value mode
func isValueModeEnable(s string) bool {
	return strings.Contains(s, "value")
}

func isSkippedField(field *ast.Field) bool {
	if field.Tag == nil {
		return false
	}
	return strings.Contains(field.Tag.Value, `newc:"-"`)
}
