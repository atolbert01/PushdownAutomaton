package main

// Defines the type PdaProcessor.
type PdaProcessor struct {

	// The replica group ID of this pda processor
	Gid int

	// The Id is used for indexing purposes when querying the database.
	Id int `json:"id"`

	// Note: field names must begin with a capital in order to be recognized by the JSON Marshaller
	Name string `json:"name"`
	States []string `json:"states"`
	InputAlphabet []string `json:"input_alphabet"`
	StackAlphabet []string `json:"stack_alphabet"`
	AcceptingStates []string `json:"accepting_states"`
	StartState string `json:"start_state"`
	Transitions [][]string `json:"transitions"`
	EosToken string `json:"eos"`

	// Holds the current state.
	CurrentState string `json:"current_state"`

	// The slice is used to hold the tokens.
	TokenStack []string `json:"token_stack"`

	// Holds the position value of the last put token.
	LastPutPosition int `json:"last_put_position"`
	
	// Holds member ids if passed in with calls to /replica_pdas/gid
	//Members []int `json:"members"`

	// Holds input tokens in order of their position.
	TokenMap map[int]string

	// Vector clock which is used to maintain consistency. The ClockMap will be of size N, where N
	// is the size of the group this pda belongs to.
	ClockMap map[int]int

	PdaCode string
}

type Snap struct {
	CurrentState string `json:"current_state"`
	QueuedTokens []string `json:"queued_tokens"`
	TopTokens []string `json:"top_tokens"`
	Cookie string `json:"cookie"`
}

type GroupMemberAddresses struct {
	Addresses []string `json:"member_addresses"`
}

type PeekResponse struct {
	Tokens string `json:"tokens"`
	Cookie string `json:"cookie"`
}