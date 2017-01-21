package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func secretHandler(next http.Handler) http.Handler {
	secret := os.Getenv("SECRET_KEY")
	signature_ip := os.Getenv("SIGNATURE_IP")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// step 1: query string is required
		if r.URL.RawQuery == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// step 2: query string show be valid
		values, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// step 3: expires and signature are required
		expires := values.Get("e")
		signature := values.Get("s")
		if expires == "" || signature == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// step 4: expires is int
		expires_at, err := strconv.ParseInt(expires, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// step 5: not expired
		now := int64(time.Now().Unix())
		if now > expires_at {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Expired"))
			return
		}

		// secret + url + expires + ip
		var base string
		if signature_ip == "" {
			base = secret + r.URL.Path + expires
		} else {
			base = secret + r.URL.Path + expires + r.RemoteAddr
		}
		sig := fmt.Sprintf("%x", md5.Sum([]byte(base)))

		if sig == signature {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Invalid signature"))
		}
	})
}

func main() {
	root := os.Args[1]
	handler := http.FileServer(http.Dir(root))
	http.Handle("/", secretHandler(handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8700"
	}

	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe("localhost:"+port, nil))
}
