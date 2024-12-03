package provisioncrawler

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func readDocx(s string) (string, error) {
	r, err := zip.OpenReader(s)
	if err != nil {
		inf, er := os.Stat(s)
		if er != nil {
			return "", fmt.Errorf("readDocx: error in os.Stat while opening %s: %s", s, er.Error())
		}
		if inf.Size() == 0 {
			return "", nil
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

func IsMoreTimesIn(str string, substr string, substrtwo string) bool {
	return strings.Count(str, substr) > strings.Count(str, substrtwo)
}

func initRozh(file File) (Rozh, error) {
	text, err := readDocx(filepath.Join(file.Path, file.File.Name()))
	if err != nil {
		return Rozh{}, err
	}
	inf, _ := file.File.Info()
	animals := Animals{
		Kone:     strings.Contains(text, "kůň") || strings.Contains(text, "koně") || strings.Contains(text, "klisna"),
		OvceKozy: strings.Contains(text, "ovce") || strings.Contains(text, "kozy"),
		Prasata:  strings.Contains(text, "prase") || strings.Contains(text, "prasečí") || strings.Contains(text, "prasat"),
		Turi:     strings.Contains(text, "tur") || strings.Contains(text, "tura") || strings.Contains(text, "turů"),
	}
	return Rozh{
		Name:        file.File.Name(),
		Path:        file.Path,
		Date:        inf.ModTime(),
		Provistions: extractExpressions(text),
		Rozhodnuti:  IsMoreTimesIn(text, "R O Z H O D N U T I", "P R I K A Z"),
		Male:        IsMoreTimesIn(text, "obviněný", "obviněná"),
		Podnikatel:  IsMoreTimesIn(text, "IČ", "RČ"),
		Animals:     animals,
	}, nil
}

func GetRozhs(dir string) []Rozh {
	files := walk(dir)

	var rozhs []Rozh

	for _, file := range files {
		if !strings.HasSuffix(file.File.Name(), ".docx") {
			continue
		}
		rozh, err := initRozh(file)
		if err != nil {
			continue
		}
		rozhs = append(rozhs, rozh)

	}
	return rozhs
}

func WriteToJson(path string) {
	rozhs := GetRozhs(path)
	o, err := json.Marshal(rozhs)
	if err != nil {
		panic(err)
	}
	os.WriteFile("output.json", o, 0644)
}
