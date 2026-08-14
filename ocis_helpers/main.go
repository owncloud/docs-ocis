package main

import (
	"fmt"
	"flag"
    "path/filepath"
	"io"
	"log"
	"os"
	"strconv"
)

// define the directories used
// note that services and output have a counterpart in .gitignore
// note that other required constants are defined in 'read_env_file.go'
const services_dir = "./services/"
const output_dir   = "./output/"
const ocis_dir     = "../../ocis/"
const folder_mode  = "0774"
const file_mode    = "0664"

func main() {

	var err error

	verboseP := flag.Bool("v", false, "Enable verbosity")
	removeP := flag.Bool("r", false, "Do not remove output directory when finished")
	helpP := flag.Bool("h", false, "Print help message")

	flag.Usage = func() {
		// if an unknown flag has been provided
		printUsage()
		os.Exit(1)
	}

	flag.Parse()

	isVerbose := *verboseP
	isRemove := !*removeP
	isHelp := *helpP

	if isHelp {
		printUsage()
	}

	// check if the ocis repo has been cloned locally
	// no colors, added in a later step
	_, err = os.Stat(ocis_dir)
	if err != nil {
		fmt.Printf("The required ocis repo cant be found locally at: %s\n", ocis_dir)
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) > 0 {
		switch args[0] {
		case "service":
			prepareDirectories(isVerbose, isRemove)
			RenderServices()
		case "rogue":
			prepareDirectories(isVerbose, isRemove)
			GetRogueEnvs()
		case "extended":
			prepareDirectories(isVerbose, isRemove)
			RenderRogueEnvs()
		case "all":
			prepareDirectories(isVerbose, isRemove)
			RenderServices()
			GetRogueEnvs()
			RenderRogueEnvs()
		case "cleanup":
			prepareDirectories(isVerbose, isRemove)
			cleanupServiceDir()
		default:
			fmt.Fprintf(os.Stderr, "Unknown task: %s\n", args[0])
			printUsage()
			os.Exit(1)
		}
	} else {
		printUsage()
	}
}

// isVerbose and isRemove are provided by flags defined above
func prepareDirectories(isVerbose bool, isRemove bool) {

	// pre-create basic directories, subdirectories will be added on the fly
	var err error

	// with the checks here, we do not need checks in follow up code
	// get permissions from string
	x := uint64(0)
	x, err = strconv.ParseUint(folder_mode, 8, 32)
	if err != nil {
		log.Fatal(err)
	}
	fm := os.FileMode(x)

	// get permissions from string - for testing purposes only if the string is valid
	y := uint64(0)
	y, err = strconv.ParseUint(file_mode, 8, 32)
	if err != nil {
		log.Fatal(err)
	}
	_ = y

	// write functional vars to a file: .env
	// required to be read by go and the templates
	// the file will be overwritten if exists
	envContent := fmt.Sprintf("IS_VERBOSE=%t\nIS_REMOVE=%t\nSERVICES_DIR=%s\nOUTPUT_DIR=%s\nOCIS_DIR=%s\nFOLDER_MOD=%s\nFILE_MOD=%s\n",
		isVerbose, isRemove, services_dir, output_dir, ocis_dir, folder_mode, file_mode)
	err = os.WriteFile(".env", []byte(envContent), fm)
	if err != nil {
		log.Fatal(err)
	}

	// read the written and provide the variables as in other go files for consistent usage
	ReadEnv()

	// create output folder if not exists
	err = os.MkdirAll(output_dir, Env.folder_mode)
	if err != nil {
		log.Fatal(err)
	}

	// create services folder if not exists
	err = os.MkdirAll(services_dir, Env.folder_mode)
	if err != nil {
		log.Fatal(err)
	}
}


func cleanupServiceDir() {

	var err error
	var folder string

	fmt.Printf(Magenta + "Remove the content of non persistent subfolders in %s \n", services_dir + Reset)

    // adoc
    folder = services_dir + adoc_files + "*"
    err = removeGlob(folder)
    if err != nil {
        log.Fatalf("Error removing files: %+v", err)
    }

    // extended
    folder = services_dir + extened_files + "*"
    err = removeGlob(folder)
    if err != nil {
        log.Fatalf("Error removing files: %+v", err)
    }

    // yaml
    folder = services_dir + yaml_files + "*"
    err = removeGlob(folder)
    if err != nil {
        log.Fatalf("Error removing files: %+v", err)
    }
}

// remove the content of a folder if exists
func removeGlob(path string) (err error) {

	contents, err := filepath.Glob(path)
	if err != nil {
		return
	}

	for _, item := range contents {
		err = os.RemoveAll(item)
		if err != nil {
			return
		}
	}
	return
}

func RemoveOutputDir() {

	var err error

	// if removal should be omitted
	if Env.isRemove {
		// first remove the output directory
		fmt.Println(Magenta + "Cleaning up (remove output directory) \n" + Reset)
		_, err = os.Stat(Env.output_dir)
		if err != nil {
			err = os.RemoveAll(Env.output_dir)
			if err != nil {
				fmt.Println(err)
			}
		}
		// then remove the .env file
		_, err = os.Stat(".env")
		if err != nil {
			err = os.Remove(".env")
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		fmt.Println(Magenta + "No cleanup (output directory is kept) \n" + Reset)
	}
}


func CopyFile(src, dst string) error {

	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create the destination file
	destinationFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destinationFile.Close()

	// Copy the content
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// set the correct permissions
	err = os.Chmod(dst, Env.file_mode)
	if err != nil {
		log.Fatal(err)
	}

	// Flush file metadata to disk
	err = destinationFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync destination file: %w", err)
	}

	return nil
}


func printUsage() {

	fmt.Printf("Usage: go run . <flags> <task>\n\n")
	fmt.Printf("Available command line flags:\n")
	flag.PrintDefaults()
	fmt.Printf("\nAvailable tasks:\n")
	fmt.Printf("  service:   generate service envvar tables\n")
	fmt.Printf("  rogue:     update/create the extended_vars.yaml file\n")
	fmt.Printf("  extended:  generate extended_configvars table\n")
	fmt.Printf("  all:       run all above at once\n")
	fmt.Printf("  cleanup:   cleanup folders in %s\n", services_dir)

}
