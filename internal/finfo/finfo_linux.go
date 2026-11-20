//go:build linux
package finfo


import (
	"os"
	"syscall"
	"strings"
	"path/filepath"
	"io/fs"
)

func fill(in *Info, st *syscall.Stat_t) {
	in.Owner = userName(st.Uid)
	in.Group = groupName(uint32(st.Gid))
}

// Stat reads path's attributes, of a link itself rather than its target.
func Stat(path string) (Info, error) {
	st, err := os.Lstat(path)
	if err != nil {
		return Info{}, err
	}
	in := Info{
		Name: filepath.Base(path), Path: path, Dir: st.IsDir(),
		Size: st.Size(), Modified: st.ModTime(), Mode: st.Mode(),
	}
	if st.Mode()&fs.ModeSymlink != 0 {
		in.Link = true
		in.Target, _ = os.Readlink(path)
		if t, err := os.Stat(path); err == nil && t.IsDir() {
			in.Dir = true
		}
	}
	if sys, ok := st.Sys().(*syscall.Stat_t); ok {
		fill(&in, sys)
	}
	in.Kind = Kind(in.Name, in.Dir, in.Link)
	in.Ext = strings.TrimPrefix(filepath.Ext(in.Name), ".")
	readExtra(&in)
	return in, nil
}
