---
layout: default
title: Chapter 1 
nav_order: 2
parent: Egg timer
has_children: false
---

# Chapter 1 - The empty window

## Goals
The intention of this section is to create a blank canvas that we later can draw upon. 

![An empty window](01_empty_window.gif)

## Outline

The code does three main things:
 - Imports Gio
 - Creates and calls a goroutine that ...
   - ... creates a new window, called `w`
   - ... starts an eternal loop that waits for Events in the window (none will come in this example)

That's it. Let's look at the code:

## Code
```go
package main

import (
	"gioui.org/app"
)

func main() {
	go func() {
		// create new window
		w := app.NewWindow()

		// listen for events in the window.
		for range w.Events() {
		}
	}()
	app.Main()
}
```

## Comments

The code looks simple enough, right? Still, let's take the tiime to to look at what's going on.

1. We import **gioui.org/app**. What's that?
   
   Looking at [the docs](https://pkg.go.dev/gioui.org/app) we find:
   > Package app provides a platform-independent interface to operating system functionality for running graphical user interfaces.
   
   This is good news. Gio takes care of all the plaform-dependent stuff for us. I routinely code on Windows and Macos. Gio just works. [GioUI.org](gioui.org) lists even more, iOS and Android included. 
   
   This is deeper than you might realize. Because, even if your app today is single-platform, your *skillset* is now multi-platform. 
   *"We should port to mac"* Consider it done! *"Hot startup seeking app- and desktop experts*" No problem. *"Who here knows tvOS"* You do.
   *"The pilot died, can anyone land this plane*" ... ok, maybe not, but the point still stand. The diversity of Gio is nothing less than amazing.
   
2. The **event loop** in the goroutine
   
   - The event.loop is the `for range w.Events()` loop that listens for events in the window. For now we just let it listen without doing anything with the events it receives. Later we'll start reacting to them.
   
     From [app.main](https://pkg.go.dev/gioui.org/app#hdr-Main) we learn:
     > Because Main is also blocking on some platforms, the event loop of a Window must run in a goroutine.

    - A goroutine with no name, i.e. an *anonymous function*, is created and runs the event loop. Since it's in a goroutine it will spin concurrently with the rest of the program.
   ```go
     go func {
       ...
     }()
   ```
  
      Jeremy Bytes [writes well about anonymous functions](https://jeremybytes.blogspot.com/2021/02/go-golang-anonymous-functions-inlining.html). They're useful in many contexts, not only with Gio.

   
   

3. Start it by calling **app.Main()**
   > The Main function must be called from a program's main function, to hand over control of the main thread to operating systems that need it.
  
   From [app.Main documentation](https://pkg.go.dev/gioui.org/app#hdr-Main)
