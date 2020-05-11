package main

import(
	"io"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"
	"github.com/gorilla/mux"
)


type Snap struct {
	CurrentState string `json:"current_state"`
	QueuedTokens []string `json:"queued_tokens"`
	TopTokens []string `json:"top_tokens"`
	Cookie string `json:"cookie"`
}

type PeekResponse struct {
	Tokens string `json:"tokens"`
	Cookie string `json:"cookie"`
}

type GroupMemberAddresses struct {
	Addresses []string `json:"member_addresses"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

// Handles requests to URL localhost:8080/pdas. Returns list of ids of PDAs available at server.
func ShowPdas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	pdas := RepoGetPdas()

	var pdaIds string

	for _, pda := range pdas {
		if pda.IsValid() {
			pdaIds += strconv.Itoa(pda.Id) + " "
		} else {
			panic("Invalid pda")
		}
	}

	w.Write([]byte(pdaIds))
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
	pda.PdaCode = string(body)

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

	// Try to parse the request
	parseErr := r.ParseMultipartForm(32 << 20) // Max memory is 32 MB
	if parseErr != nil {
		http.Error(w, "failed to parse multipart message", http.StatusBadRequest)
		return
	}

	token := r.FormValue("token_value")
	sessionCookie := r.FormValue("session_cookie")

	if len(token) < 1 || len(sessionCookie) < 1 {
		w.WriteHeader(422) // Unprocessable entity
		panic("Error, could not present token. Invalid request parameters")
	}

	clientClock := StringToClockMap(sessionCookie)
	
	var pda PdaProcessor
	pda = RepoFindPda(id)
	var updatedCookie = make(map[int]int)


	consistentId := RepoFindConsistentPda(pda, clientClock)
	fmt.Println("Last updated pda: ", consistentId)
	
	pda = RepoMakeConsistent(pda.Id, consistentId, clientClock)
	updatedCookie = pda.ClockMap

	// Finally present the token to pda and rememeber to increment the clock for this id.
	pda.Put(position, token)

	RepoUpdatePda(pda)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	
	w.Write([]byte(ClockMapToString(updatedCookie)))
}

// Handles requests to URL localhost:8080/pdas/{id}/eos/{position}. Call eos() for the given pda 
// after the given position
func PutEos(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])
	var position, _ = strconv.Atoi(vars["position"])
	
	// Read in the request body.
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	sessionCookie := string(body)

	clientClock := StringToClockMap(sessionCookie)

	var pda PdaProcessor
	pda = RepoFindPda(id)
	if !pda.IsValid() {
		panic("PDA not found.")
	}
	
	consistentId := RepoFindConsistentPda(pda, clientClock)
	pda = RepoMakeConsistent(pda.Id, consistentId, clientClock)

	var updatedCookie = make(map[int]int)
	updatedCookie = pda.ClockMap

	if pda.Eos(position) {
		// If the EOS is successful then we update the Repo with the EOSd PDA.
		//RepoCreatePda(pda)
		RepoUpdatePda(pda)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(ClockMapToString(updatedCookie)))
	} else {
		w.Write([]byte("Error: " + pda.Name + ", Eos() - " + " Invalid input stream, input " + 
		"rejected. The language is not recognized."))
	}
}

// Handles requests to URL localhost:8080/pdas/{id}/is_accepted. Call IsAccepted() for the given pda
func GetIsAccepted(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])

	// Read in the request body.
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	sessionCookie := string(body)
	clientClock := StringToClockMap(sessionCookie)

	var pda PdaProcessor
	pda = RepoFindPda(id)
	if !pda.IsValid() {
		panic("PDA not found.")
	}
	
	consistentId := RepoFindConsistentPda(pda, clientClock)
	pda = RepoMakeConsistent(pda.Id, consistentId, clientClock)

	var updatedCookie = make(map[int]int)
	updatedCookie = pda.ClockMap

	if pda.IsAccepted(){
		RepoUpdatePda(pda)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(ClockMapToString(updatedCookie)))
	} else {
		w.Write([]byte("False: " + pda.Name + " not in accepting state: " + pda.CurrentState))
	}
}
// Handles requests to URL localhost:8080/pdas/{id}/stack/top/{k}. Call and return the value of peek(k)
func GetPeek(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])
	var k, _ = strconv.Atoi(vars["k"])

	// Read in the request body.
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	sessionCookie := string(body)
	clientClock := StringToClockMap(sessionCookie)

	var pda PdaProcessor
	pda = RepoFindPda(id)
	if !pda.IsValid() {
		panic("PDA not found.")
	}
	
	consistentId := RepoFindConsistentPda(pda, clientClock)
	pda = RepoMakeConsistent(pda.Id, consistentId, clientClock)

	var updatedCookie = make(map[int]int)
	updatedCookie = pda.ClockMap

	var results = pda.Peek(k)
	var combined = ""
	for _, result := range results {
		combined += result + " "
	}
	//fmt.Println(combined)
	RepoUpdatePda(pda)

	peekResp := PeekResponse{
		combined,
		ClockMapToString(updatedCookie),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(peekResp); err != nil {
		panic(err)
	}
}

// Handles request to URL http://localhost:8080/pdas/{id}/stack/len: Return the number of tokens 
// currently in the stack.
func GetLen(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])

	// Read in the request body.
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	sessionCookie := string(body)
	clientClock := StringToClockMap(sessionCookie)

	var pda PdaProcessor
	pda = RepoFindPda(id)
	if !pda.IsValid() {
		panic("PDA not found.")
	}
	
	consistentId := RepoFindConsistentPda(pda, clientClock)
	pda = RepoMakeConsistent(pda.Id, consistentId, clientClock)
	RepoUpdatePda(pda)

	var updatedCookie = make(map[int]int)
	updatedCookie = pda.ClockMap

	result := PeekResponse{
		strconv.Itoa(len(pda.TokenStack)),
		ClockMapToString(updatedCookie),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

// Handles requests to URL http://localhost:8080/pdas/{id}/state: Call and return the value of 
// current_state
func GetState(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])

	
	// Read in the request body.
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	sessionCookie := string(body)
	clientClock := StringToClockMap(sessionCookie)

	var pda PdaProcessor
	pda = RepoFindPda(id)
	if !pda.IsValid() {
		panic("PDA not found.")
	}
	
	consistentId := RepoFindConsistentPda(pda, clientClock)
	pda = RepoMakeConsistent(pda.Id, consistentId, clientClock)
	RepoUpdatePda(pda)

	var updatedCookie = make(map[int]int)
	updatedCookie = pda.ClockMap

	result := PeekResponse{
		pda.Name + " current state: " + pda.CurrentState,
		ClockMapToString(updatedCookie),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

// Handles requests to URL http://localhost:8080/pdas/{id}/tokens: Call and return the value of 
// queued_tokens
func GetQueue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])

	// Read in the request body.
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	sessionCookie := string(body)
	clientClock := StringToClockMap(sessionCookie)

	var pda PdaProcessor
	pda = RepoFindPda(id)
	if !pda.IsValid() {
		panic("PDA not found.")
	}
	
	consistentId := RepoFindConsistentPda(pda, clientClock)
	pda = RepoMakeConsistent(pda.Id, consistentId, clientClock)
	RepoUpdatePda(pda)

	var updatedCookie = make(map[int]int)
	updatedCookie = pda.ClockMap

	var results = pda.QueuedTokens()

	var combined = ""
	for _, result := range results {
		combined += result + " "
	}

	queueResp := PeekResponse{
		pda.Name + " queued tokens: " + combined,
		ClockMapToString(updatedCookie),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(queueResp); err != nil {
		panic(err)
	}
}

// Handles requests to URL http://localhost:8080/pdas/{id}/snapshot/{k}: Return a JSON message 
// (array) with 3 components:
// 	1. current_state()
// 	2. queued_tokens()
// 	3. peek(k)
func Snapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])
	var k, _ = strconv.Atoi(vars["k"])

	// Read in the request body.
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	sessionCookie := string(body)
	clientClock := StringToClockMap(sessionCookie)

	var pda PdaProcessor
	pda = RepoFindPda(id)
	if !pda.IsValid() {
		panic("PDA not found.")
	}
	
	consistentId := RepoFindConsistentPda(pda, clientClock)
	pda = RepoMakeConsistent(pda.Id, consistentId, clientClock)
	RepoUpdatePda(pda)

	var updatedCookie = make(map[int]int)
	updatedCookie = pda.ClockMap

	var stateResult = pda.CurrentState
	var queueResults = pda.QueuedTokens()
	var peekResults = pda.Peek(k)

	snap := Snap{
		stateResult,
		queueResults,
		peekResults,
		ClockMapToString(updatedCookie),
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

/********************************** BEGIN REPLICA GROUP HANDLERS **********************************/

// Handles GET requests for http://localhost:8080/replica_pdas: Return a list of ids of replica 
// groups currently defined.
func GetGroupIds(w http.ResponseWriter, r *http.Request) {
	results := RepoGetGroupIds()
	resultsText := ""
	
	for _, result := range results {
		resultsText += strconv.Itoa(result) + " "
	}

	w.Write([]byte("Currently defined replica groups: " + resultsText))
}

// Handles PUT requests for http://localhost:8080/replica_pdas/gid: Define a new replica group.
//
// Expects two request parameters: 
//		(1) pda_code: gives the specification used to create the pdas
//		(2) group_members: gives the group member pda addresses.
//
// Create/replace the group members as needed. 
func InitGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var gid, _ = strconv.Atoi(vars["gid"])

	var pda PdaProcessor

	// Try to parse the request
	parseErr := r.ParseMultipartForm(32 << 20) // Max memory is 32 MB
	if parseErr != nil {
		http.Error(w, "failed to parse multipart message", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal([]byte(r.FormValue("pda_code")), &pda); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // Unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	memberStr := strings.Split(r.FormValue("members"), " ")
	members := make([]int, len(memberStr))
	for i := range memberStr {
		members[i], _ = strconv.Atoi(memberStr[i])
	}

	// Reset initializes the pda with starting values.
	pda.Reset()

	if !pda.IsValid() {
		panic("PDA not found.")
	}

	RepoInitGroup(gid, pda, members, string([]byte(r.FormValue("pda_code"))))
}

// Handles PUT requests for http://localhost:8080/replica_pdas/gid/reset: For each pda in the group,
// reset.
func ResetGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var gid, _ = strconv.Atoi(vars["gid"])

	if !RepoResetGroup(gid){
		panic(("Error resetting group"))
	}

	w.Write([]byte(("Group reset")))
}

// Handles GET requests for http://localhost:8080/replica_pdas/gid/members: Return a JSON array with
// the addresses of the members in the given group.
func GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var gid, _ = strconv.Atoi(vars["gid"])

	members := RepoGetGroupMembers(gid)

	if len(members) <= 0 {
		panic("Error retrieving member addresses")
	}

	resultsText := ""
	for _, result := range members {
		resultsText += ("http://localhost:8080/pdas/" + strconv.Itoa(result) + " ")
	}
	
	addresses := GroupMemberAddresses{
		strings.Split(resultsText, " "),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(addresses); err != nil {
		panic(err)
	}

}

// Handles GET requests for http://localhost:8080/replica_pdas/gid/connect: Return the address of a
// random group member that a client could connect to.
func GetConnectMemberId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var gid, _ = strconv.Atoi(vars["gid"])

	id := RepoGetRandomMember(gid)

	if id == -1 {
		panic("Error retrieving connect address")
	}

	w.Write([]byte("http://localhost:8080/pdas/" + strconv.Itoa(id)))
}

// Handles PUT requests for http://localhost:8080/replica_pdas/gid/close: Close the pdas of all
// group members.
func CloseGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	w.Write([]byte("Group closed: " + vars["gid"]))
}

// Handles DELETE requests for http://localhost:8080/replica_pdas/gid/delete: Delete the replica
// group and all its members.
func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var gid, _ = strconv.Atoi(vars["gid"])

	var msg = "Delete successful: " + vars["gid"]
	if !RepoDeleteGroup(gid) {
		msg = "Error during delete: " + vars["gid"]
	}

	w.Write([]byte(msg))
}

// Handles PUT requests for http://localhost:8080/pdas/id/join: Join the pda with the given id to 
// the replica group with group address provided as a request parameter (replica_group).
func PdaJoinGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])

	// Read in the request body.
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	tokens := strings.Split(string(body), "/")
	if len(tokens) < 1 {
		panic("Error, improper request body. Could not join pda")
	}
	gid, _ := strconv.Atoi(tokens[len(tokens) - 1])
	
	var pda PdaProcessor
	RepoJoinPda(id, gid, pda) // The last token should be the gid

	w.Write([]byte("Pda successfully joined"))
}

// Handles GET requests for http://localhost:8080/pdas/id/code: Return the JSON specification of the
// pda with given id.
func GetPdaCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])

	pdaCode := RepoGetPdaCode(id)

	if len(pdaCode) <= 0 {
		panic("Error, could not retrieve pda code")
	}
	w.Write([]byte(pdaCode))
}

// Handles GET requests for http://localhost:8080/pdas/id/c3state: Return JSON message with the
// state information maintained to support client-centric consistency.
func GetPdaStateInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // Get the variables from the request.
	var id, _ = strconv.Atoi(vars["id"])

	pdaClockMap := RepoGetClockMap(id)

	clockString := ClockMapToString(pdaClockMap)

	/*for id, timestamp := range pdaClockMap {
		clockString += strconv.Itoa(id) + ":" + strconv.Itoa(timestamp) + " "
	}*/

	w.Write([]byte(clockString))
}

/*********************************** END REPLICA GROUP HANDLERS ***********************************/

func ClockMapToString(clockMap map[int]int) (clockString string) {
	for id, timestamp := range clockMap {
		clockString += strconv.Itoa(id) + ":" + strconv.Itoa(timestamp) + " "
	}
	return clockString
}

func StringToClockMap(clockString string) (clockMap map[int]int) {
	clockMap = make(map[int]int)

	clockStringArr := strings.Split(clockString, " ")

	for _, pair := range clockStringArr {

		if(len(pair) > 1) {
			splitPair := strings.Split(pair, ":")
			var clockId, _ =  strconv.Atoi(splitPair[0])
			var clockTimestamp, _ =  strconv.Atoi(splitPair[1])
			clockMap[clockId] = clockTimestamp
		}
	}

	return clockMap
}