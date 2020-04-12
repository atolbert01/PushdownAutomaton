package main

import(
	"io"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

// Handles requests to URL localhost:8080/pdas. Returns list of names of PDAs available at server.
func ShowPdas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	//var pdas []PdaProcessor
	pdas := RepoGetPdas()

	var pdaNames string

	for _, pda := range pdas {
		if pda.IsValid() {
			pdaNames += pda.Name + " "
		} else {
			panic("Invalid pda")
		}
	}

	w.Write([]byte(pdaNames))
}

// Handles requests to URL localhost:8080/pdas/{id}. Creates PDA with given id at server.
func CreatePda(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r) // Get the variables from the request.
	id := vars["id"]
	
	var pda PdaProcessor

	// Read in the request body.
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &pda); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // Unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	// Make sure the pda is created in the repo with the correct id. This is used for retrieval.
	pda.Id, _ = strconv.Atoi(id)

	// Reset initializes the pda with starting values.
	pda.Reset()

	if !pda.IsValid() {
		panic("PDA not found.")
	}

	p := RepoCreatePda(pda)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		panic(err)
	}
}

// Handles requests to URL localhost:8080/pdas/{id}/reset. Resets the PDA with given id.
func ResetPda(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])
	
	var pda PdaProcessor

	pda = RepoFindPda(id)

	if !pda.IsValid() {
		panic("PDA not found.")
	}

	pda.CurrentState = "q2"
	pda.TokenStack = []string{"0", "0", "1", "1"}
	fmt.Println(pda)

	pda.Reset()

	fmt.Println(pda)

	p := RepoCreatePda(pda)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		panic(err)
	}
}

// Handles requests to URL localhost:8080/pdas/{id}/tokens/{position}. Present a token at the given 
// position.
func PresentToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])
	var position, _ = strconv.Atoi(vars["position"])

	// Read in the request body.
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	var token = string(body)

	fmt.Println("Token presented: ", token)

	var pda PdaProcessor

	pda = RepoFindPda(id)

	if !pda.IsValid() {
		panic("PDA not found.")
	}

	pda.Put(position, token)

	p := RepoCreatePda(pda)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		panic(err)
	}

}

// Handles requests to URL localhost:8080/pdas/{id}/eos/{position}. Call eos() for the given pda 
// after the given position
func PutEos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])
	var position, _ = strconv.Atoi(vars["position"])

	var pda PdaProcessor

	pda = RepoFindPda(id)

	if !pda.IsValid() {
		panic("PDA not found.")
	}

	if pda.Eos(position) {
		// If the EOS is successful then we update the Repo with the EOSd PDA.
		RepoCreatePda(pda)
		w.Write([]byte("Input stream is accepted. Language recognized.\n"))
		w.WriteHeader(http.StatusCreated)
	} else {
		w.Write([]byte("Error: " + pda.Name + ", Eos() - " + " Invalid input stream, input " + 
		"rejected. The language is not recognized.\n"))
	}
}

// Handles requests to URL localhost:8080/pdas/{id}/is_accepted. Call IsAccepted() for the given pda
func GetIsAccepted(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])

	var pda PdaProcessor

	pda = RepoFindPda(id)

	if !pda.IsValid() {
		panic("PDA not found.")
	}

	if pda.IsAccepted(){
		w.Write([]byte("True: " + pda.Name + " is in accepting state: " + pda.CurrentState))	
	} else {
		w.Write([]byte("False: " + pda.Name + " not in accepting state: " + pda.CurrentState))
	}
}

func GetPeek(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	id := vars["id"]
	k := vars["k"]

	_ = id
	_ = k
}

func GetLen(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	id := vars["id"]

	_ = id
}

func GetState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	id := vars["id"]

	_ = id
}

func GetQueue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	id := vars["id"]

	_ = id
}

func Snapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	id := vars["id"]
	k := vars["k"]

	_ = id
	_ = k
}

func PutClose(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	id := vars["id"]

	_ = id
}

func DeletePda(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	id := vars["id"]

	_ = id
}