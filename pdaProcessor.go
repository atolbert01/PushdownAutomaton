package main

import (
	"fmt"
)

// Used as the key in the transition function map.
type PdaConfig struct {
	State string // The state the PDA should be in to take the associated transition
	InputToken string // The input token required to take the associated transition
	TopToken string // The token that should be on the token top to take the associated transition
}

// Used as the value in the transition function map.
type PdaTransition struct {
	NextState string // The state the PDA should transition to if it is in the required config.
	PushToken string // The token that should be pushed to the stack if this transition is taken.
}

// Defines the type PdaProcessor.
type PdaProcessor struct {

	// The Id is used for indexing purposes when querying the database.
	Id int `json:"id"`

	// Note: field names must begin with a capital in order to be recognized by the JSON Marshaller
	Name string `json:"name"`
	States []string `json:"states"`
	InputAlphabet []string `json:"input_alphabet"`
	StackAlphabet []string `json:"stack_alphabet"`
	AcceptingStates []string `json:"accepting_states"`
	StartState string `json:"start_state"`
	Transitions [][]string `json:"transitions"`
	EosToken string `json:"eos"`

	// Holds the current state.
	CurrentState string `json:"current_state"`

	// The slice is used to hold the tokens.
	TokenStack []string `json:"token_stack"`

	// Holds the position value of the last put token.
	LastPutPosition int `json:"last_put_position"`
	
	// Holds input tokens in order of their position.
	TokenMap map[int]string
}

// Checks pda to make sure it has been initialized properly.
func (pda *PdaProcessor) IsValid()(bool){
	// Validate input.	
	if len(pda.Name) == 0 || len(pda.States) == 0 || len(pda.InputAlphabet) == 0 || 
	len(pda.StackAlphabet) == 0 || len(pda.AcceptingStates) == 0 || len(pda.StartState) == 0 ||
	len(pda.Transitions) == 0 || len(pda.EosToken) == 0 {

		return false
	}
	return true
}

// Find the appropriate transition given a PdaConfig as input.
func (pda *PdaProcessor) FindTransition(config PdaConfig) (PdaTransition){
	for _, t := range pda.Transitions {

		if config.State == t[0] && config.InputToken == t[1] && config.TopToken == t[2] {
			return PdaTransition{t[3], t[4]}
		}
	}
	return PdaTransition{}
}

// Sets the CurrentState to StartState and assigns TokenStack a new empty slice
func (pda *PdaProcessor) Reset(){
	pda.CurrentState = pda.StartState
	pda.TokenStack = []string{}
	pda.TokenMap = make(map[int]string)
	pda.LastPutPosition = -1
}

func (pda *PdaProcessor) Put(position int, newToken string) (numTransitions int){

	// First add the new token to the map. If this is a fresh Put call, then the token does not
	// currently exist in our map.
	pda.TokenMap[position] = newToken

	// First see if we can make a transition without using a token
	numTransitions = pda.ExecutePut(true)
	numTransitions = pda.ExecutePut(false)
	
	return numTransitions
}

func (pda *PdaProcessor) ExecutePut(isNullToken bool) (numTransitions int){
	
	numTransitions = 0
	var token string = ""
	var nextPosition = -1

	var trueNext = pda.LastPutPosition + 1
	
	if !isNullToken {
		// Now get the next token in the map.
		nextPosition = pda.NextQueuedPosition()
		token = pda.TokenMap[nextPosition]
	}

	// Compare the next queued token position to the position of the most recently Put token.
	// If the nextPosition is one more than the last put token, or equal then we will put this token
	// Otherwise we will wait for the next Put. nextPosition and LastPutPosition will be equal if we
	// get position zero on the first Put.
	if nextPosition == pda.LastPutPosition + 1 || isNullToken {

		top := pda.CurrentTop()
		currentConfig := PdaConfig{pda.CurrentState, token, top}
		fmt.Println("Current config: ", currentConfig)

		transition := pda.FindTransition(currentConfig)
		if len(transition.NextState) > 0 {
			fmt.Println("Transition found: ", transition)
			pda.CurrentState = transition.NextState;

			// Check if we are pushing a null token
			if len(transition.PushToken) == 0 {
				
				// If the top token is not null then pop it.
				if len(pda.CurrentTop()) != 0 {

					pda.Pop()
				}
			} else {

				pda.Push(transition.PushToken)
			}

			numTransitions++
		}

		if numTransitions == 0 && !isNullToken {
			
			pda.Push(token)
		}

		if nextPosition == trueNext{
			pda.LastPutPosition = nextPosition
		}
		numTransitions = pda.ExecutePut(false)
	}

	return numTransitions
}

// Checks for the top k tokens in the stack and returns them without removing them.
func (pda *PdaProcessor) Peek(k int) ([]string) {

	var tokens []string

	if pda.IsEmpty() {
		return tokens
	} else {
		stackSize := len(pda.TokenStack)

		for i := stackSize - 1; i >= stackSize - k; i-- {
			token := (pda.TokenStack)[i]
			tokens = append(tokens, token)
		}
		return tokens
	}
}

// Return the token at the top of the stack.
func (pda *PdaProcessor) CurrentTop() (string) {
	if pda.IsEmpty() {
		return ""
	} else {
		return pda.Peek(1)[0]
	}
}

func (pda *PdaProcessor) Eos(position int) (bool) {
	pda.Put(position, pda.EosToken)
	if position == pda.LastQueuedPosition(){
		for _, s := range pda.AcceptingStates {

			// We are in an accepting state.
			if pda.CurrentState == s {
				// Are we att the end of the stack?
				if pda.CurrentTop() == pda.EosToken {

					pda.Pop()
				}
				// Is the stack empty?
				if pda.IsEmpty() {
				
					return true
				}
			}
		}
	}
	return false
}


func (pda *PdaProcessor) IsAccepted() (bool) {

	for _, s := range pda.AcceptingStates {
		if pda.CurrentState == s && pda.IsEmpty() {
			return true
		}
	}
	return false
}

func (pda *PdaProcessor) QueuedTokens() ([]string) {
	var results []string
	startPos := pda.NextQueuedPosition()
	sorted := []int{startPos}

	for position, _ := range pda.TokenMap {
		if position >= startPos {
			sorted = append(sorted, position)
		}
	}
	
	for _, i := range sorted {
		results = append(results, pda.TokenMap[i])
	}

	return results
}

func (pda *PdaProcessor) NextQueuedPosition() (nextPosition int) {
	
	nextPosition = -1

	for position, _ := range pda.TokenMap {
	
		if position == pda.LastPutPosition + 1 {
			nextPosition = position
		}
	}
	return nextPosition
}

func (pda *PdaProcessor) LastQueuedPosition() (lastPosition int) {
	lastPosition = -1

	for position, _ := range pda.TokenMap {
		
		if position > lastPosition {
			lastPosition = position
		}
	}
	return lastPosition
}

/*********************************** BEGIN STACK IMPLEMENTATION ***********************************/

// Find out if the token stack is empty.
func (pda *PdaProcessor)IsEmpty() bool {
	return len(pda.TokenStack) == 0
}

// Push a new token to the stack.
func (pda *PdaProcessor)Push(str string) {
	pda.TokenStack = append(pda.TokenStack, str)
}

// Remove and return top token of stack. Return false if stack is empty.
func (pda *PdaProcessor) Pop() (string, bool) {
	if pda.IsEmpty() {
		return "", false
	} else {
		index := len(pda.TokenStack) - 1 // Get the index of top stack token.
		token := (pda.TokenStack)[index] // Index into the slice and obtain the token.
		pda.TokenStack = (pda.TokenStack)[:index] // Remove it from the stack by slicing it off
		return token, true
	}
}

/************************************ END STACK IMPLEMENTATION ************************************/