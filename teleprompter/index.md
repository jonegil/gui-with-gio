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
[View it on GitHub](https://github.com/jonegil/gui-with-gio/teleprompter/code){: .btn .fs-5 .mb-4 .mb-md-0 }

## Goals

This project continues where the [egg timer](../egg_timer/) leaves off. The timer was a good start and gave us the foundation to build an app. But we're not done. Especially we should look closer at how how to deal with user input, both keyboard and mouse. 

To do that we'll build what's known as a [teleprompter](https://en.wikipedia.org/wiki/Teleprompter). That's simply an app that displays and scrolls text, but we'll make sure it's both lively and responsive for the user to interact with. We'll make sure to look into some other new parts of Gio as well.

![Mr_Gorbachev_tear_down_this_wall](teleprompter_Mr_Gorbachev.gif)

Ready? 

Let's (sc)roll!
(sorry)

## Outline

More precisely our teleprompter should
1. Read text from a ```.txt``` file so the speaker can easily display his or her own scripts
1. Allow full flexibility to adjust **font-size** and **text width**
1. Help the speaker by displaying a **focusbar** that can be moved to where it's most useful
1. **Full control**  of manual scroll, but also **auto scroll** that's easy to start, stop, pause, speed up and slow down.
1. Easy to remember keyboard controls, and be fully controllable with only one hand. Designed for  guestilating Italians!

<p align="center">
  <img src="teleprompter_with_text.jpeg" alt="Teleprompter with text" height="250"/>
  <!--img src="teleprompter.jpeg" alt="Teleprompter and camera" height="250"/-->
</p>

## Source code

Time for code. We´ll zoom out a bit and look at the main blocks and structure of the program.


### Section 1 - Starting up

In the first part of the program we list the imports and start up the app. 

#### New imports
```go
import (
  // ... many normal imports we discussed earlier ...

  // ... plus two new interesting Gio imports
  "gioui.org/io/key"
  "gioui.org/io/pointer"
)
```
These two are new to us and gives support for keyboard and mouse events:
 - *Package [io/key](https://pkg.go.dev/gioui.org/io/key) implements key and text events and operations.*
 - *Package [io/pointer](https://pkg.go.dev/gioui.org/io/pointer) implements pointer events and operations. A pointer is either a mouse controlled cursor or a touch object such as a finger.*


#### A placeholder for the speech
```go
// A []string to hold the speech as a list of paragraphs
var paragraphList []string
```

The text to be displayed is read from ```speech.txt```. We could have stored it in one long string, but as we will see later, it will be easier to work with it as a list of paragraphs. Hence we store it as a ```[]string``` slice insteads.


#### Read the file and start the gui

```go
func main() {
  // Part 1 - Read from file
  f, err := ioutil.ReadFile("speech.txt")
  if err == nil {
    // Convert whole text into a slice of strings.
    paragraphList = strings.Split(string(f), "\n")
    // Add extra empty lines a the end. Cheap but effective trick to ensure
    // the last line of the speech scrolls out of the screen
    for i := 1; i <= 10; i++ {
      paragraphList = append(paragraphList, "")
    }
  }

  // Part 2 - Start the gui
  go func() {
    // create new window
    w := app.NewWindow(
      app.Title("Teleprompter"),
      app.Size(unit.Dp(350), unit.Dp(300)),
    )
    // draw on screen
    if err := draw(w); err != nil {
      log.Fatal(err)
    }
    os.Exit(0)
  }()

  app.Main()
}
```

Part 1 above reads the text and splits it by ```\n```, newline. To enure we have enough empty space to allow the last element to scroll off screen, we add 10 empty paragraphs at the end. That's a pragmatic but unelegant solution I think. But not worth pondering too long on either, so let's move on.

### Section 2 - Define variables to control behaviour

```go
func draw(w *app.Window) error {
  // variables that control behviour
  // y-position for text
  var scrollY int = 0

  // y-position for red focusbar
  var highlightY int = 78

  // width of text area
  var textWidth int = 300

  // fontSize
  var fontSize int = 35

  // Are we auto scrolling?
  var autoscroll bool = false
  var autospeed int = 1
  
```

Now we're getting into the meat of things. In order to control the behaviour of the program we need multiple state variables. These keep track of multiple things for us:
 - How far down have we scrolled? Defaults to starting on the top.
 - The red focusbar is there to keep the speekers attention on one single line. But it's helpful to see the last line spoken, and a few of the next ones too. So the focusbar should highlight a line quite high up, but not all the way. Hence what I feel is a sensible starting point. But remember that the should be able to adjust this, so it's a variable, not a constant.
 - Explain the textWidth ...
 - Explain the fontSize ...
 - Autoscroll ...

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
