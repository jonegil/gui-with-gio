---
layout: default
title: Chapter 1 - Window
nav_order: 2
parent: Egg timer
has_children: false
---

# Chapter 1 - An empty window

Updated Feb 23rd 2024

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
    "gioui.org/app"
)

func main() {
    go func() {
        // create new window
        w := app.NewWindow()

        // listen for events in the window
        for {
            w.NextEvent()
        }
    }()
    app.Main()
}

```

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

    - The event loop is the `for { w.NextEvent() }` loop.
      As described in the docs, [w.NextEvent](https://pkg.go.dev/gioui.org/app#Window.NextEvent) simply _blocks until an event is received from the window_. For now we just let it listen without doing anything with the events it receives. Later we'll start reacting to them.

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

[Next chapter](02_title_and_size.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[View it on GitHub](https://github.com/jonegil/gui-with-gio/tree/main/egg_timer){: .btn .fs-5 .mb-4 .mb-md-0 }
