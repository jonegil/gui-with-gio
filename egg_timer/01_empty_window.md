---
layout: default
title: Chapter 1 - Window
nav_order: 2
parent: Egg timer
has_children: false
---

# Chapter 1 - An empty window

## Goals

The intent of this section is to create a blank canvas that we later can draw upon.

![An empty window](01_empty_window.gif)

## Outline

The code does three main things:

- Imports Gio
- Creates and calls a goroutine that:
  - Creates a new window, called `w`
  - Starts a never-ending loop that waits for Events in the window (no Event will ever occur in this example)

That's it! Let's look at the code:

## Code

```go
package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
)

func main() {
	go func() {
		// create new window
		w := app.NewWindow()
		var ops op.Ops

		// listen for events in the window.
		for e := range w.Events() {

			// detect what type of event
			switch e := e.(type) {

			// Is ita FrameEvent? Those are sent when the application should re-render.
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				e.Frame(gtx.Ops)
			}
		}
		os.Exit(0)
	}()
	app.Main()
}

```
---
Out of sync
{: .label .label-red }

As of today, April 3rd, the text below is not fully in sync with the code. The code runs fine, but I need to rewrite some of the text below. Should still be helpful though. 

---

## Comments

The code looks simple enough, right? Still, let's take the time to to look at what's going on.

1.  We import `gioui.org/app`. What's that?

    Looking at [the docs](https://pkg.go.dev/gioui.org/app) we find:

    > _Package app provides a platform-independent interface to operating system functionality for running graphical user interfaces._

    This is good news. Gio takes care of all the platform dependent stuff for us. I routinely code on Windows and MacOS. Gio just works. [GioUI.org](https://gioui.org/#installation) lists even more, iOS and Android included.

    This is deeper than you might realize. Even if your app today is single-platform, your _skillset_ is now multi-platform.
    _"We should port to Mac."_ &nbsp;Consider it done! _"Hot startup seeking app and desktop experts._" &nbsp;No problem. _"Who here knows tvOS?"_ &nbsp;You do!
    _"The pilot died, can anyone land this plane?!_" &nbsp;OK, maybe not that last one but the point still stands. The diversity of Gio is nothing less than amazing.

2.  The **event loop** in the goroutine

    - The event loop is the `for range w.Events()` loop that listens for events in the window. For now we just let it listen without doing anything with the events it receives. Later we'll start reacting to them.

      From [app.main](https://pkg.go.dev/gioui.org/app#hdr-Main) we learn:

      > _Because Main is also blocking on some platforms, the event loop of a Window must run in a goroutine._

    - A goroutine with no name (i.e. an _anonymous function_) is created and runs the event loop. Since it's in a goroutine it will spin concurrently with the rest of the program.

    ```go
    go func {
    // ...
    }()
    ```

          Jeremy Bytes [writes well about anonymous functions](https://jeremybytes.blogspot.com/2021/02/go-golang-anonymous-functions-inlining.html). They're useful in many contexts, not only with Gio.

3.  Start it by calling `app.Main()`
    From [app.Main documentation](https://pkg.go.dev/gioui.org/app#hdr-Main):
    > _The Main function must be called from a program's main function, to hand over control of the main thread to operating systems that need it._

---

[Next chapter](02_title_and_size.md){: .btn .fs-5 .mb-4 .mb-md-0 }
