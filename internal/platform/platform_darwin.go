//go:build darwin

package platform

import (
	"syscall"
	"os"
)


// mntDontBrowse marks a volume that Finder does not show (MNT_DONTBROWSE):
// Recovery, Preboot, VM and the like.
const mntDontBrowse = 0x00100000

// hiddenMount tells whether p is a mount point Finder doesn't show. A
// folder that is not a mount point has its volume's flags (the hidden
// Data volume's, for one under /Volumes), so only mount points count.
func HiddenMount(p string) bool {
	var st syscall.Statfs_t
	if syscall.Statfs(p, &st) != nil || st.Flags&mntDontBrowse == 0 {
		return false
	}
	var on []byte
	for _, c := range st.Mntonname {
		if c == 0 {
			break
		}
		on = append(on, byte(c))
	}
	return string(on) == p
}

func Device(p string) (uint64, bool) {
	fi, err := os.Stat(p)
	if err != nil {
		return 0, false
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, false
	}
	return uint64(st.Dev), true
}
