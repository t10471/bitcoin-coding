package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/fatih/structtag"
	"golang.org/x/tools/go/packages"
)

type packageInfo struct {
	name  string
	files []*ast.File
}

func makePackageInfo(path string) (*packageInfo, error) {
	cfg := &packages.Config{Mode: packages.NeedName | packages.NeedSyntax, Tests: false}
	packageList, err := packages.Load(cfg, path)
	if err != nil {
		return nil, err
	}
	if len(packageList) != 1 {
		return nil, fmt.Errorf("error: %d packages found", len(packageList))
	}
	p := packageList[0]
	return &packageInfo{name: p.Name, files: p.Syntax}, nil
}

type typeInspector struct {
	typeName string
	typeSpec *ast.TypeSpec
}

func (p *packageInfo) findTypeSpec(typeName string) (*ast.TypeSpec, error) {
	ti := &typeInspector{typeName: typeName}
	for _, file := range p.files {
		if file == nil {
			continue
		}
		ast.Inspect(file, ti.inspect)
		if ti.typeSpec != nil {
			return ti.typeSpec, nil
		}
	}
	return nil, errors.New("not found type")
}

func (t *typeInspector) inspect(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.TYPE {
		return true
	}
	for _, spec := range decl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		if typeSpec.Name.String() != t.typeName {
			continue
		}
		t.typeSpec = typeSpec
		return false
	}
	return true
}

type StructureInfo struct {
	TypeName      string
	ReceiverChar  string
	FieldInfoList []*fieldInfo
}

func makeStructureInfo(ts *ast.TypeSpec, typeName string) (*StructureInfo, error) {
	structType, ok := interface{}(ts.Type).(*ast.StructType)
	if !ok {
		return nil, errors.New("not ast.StructType")
	}
	si := &StructureInfo{
		TypeName:      typeName,
		ReceiverChar:  string(strings.ToLower(typeName)[0]),
		FieldInfoList: make([]*fieldInfo, 0, len(structType.Fields.List)),
	}
	for _, fi := range structType.Fields.List {
		fi, err := makeFieldInfo(fi)
		if err != nil {
			return nil, err
		}
		si.FieldInfoList = append(si.FieldInfoList, fi)
	}
	return si, nil
}

type FieldType int

//go:generate stringer -type=FieldType
const (
	FieldTypeSlice FieldType = iota
	FieldTypeStructure
	FieldTypeBaseType
)

type fieldInfo struct {
	FieldName      string
	TypeName       string
	FieldType      FieldType
	SliceCountName string
	BaseTypeName   string
}

func makeFieldInfo(field *ast.Field) (*fieldInfo, error) {
	typeName, err := exprToTypeName(field.Type)
	if err != nil {
		return nil, err
	}
	if len(field.Names) == 0 {
		return nil, errors.New("not allow embed type")
	}
	fi := &fieldInfo{
		TypeName:  typeName,
		FieldName: field.Names[0].String(),
		FieldType: FieldTypeStructure,
	}
	tg := &tag{}
	if err := tg.parse(field.Tag); err != nil {
		msg := "parseTag error at `%s` error is %s"
		return nil, fmt.Errorf(msg, field.Names[0].String(), err.Error())
	}
	if _, isSlice := field.Type.(*ast.ArrayType); isSlice {
		fi.FieldType = FieldTypeSlice
		if tg.countName == "" {
			return nil, errors.New("slice should have coding-count tag")
		}
		fi.SliceCountName = tg.countName
	} else if isBaseType(fi.TypeName) {
		fi.FieldType = FieldTypeBaseType
		s := strings.ReplaceAll(fi.TypeName, "basetype.", "")
		fi.BaseTypeName = s
	}
	return fi, nil
}

func exprToTypeName(expr ast.Expr) (string, error) {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name, nil
	}
	if selector, ok := expr.(*ast.SelectorExpr); ok {
		sel, err := exprToTypeName(selector.Sel)
		if err != nil {
			return "", err
		}
		x, err := exprToTypeName(selector.X)
		if err != nil {
			return "", err
		}
		return x + "." + sel, nil
	}
	if array, ok := expr.(*ast.ArrayType); ok {
		x, err := exprToTypeName(array.Elt)
		if err != nil {
			return "", err
		}
		return "[]" + x, nil
	}
	return "", errors.New("invalid type name")
}

var baseType = []string{
	"basetype.Hash",
	"basetype.VarInt",
	"basetype.Uint32",
	"basetype.Int32",
	"basetype.Uint32Time",
}

func isBaseType(name string) bool {
	for _, b := range baseType {
		if name == b {
			return true
		}
	}
	return false
}

type tag struct {
	countName string
}

func (t *tag) parse(b *ast.BasicLit) error {
	if b == nil {
		return nil
	}
	tg, err := structtag.Parse(strings.Trim(b.Value, "`"))
	if err != nil {
		return err
	}
	if tg == nil {
		return nil
	}
	if c, err := tg.Get("coding-count"); err == nil && c != nil {
		t.countName = c.Value()
	}
	return nil
}
