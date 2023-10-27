package model

type Thread struct {
	Id       int
	Title    string
	Messages []Message
}

type Message struct {
	Id        *int
	Message   *string
	Thread_id *int
}
