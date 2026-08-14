package main

// note that this file will be copied to all subfolders where templates are stored automatically
// the filename is therefore crucial !

import (
//	"fmt"
	"os"
	"strings"
	"strconv"
)

// the struct that provides all envvars defined in .env
// note that the values of these variables are defined in 'main.go'
type ESt struct {
    isVerbose    bool
    isRemove     bool
    services_dir string
    output_dir   string
    ocis_dir     string
    folder_mode  os.FileMode
    file_mode    os.FileMode
}

// the variable that provides the struct
var Env ESt

// define global colors for fmt.print
const Reset   = "\033[0m"
const Red     = "\033[31m"
const Green   = "\033[32m"
const Yellow  = "\033[33m"
const Blue    = "\033[34m"
const Magenta = "\033[35m"
const Cyan    = "\033[36m"
const Gray    = "\033[37m"
const White   = "\033[97m"


// provide common paths and file names
const persistent_files    = "persistent_files/"
const extened_files       = "extended/"
const adoc_files          = "adoc/"
const yaml_files          = "yaml/"

// there is another folder: "env_var_deltas/"
// the folder and its content is managed by a Python program (changed_envvars.py)

const yamlServiceSource   = "env_vars.yaml"
const yamlExtendedSource  = "extended_vars.yaml"
const backupFile          = "extended_vars.yaml.do_not_delete"

const adoc_configvars     = "_configvars.adoc"
const adoc_deprecation    = "_deprecation.adoc"
const adoc_global         = "global_configvars.adoc"
const yaml_example        = "-config-example.yaml"

// Read the .env file and polulate variables
func ReadEnv() {

	// get informations from envvar file to populate variables
	// for ease of handling, all of the same look identical in different files
	x            := ""
	y            := uint64(0)
	if envData, err := os.ReadFile(".env"); err == nil {
		for _, line := range strings.Split(string(envData), "\n") {
			if strings.HasPrefix(line, "IS_VERBOSE=") {
				Env.isVerbose = strings.TrimPrefix(line, "IS_VERBOSE=") == "true"
				_ = Env.isVerbose
				continue
			}
			if strings.HasPrefix(line, "IS_REMOVE=") {
				Env.isRemove = strings.TrimPrefix(line, "IS_REMOVE=") == "true"
				_ = Env.isRemove
				continue
			}
			if strings.HasPrefix(line, "SERVICES_DIR=") {
				Env.services_dir = strings.TrimPrefix(line, "SERVICES_DIR=")
				_ = Env.services_dir
				continue
			}
			if strings.HasPrefix(line, "OUTPUT_DIR=") {
				Env.output_dir = strings.TrimPrefix(line, "OUTPUT_DIR=")
				_ = Env.output_dir
				continue
			}
			if strings.HasPrefix(line, "OCIS_DIR=") {
				Env.ocis_dir = strings.TrimPrefix(line, "OCIS_DIR=")
				_ = Env.ocis_dir
				continue
			}
			if strings.HasPrefix(line, "FOLDER_MOD=") {
				x = strings.TrimPrefix(line, "FOLDER_MOD=")
				y, _ = strconv.ParseUint(x, 8, 32)
				Env.folder_mode = os.FileMode(y)
				_ = Env.folder_mode
				continue
			}
			if strings.HasPrefix(line, "FILE_MOD=") {
				x = strings.TrimPrefix(line, "FILE_MOD=")
				y, _ = strconv.ParseUint(x, 8, 32)
				Env.file_mode = os.FileMode(y)
				_ = Env.file_mode
			}
		}
	}
}

