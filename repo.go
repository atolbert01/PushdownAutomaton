package main

import (
	"math/rand"
	"time"
)

// A map which holds all replica groups where the key is the gid and the value is a 
// map[int]PdaProcessor where the key is the id of the pda and all PdaProcessors share the same
// spec.
var replicaGroups map[int]map[int]PdaProcessor//map[int]string

// Stores the base implementation of each pda assigned to a group. The key is the group id.
var pdaCodes map[int]PdaProcessor

// Stores the pdas at a given id.
var	pdas map[int]PdaProcessor

func RepoInit() {
	pdas = make(map[int]PdaProcessor)
	replicaGroups = make(map[int]map[int]PdaProcessor)
	pdaCodes = make(map[int]PdaProcessor)
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

/********************************* BEGIN REPLICA GROUP FUNCTIONS **********************************/

func RepoInitGroup(gid int, pda PdaProcessor, members []int){
	pdaCodes[gid] = pda // Store the pda code so we can quickly retrieve it later.
	
	var newGroup = make(map[int]PdaProcessor)

	for _, m := range members {
		var newPda = pda
		newPda.Gid = gid
		newPda.Id = m
		newGroup[m] = newPda
	}
	replicaGroups[gid] = newGroup
	InitClocks(gid)


	// for _, p := range replicaGroups[gid] {
	// 	fmt.Println(p)
	// }
}

// Function to reset all the clocks in a group to zero.
func InitClocks(gid int) {
	var length = len(replicaGroups[gid])

	// for {key}, {value} := range {list}
	for _, p := range replicaGroups[gid] {
		p.ResetClock(length)
		
		for _, m := range replicaGroups[gid] {
		
			p.SetClock(m.Id, 0) // Set every timestamp in the map to 0
		}
		replicaGroups[gid][p.Id] = p
	}
}

func RepoGetGroupIds() (gids []int) {

	for gid, _ := range replicaGroups {
		gids = append(gids, gid)
	}

	return gids
}

func RepoResetGroup(gid int) (bool) {

	var length = len(replicaGroups[gid])

	for _, p := range replicaGroups[gid] {
		p.Reset()
		p.ResetClock(length)
		replicaGroups[gid][p.Id] = p
	}

	return true
}

func RepoGetGroupMembers(gid int) (ids []int) {
	for _, p := range replicaGroups[gid] {
		ids = append(ids, p.Id)
	}

	return ids
}

func RepoGetRandomMember(gid int) (id int) {

	id = -1
	rand.Seed(time.Now().Unix()) // Initialize global pseudo random generator
	










	// NOT WORKING




	idx := rand.Intn(len(replicaGroups[gid]))

	for key, _ := range replicaGroups[gid] {
		if key == replicaGroups[gid][idx].Id {
			id = key
		}
	}

	return id
}

/********************************** END REPLICA GROUP FUNCTIONS ***********************************/