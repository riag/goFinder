//go:build windows

package platform


func HiddenMount(p string) bool {
	return false
}

func Device(p string) (uint64, bool) {
	return 0, false
}
