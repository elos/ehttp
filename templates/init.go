package templates

import "path/filepath"

/*
	joinDir prepends a full path to a slice of relative paths

	i.e., base = "/root/here/", files  ["1.go", "2.go"]
	=> ["/root/here/1.go", "/root/here/2.go"]

	Useful for building templateSet paths
*/
func JoinDir(base string, files []string) []string {
	r := make([]string, len(files))
	for i := range files {
		r[i] = filepath.Join(base, files[i])
	}
	return r
}
