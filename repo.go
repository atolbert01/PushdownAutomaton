package main

import (
	"math/rand"
	"time"
	"encoding/json"
	"fmt"
)

// A map which holds all replica groups where the key is the gid and the value is a 
// map[int]PdaProcessor where the key is the id of the pda and all PdaProcessors share the same
// spec.
var replicaGroups map[int]map[int]PdaProcessor//map[int]string

// Stores the base implementation of each pda assigned to a group. The key is the group id.
var pdaCodes map[int]string

// Master list of all pdas currently in existence. Stores the pdas at a given id.
var	pdas map[int]PdaProcessor

func RepoInit() {
	pdas = make(map[int]PdaProcessor)
	replicaGroups = make(map[int]map[int]PdaProcessor)
	pdaCodes = make(map[int]string)
}

func RepoCreatePda(pda PdaProcessor) PdaProcessor {
	//RepoRemovePda(pda.Id) // If a pda with this id already exists in a group or alone, delete it.
	
	pdas[pda.Id] = pda // Add the pda to the master list, but no group.

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

	// Also delete this pda from any group it may be a part of.
	for key, _ := range replicaGroups {
		delete(replicaGroups[key], id)
	}

	return pdas
}

/********************************* BEGIN REPLICA GROUP FUNCTIONS **********************************/

func RepoInitGroup(gid int, pda PdaProcessor, members []int, code string){
	pdaCodes[gid] = code // Store the pda code so we can quickly retrieve it later.
	
	var newGroup = make(map[int]PdaProcessor)

	for _, m := range members {
		var newPda = pda
		newPda.Gid = gid
		newPda.Id = m
		newPda.PdaCode = code

		// Add to the newly created group.
		newGroup[m] = newPda

		// Also add to our master list of pdas.
		RepoCreatePda(newPda)
	}
	replicaGroups[gid] = newGroup
	InitClocks(gid)
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
		pdas[p.Id] = p
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

	var keys []int
	for key, _ := range replicaGroups[gid] {
		keys = append(keys, key)
	}

	id = keys[rand.Intn(len(keys) - 1)]

	return id
}

func RepoDeleteGroup(gid int) (bool) {
	if _, ok := replicaGroups[gid]; ok {
		for _, p := range replicaGroups[gid] {
			RepoRemovePda(p.Id) // Make sure we remove the pda from master list
		}
		delete(replicaGroups, gid)
		return true
	}
	return false
}

func RepoJoinPda(id int, gid int, pda PdaProcessor) (bool) {
	if _, ok := pdas[id]; ok { // if the pda exists join, else return false

		if _, ok := replicaGroups[gid]; ok { // if the replica group exists, then join
			
			jsonData := []byte(pdaCodes[gid])

			if err := json.Unmarshal(jsonData, &pda); err != nil {
				panic("Could not read pda_code during join")
			}

			pda.Reset()
			
			pda.Gid = gid
			pda.Id = id
			pda.PdaCode = pdaCodes[gid]
			
			pda.ResetClock(len(replicaGroups[gid]))
			RepoCreatePda(pda)
			replicaGroups[gid][id] = pda

			// Init Clocks routine			
			for _, m := range replicaGroups[gid] {
				
				pda.SetClock(m.Id, 0) // Set every timestamp in the map to 0
				m.SetClock(id, 0) // Add the new pda to every clock in the group

				replicaGroups[gid][m.Id] = m
				pdas[m.Id] = m			
			}
			pda.SetClock(pda.Id, 0) // Add self to clock map

			replicaGroups[gid][id] = pda
			pdas[id] = pda

			fmt.Println()

		} else { // Create a new group with just this pda in it.
			RepoInitGroup(gid, pdas[id], []int{id}, pdas[id].PdaCode)
		}
		return true
	}
	return false
}

func RepoGetPdaCode(id int) (string) {
	pda := pdas[id]
	return pda.PdaCode
}

func RepoGetClockMap(id int) (map[int]int) {
	/*if len(pdas[id].ClockMap) > 1 {

	} else {
		return pdas[id].ClockMap
	}*/
	return pdas[id].ClockMap
}

func RepoFindConsistentPda(pda PdaProcessor, clientClock map[int]int) (consistentId int) {
	
	highest := -1
	consistentId = -1
	
	// Get the highest timestamp
	for id, ts := range clientClock {
		if ts >= highest {
			highest = ts
			consistentId = id
		}
	}

	if highest == 0 {
		return pda.Id
	}

	return consistentId
}

func RepoMakeConsistent(idToUpdate int, consistentId int, clientClock map[int]int) (PdaProcessor) {
	
	pdaToUpdate := pdas[idToUpdate]
	
	pdaToUpdate.ClockMap[idToUpdate] = clientClock[consistentId] + 1
		
	consistentPda := pdas[consistentId]

	pdaToUpdate.TokenMap = consistentPda.TokenMap
	pdaToUpdate.LastPutPosition = consistentPda.LastPutPosition
	pdaToUpdate.TokenStack = consistentPda.TokenStack
	pdaToUpdate.CurrentState = consistentPda.CurrentState

	RepoUpdatePda(consistentPda)

	for id, ts := range clientClock {
		if ts > pdaToUpdate.ClockMap[id] {
			pdaToUpdate.ClockMap[id] = ts
		}
	}
	return pdaToUpdate
}

func RepoUpdatePda(pda PdaProcessor) PdaProcessor {
	pdas[pda.Id] = pda // Add the pda to the master list, but no group.
	if (len(pda.ClockMap) > 1) {
		replicaGroups[pda.Gid][pda.Id] = pda
	}

	return pda
}
/********************************** END REPLICA GROUP FUNCTIONS ***********************************/