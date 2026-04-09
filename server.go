package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const GREEN_BASE = "https://api.green-api.com"

type ErrorResp struct {
	Error string `json:"error"`
}

func main() {
	mux := http.NewServeMux()

	origin := func(r *http.Request) string {
		return "*"
	}

	withCORS := func(h http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin(r))
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			h(w, r)
		})
	}
	mux.HandleFunc("/api/getSettings", withCORS(getSettings))
	mux.HandleFunc("/api/getStateInstance", withCORS(getStateInstance))
	mux.HandleFunc("/api/sendMessage", withCORS(sendMessage))
	mux.HandleFunc("/api/sendFileByUrl", withCORS(sendFileByUrl))

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	fmt.Println("Listening on port 8080")
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func proxyJSON(w http.ResponseWriter, r *http.Request, greenUrl string, method string, body io.Reader) {
	client := &http.Client{
		Timeout: time.Second * 15,
	}
	req, err := http.NewRequest(method, greenUrl, body)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResp{err.Error()})
		return
	}
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := client.Do(req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResp{err.Error()})
		return
	}
	defer resp.Body.Close()

	respText, _ := io.ReadAll(resp.Body)

	var js any
	if err := json.Unmarshal(respText, &js); err == nil {
		writeJSON(w, resp.StatusCode, js)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(resp.StatusCode)
	_, _ = w.Write(respText)
}

func getSettings(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	idInstance := q.Get("idInstance")
	apiToken := q.Get("apiToken")
	if idInstance == "" || apiToken == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResp{"idInstance or apiToken is empty"})
		return
	}

	greenUrl := fmt.Sprintf("%s/waInstance%s/getSettings/%s", GREEN_BASE, idInstance, apiToken)
	proxyJSON(w, r, greenUrl, http.MethodGet, nil)
}

func getStateInstance(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	idInstance := q.Get("idInstance")
	apiToken := q.Get("apiToken")
	if idInstance == "" || apiToken == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResp{"idInstance or apiToken is empty"})
		return
	}

	greenUrl := fmt.Sprintf("%s/waInstance%s/getSettings/%s", GREEN_BASE, idInstance, apiToken)
	proxyJSON(w, r, greenUrl, http.MethodGet, nil)
}

type sendMessageReq struct {
	IdInstance string `json:"idInstance"`
	ApiToken   string `json:"apiToken"`
	ChatID     string `json:"chatId"`
	Message    string `json:"message"`
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	var body sendMessageReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResp{err.Error()})
		return
	}
	if body.IdInstance == "" || body.ApiToken == "" || body.ChatID == "" || strings.TrimSpace(body.Message) == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResp{"idInstance or apiToken is empty"})
		return
	}
	greenURL := fmt.Sprintf("%s/waInstance%s/sendMessage/%s", GREEN_BASE, body.IdInstance, body.ApiToken)

	payload := map[string]string{
		"chatId":  body.ChatID,
		"message": body.Message,
	}
	b, _ := json.Marshal(payload)
	proxyJSON(w, r, greenURL, http.MethodPost, strings.NewReader(string(b)))
}

type sendFileByUrlReq struct {
	IdInstance string `json:"idInstance"`
	ApiToken   string `json:"apiToken"`
	ChatID     string `json:"chatId"`
	UrlFile    string `json:"urlFile"`
}

func sendFileByUrl(w http.ResponseWriter, r *http.Request) {
	var body sendFileByUrlReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResp{err.Error()})
		return
	}
	if body.IdInstance == "" || body.ApiToken == "" || body.ChatID == "" || body.UrlFile == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResp{"idInstance or apiToken is empty"})
		return
	}
	if u, err := url.Parse(body.UrlFile); err != nil || u.Scheme == "" || u.Host == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResp{"urlFile is invalid"})
		return
	}
	greenURL := fmt.Sprintf("%s/waInstance%s/sendFileByUrl/%s", GREEN_BASE, body.IdInstance, body.ApiToken)

	payload := map[string]string{
		"chatId": body.ChatID,
		"url":    body.UrlFile,
	}

	b, _ := json.Marshal(payload)
	proxyJSON(w, r, greenURL, http.MethodPost, strings.NewReader(string(b)))
}
