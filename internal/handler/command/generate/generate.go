package generate

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fanky5g/ponzu/content"
	contentGenerator "github.com/fanky5g/ponzu/content/generator"
	"github.com/fanky5g/ponzu/content/generator/parser"
	generatorTypes "github.com/fanky5g/ponzu/content/generator/types"
	log "github.com/sirupsen/logrus"
	"go/format"
	"html/template"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var (
	rootPath string
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	rootPath = filepath.Join(filepath.Dir(b), "../../../..")
}

type GeneratorArgs struct {
	Arguments   []string
	Target      contentGenerator.Target
	ContentType content.Type
}

type Generator interface {
	CopyFiles(src, target string) error
	WriteTemplate(name string, w io.Writer) error
	Generate() error
}

type generator struct {
	args             GeneratorArgs
	templateRoot     string
	hasExistingTypes bool
}

func (g *generator) mapContentType(contentType content.Type) string {
	switch contentType {
	case content.TypePlain:
		return "content.TypePlain"
	case content.TypeContent:
		return "content.TypeContent"
	case content.TypeFieldCollection:
		return "content.TypeFieldCollection"
	}

	return ""
}

func (g *generator) WriteTemplate(name string, w io.Writer) error {
	templ := template.Must(template.New(name).ParseFiles(filepath.Join(g.templateRoot, name)))
	buf := &bytes.Buffer{}
	if err := templ.Execute(buf, map[string]interface{}{
		"Arguments":   g.args.Arguments,
		"Target":      g.args.Target,
		"ContentType": g.mapContentType(g.args.ContentType),
		"ModuleName":  "content-generator",
		"GoVersion":   strings.TrimPrefix(runtime.Version(), "go"),
	}); err != nil {
		return err
	}

	var err error
	if strings.HasSuffix(strings.TrimSuffix(name, ".tmpl"), ".go") {
		var fmtBuf []byte
		fmtBuf, err = format.Source(buf.Bytes())
		if err != nil {
			return fmt.Errorf("failed to format template: %s", err.Error())
		}

		_, err = io.Copy(w, bytes.NewBuffer(fmtBuf))
	} else {
		_, err = io.Copy(w, buf)
	}

	return err
}

func NewGenerator(args GeneratorArgs) (Generator, error) {
	hasExistingTypes := true
	_, err := os.Stat(filepath.Join(args.Target.Path.Root, args.Target.Path.Base, "ponzu_types.go"))
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		hasExistingTypes = false
	}

	return &generator{
		args:             args,
		templateRoot:     fmt.Sprintf("%s/internal/handler/command/generate/templates", rootPath),
		hasExistingTypes: hasExistingTypes,
	}, nil
}

func (g *generator) CopyFiles(src, target string) error {
	dirName := filepath.Clean(src)
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			targetPath := filepath.Join(target, strings.TrimPrefix(path, dirName))
			return g.copyFile(path, targetPath)
		}

		return nil
	})
}

func (g *generator) writeTemplateFile(templateName, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: go.mod. Error: %v", err)
	}

	defer func() {
		if err = f.Close(); err != nil {
			log.WithFields(log.Fields{
				"Error": err,
				"File":  f.Name(),
			}).Warning("Failed to close file")
		}
	}()

	return g.WriteTemplate(templateName, f)
}

func (g *generator) Generate() error {
	if g.hasExistingTypes {
		workingDir, err := os.MkdirTemp(os.TempDir(), "generator")
		if err != nil {
			return fmt.Errorf("failed to create temp dir: %v", err)
		}

		defer func() {
			if err = os.RemoveAll(workingDir); err != nil {
				log.WithFields(log.Fields{
					"Error":     err,
					"Directory": workingDir,
				}).Warning("Failed to delete temp working dir")
			}
		}()

		if err = g.writeTemplateFile("go.mod.tmpl", filepath.Join(workingDir, "go.mod")); err != nil {
			return fmt.Errorf("failed to write template: %v", err)
		}

		err = g.CopyFiles(
			filepath.Join(g.args.Target.Path.Root, g.args.Target.Path.Base),
			filepath.Join(workingDir, g.args.Target.Path.Base),
		)

		if err != nil {
			return fmt.Errorf("failed to copy entities: %v", err)
		}

		if err = g.writeTemplateFile("main.go.tmpl", filepath.Join(workingDir, "main.go")); err != nil {
			return fmt.Errorf("failed to write template: %v", err)
		}

		var output string
		output, err = g.exec("go run main.go", workingDir)
		if err != nil {
			errorMessage := err.Error()
			if output != "" {
				errorMessage = output
			}

			return fmt.Errorf("failed to generate content: %v", errorMessage)
		}

		if output != "" {
			log.Debug(output)
		}

		return nil
	}

	contentTypes := content.Types{
		Content:          make(map[string]content.Builder),
		FieldCollections: make(map[string]content.Builder),
		Definitions:      make(map[string]generatorTypes.TypeDefinition),
	}

	generatorInstance, err := contentGenerator.New(contentGenerator.Config{
		Types: contentTypes,
		Target: contentGenerator.Target{
			Path: contentGenerator.Path{
				Root: g.args.Target.Path.Root,
				Base: g.args.Target.Path.Base,
			},
			Package: g.args.Target.Package,
		},
	})

	if err != nil {
		return err
	}

	p, err := parser.New(contentTypes)
	if err != nil {
		return err
	}

	typeDefinition, err := p.ParseTypeDefinition(g.args.ContentType, g.args.Arguments)
	if err != nil {
		return err
	}

	return generatorInstance.Generate(g.args.ContentType, typeDefinition)
}

func (g *generator) exec(command, workingDir string) (string, error) {
	re := regexp.MustCompile(`\s+`)
	tokens := re.Split(command, -1)

	if len(tokens) < 1 {
		return "", fmt.Errorf("invalid command: %v", command)
	}

	commandName := tokens[0]
	args := tokens[1:]

	log.Info("Running command", commandName, args)
	cmd := exec.Command(commandName, args...)
	cmd.Dir = workingDir
	output, errCommandExec := cmd.CombinedOutput()
	if errCommandExec != nil {
		return formatOutput(output), errCommandExec
	}

	return formatOutput(output), nil
}

func formatOutput(output []byte) string {
	if len(output) == 0 {
		return ""
	}

	return string(output)
}

func (g *generator) copyFile(src, target string) error {
	targetDir := filepath.Dir(target)

	_, err := os.Stat(targetDir)
	if errors.Is(err, fs.ErrNotExist) {
		err = os.MkdirAll(targetDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	var srcFile *os.File
	srcFile, err = os.Open(src)
	if err != nil {
		return err
	}

	defer func() {
		if err = srcFile.Close(); err != nil {
			log.WithFields(log.Fields{
				"Error": err,
				"Path":  src,
			}).Warning("Failed to close src file")
		}
	}()

	var targetFile *os.File
	targetFile, err = os.Create(target)
	if err != nil {
		return err
	}

	defer func() {
		if err = targetFile.Close(); err != nil {
			log.WithFields(log.Fields{
				"Error": err,
				"Path":  target,
			}).Warning("Failed to close target file")
		}
	}()

	_, err = io.Copy(targetFile, srcFile)
	return err
}
