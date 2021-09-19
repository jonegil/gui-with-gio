---
layout: default
title: Teleprompter
nav_order: 3
has_children: true
has_toc: false
---


# Teleprompter - Animation and interaction

```

Hi mate!

The code for this part is done - but not the text. 
Please look around while I continue to write this chapter. 

As always, please pull the repo and play around.

Cheers

```

## Goals

This project continues where the [egg timer](../egg_timer/) leaves off. The timer was a good start and gave us the foundation to build an app. But we're not done. Especially we should look closer at how how to deal with user input, both keyboard and mouse. 

To do that we'll build what's known as a [teleprompter](https://en.wikipedia.org/wiki/Teleprompter). That's simply an app that displays and scrolls text, but we'll make sure it's both lively and responsive for the user to interact with. We'll make sure to look into some other new parts of Gio as well.

![Mr_Gorbachev_tear_down_this_wall](teleprompter_Mr_Gorbachev.gif)

Ready? 

Let's (sc)roll!
(sorry)

[View it on GitHub](https://github.com/jonegil/gui-with-gio/teleprompter/code){: .btn .fs-5 .mb-4 .mb-md-0 }

## Outline

More precisely our teleprompter should
1. Read text from a ```.txt``` file so the speaker can easily display his or her own scripts
1. Allow full flexibility to adjust **font-size** and **text width**
1. Help the speaker by displaying a **focus bar** that can be moved to where it's most useful
1. **Full control**  of manual scroll, but also **auto scroll** that's easy to start, stop, pause, speed up and slow down.
1. Easy to remember keyboard controls, and be fully controllable with only one hand. Designed for  guestilating Italians!

<p align="center">
  <img src="teleprompter_with_text.jpeg" alt="Teleprompter with text" height="250"/>
  <!--img src="teleprompter.jpeg" alt="Teleprompter and camera" height="250"/-->
</p>

## Source code

Time for code. We´ll zoom out a bit and look at the main blocks and structure of the program.


### Section 1 - Starting up

In the first part of the program we list the imports and start up the app. A slightly shortened version goes like this:

```go
import (
  // Two new interesting Gio imports
  "gioui.org/io/key"
	"gioui.org/io/pointer"
)

// the []string to hold the speech as a list of paragraphs
var paragraphList []string

func main() {
	// Part 1 - Read the speech from file
  f, err := ioutil.ReadFile("speech.txt")

  // Part 2 - Start the GUI
	go func() {
		// draw on screen
    w := app.NewWindow
    draw(w) 
	}()
}
```

[gioui.org/io/key]() and [gioui.org/io/pointer]() are new to us, and as their names suggest give us support for processing keyboard and mouse events.

We then define a ```[]string``` that we'll later fill with contents from the ```.txt``` we read from disk.

Finally create a new window and start the app.

### Section 2 - Define variables to control behaviour

```go
func draw(w *app.Window) error {
  // variables that control behviour
	// y-position for text
	var scrollY int = 0

	// y-position for red highlight bar
	var highlightY int = 78

	// width of text area
	var textWidth int = 300

	// fontSize
	var fontSize int = 35

	// Are we auto scrolling?
	var autoscroll bool = false
	var autospeed int = 1
  
```

### Section 3 - Listen for events
In the second part we listen for events and draw on screen.

```go

	// listen for events in the window.
	for e := range w.Events() {

		// Detect what type of event
		switch e := e.(type) {

		// A keypress?
		case key.Event:
	
		// A mouse event?
		case pointer.Event:
	
		// A re-render request?
		case system.FrameEvent:
	
			// Layout the speech as a list of paragraphs
			// 1) First the margins ...

        // 2) ... then the list inside those margins ...
				
          // 3) ... where each paragraph is it's separate item
			
			// Draw a transparent red focusbar.

		// Shutdown?
		case system.DestroyEvent:
	
    }
  }

```


All the source-code is in this repo, in the ```teleprompter/code``` folder.
