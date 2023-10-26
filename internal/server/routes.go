package server

import "net/http"

// Add routes here. Works best with libraries which conforms to the standars net
// library. Eg. https://github.com/go-chi/chi and https://github.com/connectrpc/connect-go
// will work.
func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloWorld())
	return mux
}

// Just a simple example of a handler.
func helloWorld() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if _,err := w.Write([]byte("Hello world")); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
