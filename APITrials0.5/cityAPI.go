package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type city struct {
	Name string
	Area uint64
}

func mainLogic(w http.ResponseWriter, r *http.Request) {
	// Check if method is POST
	if r.Method == "POST" {
		var tempCity city
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&tempCity)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()
		// no se que mierda de la resource creation logic va a aca
		fmt.Printf("Got %s city with area of %d sq miles!\n", tempCity.Name, tempCity.Area)
		// deci que esta todo bien
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("201 - Created"))
	} else {
		// deci que no esta permitido
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method not allowed"))
	}

}

func main() {
	http.HandleFunc("/city", mainLogic)
	http.ListenAndServe(":8000", nil)
}

numGenerator := generator()
for i := 0; i < 5; i++ {
	fmt.Print(numGenerator(), "\t")

}
//HandlerFunc retorna un handler de http
mainLogicHandler := http.HandlerFunc(mainLogic)
http.Handle("/", middleware(mainLogicHandler))
http.ListenAndServe(":5500", nil)
// this function returnas another function
func generator() func() int {
	var i = 0
	return func() int {
		i++
		return i
	}
}

func middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(escritor http.ResponseWriter, request *http.Request) {
		fmt.Printf("Software ejecutado antes de que la request llegue a la API")
		// Devolviendo el control al handler
		handler.ServeHTTP(escritor, request)
		fmt.Println("Y esto dice que ejecuta el software despues de la response phase")
	})
}

func mainLogic(escritor http.ResponseWriter, request *http.Request) {
	// El business logic va aqui....
	fmt.Println("Este es el manejador principal...")
	escritor.Write([]byte("OK"))
}
