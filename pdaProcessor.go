package main

import (
	"encoding/json"
	//"fmt"
)

// A struct that defines values for a State. isCurrent is true when the State is the active state of
// the pda. isAccepting is true if this State is one of the accepting states defined in the json.
type State struct {
	isCurrent bool
	isAccepting bool
}

// Defines the type PdaProcessor.
type PdaProcessor struct {
	Name string `json:"name"`
	TextStates []string `json:"states"`
	InputAlphabet []string `json:"input_alphabet"`
	StackAlphabet []string `json:"stack_alphabet"`
	AcceptingStates []string `json:"accepting_states"`
	StartState string `json:"start_state"`
	Transitions [][]string `json:"transitions"`
	Eos string `json:"eos"`

	// Stores all the States defined in the input text.
	States map[string]State
}

// Unmarshals the jsonText string. Returns true if it succeeds.
func (pda *PdaProcessor) Open(jsonText string) (bool){
	json.Unmarshal([]byte(jsonText), &pda)
	
	// Validate input.	
	if len(pda.Name) == 0 || len(pda.TextStates) == 0 || len(pda.InputAlphabet) == 0 || 
	len(pda.StackAlphabet) == 0 || len(pda.AcceptingStates) == 0 || len(pda.StartState) == 0 ||
	len(pda.Transitions) == 0 || len(pda.Eos) == 0{
		return false
	}

	pda.States = make(map[string]State)

	return true

	// Create State 'enumerator'. Go doesn't have actual enumerators, so we must improvise.
	/*for i, v := range pda.TextStates {
		fmt.Println(i)
		pda.States[v] = State{false, false}
	}

	for k, v := range pda.States {
		fmt.Println("%s -> %s", k, v)
	}*/

	return true
}

func (pda *PdaProcessor) Reset(){

}