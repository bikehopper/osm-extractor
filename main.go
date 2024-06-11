package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type Extracts struct {
	Extracts []Extract `json:"Extracts"`
}

type Extract struct {
	Output      string  `json:"output"`
	Directory   string  `json:"directory"`
	Description string  `json:"description"`
	Polygon     polygon `json:"polygon"`
}

type polygon struct {
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
}

func osmium_extract() {
	cmd := exec.Command("osmium", "extract", "-d", "./volumes/output", "-c", "./config.json", "./volumes/input/latest.osm.pbf")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Print("error getting osmium standerr pipe")
	}
	err = cmd.Start()
	if err != nil {
		log.Printf("Osmium extract cmd.Start() failed with %s\n", err)
		log.Fatal(err)
	}

	stderrin := bufio.NewScanner(stderr)
	for stderrin.Scan() {
		fmt.Println(stderrin.Text())
	}
	err = cmd.Wait()
	if err != nil {
		log.Printf("Osmium failed: %s\n", err)
		log.Fatal(err)
	}
}

func move_extracts(src string, dest string) {
	err := os.MkdirAll(filepath.Dir(dest), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	err = os.Rename(src, dest)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	osmium_extract()

	jsonFile, err := os.Open("config.json")

	if err != nil {
		log.Fatalf("Failed to read config JSON: %s\n", err)
	}

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var extracts Extracts

	err = json.Unmarshal(byteValue, &extracts)
	if err != nil {
		log.Fatalf("Error nnmarshaling JSON: %s\n", err)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(extracts.Extracts); i++ {
		sourcePath := filepath.Join(currentDir, extracts.Extracts[i].Output)
		destDir := filepath.Join(currentDir, extracts.Extracts[i].Directory)
		destPath := filepath.Join(destDir, extracts.Extracts[i].Output)

		move_extracts(sourcePath, destPath)

		fmt.Println("File moved successfully.")
	}
}
