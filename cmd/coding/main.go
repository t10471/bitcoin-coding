package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/stoewer/go-strcase"
)

var pathValue = "."
var typeName *string

func main() {
	typeName = flag.String("t", "", "must be set")
	flag.Parse()
	if *typeName == "" {
		log.Fatal("type must be set")
	}
	if err := realMain(*typeName); err != nil {
		log.Fatal(err.Error())
	}
}

func realMain(typeName string) error {
	pkg, err := makePackageInfo(pathValue)
	if err != nil {
		return err
	}
	ts, err := pkg.findTypeSpec(typeName)
	if err != nil {
		return err
	}
	si, err := makeStructureInfo(ts, typeName)
	if err != nil {
		return err
	}
	dir, _ := filepath.Abs(pathValue)
	err = outputToFile(dir, pkg.name, si)
	if err != nil {
		return err
	}
	return nil
}

func outputToFile(dir string, pkgName string, si *StructureInfo) error {
	t, err := template.New("coding").Parse(tpl)
	if err != nil {
		return fmt.Errorf("template Parse: %w", err)
	}
	src := new(bytes.Buffer)
	st := struct {
		PackageName   string
		StructureInfo *StructureInfo
	}{pkgName, si}
	if err := t.Execute(src, st); err != nil {
		return fmt.Errorf("template Execute: %w", err)
	}
	name := fmt.Sprintf("%s_coding.go", strcase.SnakeCase(si.TypeName))
	outputName := filepath.Join(dir, strings.ToLower(name))
	if err := ioutil.WriteFile(outputName, src.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}
	return nil
}
