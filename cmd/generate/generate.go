package generate

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fanky5g/ponzu/content"
	"github.com/fanky5g/ponzu/generator"
	contentGenerator "github.com/fanky5g/ponzu/generator/content"
	modelGenerator "github.com/fanky5g/ponzu/generator/models"
	"github.com/fanky5g/ponzu/util"
	log "github.com/sirupsen/logrus"
)

var rootPath string

func init() {
	_, b, _, _ := runtime.Caller(0)
	rootPath = filepath.Join(filepath.Dir(b), "../..")
}

func RunGenerator(contentType generator.Type, arguments []string) error {
	contentConfig := getContentConfig(contentType)
	modelConfig := getModelConfig()

	hasExistingTypes := true
	_, err := os.Stat(filepath.Join(contentConfig.Target.Path.Root, contentConfig.Target.Path.Base, "ponzu_types.go"))
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		hasExistingTypes = false
	}

	if hasExistingTypes {
		return RunGeneratorWithExistingTypes(contentType, arguments, map[string]generator.Config{
			"Content": contentConfig,
			"Models":  modelConfig,
		})
	}

	contentGeneratorInstance, err := contentGenerator.New(contentConfig, content.Types{
		Content:          make(map[string]content.Builder),
		FieldCollections: make(map[string]content.Builder),
		Definitions:      make(map[string]generator.TypeDefinition),
	})

	if err != nil {
		return err
	}

	generator.Register(contentGeneratorInstance)

	modelGeneratorInstance, err := modelGenerator.New(modelConfig, contentConfig)
	if err != nil {
		return err
	}

	generator.Register(modelGeneratorInstance)

	return generator.Run(contentType, arguments)
}

func RunGeneratorWithExistingTypes(contentType generator.Type, arguments []string, configs map[string]generator.Config) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

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

	for _, cfg := range configs {
		err = copyFiles(
			filepath.Join(cfg.Target.Path.Root, cfg.Target.Path.Base),
			filepath.Join(workingDir, cfg.Target.Path.Base),
		)

		if err != nil {
			return fmt.Errorf("failed to copy files for generation: %v", err)
		}
	}

	targetModule, err := util.GetModulePath()
	if err != nil {
		return err
	}

	ponzuVersion, err := util.GetPonzuVersion()
	if err != nil {
		return err
	}

	if err = writeTemplateFile(filepath.Join(workingDir, "go.mod"), "go.mod.tmpl", map[string]interface{}{
		"GoVersion":    strings.TrimPrefix(runtime.Version(), "go"),
		"PonzuVersion": ponzuVersion,
		"ModulePath":   targetModule,
		"WorkingDir":   cwd,
	}); err != nil {
		return fmt.Errorf("failed to write template: %v", err)
	}

	if err = writeTemplateFile(filepath.Join(workingDir, "main.go"), "main.go.tmpl", map[string]interface{}{
		"Arguments":   arguments,
		"ContentType": contentType,
		"ModulePath":  targetModule,
		"Config":      configs,
	}); err != nil {
		return fmt.Errorf("failed to write template: %v", err)
	}

	var output string
	output, err = exec("go get ./...", workingDir)
	if err != nil {
		errorMessage := err.Error()
		if output != "" {
			errorMessage = output
		}

		return fmt.Errorf("failed to generate content: %v", errorMessage)
	}

	output, err = exec("go run main.go", workingDir)
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
