package main

import (
	"encoding/json"
	//"fmt"
)

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
	Eos string `json:"eos"`

	// Holds the current state.
	CurrentState string

	// Token at the top of the stack.
	CurrentTop string

	// The slice is used to hold the tokens.
	TokenStack []string

}

// Unmarshals the jsonText string. Returns true if it succeeds.
func (pda *PdaProcessor) Open(jsonText string) (bool){

	if err := json.Unmarshal([]byte(jsonText), &pda); err != nil {
		check(err)
	}

	// Validate input.	
	if len(pda.Name) == 0 || len(pda.States) == 0 || len(pda.InputAlphabet) == 0 || 
	len(pda.StackAlphabet) == 0 || len(pda.AcceptingStates) == 0 || len(pda.StartState) == 0 ||
	len(pda.Transitions) == 0 || len(pda.Eos) == 0 {
		return false
	}

	return true
}

// Sets the CurrentState to StartState and assigns TokenStack a new empty slice
func (pda *PdaProcessor) Reset(){
	pda.CurrentState = pda.StartState
	pda.TokenStack = []string{}
}

func (pda *PdaProcessor) Put(token string){
	
}
