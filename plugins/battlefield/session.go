package battlefield

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/StalkR/goircbot/lib/transport"
)

// A Session represents a session to EA and Battlefield companion.
type Session struct {
	email    string
	password string

	fid         string
	executionID string
	jsessionID  string
	sid         string

	code             string
	gatewaySessionID string
}

// NewSession creates a new Session but does not perform any requests.
func NewSession(email, password string) *Session {
	return &Session{email: email, password: password}
}

// Login logs in on EA then to Battlefield companion.
// It obtains the gateway session ID needed to talk to the companion API.
func (s *Session) Login() error {
	if err := s.connectAuthInit(); err != nil {
		return err
	}
	if err := s.webLoginInit(); err != nil {
		return err
	}
	if err := s.webLogin(); err != nil {
		return err
	}
	if err := s.connectAuth(); err != nil {
		return err
	}
	if err := s.connectAuthCode(); err != nil {
		return err
	}
	return s.loginFromAuthCode()
}

func (s *Session) connectAuthInit() error {
	c, err := transport.Client("https://accounts.ea.com")
	if err != nil {
		return err
	}
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse // do not follow redirects
	}
	query := url.Values{}
	query.Set("response_type", "code")
	query.Set("client_id", "Battlefield-CoreWeb")
	resp, err := c.Get(fmt.Sprintf("https://accounts.ea.com/connect/auth?%v", query.Encode()))
	if err != nil && err != http.ErrUseLastResponse {
		return err
	}
	defer resp.Body.Close()
	// https://signin.ea.com/p/web/login?fid=xxxxxxxx
	u, err := resp.Location()
	if err != nil {
		return err
	}
	target := u.Host + u.Path
	switch target {
	case "signin.ea.com/p/web/login", "signin.ea.com/p/web2/login":
	default:
		return fmt.Errorf("battlefield: expected redirect to signin.ea.com/p/web/login, got %s%s", u.Host, u.Path)
	}
	fid := u.Query().Get("fid")
	if fid == "" {
		return fmt.Errorf("battlefield: EA login has no fid")
	}
	s.fid = fid
	return nil
}

func (s *Session) webLoginInit() error {
	c, err := transport.Client("https://signin.ea.com")
	if err != nil {
		return err
	}
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse // do not follow redirects
	}
	query := url.Values{}
	query.Set("fid", s.fid)
	resp, err := c.Get(fmt.Sprintf("https://signin.ea.com/p/web/login?%v", query.Encode()))
	if err != nil && err != http.ErrUseLastResponse {
		return err
	}
	defer resp.Body.Close()
	// https://signin.ea.com/p/web/login?execution=xxxxx&initref=https%3A%2F%2Faccounts.ea.com%3A443%2Fconnect%2Fauth%3Fresponse_type%3Dcode%26client_id%3DBattlefield-CoreWeb
	u, err := resp.Location()
	if err != nil {
		return err
	}
	if u.Host+u.Path != "signin.ea.com/p/web/login" {
		return fmt.Errorf("battlefield: expected redirect to signin.ea.com/p/web/login, got %s%s", u.Host, u.Path)
	}
	executionID := u.Query().Get("execution")
	if executionID == "" {
		return fmt.Errorf("battlefield: EA login has no execution id")
	}
	jsessionID := func() string {
		for _, c := range resp.Cookies() {
			if c.Name == "JSESSIONID" {
				return c.Value
			}
		}
		return ""
	}()
	if jsessionID == "" {
		return fmt.Errorf("battlefield: EA login has no jsessionid cookie")
	}
	s.executionID = executionID
	s.jsessionID = jsessionID
	return nil
}

func (s *Session) webLogin() error {
	c, err := transport.Client("https://signin.ea.com")
	if err != nil {
		return err
	}
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse // do not follow redirects
	}
	query := url.Values{}
	query.Set("execution", s.executionID)
	data := url.Values{}
	data.Set("email", s.email)
	data.Set("password", s.password)
	data.Set("_rememberMe", "on")
	data.Set("rememberMe", "on")
	data.Set("_eventId", "submit")
	data.Set("gCaptchaResponse", "")
	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://signin.ea.com/p/web/login?%v", query.Encode()),
		strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{
		Name:   "JSESSIONID",
		Value:  s.jsessionID,
		Path:   "/",
		Domain: "accounts.ea.com",
		Secure: true,
	})
	resp, err := c.Do(req)
	if err != nil && err != http.ErrUseLastResponse {
		return err
	}
	defer resp.Body.Close()
	// https://accounts.ea.com:443/connect/auth?response_type=code&client_id=Battlefield-CoreWeb&fid=xxxxx
	u, err := resp.Location()
	if err != nil {
		return err
	}
	if u.Host+u.Path != "accounts.ea.com:443/connect/auth" {
		return fmt.Errorf("battlefield: expected redirect to accounts.ea.com:443/connect/auth, got %s%s", u.Host, u.Path)
	}
	if u.Query().Get("fid") != s.fid {
		return fmt.Errorf("battlefield: EA login did not return the fid")
	}
	return nil
}

func (s *Session) connectAuth() error {
	c, err := transport.Client("https://accounts.ea.com")
	if err != nil {
		return err
	}
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse // do not follow redirects
	}
	query := url.Values{}
	query.Set("response_type", "code")
	query.Set("client_id", "Battlefield-CoreWeb")
	query.Set("fid", s.fid)
	resp, err := c.Get(fmt.Sprintf("https://accounts.ea.com/connect/auth?%v", query.Encode()))
	if err != nil && err != http.ErrUseLastResponse {
		return err
	}
	defer resp.Body.Close()
	// http://www.battlefield.com?code=xxxxx
	u, err := resp.Location()
	if err != nil {
		return err
	}
	if u.Host+u.Path != "www.battlefield.com" {
		return fmt.Errorf("battlefield: expected redirect to www.battlefield.com, got %s%s", u.Host, u.Path)
	}
	sid := func() string {
		for _, c := range resp.Cookies() {
			if c.Name == "sid" {
				return c.Value
			}
		}
		return ""
	}()
	if sid == "" {
		return fmt.Errorf("battlefield: signin has no sid")
	}
	s.sid = sid
	return nil
}

func (s *Session) connectAuthCode() error {
	c, err := transport.Client("https://accounts.ea.com")
	if err != nil {
		return err
	}
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse // do not follow redirects
	}
	query := url.Values{}
	query.Set("client_id", "sparta-companion-web")
	query.Set("response_type", "code")
	query.Set("prompt", "none")
	query.Set("redirect_uri", "nucleus:rest")
	req, err := http.NewRequest("GET",
		fmt.Sprintf("https://accounts.ea.com/connect/auth?%v", query.Encode()),
		nil)
	if err != nil {
		return err
	}
	req.AddCookie(&http.Cookie{
		Name:   "sid",
		Value:  s.sid,
		Path:   "/",
		Domain: "accounts.ea.com",
		Secure: true,
	})
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("battlefield: expected 200 got %d", resp.StatusCode)
	}
	var result struct {
		Code string
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if result.Code == "" {
		return fmt.Errorf("battlefield: no code received")
	}
	s.code = result.Code
	return nil
}

func (s *Session) loginFromAuthCode() error {
	c, err := transport.Client("https://companion-api.battlefield.com")
	if err != nil {
		return err
	}
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse // do not follow redirects
	}
	body := strings.NewReader(fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"method": "Companion.loginFromAuthCode",
		"params": {
			"code": "%s",
			"redirectUri": "nucleus:rest"
		}
	}`, s.code))
	resp, err := c.Post(
		"https://companion-api.battlefield.com/jsonrpc/web/api?Companion.loginFromAuthCode",
		"application/json",
		body)
	if err != nil && err != http.ErrUseLastResponse {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("battlefield: expected 200 got %d", resp.StatusCode)
	}
	var result struct {
		Result struct {
			Id string
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if result.Result.Id == "" {
		return fmt.Errorf("battlefield: no id received")
	}
	s.gatewaySessionID = result.Result.Id
	return nil
}

func (s *Session) Stats(id uint64, name string) (*Stats, error) {
	json, err := s.getCareer(id)
	if err != nil {
		return nil, err
	}
	return parseStats(id, name, json)
}

func (s *Session) getCareer(personaID uint64) ([]byte, error) {
	uri := "https://companion-api.battlefield.com/jsonrpc/web/api?Stats.getCareerForOwnedGamesByPersonaId"
	c, err := transport.Client(uri)
	if err != nil {
		return nil, err
	}
	var resp *http.Response
	for attempt := 0; ; attempt++ {
		body := strings.NewReader(fmt.Sprintf(`{
		"jsonrpc": "2.0",
		"method": "Stats.getCareerForOwnedGamesByPersonaId",
		"params": {
			"personaId": "%d"
		},
		"id":""
	}`, personaID))
		req, err := http.NewRequest("POST", uri, body)
		if err != nil {
			return nil, err
		}
		req.Header.Add("X-GatewaySession", s.gatewaySessionID)
		resp, err = c.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusForbidden {
			break
		}
		// login and retry, but only one time
		if attempt == 1 {
			return nil, fmt.Errorf("battlefield: still forbidden after login")
		}
		if err := s.Login(); err != nil {
			return nil, fmt.Errorf("battlefield: login failed")
		}
		continue
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("battlefield: expected status 200, got %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}
