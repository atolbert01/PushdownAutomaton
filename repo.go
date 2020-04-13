package main

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
}

func RepoGetPdas() map[int]PdaProcessor {
	return pdas
}

func RepoRemovePda(id int) map[int]PdaProcessor {
	delete(pdas, id)
	return pdas
}