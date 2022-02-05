package api

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"syscall"
	"time"
)

const (
	PathPostHash = "/hash"
	PathGetHash  = "/hash/"
	PathGetStats = "/stats"
	PathShutdown = "/shutdown"

	HashDelayMs = 5000
)

type HashItem struct {
	ID     int64
	Base64 string
	UnixMs int64
}

func NewHashItem(s string) *HashItem {
	hasher := sha512.New()
	hasher.Write([]byte(s))
	hi := &HashItem{
		UnixMs: time.Now().UnixMilli(),
		Base64: base64.URLEncoding.EncodeToString(hasher.Sum(nil)),
	}
	return hi
}

type HashServer struct {
	*http.ServeMux
	Items         *HashRepository
	HashAverageNs *AtomicAverage
	PID           int
}

func NewHashServer(pid int) *HashServer {
	hs := &HashServer{
		ServeMux:      http.NewServeMux(),
		Items:         NewHashRepository(),
		HashAverageNs: NewAtomicAverage(),
		PID:           pid,
	}
	hs.routes()
	return hs
}

func (hs *HashServer) routes() {
	hs.HandleFunc(PathPostHash, hs.addHashItem())
	hs.HandleFunc(PathGetHash, hs.getHashItem())
	hs.HandleFunc(PathGetStats, hs.getStats())
	hs.HandleFunc(PathShutdown, hs.shutdown())
}

func (hs *HashServer) addHashItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tns := time.Now().UnixNano()
		defer hs.HashAverageNs.Add(time.Now().UnixNano() - tns)

		if r.Method != http.MethodPost {
			http.Error(w, fmt.Sprintf("ERROR: Expected: POST, actual: %s", r.Method), http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, fmt.Sprintf("ERROR: %s", err), http.StatusBadRequest)
			return
		}

		pw := r.FormValue("password")
		if pw == "" {
			http.Error(w, "ERROR: Form field 'password' not found", http.StatusBadRequest)
			return
		}

		id := hs.Items.Add(NewHashItem(pw))
		w.Header().Set("Content-Type", "text/plain")
		out := fmt.Sprintf("%s\n", strconv.Itoa(int(id)))
		if _, err := w.Write([]byte(out)); err != nil {
			http.Error(w, fmt.Sprintf("ERROR: %s", err), http.StatusInternalServerError)
			return
		}
	}
}

func (hs *HashServer) getHashItem() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, fmt.Sprintf("ERROR: Expected: GET, actual: %s", r.Method), http.StatusBadRequest)
			return
		}

		// extract id from '/hash/{id}'
		sid := r.URL.Path[len(PathGetHash):]
		if sid == "" {
			http.Error(w, fmt.Sprintf("ERROR: Missing id %s{id}", PathGetHash), http.StatusBadRequest)
			return
		}

		id, err := strconv.ParseInt(sid, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("ERROR: %s", err), http.StatusBadRequest)
			return
		}

		item, err := hs.Items.Get(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("ERROR: %s", err), http.StatusBadRequest)
			return
		}

		// ensure HashDelayMs
		if ms := time.Now().UnixMilli() - item.UnixMs; ms < HashDelayMs {
			http.Error(w, fmt.Sprintf("ETA: %d", HashDelayMs-ms), http.StatusAccepted)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		out := fmt.Sprintf("%s\n", item.Base64)
		if _, err = w.Write([]byte(out)); err != nil {
			http.Error(w, fmt.Sprintf("ERROR: %s", err), http.StatusInternalServerError)
			return
		}
	}
}

func (hs *HashServer) getStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, fmt.Sprintf("ERROR: Expected: GET, actual: %s", r.Method), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		ss := NewStats(hs.HashAverageNs.Get())
		if err := json.NewEncoder(w).Encode(ss); err != nil {
			http.Error(w, fmt.Sprintf("ERROR: %s", err), http.StatusInternalServerError)
			return
		}
	}
}

func (hs *HashServer) shutdown() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		out := fmt.Sprintf("PID: %s Shutdown in progress\n", strconv.Itoa(hs.PID))
		if _, err := w.Write([]byte(out)); err != nil {
			http.Error(w, fmt.Sprintf("ERROR: %s", err), http.StatusInternalServerError)
			return
		}
		syscall.Kill(hs.PID, syscall.SIGINT)
	}
}
