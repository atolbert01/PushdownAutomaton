package main

import (
	"encoding/json"
)

type PdaProcessor struct {
	Name string `json:"name"`
	States []string `json:"states"`
	InputAlphabet []string `json:"input_alphabet"`
	StackAlphabet []string `json:"stack_alphabet"`
	AcceptingStates []string `json:"accepting_states"`
	StartState string `json:"start_state"`
	Transitions [][]string `json:"transitions"`
	Eos string `json:"eos"`
}

// Unmarshals the jsonText string. Returns true if it succeeds.
func (pda *PdaProcessor) Open(jsonText string) (bool){
	json.Unmarshal([]byte(jsonText), &pda)
	
	return true
}