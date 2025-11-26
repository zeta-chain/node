// This program generates chaos-wrappers for client interfaces using reflection.
// It places the generated file in /zetaclient/mode/chaos/generated.go.
//
// This program also creates a sample JSON percentages file and places it in
// /zetaclient/mode/chaos/generate/sample.go.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"text/template"

	"github.com/zeta-chain/node/zetaclient/chains/bitcoin"
	"github.com/zeta-chain/node/zetaclient/chains/evm"
	solana "github.com/zeta-chain/node/zetaclient/chains/solana/repo"
	"github.com/zeta-chain/node/zetaclient/chains/sui"
	"github.com/zeta-chain/node/zetaclient/chains/ton"
	"github.com/zeta-chain/node/zetaclient/chains/tssrepo"
	"github.com/zeta-chain/node/zetaclient/chains/zrepo"
)

type Data struct {
	Interfaces []*Interface
	Imports    map[string]string

	sample map[string](map[string]uint)
}

type Interface struct {
	Name    string
	Type    string
	Methods []*Method
}

type Method struct {
	Name string
	Ins  []*Pair
	Outs []*Pair
}

type Pair struct {
	First   string
	Second  string
	IsError bool
}

func main() {
	// Read the template file.
	templateFile, err := os.ReadFile("generate/base.go.templ")
	if err != nil {
		panic(fmt.Sprintf("failed to read the template file: %v", err))
	}

	// Initialize the template and set auxiliary functions.
	tmpl := template.New("chaos")
	tmpl.Funcs(map[string]any{
		// sep returns the separator string when the element of the slice being processed is not
		// the last element of the slice.
		"sep": func(slice reflect.Value, i int, sep string) string {
			if slice.Kind() != reflect.Slice {
				panic("invalid value")
			}
			if i < slice.Len()-1 {
				return sep
			}
			return ""
		},
		"canfail": func(mthd *Method) bool { return mthd.canFail() },
		"space":   func() string { return " " },
		"tab":     func() string { return "\t" },
		"newline": func() string { return "\n" },
	})

	// Parse the template.
	tmpl, err = tmpl.Parse(string(templateFile))
	if err != nil {
		panic(fmt.Sprintf("failed to parse the template: %v", err))
	}

	data := &Data{
		Imports:    make(map[string]string),
		Interfaces: make([]*Interface, 0),
		sample:     make(map[string](map[string]uint)),
	}

	// Process the client interfaces.
	processInterface[zrepo.ZetacoreClient](data)
	processInterface[bitcoin.BitcoinClient](data)
	processInterface[evm.EVMClient](data)
	processInterface[solana.SolanaClient](data)
	processInterface[sui.SuiClient](data)
	processInterface[ton.TONClient](data)
	processInterface[tssrepo.TSSClient](data)

	// Create the generated file.
	generatedFile, err := os.Create("generated.go")
	if err != nil {
		panic(fmt.Sprintf("failed to create the generated file: %v", err))
	}
	defer generatedFile.Close()

	// Write to the generated file.
	err = tmpl.Execute(generatedFile, data)
	if err != nil {
		panic(fmt.Sprintf("failed to execute the template: %v", err))
	}

	// Create the sample JSON file.
	err = data.createSample()
	if err != nil {
		panic(fmt.Sprintf("failed to create the sample file: %v", err))
	}
}

// ------------------------------------------------------------------------------------------------

// processInterface processes an interface type.
// It fills the Data object with relevant information for the code generation.
func processInterface[T any](data *Data) {
	typ := reflect.TypeFor[T]()

	if typ.Kind() != reflect.Interface {
		panic("type is not an interface")
	}

	numMethod := typ.NumMethod()

	itfc := &Interface{
		Name:    typ.Name(),
		Type:    data.stringify(typ),
		Methods: make([]*Method, 0, numMethod),
	}

	data.Interfaces = append(data.Interfaces, itfc)
	data.sample[itfc.Name] = make(map[string]uint)

	for i := range numMethod {
		method := typ.Method(i)
		numIn := method.Type.NumIn()
		numOut := method.Type.NumOut()

		mthd := &Method{
			Name: method.Name,
			Ins:  make([]*Pair, 0, numIn),
			Outs: make([]*Pair, 0, numOut),
		}
		itfc.Methods = append(itfc.Methods, mthd)

		for i := range numIn {
			in := method.Type.In(i)
			mthd.Ins = append(mthd.Ins, &Pair{
				First:  fmt.Sprintf("in%d", i),
				Second: data.stringify(in),
			})
		}

		for i := range numOut {
			out := method.Type.Out(i)
			mthd.Outs = append(mthd.Outs, &Pair{
				First:   fmt.Sprintf("out%d", i),
				Second:  data.stringify(out),
				IsError: out.Name() == "error",
			})
		}

		if mthd.canFail() {
			data.sample[itfc.Name][mthd.Name] = 50
		}
	}
}

func (data *Data) stringify(typ reflect.Type) string {
	if path := typ.PkgPath(); path != "" {
		return data.typeWithPath(typ)
	}

	switch typ.Kind() {
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", typ.Len(), data.stringify(typ.Elem()))
	case reflect.Chan:
		return fmt.Sprintf("%s %s", typ.ChanDir().String(), data.stringify(typ.Elem()))
	case reflect.Func:
		panic("unimplemented")
	case reflect.Interface:
		return data.typeWithPath(typ)
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s", data.stringify(typ.Key()), data.stringify(typ.Elem()))
	case reflect.Pointer:
		return fmt.Sprintf("*%s", data.stringify(typ.Elem()))
	case reflect.Slice:
		return fmt.Sprintf("[]%s", data.stringify(typ.Elem()))
	case reflect.Struct:
		return data.typeWithPath(typ)
	case reflect.UnsafePointer:
		panic("unimplemented")
	default:
		return typ.Name()
	}
}

// typeWithPath returns a type with its package identifier, like "m0.Context".
func (data *Data) typeWithPath(typ reflect.Type) string {
	name := typ.Name()
	path := typ.PkgPath()

	if path == "" {
		return name
	}

	alias, ok := data.Imports[path]
	if !ok {
		alias = fmt.Sprintf("m%d", len(data.Imports))
		data.Imports[path] = alias
	}
	return fmt.Sprintf("%s.%s", alias, typ.Name())
}

// createSample creates a sample file with the percentages.
func (data *Data) createSample() error {
	s, err := json.MarshalIndent(data.sample, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile("generate/sample.json", s, 0600)
	if err != nil {
		return err
	}

	return nil
}

// canFail returns true if Method returns an error, and false otherwise.
func (mthd *Method) canFail() bool {
	for _, out := range mthd.Outs {
		if out.IsError {
			return true
		}
	}
	return false
}
