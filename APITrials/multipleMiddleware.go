package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type city struct {
	Name string
	Area uint64
}

// Middleware para chequear que el contenido sea JSON
func filterContentType(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Middleware de chequeo de contenido")
		// Filtering requests by MIME type
		if r.Header.Get("Content-type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("415 - Unsupported media type. Please send JSON"))
			return
		}
		handler.ServeHTTP(w, r)
	})
}

// Middleware para adherir un timestamp en la response cookie
func setServerTimeCookie(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
		// Seteando cookie para cada response
		cookie := http.Cookie{Name: "Server-Time(UTC)", Value: strconv.FormatInt(time.Now().Unix(), 10)}
		http.SetCookie(w, &cookie)
		log.Println("Currently in the set server time middleware")
	})
}

func mainLogic(w http.ResponseWriter, r *http.Request) {
	// Chequeamos que el metodo sea POST
	if r.Method == "POST" {
		var tempCity city
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&tempCity)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()
		//resource creation que dice que es solo print pero por ahora
		log.Printf("Got %s city with area of %d sq miles!\n", tempCity.Name, tempCity.Area)
		// decir tooodo bienn
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("201 - Created"))
	} else {
		// decir que el metodo no es
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - metodo not allowed"))
	}
}

func main() {
	mainLogicHandler := http.HandlerFunc(mainLogic)
	http.Handle("/city",
		filterContentType(setServerTimeCookie(mainLogicHandler)))

	http.ListenAndServe(":8000", nil)
}
