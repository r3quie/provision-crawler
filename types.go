package crawler

import (
	"io/fs"
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
	Name        string
	Path        string
	Date        time.Time
	Podnikatel  bool
	Male        bool
	Animals     Animals
	Provisions []string
	Rozhodnuti  bool
}

type File struct {
	File fs.DirEntry
	Path string
}
