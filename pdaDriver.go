package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"bufio"
	"strings"
)

// The main method of the pdaDriver program. This will check the command-line args, import the json,
// and pass the input string and json spec to an instance of pdaProcessor.
func main(){
	// Check to make sure the user has provided a path for the json spec.
	if len(os.Args) < 2{
		fmt.Println("Error: command-line args must include json spec file path")
		os.Exit(0)
	}
	jsonFilename := string(os.Args[1])
	jsonText, err := ioutil.ReadFile(jsonFilename)
	
	check(err)

	pda := new(PdaProcessor)
	if pda.Open(string(jsonText)){

		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Enter input text: ")

		inputText, _ := reader.ReadString('\n')
		inputTokens := strings.Split(strings.Replace(inputText, string('\n'), "", -1), " ")
		
		// Loop until validateTokens returns true.
		for !validateTokens(inputTokens, pda.InputAlphabet) {
			fmt.Println("Error: input text invalid. Input must contain only the following: ", 
				pda.InputAlphabet)
			fmt.Println("Enter input text: ")
			
			inputText, _ = reader.ReadString('\n')
			inputTokens = strings.Split(strings.Replace(inputText, string('\n'), "", -1), " ")
		}

		fmt.Println("SUCCESS! Your input is accepted: ", inputTokens)

		// Iterate over all input tokens. First, check the current state, input token and top token 
		// then determine whether we can take a transition. If we can make a transition WITHOUT
		// consuming a token, then we will do that. Otherwise we consume a token and make the
		// appropriate transition.
		for _, t := range inputTokens {
			fmt.Println(t)
		}

	} else {
		fmt.Println("Error: could not open json spec")
	}
}

// Calls panic if it detects an error.
func check(e error){
	if e != nil{
		panic(e)
	}
}

// Checks input text to see if the provided tokens match the pda's InputAlphabet
func validateTokens(inputTokens []string, inputAlphabet []string)(bool){

	for _, t := range inputTokens {

		exists := false
		
		for _, s := range inputAlphabet {
			if string(t) == s {
				exists = true
			}
		}
		if !exists {
			return false
		}
	}
	return true
}