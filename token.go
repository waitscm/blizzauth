package blizzauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	tokenURL      = "https://us.battle.net/oauth/token?grant_type=client_credentials"
	checkTokenURL = "https://us.battle.net/oauth/check_token?token=%v"
)

type Token struct {
	Request    string
	Expiration int64
}

func (t *Token) IsExpired() bool {
	return t.Expiration < time.Now().Unix()
}

// tokenResp is the response when requesting a token.
type tokenResp struct {
	AccessToken string `json:"access_token"`
	Type        string `json:"token_type"`
	Expires     int    `json:"expires_in"`
}

// tokenStatus is the requested token status
type tokenStatus struct {
	Expiration int64  `json:"exp"`
	ClientID   string `json:"client_id"`
}

var apiTokenLock map[string]sync.Mutex

var cachedToken sync.Map

func init() {
	apiTokenLock = map[string]sync.Mutex{}
}

func getCachedToken(apiName string) (t *Token) {
	ct, ok := cachedToken.Load(apiName)
	if ok {
		ctok, _ := ct.(Token)
		return &ctok
	}
	return nil
}

// GetToken based on the name of the API key filenames you want to get a token from Blizz for.
func GetToken(apiName string) (t Token, err error) {

	ct := getCachedToken(apiName)
	if ct != nil && !ct.IsExpired() {
		return *ct, nil
	}

	// no cached token get a new one
	lock, _ := apiTokenLock[apiName]
	lock.Lock()
	defer lock.Unlock()
	// we obtained the lock make sure that someone else didn't generate the token
	ct = getCachedToken(apiName)
	if ct != nil && !ct.IsExpired() {
		return *ct, nil
	}

	log.Println("GetToken Request new token")

	keys := newKeys(apiName)
	if keys == nil {
		log.Println("Unable to get keys")
		return t, errors.New("cannot load key for token request")
	}

	req, err := http.NewRequest("POST", tokenURL, nil)
	if err != nil {
		log.Fatalln("GetToken can't create new request", err)
	}
	req.SetBasicAuth(keys.id, keys.secret)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("GetToken request failed!!!", err)
		return t, errors.New("unable to request token")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("GetToken unable to read response body", err)
		return t, errors.New("can't open GetToken response")
	}

	var token tokenResp
	err = json.Unmarshal(body, &token)
	if err != nil {
		log.Println("GetToken response is not a TokenResp", err)
		return t, errors.New("update the code! response wrong format:" + string(body))
	}

	t.Request = token.AccessToken
	log.Println("GetToken new token:", t.Request)

	// now check when the token expires
	checkURI := fmt.Sprintf(checkTokenURL, t.Request)

	resp, err := http.Get(checkURI)

	if err != nil {
		log.Println("GetToken check token req failed:", checkURI, err)
		return t, errors.New("token expiration time not obtained")
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("GetAuctionRespFile bad resp:", checkURI, string(body), err)
		return t, errors.New("check token bad response")
	}
	var status tokenStatus
	err = json.Unmarshal(body, &status)
	if err != nil {
		log.Println("Unsupported Token Status", string(body))
		return t, errors.New("update the code Token Status resp body unknown")
	}

	t.Expiration = status.Expiration
	expiryTime := time.Unix(t.Expiration, 0)
	log.Println("GetToken token expires at ", expiryTime, "(", t.Expiration, ")")

	return t, nil
}
