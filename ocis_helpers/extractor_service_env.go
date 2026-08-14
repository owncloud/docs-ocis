package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

// Render the envvars that origin in services 
func RenderServices() {

	ReadEnv()
	fmt.Printf(Green + "Generate envvar .adoc tables and yaml examples from services and the %s file\n\n" + Reset, yamlServiceSource)
	doTemplates()
	RemoveOutputDir()

}

// finally do the templating stuff
func doTemplates() {

	/*
	  renders:
	  1. adoc files for envvars in services
	  2. yaml files for envvars in services
	  3. the env_var.yaml file
	*/
	var targets = map[string]string{
		"templates/adoc-generator.go.tmpl":                 output_dir + "adoc/adoc-generator.go",
		"templates/example-yaml-config-generator.go.tmpl":  output_dir + "exampleconfig/example-yaml-config-generator.go",
		"templates/envar-db-table.go.tmpl":                 output_dir + "env/envvar-db-table.go",
	}

	paths, err := filepath.Glob(Env.ocis_dir + "services/*/pkg/config/defaults/defaultconfig.go")
	if err != nil {
		log.Fatal(err)
	}
	replacer := strings.NewReplacer(
		ocis_dir, "github.com/owncloud/ocis/v2/",
		"/defaultconfig.go", "",
	)
	for i := range paths {
		paths[i] = replacer.Replace(paths[i])
	}

	for template, output := range targets {
		generateIntermediateCode(template, output, paths)
		runIntermediateCode(output)
	}

}

func generateIntermediateCode(templatePath string, intermediateCodePath string, paths []string) {

	var err error

	content, err := os.ReadFile(templatePath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Generating intermediate go code for " + intermediateCodePath + " using template " + templatePath)
	tpl := template.Must(template.New("").Parse(string(content)))

	err = os.MkdirAll(path.Dir(intermediateCodePath), Env.folder_mode)
	if err != nil {
		log.Fatal(err)
	}
	runner, err := os.Create(intermediateCodePath)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Chmod(intermediateCodePath, Env.file_mode)
	if err != nil {
		log.Fatal(err)
	}
	err = tpl.Execute(runner, paths)
	if err != nil {
		log.Fatal(err)
	}
}

func runIntermediateCode(intermediateCodePath string) {
	fmt.Println("Running intermediate go code for " + intermediateCodePath)
	defaultConfigPath := "/etc/ocis"
	defaultDataPath := "/var/lib/ocis"
	os.Setenv("OCIS_BASE_DATA_PATH", defaultDataPath)
	os.Setenv("OCIS_CONFIG_DIR", defaultConfigPath)

	// Set AUTOMEMLIMIT_EXPERIMENT=system on non-Linux systems to avoid cgroups errors
	if runtime.GOOS != "linux" {
		os.Setenv("AUTOMEMLIMIT_EXPERIMENT", "system")
	}

	// copy the required go file to the template code
	envFileDir := path.Dir(intermediateCodePath)
	envFileDst := path.Join(envFileDir, "read_env_file.go")
	envFileSrc, err := os.ReadFile("read_env_file.go")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(envFileDst, envFileSrc, Env.file_mode)
	if err != nil {
		log.Fatal(err)
	}

	out, err := exec.Command("go", "run", intermediateCodePath, envFileDst).CombinedOutput()
	if err != nil {
		log.Fatal(string(out), err)
	}
	fmt.Println(string(out))
}
