package main

type Clip interface {
	Inform() <-chan string
	Update(string)
}
