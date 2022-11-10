package main

type ChangeStateMessage struct {
	Option string
	Key    string
	Val    interface{}
}
