package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func printArray(arr []string) {
	for i := 0; i < len(arr); i++ {
		fmt.Println(arr[i])
	}
}

func readDocx(s string) (string, error) {
	r, err := zip.OpenReader(s)
	if err != nil {
		inf, er := os.Stat(s)
		if er != nil {
			return "", fmt.Errorf("readDocx: error in os.Stat while opening %s: %s", s, er.Error())
		}
		if inf.Size() == 0 {
			return "", fmt.Errorf("readDocx: 0 size file")
		}
		return "", fmt.Errorf("readDocx: error in zip.OpenReader while opening %s: %s", s, err.Error())
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			// fmt.Println("found")
			doc, errdoc := f.Open()
			if errdoc != nil {
				return "", errdoc
			}
			buf := new(strings.Builder)
			_, errd := io.Copy(buf, doc)
			if errd != nil {
				return "", errd
			}
			reg, _ := regexp.Compile(`\<.*?\>`)
			return reg.ReplaceAllString(buf.String(), ""), nil

		}
	}
	return "", nil
}

func extractExpressions(a string) []string {
	re := regexp.MustCompile(`§\s\d+(\sodst\.\s\d+)?(\spísm\.\s[a-z]\))?(\sbod\s\d+)?`)
	return re.FindAllString(a, -1)
}

func walk(d string) []File {
	var files []File
	dir, _ := os.ReadDir(d)
	for _, entry := range dir {
		if entry.IsDir() {
			files = append(files, walk(d+"/"+entry.Name())...)
		} else {
			files = append(files, File{entry, d})
		}
	}
	return files
}

func GetProvisons(dir string) [][]string {
	files := walk(dir)
	var provs [][]string
	for _, file := range files {
		if !strings.HasSuffix(file.File.Name(), ".docx") {
			continue
		}
		text, err := readDocx(filepath.Join(file.Path, file.File.Name()))
		if err != nil {
			if err.Error() == "readDocx: 0 size file" {
				continue
			}
			panic(err)
		}
		provs = append(provs, extractExpressions(text))
	}
	return provs
}

func main() {
	var input string
	fmt.Println("Enter the path to the directory with the files:")
	fmt.Scanln(&input)
	provs := GetProvisons(input)
	for _, prov := range provs {
		printArray(prov)
	}
}
