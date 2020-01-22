package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	// InstMainURL ...
	InstMainURL = "https://www.instagram.com"
	// MyReg0 ...
	MyReg0 = `<link rel="preload" href="\/static\/bundles\/metro\/Consumer\.js\/[a-zA-Z0-9]+.js"`
	// MyReg1 ...
	MyReg1 = `<link rel="preload" href="\/static\/bundles\/es6\/Consumer\.js\/[a-zA-Z0-9]+.js"`

	// MyReg01 ...
	MyReg01 = `\/static\/bundles\/es6\/Consumer\.js\/[a-zA-Z0-9]+.js`
	// MyReg02 ...
	MyReg02 = `\/static\/bundles\/es6\/ConsumerLibCommons\.js\/[a-zA-Z0-9]+.js"`
	// MyReg002 ...
	MyReg002 = `e.instagramWebFBAppId='([0-9]{16})'`
	// MyReg2 query_hash ПОДПИСКИ
	MyReg2 = `n="([a-zA-z0-9]{32})"`
	// MyReg3 ...csrf_token
	MyReg3 = `"csrf_token":"([a-zA-Z0-9]{32})"`
)

var (
	client = &http.Client{
		Timeout: time.Second * 10,
	}
)

// UserJSON ...
type UserJSON struct {
	LoggingPageID string `json:"logging_page_id"` // logging_page_id ...
	Graphql       map[string]struct {
		ID string `json:"id"` // id ...
	} `json:"graphql"` // graphql ...
}

// User ...
type User struct {
	UserName        string
	Queryhashfollow string         // Queryhashfollow ...
	ID              string         //FolowID...
	CsrfToken       string         //csrf_token
	XIgAppID        string         //x-ig-app-id
	Cookies         []*http.Cookie //Cookies
}

// NewUser ...
func NewUser() *User {
	return &User{
		UserName: "beatbot_",
	}
}

// GetBodyByte ...
func GetBodyByte(u string) ([]byte, error) {
	resp, err := client.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, err
}

// MyReg ...
func MyReg(b []byte, reg string) [][]byte {
	re, _ := regexp.Compile(reg)
	u := re.FindSubmatch(b)
	return u
}

// GetBody2 ...
func GetBody2(u string, us *User) ([]byte, error) {
	url, err := url.Parse(InstMainURL + "/" + u)
	request := &http.Request{
		Method: "GET",
		URL:    url,
		Header: http.Header{
			"Origin":     []string{InstMainURL},
			"Referer":    []string{InstMainURL + "/" + u},
			"User-agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36"},
		},
	}
	var res1 *http.Response
	res1, err = client.Do(request)
	if err != nil {
		return nil, err
	}
	us.Cookies = res1.Cookies()
	// res1.Cookies := res1.Cookies()
	// for _, cookie := range cookies {
	// 	fmt.Println(cookie.Name)
	// 	fmt.Println(cookie.Value)
	// }

	defer res1.Body.Close()
	b, err := ioutil.ReadAll(res1.Body)
	if err != nil {
		return nil, err
	}
	return b, err
}

// GetBody3 ...
func GetBody3(u string, us *User) ([]byte, error) {
	url, err := url.Parse(u)
	request := &http.Request{
		Method: "GET",
		URL:    url,
		Header: http.Header{
			// ":authority":       []string{"www.instagram.com"},
			// ":method":          []string{"GET"},
			// ":scheme":          []string{"https"},
			"accept":           []string{"*/*"},
			"accept-encoding":  []string{"gzip, deflate, br"},
			"accept-language":  []string{"ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7"},
			"Referer":          []string{InstMainURL + "/" + u + "/following/"},
			"sec-fetch-mode":   []string{"cors"},
			"sec-fetch-site":   []string{"same-origin"},
			"User-agent":       []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36"},
			"x-csrftoken":      []string{us.CsrfToken},
			"x-ig-app-id":      []string{us.XIgAppID},
			"x-requested-with": []string{"XMLHttpRequest"},
		},
	}
	for _, cookie := range us.Cookies {
		fmt.Println(cookie.Name)
		fmt.Println(cookie.Value)
		request.AddCookie(cookie)
	}

	var res1 *http.Response
	res1, err = client.Do(request)
	if err != nil {
		return nil, err
	}
	defer res1.Body.Close()
	fmt.Println(res1.Header)
	fmt.Println(res1.StatusCode)
	if res1.StatusCode != 200 {

		return nil, err
	}
	b, err := ioutil.ReadAll(res1.Body)
	if err != nil {
		return nil, err
	}
	return b, err
}

// Marsh ...
func Marsh(b []byte) map[string]interface{} {
	var myjson map[string]interface{}
	if err := json.Unmarshal(b, &myjson); err != nil {
		panic(err)
	}
	return myjson
}

// Marsh2 ...
func Marsh2(b []byte, uj *UserJSON) {
	// var myjson map[string]interface{}

	if err := json.Unmarshal(b, &uj); err != nil {
		panic(err)
	}
}

// GetFollowList ...
func GetFollowList(u *User) {
	url1 := InstMainURL + "/graphql/query/?query_hash=" + u.Queryhashfollow + `&variables=%7B"id"%3A"` + u.ID + `"%2C"include_reel"%3Afalse%2C"fetch_mutual"%3Afalse%2C"first"%3A20%7D`
	// url2 := `https://www.instagram.com/graphql/query/?query_hash=d04b0a864b4b54837c0d870b0e77e076&variables=%7B%22id%22%3A%228655437174%22%2C%22include_reel%22%3Afalse%2C%22fetch_mutual%22%3Afalse%2C%22first%22%3A51%7D`
	fmt.Println(url1)
	// b, _ := GetBodyByte(url2)
	// fmt.Println(string(b))
	b1, _ := GetBody3(url1, u)
	fmt.Println(string(b1))

}
func main() {
	fmt.Println(Testing("Hi"))
	// testGet()
	user := NewUser()
	var UJson UserJSON

	body0, _ := GetBodyByte(InstMainURL + "/" + user.UserName + "/?__a=1")

	Marsh2(body0, &UJson)
	user.ID = UJson.Graphql["user"].ID
	fmt.Println(user.ID)

	body, _ := GetBody2(user.UserName, user)

	u1 := MyReg(body, MyReg01)
	u12 := MyReg(body, MyReg3)
	user.CsrfToken = string(u12[1])
	fmt.Println(user.CsrfToken)
	u21 := MyReg(body, MyReg02)
	fmt.Println(string(u21[0]))
	body02, _ := GetBodyByte(InstMainURL + string(u21[0]))
	u02 := MyReg(body02, MyReg002)
	if len(u02) >= 2 {
		user.XIgAppID = string(u02[1])
	}
	fmt.Println(user.XIgAppID)
	body2, _ := GetBodyByte(InstMainURL + string(u1[0]))
	u2 := MyReg(body2, MyReg2)

	if len(u2) >= 2 {
		user.Queryhashfollow = string(u2[1])
	}
	fmt.Println(user.Queryhashfollow)
	GetFollowList(user)
	// TESTRES(user)
}

// MyRegMy ...
func MyRegMy(re string, b string) [][]string {
	r := regexp.MustCompile(re)
	return r.FindAllStringSubmatch(b, -1)
}

// ReadToN ...
func ReadToN(sc string) {
	configLines := strings.Split(sc, "\n")
	for i := 0; i < len(configLines); i++ {
		fmt.Println(i)
		if configLines[i] != "" {
			res1 := MyRegMy(MyReg2, configLines[i])
			if len(res1) >= 1 {
				fmt.Println(res1)
				break
			}

		}
	}
}

// TESTRES ...
func TESTRES(us *User) {
	url, err := url.Parse("https://httpbin.org/get")
	request := &http.Request{
		Method: "GET",
		URL:    url,
		Header: http.Header{
			// ":authority":       []string{"www.instagram.com"},
			// ":method":          []string{"GET"},
			// ":scheme":          []string{"https"},
			"accept":           []string{"*/*"},
			"accept-encoding":  []string{"gzip, deflate, br"},
			"accept-language":  []string{"ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7"},
			"Referer":          []string{InstMainURL + "/" + "u" + "/following/"},
			"sec-fetch-mode":   []string{"cors"},
			"sec-fetch-site":   []string{"same-origin"},
			"User-agent":       []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36"},
			"x-csrftoken":      []string{us.CsrfToken},
			"x-ig-app-id":      []string{us.XIgAppID},
			"x-requested-with": []string{"XMLHttpRequest"},
		},
	}
	for _, cookie := range us.Cookies {
		fmt.Println(cookie.Name)
		fmt.Println(cookie.Value)
		request.AddCookie(cookie)
	}

	var res1 *http.Response
	res1, err = client.Do(request)
	if err != nil {
		return
	}
	defer res1.Body.Close()

	fmt.Println(res1.StatusCode)
	if res1.StatusCode != 200 {

		return
	}

	b, err := ioutil.ReadAll(res1.Body)
	if err != nil {
		return
	}
	fmt.Println(res1.Header)
	fmt.Println(string(b))
	return
}
