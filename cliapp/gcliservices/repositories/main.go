package repositories

import "strings"

func reduceRepoURL(path string) string {
	index := strings.Index(path, "://")
	if index != -1 {
		path = path[index+3:]
	}
	if strings.HasSuffix(path, ".git") {
		path = path[:len(path)-4]
	}
	return path
}
