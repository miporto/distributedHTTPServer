package httpparser

import "testing"

type testpair struct {
	val string
	exp string
}

func TestGetMethod(t *testing.T) {
	tests := []testpair{
		{"GET /test/example/1 HTTP/1.1\nHost: example.net\n\n", "GET"},
		{"DELETE /test/example/1 HTTP/1.1\nHost: localhost:8080\n\n", "DELETE"},
		{"POST /test/example/1 HTTP/1.1\nHost: localhost:8080\nContent-Type: application/json\n\n", "POST"},
		{"PUT /test/example/1 HTTP/1.1\nHost: localhost:8080\nContent-Type: application/json\n\n", "PUT"},
	}
	for _, pair := range tests {
		res := GetMethod(pair.val)
		if res != pair.exp {
			t.Error("For ", pair.val, " expected ", pair.exp, "got ", res)
		}
	}
}

func TestGetURI(t *testing.T) {
	tests := []testpair{
		{"GET /test/example/1 HTTP/1.1\nHost: example.net\n\n", "/test/example/1"},
		{"DELETE /inv4lid/example/1 HTTP/1.1\nHost: localhost:8080\n\n", ""},
		{"POST /test/newexample/987 HTTP/1.1\nHost: localhost:8080\nContent-Type: application/json\n\n",
			"/test/newexample/987"},
		{"PUT /test/example/s HTTP/1.1\nHost: localhost:8080\nContent-Type: application/json\n\n", ""},
	}
	for _, pair := range tests {
		res := GetURI(pair.val)
		if res != pair.exp {
			t.Error("For ", pair.val, " expected ", pair.exp, "got ", res)
		}
	}
}

func TestGetContentLength(t *testing.T) {
	test := `POST /bin/login HTTP/1.1
	Host: 127.0.0.1:8000
	Content-Type: application/x-www-form-urlencoded
	Content-Length: 37
	   
	User=Peter+Lee&pw=123456&action=login`

	res := GetContentLength(test)
	if res != 37 {
		t.Error("Expected ", 37, "got ", res)
	}
}
