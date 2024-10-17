package main

import (
	"fmt"

	"github.com/riadafridishibly/go-events"
)

func main() {
	eventHandler := events.NewEventHandler[int]()

	eventHandler.On("foobar", func(i int) {
		fmt.Println("Number 1:", i)
	})

	eventHandler.On("foobar", func(i int) {
		fmt.Println("Number 2:", i)
	})

	eventHandler.On("baz", func(i int) {
		eventHandler.Emit("foobar", i+1)
	})

	eventHandler.Emit("foobar", 42)
	eventHandler.Emit("baz", 43)

	eventHandler.Wait()
}
