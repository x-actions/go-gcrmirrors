package github

import "strings"

func parseTxtRelativePath(str string) string {
	// str like "https://github.com/x-mirrors/gcr.io/blob/main/k8s.gcr.io/k8s.txt"
	// return "k8s.gcr.io/k8s.txt"
	temp := strings.Split(str, "/")
	return strings.Join(temp[len(temp)-2:], "/")
}
