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


type Snap struct {
	CurrentState string `json:"current_state"`
	QueuedTokens []string `json:"queued_tokens"`
	TopTokens []string `json:"top_tokens"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

// Handles requests to URL localhost:8080/pdas. Returns list of names of PDAs available at server.
func ShowPdas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

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
		w.Write([]byte("Input stream is accepted. Language recognized."))
		w.WriteHeader(http.StatusCreated)
	} else {
		w.Write([]byte("Error: " + pda.Name + ", Eos() - " + " Invalid input stream, input " + 
		"rejected. The language is not recognized."))
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
// Handles requests to URL localhost:8080/pdas/{id}/stack/top/{k}. Call and return the value of peek(k)
func GetPeek(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])
	var k, _ = strconv.Atoi(vars["k"])

	var pda PdaProcessor

	pda = RepoFindPda(id)

	if !pda.IsValid() {
		panic("PDA not found.")
	}

	var results = pda.Peek(k)
	var combined = ""
	for _, result := range results {
		combined += result + " "
	}

	if (combined == "") {
		w.Write([]byte(pda.Name + " peek returned no results."))
	} else {
		w.Write([]byte(pda.Name + " top " + strconv.Itoa(k) + " elements: " + combined))
	}
}

// Handles request to URL http://localhost:8080/pdas/{id}/stack/len: Return the number of tokens 
// currently in the stack.
func GetLen(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])

	var pda PdaProcessor

	pda = RepoFindPda(id)

	if !pda.IsValid() {
		panic("PDA not found.")
	}

	var result = pda.Name + " stack size: " + strconv.Itoa(len(pda.TokenStack))

	w.Write([]byte(result))
}

// Handles requests to URL http://localhost:8080/pdas/{id}/state: Call and return the value of 
// current_state
func GetState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])

	var pda PdaProcessor

	pda = RepoFindPda(id)

	if !pda.IsValid() {
		panic("PDA not found.")
	}

	var result = pda.Name + " current state: " + pda.CurrentState

	w.Write([]byte(result))
}

// Handles requests to URL http://localhost:8080/pdas/{id}/tokens: Call and return the value of 
// queued_tokens
func GetQueue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])

	var pda PdaProcessor

	pda = RepoFindPda(id)

	if !pda.IsValid() {
		panic("PDA not found.")
	}

	var results = pda.QueuedTokens()

	var combined = ""
	for _, result := range results {
		combined += result + " "
	}

	w.Write([]byte(pda.Name + " queued tokens: " + combined))
}

// Handles requests to URL http://localhost:8080/pdas/{id}/snapshot/{k}: Retrun a JSON message 
// (array) with 3 components:
// 	1. current_state()
// 	2. queued_tokens()
// 	3. peek(k)
func Snapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])
	var k, _ = strconv.Atoi(vars["k"])

	var pda PdaProcessor

	pda = RepoFindPda(id)

	if !pda.IsValid() {
		panic("PDA not found.")
	}

	var stateResult = pda.CurrentState
	var queueResults = pda.QueuedTokens()
	var peekResults = pda.Peek(k)

	snap := Snap{
		stateResult,
		queueResults,
		peekResults,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(snap); err != nil {
		panic(err)
	}

}

func PutClose(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	id := vars["id"]

	_ = id
}

// Handles requests for http://localhost:8080/pdas/{id}/delete: Delete the PDA with name from the 
// server.
func DeletePda(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])

	var pda PdaProcessor

	pda = RepoFindPda(id)

	if !pda.IsValid() {
		panic("PDA not found.")
	}

	pdas := RepoRemovePda(id)

	var pdaNames string

	for _, p := range pdas {
		if p.IsValid() {
			pdaNames += p.Name + " "
		} else {
			panic("Invalid pda")
		}
	}

	w.Write([]byte(pda.Name + " deleted. Remaining pdas: " + pdaNames))

}