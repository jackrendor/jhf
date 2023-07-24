package resources

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const useragent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36"

type Hashes struct {
	Hash   string
	Value  string
	Solved bool
}

var LoadedHashes []Hashes

func LoadHashes(hashes []string) {
	for _, hash := range hashes {
		LoadedHashes = append(LoadedHashes, Hashes{
			Hash:   hash,
			Solved: false,
		})
	}
}

func SolveHash(hash string, value string) bool {
	if value == "" {
		return false
	}
	for i, hashElem := range LoadedHashes {
		if hashElem.Hash == hash || hashElem.Hash == "*"+strings.ToUpper(hash) {
			LoadedHashes[i].Solved = true
			LoadedHashes[i].Value = value
			return true
		}
	}
	return false
}

func GetUnsolved() []string {
	var res []string
	for _, hash := range LoadedHashes {
		if !hash.Solved {
			res = append(res, hash.Hash)
		}
	}
	return res
}

// Decrypt MD5, SHA1, MySQL, NTLM, SHA256, MD5 Email, SHA256 Email, SHA512 hashes
func Hashes_Com(hashes []string) []Hashes {
	//build the request
	res := []Hashes{}
	reqGet, reqGetErr := http.NewRequest(http.MethodGet, "https://hashes.com/en/decrypt/hash", nil)
	if reqGetErr != nil {
		log.Println("[hashes.com]:", reqGetErr.Error())
		return res
	}
	reqGet.Header.Set("User-Agent", useragent)

	// perform the request
	respGet, respGetErr := http.DefaultClient.Do(reqGet)
	if reqGetErr != nil {
		log.Println("[hashes.com]:", respGetErr.Error())
		return res
	}

	cookies := respGet.Cookies()
	//Parse html so it can be analyzed later
	nodeGet, parseErr := html.Parse(respGet.Body)
	if parseErr != nil {
		log.Println("[hashes.com] [html.Parse]:", parseErr.Error())
		return res
	}

	//xpath to get csrf token
	csrfTokenNode, csrfTokenNodeErr := htmlquery.Query(nodeGet, "//input[@name='csrf_token']")
	if csrfTokenNodeErr != nil {
		log.Println("Couldn't get csrf_token:", csrfTokenNodeErr.Error())
		return res
	}
	if csrfTokenNode == nil {
		log.Println("Couldn't get csrf_token")
		return res
	}
	csrfTokenString := htmlquery.SelectAttr(csrfTokenNode, "value")

	//Setup post values with hash and csrf token
	data := url.Values{}
	data.Set("csrf_token", csrfTokenString)
	data.Set("hashes", strings.Join(hashes, "\r\n"))
	data.Set("submitted", "true")
	// Create request for lookup
	reqPost, reqPostErr := http.NewRequest(http.MethodPost, "https://hashes.com/en/decrypt/hash", strings.NewReader(data.Encode()))
	if reqPostErr != nil {
		log.Println("[hashes.com] [http.NewRequest]:", reqPostErr.Error())
		return res
	}

	// Adding headers
	reqPost.Header.Set("User-Agent", useragent)
	reqPost.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Parsing Form
	if err := reqPost.ParseForm(); err != nil {
		log.Println("[hashes.com] [reqPost.ParseForm]:", err.Error())
		return res
	}

	// Adding cookies
	for _, cookie := range cookies {
		reqPost.AddCookie(cookie)
	}
	reqPost.AddCookie(&http.Cookie{
		Name:  "csrf_cookie",
		Value: csrfTokenString,
	})

	// Perform lookup http post request
	respPost, respPostErr := http.DefaultClient.Do(reqPost)
	if respPostErr != nil {
		log.Println("[hashes.com] [http.Do]:", reqPostErr.Error())
		return res
	}

	// parsing html response
	nodePost, nodePostErr := html.Parse(respPost.Body)
	if nodePostErr != nil {
		log.Println("[hashes.com] [http.Do]:", nodePostErr.Error())
	}
	hashValueNodes, hashValueNodeErr := htmlquery.QueryAll(nodePost, "//div[@class='py-1']")
	if hashValueNodeErr != nil {
		log.Println("[hashes.com] [htmlquery.Query] [hashValueNode]:", hashValueNodeErr.Error())
		return res
	}
	for _, hashValueNode := range hashValueNodes {
		tmp := strings.Split(htmlquery.InnerText(hashValueNode), ":")
		res = append(res, Hashes{Hash: tmp[0], Value: tmp[1], Solved: true})
	}

	return res
}

func HashToolkit_Com(hash string) Hashes {
	res := Hashes{
		Hash:   hash,
		Solved: false,
	}

	// Preparing request
	reqGet, reqGetErr := http.NewRequest(http.MethodGet, "https://hashtoolkit.com/decrypt-hash/", nil)

	if reqGetErr != nil {
		log.Println("[hashtoolkit.com] [http.NewRequest]:", reqGetErr.Error())
		return res
	}

	// Add Query parameters
	getQuery := reqGet.URL.Query()
	getQuery.Add("hash", hash)
	reqGet.URL.RawQuery = getQuery.Encode()

	// Set headers
	reqGet.Header.Set("User-Agent", useragent)

	resGet, resGetErr := http.DefaultClient.Do(reqGet)
	if resGetErr != nil {
		log.Println("[hashtoolkit.com] [http.Do]", resGetErr.Error())
		return res
	}

	nodeGet, parseErr := html.Parse(resGet.Body)
	if parseErr != nil {
		log.Println("[hashtoolkit.com] [html.Parse]:", parseErr.Error())
		return res
	}

	if strings.Contains(htmlquery.InnerText(nodeGet), "No hashes found for") {
		return res
	}

	hashValueNode, hashValueNodeErr := htmlquery.Query(nodeGet, "//td[@class='res-text']")
	if hashValueNodeErr != nil {
		log.Println("[hashtoolkit.com] [htmlquery.Query]:", hashValueNodeErr.Error())
		return res
	}

	// Skip if no node found
	if hashValueNode == nil {
		return res
	}

	// Skip if result is empty
	hashValue := htmlquery.InnerText(hashValueNode)
	if hashValue == "" {
		return res
	}
	res.Solved = true
	res.Value = strings.TrimSpace(hashValue)

	return res
}

func Crack(hashes []string) []Hashes {
	LoadHashes(hashes)
	solved := Hashes_Com(GetUnsolved())
	for _, elem := range solved {
		SolveHash(elem.Hash, elem.Value)
	}
	for _, elem := range GetUnsolved() {
		hashValue := HashToolkit_Com(elem)
		SolveHash(hashValue.Hash, hashValue.Value)
	}
	return LoadedHashes
}
