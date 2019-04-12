package main

import(
	"log"
	"net/http"
	"text/template"
	"path/filepath"
	"fmt"
	"os"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"time"
)

const (
	FormNameEmail = "email"
	FormNamePassword = "password"
	FormNameRecaptcha = "recaptcha_token"
	EnvRecaptchaSiteKey = "RECAPCHA_SITE_KEY"
	EnvRecaptchaSiteSecret = "RECAPCHA_SITE_SECRET"
)

type loginHandler struct {
	filename string
}

func (l *loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("LoginHandler %#v", r)

	l.filename = "login.html"
	data := map[string]interface{}{
		"FormName": map[string]interface{}{
			"Email": FormNameEmail,
			"Password": FormNamePassword,
			"Recaptcha": FormNameRecaptcha,
		},
		"Recaptcha": map[string]interface{}{
			"SiteKey": os.Getenv(EnvRecaptchaSiteKey),
		},
	}
	if r.Method == http.MethodPost {
		r.ParseForm()
		if err := loginVerify(r.FormValue(FormNameEmail), r.FormValue(FormNamePassword)); err != nil {
			data["Error"] = err
		} else {
			if r, err := reCaptchaVerify(r.FormValue(FormNameRecaptcha)); err != nil {
				data["Error"] = err
			} else{
				if r.Success { // && r.Score > 0.7 {
					// Success reCAPTCHA
					data["RecaptchaResponse"] = r
					l.filename = "index.html"
				} else {
					data["Error"] = fmt.Sprintf("reCAPTCHA failed:%s", r.ErrorCodes)
				}
			}
		}
	}

	t := template.Must(template.ParseFiles(filepath.Join("templates", l.filename)))
	t.Execute(w, data)
}

func loginVerify(email string, passowrd string) error {
	// something verify
	if "" == email || "" == passowrd {
		return fmt.Errorf("Login failed")
	}
	return nil
}

type RecaptchaResponse struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

func reCaptchaVerify(token string) (r RecaptchaResponse, err error) {
	const recaptchaServerName = "https://www.google.com/recaptcha/api/siteverify"
	resp, err := http.PostForm(recaptchaServerName,
		url.Values{"secret": {os.Getenv(EnvRecaptchaSiteSecret)}, "remoteip": {"127.0.0.1"}, "response": {token}})
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return
	}
	return
}

func main() {
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		fmt.Fprintln(w, "data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQEAYAAABPYyMiAAAABmJLR0T///////8JWPfcAAAACXBIWXMAAABIAAAASABGyWs+AAAAF0lEQVRIx2NgGAWjYBSMglEwCkbBSAcACBAAAeaR9cIAAAAASUVORK5CYII=")
	})
	http.Handle("/", &loginHandler{})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}