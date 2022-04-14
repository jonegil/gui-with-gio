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

1. Imports Gio
1. Creates and calls a goroutine that open a new window, called `w`
1. Starts a never-ending loop that waits for Events in the window (no Event will ever occur in this example)

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

			// detect which type of event
			switch e := e.(type) {

			// is it a FrameEvent? Those are sent when the application should re-render.
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

## Comments

The code looks simple enough, right? Still, let's take the time to to look at what's going on.

1.  We have some imports. What are those?

 - [os](https://pkg.go.dev/os) - The docs telll us:

    > _Package os provides a platform-independent interface to operating system functionality._

    That's convenient. We'll use it to gracefully exit the program, to ensure a smooth experience for users.

  - [gioui.org/app](https://pkg.go.dev/gioui.org/app) - Specific GUI functionality, cross platform.

    > _Package app provides a platform-independent interface to operating system functionality for running graphical user interfaces._

    This is good news. Gio takes care of all the platform dependent stuff for us. I routinely code on Windows and MacOS. Gio just works. [GioUI.org](https://gioui.org/#installation) lists even more, iOS and Android included.

    This is deeper than you might realize. Even if your app today is single-platform, your _skillset_ is now multi-platform.
    _"We should port to Mac."_ &nbsp;Consider it done! _"Hot startup seeking app and desktop experts._" &nbsp;No problem. _"Who here knows tvOS?"_ &nbsp;You do!
    _"The pilot died, can anyone land this plane?!_" &nbsp;OK, maybe not that last one but the point still stands. The diversity of Gio is nothing less than amazing.

  - [gioui.org/io/system](https://pkg.go.dev/gioui.org/io/system) - Provides high-level events that are sent from the window. Most important is the `system.FrameEvent` which requests a new frame. The new frame is defined through a list of operations. The operations detail *what* to display and *how* to handle input. What and how. That's it.

  - [gioui.org/layout](https://pkg.go.dev/gioui.org/layout) - Defines useful parts of a layout, such as _dimensions_, _constraints_ and _directions_. Also, it includes the layout-concept known as [Flexbox](https://pkg.go.dev/gioui.org/layout#Flex). It's widely used web and user interface development. Among the many introductions, I recommend the one from [Mozilla](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Flexible_Box_Layout/Basic_Concepts_of_Flexbox).

  - [gioui.org/op](https://pkg.go.dev/gioui.org/op) - Operations, or ops, are central in Gio. They are used to update the user interface. There are operations used to draw, handle input, change window properties, scale, rotate and more. Interestingly there are also [macros](https://pkg.go.dev/gioui.org/op#MacroOp), making it possible to record opertions to be executed later. Taken together this means a list of opererations is a _mutable stack_, where you can control the flow.


2. The **event loop** in the goroutine

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
