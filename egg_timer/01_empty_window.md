---
layout: default
title: Chapter 1 - Window
nav_order: 2
parent: Egg timer
has_children: false
---

# Chapter 1 - An empty window

## Goals

The intent of this section is to create an empty window that we later can fill with our application.

![An empty window](01_empty_window.gif)

## Outline

The intention of this chapter is to build the simplest possible application that still runs. 

Our code will do three main things:

1. Import the packages needed for a bare-bones application
1. Open a new window
1. Listen for events, although we don't yet handle them

Although this on purpose is a very limited ambition, we still need to cover some ground. When learning I'm a breadth-first kind of guy, so I'll focus on conveying the overall structure first. Once that's in place, we'll go deeper in later chapters, but not here. Still we need to cover quite some ground, so hang with me, and as long as you can see the forest and not get lost among the trees, we're doing great. 

All good? Let's look at the code:

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

**Imports**

We have some imports. To give you a feel for what's going on, we'll look at each one. Later we will go into each in more detail, and discuss how they are used. Here we only touch on each one but you don't need to drill deep into each. Stay with me on the big picture as we start dipping our toes line by line.

 - [os](https://pkg.go.dev/os) - The docs tell us:

    > _Package os provides a platform-independent interface to operating system functionality._

    That's convenient. We'll only use it to gracefully exit the program. Without it some people reported a non-smooth experience in this Chapter. And we want smooooooth.

  - [gioui.org/app](https://pkg.go.dev/gioui.org/app) - Specific GUI functionality, cross platform.

    > _Package app provides a platform-independent interface to operating system functionality for running graphical user interfaces._

    This is good news. Gio takes care of all the platform dependent stuff for us. I routinely code on Windows and MacOS. Gio just works. [GioUI.org](https://gioui.org/#installation) lists even more, iOS and Android included.

    This is deeper than you might realize. Even if your app today is single-platform, your _skillset_ is now multi-platform.
    _"We should port to Mac."_ &nbsp;Consider it done! _"Hot startup seeking app and desktop experts._" &nbsp;No problem. _"Who here knows tvOS?"_ &nbsp;You do!
    _"The pilot died, can anyone land this plane?!_" &nbsp;OK, maybe not that last one but the point still stands. The diversity of Gio is nothing less than amazing.

  - [gioui.org/io/system](https://pkg.go.dev/gioui.org/io/system) - Provides high-level events that are sent from the window. Most important is the `system.FrameEvent` which requests a new frame to be drawn on your screen. The new frame is defined through a list of operations. The operations detail *what* to display and *how* to handle input. What and how. That's it.

  - [gioui.org/layout](https://pkg.go.dev/gioui.org/layout) - Defines useful parts of a layout, such as **dimensions**, **constraints** and **directions**. Also, it includes the layout-concept known as [Flexbox](https://pkg.go.dev/gioui.org/layout#Flex) which is widely used for both web and user interfaces development. Among the many introductions, I recommend the one from [Mozilla](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Flexible_Box_Layout/Basic_Concepts_of_Flexbox).

  - [gioui.org/op](https://pkg.go.dev/gioui.org/op) - Operations, or ops, are central in Gio. They are used to update the user interface. There are operations used to draw, handle input, change window properties, scale, rotate and more. Interestingly there are also [macros](https://pkg.go.dev/gioui.org/op#MacroOp), making it possible to record opertions to be executed later. Taken together this means a list of opererations is a _mutable stack_, where you can control the flow.

**The main() function**

The `main()` function opens with a goroutine with no name. What's up with that? Here's the structure:

```go
go func {
  // ... gui stuff
}()
app.Main()
```

A goroutine with no name is known as a an _anonymous function_. It is created to run a loop that listens for events from the window to react to. Since it's in a goroutine it will spin concurrently with the rest of any other logic your program might have. 

Jeremy Bytes [writes well about anonymous functions](https://jeremybytes.blogspot.com/2021/02/go-golang-anonymous-functions-inlining.html). They're useful in many contexts, not only with Gio.

Ok, so that's what it _is_. But _why_ do we need it? Turns out, from [app.main](https://pkg.go.dev/gioui.org/app#hdr-Main) we learn:

> _Because Main is also blocking on some platforms, the event loop of a Window must run in a goroutine._

Allright. Goroutine it is then. Start it by calling `app.Main()`
From [app.Main documentation](https://pkg.go.dev/gioui.org/app#hdr-Main):
> _The Main function must be called from a program's main function, to hand over control of the main thread to operating systems that need it._



**The event loop**
  Okay, with imports and `main()` setup out of the way, let's look at the meat of things, the event loop. The event loop is simply a _loop_ that listens for _events_. TaDa! (Sorry, I couldn't help myself).
  
  - `for e := range w.Events()` spins forever. `w.Events()` gets us the _channel_ through which events are delivered, and thus we simply listen to whatever comes through it.

  - Then we need to react. But ... what's this `e := e.(type)` line? It's actually a neat thing, known as a [type switch](https://tour.golang.org/methods/16) that allows us to take different actions depending on the `type` of event that's being processed.

  - In our case, we're only interested if the event is a `system.FrameEvent`. If so:
  
    - We define a new **graphical context**, or `gtx`. The context receives the pointer to `ops` as well as the FrameEvent `e`.

    - Later we will add a lot of operations to write and draw and lay out stuff, but for now we skip that.

    - The end the FrameEvent by drawing the grapical operations from ops into the window. That's done by sending the operations `ops` from the context `gtx` to the FrameEvent `e`.

**Leaving the event loop**
Once the event loop is complete, the program has done what it should and we can gracefully close it. For that we use [os.Exit()](https://pkg.go.dev/os?utm_source=gopls#Exit). 

The convention is that a zero exit code indicates success, which is what we send from `os.Exit(0)`if _err_ is nil. If not, we call `log.Fatal(err)` which prints the error message and exits with `os.Exit(1)`. There's no logic so far to report various errors though.


## Well done!

Well done! We've covered lots of ground and laid down some solid foundations. From here we will buid on those, and go deeper into the details where needed. And if it feels a bit much, don't sweat details, as long as you have a feel for the the big picture. Telescope, not microscope.

Ready for the next one? Let's roll!

---

[Next chapter](02_title_and_size.md){: .btn .fs-5 .mb-4 .mb-md-0 }
