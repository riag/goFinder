//go:build windows

package finfo

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

const (
	faHidden = 0x2 // FILE_ATTRIBUTE_HIDDEN
	// ftEpochOffset 是 1601-01-01 到 1970-01-01 之间的 100 纳秒间隔数
	ftEpochOffset = 11644473600 * 10000000 // = 116444736000000000
)

func ftToTime(ft syscall.Filetime) time.Time {
	// FILETIME 是 LowDateTime|HighDateTime 拼成的 100-纳秒间隔数（从 1601-01-01 起）
	intervals := int64(ft.HighDateTime)<<32 + int64(ft.LowDateTime)
	intervals -= ftEpochOffset // 转成从 1970-01-01 起的 100-纳秒间隔
	return time.Unix(intervals/10000000, (intervals%10000000)*100)
}

func fill(in *Info, st *syscall.Win32FileAttributeData) {
	in.Created = ftToTime(st.CreationTime)
	in.Accessed = ftToTime(st.LastAccessTime)
	in.Hidden = st.FileAttributes&faHidden != 0
	// Owner/Group: Windows Lstat 不返回 SID，留空
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
	if sys, ok := st.Sys().(*syscall.Win32FileAttributeData); ok {
		fill(&in, sys)
	}
	in.Kind = Kind(in.Name, in.Dir, in.Link)
	in.Ext = strings.TrimPrefix(filepath.Ext(in.Name), ".")
	readExtra(&in)
	return in, nil
}
