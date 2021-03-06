//+build js,wasm

package main

import (
	"fmt"
	. "github.com/siongui/godom/wasm"
	"syscall/js"
	"time"
)

var (
	mainButton    = Document.GetElementById("main-button")
	mainText      = Document.GetElementById("main-text")
	addTextButton = Document.CreateElement("button")
	exitButton    = Document.CreateElement("button")
	textCounter   = 0
)

const (
	addMoreText       = "Should I say hello world?"
	exitText          = "quit this nonsense"
	fadeTimeInSeconds = 1
)

func setHeaderStuff() {
	favLink := Document.QuerySelector("link[rel~='icon']")
	favLink.Set("href", "favicon.ico")
}

func setupScreen(c chan bool) {
	Document.Set("title", "Go WASM demo fun!")
	addTextButton.Set("textContent", addMoreText)
	addTextButton.Get("style").Set("marginRight", "50px")
	mainButton.Get("style").Set("marginTop", "10px")

	mainButton.Call("appendChild", addTextButton)
	addTextButton.Call("addEventListener", "click",
		// this function doesn't need the channel
		//therefore no requirement for function to be inlined
		js.FuncOf(addText))

	exitButton.Set("textContent", exitText)
	mainButton.Call("appendChild", exitButton)
	exitCallBack := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		transition := fmt.Sprintf("%ds ease", fadeTimeInSeconds)

		addTextButton.Get("style").Set("transition", transition)
		addTextButton.Get("style").Set("opacity", "0")

		exitButton.Get("style").Set("transition", transition)
		exitButton.Get("style").Set("opacity", "0")
		go func() {
			time.Sleep(time.Duration(fadeTimeInSeconds) * time.Second)
			mainButton.Call("removeChild", exitButton)
			mainButton.Call("removeChild", addTextButton)
			c <- true
		}()
		return nil
	})
	exitButton.Call("addEventListener", "click", exitCallBack)
}

func addText(_ js.Value, _ []js.Value) interface{} {
	textCounter += 1
	textNode := Document.CreateTextNode(fmt.Sprintf("Hello world x %d", textCounter))
	br := Document.CreateElement("br")

	mainText.Call("appendChild", textNode)
	mainText.Call("appendChild", br)

	return nil
}

func orchestrator(c chan bool) {
	setHeaderStuff()
	setupScreen(c)
}

func main() {
	c := make(chan bool)
	go orchestrator(c)
	<-c
}
