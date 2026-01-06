package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/emicklei/proto"
)

// ModuleData holds the documentation for a single module
type ModuleData struct {
	PackageName string
	Messages    []MessageData
	Files       []FileData
}

// MessageData holds the documentation for a single message
type MessageData struct {
	Name         string
	Comment      string
	ProtoMessage string
	HasComment   bool
}

// FileData holds information about files from docs/spec
type FileData struct {
	Name     string
	Content  string
	FilePath string
}

func main() {
	var startDir string
	var outputFile string
	var docsSpecDir string

	flag.StringVar(&startDir, "input", "./proto", "Input directory to start searching for proto files")
	flag.StringVar(&outputFile, "output", "./docs/spec/generated.md", "Output file for generated markdown")
	flag.StringVar(&docsSpecDir, "docs-spec", "./docs/spec", "Directory containing existing spec documentation")
	flag.Parse()

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputFile)
	err := os.MkdirAll(outputDir, 0750)
	if err != nil {
		fmt.Printf("Error creating output directory %q: %v\n", outputDir, err)
		return
	}

	// Collect all module data
	var modules []ModuleData
	err = filepath.Walk(startDir, func(path string, f os.FileInfo, err error) error {
		return visit(path, f, err, &modules, docsSpecDir)
	})
	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", startDir, err)
		return
	}

	// Write all data to single file
	err = writeToSingleFile(outputFile, modules)
	if err != nil {
		fmt.Printf("Error writing to file %q: %v\n", outputFile, err)
	}
}

func visit(path string, f os.FileInfo, err error, modules *[]ModuleData, docsSpecDir string) error {
	if err != nil {
		fmt.Printf("Error visiting %q: %v\n", path, err)
		return err
	}

	if f.IsDir() {
		return nil
	}

	if filepath.Ext(path) == ".proto" {
		err := processProtoFile(path, modules, docsSpecDir)
		if err != nil {
			fmt.Printf("Error processing proto file %q: %v\n", path, err)
			return err
		}
	}

	return nil
}

func processProtoFile(path string, modules *[]ModuleData, docsSpecDir string) error {
	safePath := filepath.Clean(path)
	reader, err := os.Open(safePath)
	if err != nil {
		fmt.Printf("Error opening proto file %q: %v\n", path, err)
		return err
	}
	/* #nosec G307 */
	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		fmt.Printf("Error parsing proto file %q: %v\n", path, err)
		return err
	}

	var packageName string
	var msgServices []*proto.Service
	messageMap := make(map[string]*proto.Message)

	proto.Walk(definition,
		proto.WithPackage(func(p *proto.Package) {
			packageName = p.Name
		}),
		proto.WithService(func(s *proto.Service) {
			if s.Name == "Msg" {
				msgServices = append(msgServices, s)
			}
		}),
		proto.WithMessage(func(m *proto.Message) {
			messageMap[m.Name] = m
		}),
	)

	if len(messageMap) == 0 {
		return nil
	}

	if packageName != "" && len(msgServices) > 0 {
		moduleName := getLastSegmentOfPackageName(packageName)

		// Read all files from docs/spec for this module
		files := readAllFilesFromDocsSpec(docsSpecDir, moduleName)

		moduleData := ModuleData{
			PackageName: packageName,
			Messages:    []MessageData{},
			Files:       files,
		}

		// Get the current working directory
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error getting current working directory: %v\n", err)
			return err
		}

		// Search for the corresponding Go function
		goFunctionPath := filepath.Join(currentDir, "x")

		for _, service := range msgServices {
			for _, element := range service.Elements {
				if rpc, ok := element.(*proto.RPC); ok {
					functionName := strings.TrimPrefix(rpc.RequestType, "Msg")
					if message, found := messageMap[rpc.RequestType]; found {
						functionComment, functionFound := findFunctionInGoFiles(functionName, goFunctionPath)

						messageData := MessageData{
							Name:         rpc.RequestType,
							Comment:      functionComment,
							ProtoMessage: messageToString(message),
							HasComment:   functionFound && functionComment != "",
						}
						moduleData.Messages = append(moduleData.Messages, messageData)
					} else {
						messageData := MessageData{
							Name:         rpc.RequestType,
							Comment:      "",
							ProtoMessage: rpc.RequestType,
							HasComment:   false,
						}
						moduleData.Messages = append(moduleData.Messages, messageData)
					}
				}
			}
		}

		*modules = append(*modules, moduleData)
	}
	return nil
}

func readAllFilesFromDocsSpec(docsSpecDir, moduleName string) []FileData {
	var files []FileData
	moduleDir := filepath.Join(docsSpecDir, moduleName)

	// Check if module directory exists
	if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
		return files
	}

	// Read all files in the module directory
	err := filepath.Walk(moduleDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(moduleDir, path)
		if err != nil {
			fmt.Printf("Error getting relative path for %q: %v\n", path, err)
			return nil
		}
		if strings.HasPrefix(relPath, "..") {
			fmt.Printf("Skipping suspicious file path: %q\n", path)
			return nil
		}

		// #nosec G304 -- path is validated to be within moduleDir above
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Printf("Error reading file %q: %v\n", path, err)
			return nil
		}

		files = append(files, FileData{
			Name:     f.Name(),
			Content:  string(content),
			FilePath: relPath,
		})

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking module directory %q: %v\n", moduleDir, err)
	}

	return files
}

func writeToSingleFile(outputFile string, modules []ModuleData) error {
	// #nosec G304
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating file %q: %v\n", outputFile, err)
		return err
	}
	/* #nosec G307 */
	defer file.Close()

	_, err = file.WriteString("---\ntitle: Node Modules Specification\n---\n\n")
	if err != nil {
		fmt.Printf("Error writing to file %q: %v\n", outputFile, err)
		return err
	}

	for _, module := range modules {
		// Write module heading
		moduleName := getLastSegmentOfPackageName(module.PackageName)
		_, err = fmt.Fprintf(file, "## %s\n\n", moduleName)
		if err != nil {
			fmt.Printf("Error writing module heading to file %q: %v\n", outputFile, err)
			return err
		}

		// Write all files from docs/spec first
		for _, fileData := range module.Files {
			_, err = file.WriteString(fileData.Content)
			if err != nil {
				fmt.Printf("Error writing file content to file %q: %v\n", outputFile, err)
				return err
			}
			_, err = file.WriteString("\n\n")
			if err != nil {
				fmt.Printf("Error writing newlines to file %q: %v\n", outputFile, err)
				return err
			}
		}

		// Write messages section
		if len(module.Messages) > 0 {
			_, err = fmt.Fprintf(file, "### Messages\n\n")
			if err != nil {
				fmt.Printf("Error writing messages section to file %q: %v\n", outputFile, err)
				return err
			}

			// Write messages for this module
			for _, message := range module.Messages {
				_, err = fmt.Fprintf(file, "#### %s\n\n", message.Name)
				if err != nil {
					return err
				}
				if message.HasComment {
					_, err = fmt.Fprintf(file, "%s\n", message.Comment)
					if err != nil {
						return err
					}
				}
				_, err = fmt.Fprintf(file, "```proto\n")
				if err != nil {
					fmt.Printf("Error writing to file %q: %v\n", outputFile, err)
					return err
				}
				_, err = file.WriteString(message.ProtoMessage)
				if err != nil {
					fmt.Printf("Error writing to file %q: %v\n", outputFile, err)
					return err
				}
				_, err = file.WriteString("```\n\n")
				if err != nil {
					fmt.Printf("Error writing to file %q: %v\n", outputFile, err)
					return err
				}
			}
		}
	}

	return nil
}

func messageToString(message *proto.Message) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("message %s {\n", message.Name))

	for _, element := range message.Elements {
		switch field := element.(type) {
		case *proto.NormalField:
			sb.WriteString(fmt.Sprintf("\t%s %s = %d;\n", field.Type, field.Name, field.Sequence))
		case *proto.MapField:
			sb.WriteString(fmt.Sprintf("\tmap<%s, %s> %s = %d;\n", field.KeyType, field.Type, field.Name, field.Sequence))
		case *proto.Oneof:
			sb.WriteString(fmt.Sprintf("\toneof %s {\n", field.Name))
			for _, of := range field.Elements {
				if oneOfField, ok := of.(*proto.OneOfField); ok {
					sb.WriteString(fmt.Sprintf("\t\t%s %s = %d;\n", oneOfField.Type, oneOfField.Name, oneOfField.Sequence))
				}
			}
			sb.WriteString("\t}\n")
		}
	}

	sb.WriteString("}\n")

	return sb.String()
}

func findFunctionInGoFiles(functionName string, startDir string) (string, bool) {
	var functionDoc string
	found := false

	fileSet := token.NewFileSet()
	err := filepath.Walk(startDir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		node, err := parser.ParseFile(fileSet, path, nil, parser.ParseComments)
		if err != nil {
			fmt.Printf("Error parsing Go file %q: %v\n", path, err)
			return err
		}

		for _, decl := range node.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {
				if funcDecl.Name.Name == functionName {
					if funcDecl.Doc != nil {
						functionDoc = funcDecl.Doc.Text()
					}
					found = true
					return filepath.SkipDir
				}
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", startDir, err)
	}

	return functionDoc, found
}

func getLastSegmentOfPackageName(packageName string) string {
	segments := strings.Split(packageName, ".")
	return segments[len(segments)-1]
}
