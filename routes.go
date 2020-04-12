package main

import "net/http"

type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes {
	// Test route which responds to any request by outputting the request URL
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	// base/pdas: List of names of PDAs available at the server.
	Route{
		"ShowPdas",
		"GET",
		"/pdas",
		ShowPdas,
	},

	// base/pdas/{id}: Create PDA at server with given id and specification provided in request 
	// body. Calls open() method of PDA processor.
	Route{
		"CreatePda",
		"PUT",
		"/pdas/{id}",
		CreatePda,
	},

	// base/pdas/{id}/reset: Call reset() method for pda with given id.
	Route{
		"ResetPda",
		"PUT",
		"/pdas/{id}/reset",
		ResetPda,
	},

	// base/pdas/{id}/tokens/{position}: Present a token at the given position.
	Route{
		"PresentToken",
		"PUT",
		"/pdas/{id}/tokens/{position}",
		PresentToken,
	},

	// base/pdas/{id}/eos/{position}: Call eos() for the given pda after the given position.
	Route{
		"PutEos",
		"PUT",
		"/pdas/{id}/eos/{position}",
		PutEos,
	},

	// base/pdas/{id}/is_accepted: Call and return the value of is_accepted()
	Route{
		"GetIsAccepted",
		"GET",
		"/pdas/{id}/is_accepted",
		GetIsAccepted,
	},

	// base/pdas/{id}/stack/top/{k}: Call and return the value of peek(k)
	Route{
		"GetPeek",
		"GET",
		"/pdas/{id}/stack/top/{k}",
		GetPeek,
	},

	// base/pdas/{id}/stack/len: Return the number of tokens currently in the stack.
	Route{
		"GetLen",
		"GET",
		"/pdas/{id}/stack/len",
		GetLen,
	},

	// base/pdas/{id}/state: Call and return the value of current_state
	Route{
		"GetState",
		"GET",
		"/pdas/{id}/state",
		GetState,
	},

	// base/pdas/{id}/tokens: Call and return the value of queued_tokens
	Route{
		"GetQueue",
		"GET",
		"/pdas/{id}/tokens",
		GetQueue,
	},

	// base/pdas/{id}/snapshot/{k}: Retrun a JSON message (array) with 3 components:
	// 	1. current_state()
	// 	2. queued_tokens()
	// 	3. peek(k)
	Route{
		"Snapshot",
		"GET",
		"/pdas/{id}/snapshot/{k}",
		Snapshot,
	},

	// base/pdas/{id}/close: Call close()
	Route{
		"PutClose",
		"PUT",
		"/pdas/{id}/close",
		PutClose,
	},

	// base/pdas/{id}/delete: Delete the PDA with name from the server.
	Route{
		"DeletePda",
		"DELETE",
		"/pdas/{id}/delete",
		DeletePda,
	},
}