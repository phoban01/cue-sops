package main

import (
	"os"
	"strings"
	"time"

	"go.mozilla.org/sops/v3"
	"go.mozilla.org/sops/v3/aes"
	"go.mozilla.org/sops/v3/cmd/sops/common"
	"go.mozilla.org/sops/v3/cmd/sops/formats"
	"go.mozilla.org/sops/v3/config"
	"go.mozilla.org/sops/v3/keyservice"
	"go.mozilla.org/sops/v3/version"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/ast/astutil"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/parser"
	"cuelang.org/go/cue/token"
	"cuelang.org/go/encoding/json"
	cuejson "cuelang.org/go/pkg/encoding/json"
)

type Operation int

const (
	encrypt Operation = iota
	decrypt
)

var (
	inplace   bool = true
	operation Operation
)

func main() {
	if len(os.Args) != 3 {
		panic("must provide encrypt or decrypt argument and a cue file")
	}

	switch os.Args[1] {
	case "encrypt":
		operation = encrypt
	case "decrypt":
		operation = decrypt
	default:
		panic("must provide encrypt or decrypt")
	}

	filename := os.Args[2]

	tree, err := parser.ParseFile(filename, nil)
	if err != nil {
		panic(err)
	}

	ctx := cuecontext.New()

	cueData := ctx.BuildFile(tree)

	sopsStruct := cueData.LookupPath(cue.ParsePath("sops"))
	if sopsStruct.Exists() && operation == encrypt {
		os.Exit(0)
	} else if !sopsStruct.Exists() && operation == decrypt {
		os.Exit(1)
	}

	fields := make([]string, 0)
	for _, v := range getFieldsToEncode(tree) {
		name := v.(*ast.Field).Label.(*ast.Ident).Name
		fields = append(fields, name)
	}

	var output []byte

	switch operation {
	case encrypt:
		output, err = encryptCue(ctx, tree, cueData, fields)
		if err != nil {
			panic(err)
		}
	case decrypt:
		output, err = decryptCue(ctx, tree, cueData)
		if err != nil {
			panic(err)
		}
	}

	if err := os.WriteFile(filename, output, os.ModePerm); err != nil {
		panic(err)
	}
}

func encryptCue(ctx *cue.Context, tree *ast.File, v cue.Value, fields []string) ([]byte, error) {
	jsonData, err := cuejson.Marshal(v)
	if err != nil {
		return nil, err
	}

	encData, err := encryptData([]byte(jsonData), strings.Join(fields, "|"))
	if err != nil {
		return nil, err
	}

	encDataCue, err := json.Extract("encoded.cue", encData)
	if err != nil {
		return nil, err
	}

	expr := ctx.BuildExpr(encDataCue)

	metadata := &ast.Field{
		Label: ast.NewIdent("sops"),
		Value: ast.Embed(expr.LookupPath(cue.ParsePath("sops")).Syntax().(*ast.StructLit)).Expr,
	}

	result, err := insertEncodedValues(tree, expr, operation)
	if err != nil {
		return nil, err
	}

	comment := ast.NewLit(token.COMMENT, "\n// DO NOT EDIT: auto-generated by cue-sops")

	result.(*ast.File).Decls = append(result.(*ast.File).Decls, comment, metadata)

	return format.Node(result)
}

func decryptCue(ctx *cue.Context, tree *ast.File, v cue.Value) ([]byte, error) {
	jsonData, err := cuejson.Marshal(v)
	if err != nil {
		return nil, err
	}

	decData, err := decryptData([]byte(jsonData))
	if err != nil {
		return nil, err
	}

	decDataCue, err := json.Extract("decoded.cue", decData)
	if err != nil {
		return nil, err
	}

	expr := ctx.BuildExpr(decDataCue)

	result, err := insertEncodedValues(tree, expr, operation)
	if err != nil {
		return nil, err
	}

	// remove generated fields
	decls := result.(*ast.File).Decls
	for i, j := range decls {
		switch j.(type) {
		case *ast.Field:
			if j.(*ast.Field).Label.(*ast.Ident).Name == "sops" {
				decls = append(decls[:i], decls[i+1:]...)
			}
		}
	}

	result.(*ast.File).Decls = decls

	return format.Node(result)
}

func encryptData(data []byte, regex string) ([]byte, error) {
	sopsConf, err := config.LoadCreationRuleForFile(".sops.yaml", "*.cue", nil)
	if err != nil {
		return nil, err
	}

	store := common.StoreForFormat(formats.Json)

	branches, err := store.LoadPlainFile(data)
	if err != nil {
		return nil, err
	}

	tree := sops.Tree{
		Branches: branches,
		Metadata: sops.Metadata{
			UnencryptedSuffix: sopsConf.UnencryptedSuffix,
			EncryptedSuffix:   sopsConf.EncryptedSuffix,
			UnencryptedRegex:  sopsConf.UnencryptedRegex,
			EncryptedRegex:    regex,
			KeyGroups:         sopsConf.KeyGroups,
			ShamirThreshold:   sopsConf.ShamirThreshold,
			Version:           version.Version,
		},
	}

	svcs := []keyservice.KeyServiceClient{
		keyservice.NewLocalClient(),
	}

	dataKey, errs := tree.GenerateDataKeyWithKeyServices(svcs)
	if len(errs) > 0 {
		return nil, err
	}

	cipher := aes.NewCipher()
	unencryptedMac, err := tree.Encrypt(dataKey, cipher)
	if err != nil {
		return nil, err
	}

	tree.Metadata.LastModified = time.Now().UTC()
	tree.Metadata.MessageAuthenticationCode, err = cipher.Encrypt(unencryptedMac, dataKey, tree.Metadata.LastModified.Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	return common.StoreForFormat(formats.Json).EmitEncryptedFile(tree)
}

func decryptData(data []byte) ([]byte, error) {
	store := common.StoreForFormat(formats.Json)

	tree, err := store.LoadEncryptedFile(data)
	if err != nil {
		return nil, err
	}

	svcs := []keyservice.KeyServiceClient{
		keyservice.NewLocalClient(),
	}

	dataKey, err := tree.Metadata.GetDataKeyWithKeyServices(svcs)
	if err != nil {
		return nil, err
	}

	cipher := aes.NewCipher()
	_, err = tree.Decrypt(dataKey, cipher)
	if err != nil {
		return nil, err
	}

	res, err := common.StoreForFormat(formats.Json).EmitPlainFile(tree.Branches)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func getFieldsToEncode(tree *ast.File) []ast.Node {
	fields := make([]ast.Node, 0)
	ast.Walk(tree, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.Field:
			field := n.(*ast.Field)
			if len(field.Attrs) == 0 {
				return true
			}

			for _, a := range field.Attrs {
				ident, kind := a.Split()
				if ident != "secret" && kind != "sops" {
					continue
				}
				break
			}

			fields = append(fields, n)
		}
		return true
	}, nil)
	return fields
}

func insertEncodedValues(tree *ast.File, v cue.Value, op Operation) (ast.Node, error) {
	f := func(c astutil.Cursor) bool {
		n := c.Node()
		switch n.(type) {
		case *ast.Field:
			field := n.(*ast.Field)
			if len(field.Attrs) == 0 {
				return true
			}

			var ident, kind string
			for _, a := range field.Attrs {
				ident, kind = a.Split()
				if ident != "secret" && kind != "sops" {
					continue
				}
				break
			}

			if ident == "" {
				return true
			}

			pathItems := make([]string, 0)
			parent := c.Parent()
			for parent != nil {
				switch parent.Node().(type) {
				case *ast.Field:
					l := parent.Node().(*ast.Field).Label.(*ast.Ident).Name
					pathItems = append(pathItems, l)
				}
				parent = parent.Parent()
			}

			pathItems = append(pathItems, field.Label.(*ast.Ident).Name)

			updatedValue := v.LookupPath(cue.ParsePath(strings.Join(pathItems, "."))).Syntax()

			field.Value = updatedValue.(ast.Expr)

			c.Replace(field)
		}
		return true
	}

	return astutil.Apply(tree, f, nil), nil
}
