package httpparser

import "regexp"

func find(pattern string, s string) string {
	r, _ := regexp.Compile(pattern)
	return r.FindString(s)
}

func findSubmatch(pattern string, s string) string {
	r := regexp.MustCompile(pattern)
	return r.FindStringSubmatch(s)[0]
}
func GetMethod(s string) string {
	return find(`^[A-Z]+\b`, s)
}

func GetURI(s string) string {
	return find(`/[a-z]+/[a-z]+/[0-9]+`, s)
}

func GetHost(s string) string {
	return find(`^Host:\s([a-z]+:*[0-9]*)`, s)
}

func GetContentLength(s string) string {
	return find("", s)
}
