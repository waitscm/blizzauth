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
	ErrNoKeys     = "no keys found"
)

type Auth struct {
	request    string
	expiration int64
	clientName string
	keys       keys
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

var apiTokenLock sync.Mutex

var cachedToken sync.Map

func getCachedAuth(apiName string) (t *Auth) {
	ct, ok := cachedToken.Load(apiName)
	if ok {
		ctok, _ := ct.(*Auth)
		return ctok
	}
	return nil
}

func setCachedAuth(t *Auth) {
	cachedToken.Store(t.clientName, t)
}

// GetAuth based on the name of the API key filenames in the .blizzard directory that you want to get authorization for.
func GetAuth(apiName string) (t *Auth, err error) {

	// no cached token get a new one
	apiTokenLock.Lock()
	defer apiTokenLock.Unlock()
	// we obtained the lock make sure that someone else didn't generate the token
	ct := getCachedAuth(apiName)
	if ct != nil {
		return ct, nil
	}

	// create a new one
	var auth Auth
	// cache the auth
	auth.clientName = apiName
	keys := newKeys(auth.clientName)
	if keys == nil {
		log.Println("Unable to get keys. Check that you have keys <api_name>.id and <api_name>.secret in $HOME/.blizzard")
		return nil, errors.New(ErrNoKeys)
	}
	auth.keys = *keys
	setCachedAuth(&auth)

	return &auth, nil
}

//clearToken if there is an error
func (a *Auth) clearToken() {
	a.request = ""
}
func (a *Auth) isExpired() bool {
	return a.expiration < time.Now().Unix()
}

func (a *Auth) needNewToken() bool {

	if a.request == "" {
		return true
	}

	if a.isExpired() {
		return true
	}

	return false
}

// GetAccessToken get an access token or request a new one if current is expired
func (a *Auth) GetAccessToken() (string, error) {
	apiTokenLock.Lock()
	defer apiTokenLock.Unlock()
	if !a.needNewToken() {
		return a.request, nil
	}
	a.clearToken()
	log.Println("GetAccessToken Request new token")

	req, err := http.NewRequest("POST", tokenURL, nil)
	if err != nil {
		log.Fatalln("GetToken can't create new request", err)
	}
	req.SetBasicAuth(a.keys.id, a.keys.secret)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("GetToken request failed!!!", err)
		return "", errors.New("unable to request token")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("GetToken unable to read response body", err)
		return "", errors.New("can't open GetToken response")
	}

	var token tokenResp
	err = json.Unmarshal(body, &token)
	if err != nil {
		log.Println("GetToken response is not a TokenResp", err)
		return "", errors.New("update the code! response wrong format:" + string(body))
	}

	a.request = token.AccessToken
	log.Println("GetToken new token:", a.request)

	// now check when the token expires
	checkURI := fmt.Sprintf(checkTokenURL, a.request)

	resp, err := http.Get(checkURI)

	if err != nil {
		log.Println("GetToken check token req failed:", checkURI, err)
		return a.request, errors.New("token expiration time not obtained")
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("GetAuctionRespFile bad resp:", checkURI, string(body), err)
		return a.request, errors.New("check token bad response")
	}
	var status tokenStatus
	err = json.Unmarshal(body, &status)
	if err != nil {
		log.Println("Unsupported Auth Status", string(body))
		return a.request, errors.New("update the code Auth Status resp body unknown")
	}

	a.expiration = status.Expiration
	expiryTime := time.Unix(a.expiration, 0)
	log.Println("GetToken token expires at ", expiryTime, "(", a.expiration, ")")
	return a.request, nil
}
