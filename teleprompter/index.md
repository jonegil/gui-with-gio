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

To do that we'll build what's known as a [teleprompter](https://en.wikipedia.org/wiki/Teleprompter). A teleprompter is simply a device that displays and scrolls text. Sophisticated and expensive equipment exists, but it can just as easily be done with an app that displays and scrolls text. And that's the version we will build here. Since it needs to be lively and responsive for the user, it it's a great example for how we can react to keypresses and mouse scrolls. We'll make sure to look into some other new parts of Gio as well.

![Mr_Gorbachev_tear_down_this_wall](teleprompter_Mr_Gorbachev.gif)

Ready? 

Let's (sc)roll!
(sorry)

[View it on GitHub](https://github.com/jonegil/gui-with-gio/teleprompter/code){: .btn .fs-5 .mb-4 .mb-md-0 }

## Outline

The goals for our teleprompter are to:
1. Read text from a ```.txt``` file so the speaker can display personal scripts.
1. **Full control** of manual scroll, as well as an **auto scroll** that is easy to start, stop, pause, speed up and slow down.
1. Allow full flexibility to adjust **font-size** and **text width**
1. Help the speaker by displaying a **focusbar** that can be moved to where it's most useful
1. Easy to remember keyboard controls, fully controllable with only one hand. Designed for gesticulation!

<p align="center">
  <img src="teleprompter_with_text.jpeg" alt="Teleprompter with text" height="250"/>
</p>

## Source code

Time for code. We´ll take one step back and start with the overall structure of the program.

### Section 1 - Starting up

In the first part of the program we list the imports and start up the app. 

#### New imports
```go
import (
  // Many normal imports we discussed earlier ...
  // ... plus two new interesting Gio imports
  "gioui.org/io/key"
  "gioui.org/io/pointer"
)
```
These two are new to us and gives support for keyboard and mouse events:
 - *Package [io/key](https://pkg.go.dev/gioui.org/io/key) implements key and text events and operations.*
 - *Package [io/pointer](https://pkg.go.dev/gioui.org/io/pointer) implements pointer events and operations. A pointer is either a mouse controlled cursor or a touch object such as a finger.*

Notice how pointer supports both mouse gestures on a desktop/laptop and fingers on a screen. Nice, again an example of how learning a cross-platform framework gives you tools to master multiple devices.

#### A placeholder for the speech
```go
// A []string to hold the speech as a list of paragraphs
var paragraphList []string
```

The text to be displayed is read from ```speech.txt```. The text in that file will be parsed and stored as a collection of paragraphs. As we will see later, that will be easier to work with. Hence we store it as a ```[]string``` slice directly.


#### Read the file and start the gui

```go
func main() {
  // Part 1 - Read from file
  f, err := ioutil.ReadFile("speech.txt")
  if err == nil {
    // Convert whole text into a slice of strings.
    paragraphList = strings.Split(string(f), "\n")
    // Add extra empty lines a the end. Simple trick to ensure
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

**Part 1** above reads the text and splits it by ```\n```, newline. Each of those parts, or parapgraphs, are then appended to the ```[]paragraphList```. One small issue I experienced was that the final paragraph tended to be shown on screen even when the speech was complete. To solve that, a very simple solution is to add empty paragraphs at the end of the list. 

**Part 2** starts the Gui in a normal manner. 

### Section 2 - Define variables to control behaviour

```go
func draw(w *app.Window) error {
  // variables that control behviour
  // y-position for text
  var scrollY int = 0

  // y-position for red focusBar
  var focusBarY int = 78

  // width of text area
  var textWidth int = 300

  // fontSize
  var fontSize int = 35

  // Are we auto scrolling?
  var autoscroll bool = false
  var autospeed int = 1
  
```

Now we're getting into the meat of things. In order to control the behaviour of the program we need multiple state variables. The user will adjust all of these while using the program, so we can't have them hard coded into the various portions of the visualisation. Instead we collect them here to keep the program tidy.

 - ```scrollY``` - How far into the speech are we? It starts on top of course, and its value increments as the speaker scrolls down into the speech. This is the variable we adjust either manually or automatically to progress the speech. 
   - Move text with ```two-finger trackpad```,  ```mouse scrollwheel```, ```arrow keys```, ```j``` and ```k``` (vim)
 - ```focusBarY``` - A red focusbar heps keep the speeker's attention on one single line. At the same time it's helpful to see some context, so some of the words just spoken are above it, and more of the words to come are below. Therefore I like the focusbar to be a few lines from the top. But most importantly, it should be easy to adjust, so it's a variable that can be controlled.
   - Move it up with ```u``` and down with ```d```
 - ```textWidth``` - Should the speech fill the full width of the window, or a narrower portion of it. I prefer the speech to be fairly narrow, but that all depends on screen-setup, distance to screen, where the camera is, or if there even is one. On a laptop, the camera is very close to your face, so narrow text will not create too much eye-movement. Experiment and find what works.
   - Make the text ```w```ider or ```n```arrower 
 - ```fontSize``` - The size of the font, obviously 
   - Tune it with ```+``` and ```-```
 - ```autoscroll``` and ```autospeed``` - Should we scroll automatically? And if so, how fast? 
   - Start and stop with ```space```. Make it ```f```aster or ```s```lower

### Section 3 - Listen for events

Finally, we get to listen for events. As outlined above, there are quite a few inputs here, and they can have mutual impact on each other. For example, if ```textWdith```increases, the line breaks will adjust since there are now space for more words on each line. But if the user increases ```fontSize```, each word requires more space and line break changes again. Luckily for us Gio takes care of all of the underlying details, as long as we're keeping track of the value of those state variables. 

The structure is quite rich:

```go

  // listen for events in the window.
  for e := range w.Events() {

    // Detect what type of event
    switch e := e.(type) {

    // A keypress?
    case key.Event:
      // Update and store state for size, width and positioning

    // A mouse event?
    case pointer.Event:
      // Update and store positioning state
      
    // A re-render request?
    case system.FrameEvent:
      // Layout the speech as a list 

    // Shutdown?
    case system.DestroyEvent:
      // Sayonara

    }
  }

```


All the source-code is in this repo, in the ```teleprompter/code``` folder.
