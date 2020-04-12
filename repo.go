package main

//var currentId int = 0

//var pdas []PdaProcessor

// Stores the pdas at a given id.
var	pdas map[int]PdaProcessor

func RepoInit() {
	pdas = make(map[int]PdaProcessor)
}

func RepoCreatePda(pda PdaProcessor) PdaProcessor {
	pdas[pda.Id] = pda

	return pda
}

func RepoFindPda(id int) PdaProcessor {
	return pdas[id]
	/*for _, pda := range pdas {
		if pda.Id == id {
			return pda
		}
	}
	return PdaProcessor{}*/
}

func RepoGetPdas() map[int]PdaProcessor {
	return pdas
}