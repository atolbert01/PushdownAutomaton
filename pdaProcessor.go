package main

import (
	"encoding/json"
	"fmt"
)

type Stack []string

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
	CurrentState string

	// The slice is used to hold the tokens.
	TokenStack Stack

	// A map that holds the transition functions
	TransitionMap map[PdaConfig]PdaTransition

}

// Unmarshals the jsonText string. Returns true if it succeeds.
func (pda *PdaProcessor) Open(jsonText string) (bool){

	if err := json.Unmarshal([]byte(jsonText), &pda); err != nil {
		check(err)
	}

	// Validate input.	
	if len(pda.Name) == 0 || len(pda.States) == 0 || len(pda.InputAlphabet) == 0 || 
	len(pda.StackAlphabet) == 0 || len(pda.AcceptingStates) == 0 || len(pda.StartState) == 0 ||
	len(pda.Transitions) == 0 || len(pda.EosToken) == 0 {

		return false
	}

	// Load TransitionMap with values from the json
	pda.TransitionMap = make(map[PdaConfig]PdaTransition)
	for _, t := range pda.Transitions {

		newConfig := PdaConfig{t[0], t[1], t[2]}
		newTransition := PdaTransition{t[3], t[4]}

		pda.TransitionMap[newConfig] = newTransition
	}
	return true
}

// Sets the CurrentState to StartState and assigns TokenStack a new empty slice
func (pda *PdaProcessor) Reset(){
	pda.CurrentState = pda.StartState
	pda.TokenStack = []string{}
}

func (pda *PdaProcessor) Put(token string) (numTransitions int){

	numTransitions = 0

	top := pda.CurrentTop()

	currentConfig := PdaConfig{pda.CurrentState, token, top}

	//var transition PdaTransition
	fmt.Println("Current config: ", currentConfig)
	if transition, ok := pda.TransitionMap[currentConfig]; ok {

		pda.CurrentState = transition.NextState;

		// Check if we are pushing a null token
		if len(transition.PushToken) == 0 {
			
			// If the top token is not null then pop it.
			if len(pda.CurrentTop()) != 0 {
				
				pda.TokenStack.Pop()
			}
			//pda.TokenStack.Push(transition.PushToken)
		} else {
			pda.TokenStack.Push(transition.PushToken)
		}

		numTransitions++
	}

	if numTransitions == 0 && len(token) > 0 {
		pda.TokenStack.Push(token)
	}


	return numTransitions
}

// Checks for the top k tokens in the stack and returns them without removing them.
func (pda *PdaProcessor) Peek(k int) ([]string) {

	var tokens []string

	if pda.TokenStack.IsEmpty() {
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
	if pda.TokenStack.IsEmpty() {
		return ""
	} else {
		return pda.Peek(1)[0]
	}
}

func (pda *PdaProcessor) Eos() {
	
	for _, s := range pda.AcceptingStates {

		// We are in an accepting state.
		if pda.CurrentState == s {
			// Are we att the end of the stack?
			if pda.CurrentTop() == pda.EosToken {
				pda.TokenStack.Pop()
			}
			// Is the stack empty?
			if pda.TokenStack.IsEmpty() {
				fmt.Println("Input stream is accepted. Language recognized.")
				return
			}
		}
	}

	fmt.Println("Error: Invalid input stream, input rejected. The language is not recognized.")
}


func (pda *PdaProcessor) IsAccepted() (bool) {

	for _, s := range pda.AcceptingStates {
		if pda.CurrentState == s && pda.TokenStack.IsEmpty() {
			return true
		}
	}
	return false
}
/*********************************** BEGIN STACK IMPLEMENTATION ***********************************/

// Find out if the token stack is empty.
func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

// Push a new token to the stack.
func (s *Stack) Push(str string) {
	*s = append(*s, str)
}

// Remove and return top token of stack. Return false if stack is empty.
func (s *Stack) Pop() (string, bool) {
	if s.IsEmpty() {
		return "", false
	} else {
		index := len(*s) - 1 // Get the index of top stack token.
		token := (*s)[index] // Index into the slice and obtain the token.
		*s = (*s)[:index] // Remove it from the stack by slicing it off
		return token, true
	}
}

/************************************ END STACK IMPLEMENTATION ************************************/