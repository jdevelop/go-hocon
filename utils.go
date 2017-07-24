package hocon

import "strings"

func stripStringQuotas(src string) string {
	res := src
	if res[0] == '"' || res[0] == '\'' {
		res = res[1: len(res)-1]
	}
	return res
}

func splitPath(path string) []string {
	return strings.Split(path, ".")
}

func pathPrefix(path []string) ([]string, string) {
	length := len(path)
	if length == 1 {
		return []string{}, (path)[0]
	} else {
		return path[:len(path)-1], (path)[length-1]
	}
}
