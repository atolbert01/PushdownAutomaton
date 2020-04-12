package main

import (
	//"fmt"
	//"os"
	//"io/ioutil"
	//"bufio"
	//"strings"
	//"net/http"
	//"log"
	//"github.com/gorilla/mux"
)

// The main method of the pdaDriver program. This will check the command-line args, import the json,
// and pass the input string and json spec to an instance of pdaProcessor.
/*func main(){
	
	// Initialize router
	router := NewRouter()

	fmt.Println("Listening on port 8080...")

	log.Fatal(http.ListenAndServe(":8080", router))










	
	// Check to make sure the user has provided a path for the json spec.
	if len(os.Args) < 2{
		fmt.Println("Error: pdaDriver.main() - command-line args must include json spec file path")
		os.Exit(0)
	}
	jsonFilename := string(os.Args[1])
	jsonText, err := ioutil.ReadFile(jsonFilename)
	
	check(err)

	reader := bufio.NewReader(os.Stdin)

	inputText := ""

	// The user passed in input as a file
	if len(os.Args) > 2 {
		commandFilename := string(os.Args[2])
		commandText, err := ioutil.ReadFile(commandFilename)
		check(err)

		inputText = string(commandText)
	} else {

		fmt.Println("Enter input text: ")

		inputText, _ = reader.ReadString('\n')
	}

	inputTokens := strings.Split(strings.Replace(inputText, string('\n'), "", -1), " ")

	pda := new(PdaProcessor)
	if pda.Open(string(jsonText)){

		pda.Reset()
			
		// Loop until validateTokens returns true.
		for !validateTokens(inputTokens, pda.InputAlphabet) {

			fmt.Println("Error: " + pda.Name + ", pdaDriver.main() - input text invalid. " + 
				"Input must contain only the following: ", pda.InputAlphabet)
			fmt.Println("Enter input text: ")
			
			inputText, _ = reader.ReadString('\n')

			inputTokens = strings.Split(strings.Replace(inputText, string('\n'), "", -1), " ")
		}

		fmt.Println("Your input is valid: ", inputTokens)

		// Add the '$' token to signify the end of the input stream
		inputTokens = append(inputTokens, pda.EosToken)

		// Iterate over all input tokens. First, check the current state, input token and top token 
		// then determine whether we can take a transition. If we can make a transition WITHOUT
		// consuming a token, then we will do that. Otherwise we consume a token and make the
		// appropriate transition.
		numberTransitions := 0
		i := 0
		for ok := true; ok; ok = i < len(inputTokens) {

			fmt.Println("Input Token: ", inputTokens[i])

			// First see if a transition can be taken without consuming an input token.
			numTrans := pda.Put(0, "")

			// If not, then consume input token.
			if numTrans == 0 { 
				numTrans = pda.Put(0, inputTokens[i])
				i++
			}

			numberTransitions += numTrans
			fmt.Println("Number of transitions: ", numTrans)
			fmt.Println("Stack Size: ", len(pda.TokenStack))
			fmt.Println()
		}
		// We reached the Eos Token so now call Eos()
		pda.Eos()
		fmt.Println("Total transitions: ", numberTransitions)
		fmt.Println("\n\n")
	} else {
		fmt.Println("Error: " + pda.Name + ", pdaDriver.main() - could not open json spec")
	}
	
	close()
	
}*/

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

// Clears garbage.
func close() {
	// The garbage is collected
}