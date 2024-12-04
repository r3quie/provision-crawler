package crawler

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"
)

type Animals struct {
	Kone     bool
	OvceKozy bool
	Prasata  bool
	Turi     bool
	Jine     bool
}

type Rozh struct {
	Name       string
	Path       string
	Date       time.Time
	Podnikatel bool
	Male       bool
	Animals    Animals
	Provisions []string
	Rozhodnuti bool
}

type File struct {
	File fs.DirEntry
	Path string
}

type Rozhs []Rozh

// Returns a string representation of a Found struct (struct{subdir, filename, modtime})
func (f Rozh) String() string {
	if len(f.Name) > 63 {
		return fmt.Sprintf("%-63s %s", f.Name[:58]+"...", f.Date.Format("02.01.2006"))
	}
	return fmt.Sprintf("%-63s %s", f.Name, f.Date.Format("02.01.2006"))
}

// Sorts the slice of Found structs ([]struct{subdir, filename, modtime}) by modtime
func (f Rozhs) Sort() {
	sort.Slice(f, func(i, j int) bool {
		return f[i].Date.After(f[j].Date)
	})
}

// Returns a string representation of FoundSlice ([]struct{subdir, filename, modtime})
func (f Rozhs) WidgetText() string {
	var text strings.Builder
	for _, x := range f {
		text.WriteString(x.String() + "\n")
	}
	return text.String()
}

func (f Rozhs) Options() []string {
	s := make([]string, len(f))
	for i, x := range f {
		s[i] = x.Name
	}
	return s
}
