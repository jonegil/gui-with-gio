---
layout: default
title: Teleprompter
nav_order: 3
has_children: false
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

![Mr_Gorbachev_tear_down_this_wall](teleprompter_Mr_Gorbachev.gif)
[View it on GitHub](https://github.com/jonegil/gui-with-gio/tree/main/teleprompter){: .btn .fs-5 .mb-4 .mb-md-0 }

To do that we'll build what's known as a [teleprompter](https://en.wikipedia.org/wiki/Teleprompter). A teleprompter is simply a device that displays and scrolls text. Sophisticated and expensive equipment exists, but it can just as easily be done with an app that displays and scrolls text. And that's the version we will build here. Since it needs to be lively and responsive for the user, it it's a great example for how we can react to keypresses and mouse scrolls. We'll make sure to look into some other new parts of Gio as well.

---
**Please note:** We're not doing politics here, but when presenting this piece, it's more fun with famous speeches. And most famous speaches are given by, well, famous politicians. No matter your origin or partisanship, [Ronald Reagan's 1987 Berlin Wall Speech](https://www.nytimes.com/video/world/europe/100000006815729/reagan-berlin-wall.html) is among those, playing it's part in ending the cold war without the massive bloodshed the world feared. Both sides, East and West, deserve respect and praise for seeking and finding a mostly peaceful solution. Let us together repeat the best of our past, and learn from the worst. The future needs it.

---

On that note, let's (sc)roll!
(sorry)

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

To structure the walkthrough of the code, it's broken into the following main sections:

1. Introduce new imports to handle user input
1. Read the ```.txt``` file into a ```[]string``` slice
1. Start the application
1. Define state variables to control behaviour
1. Listen for events from the user. 

Of these, the four first are relatively straight forward, while the final one on events deserves some extra attention. That's where the we actually will deal with the various inputs from the user, and visualise the application. 

### Section 1 - New imports

Many imports are well known, but some are new:

```go
import (
  // Many normal imports we discussed earlier ...
  // ... plus two new interesting Gio imports
  "gioui.org/io/key"
  "gioui.org/io/pointer"
)
```
These two are new to us and gives support for keyboard and mouse events:
 - Package [io/key](https://pkg.go.dev/gioui.org/io/key) implements key and text events and operations.
 - Package [io/pointer](https://pkg.go.dev/gioui.org/io/pointer) implements pointer events and operations. A pointer is either a mouse controlled cursor or a touch object such as a finger.

Notice how pointer supports both mouse gestures on a desktop/laptop and fingers on a screen. Nice, again an example of how learning a cross-platform framework gives you tools to multiple devices.

### Section 2 - Read the speach into a slice

First we define variables for the program, inlcuding a slice to keep the speech in.

```go
// A []string to hold the speech as a list of paragraphs
var paragraphList []string
```

In the first part of main we actually read the speech from ```speech.txt```. The text in that file will be parsed and stored in a slice of strings. That's done inside ```main()```:

```go
func main() {
  // Read from file
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
```
The first sectioin of ```main``` reads the text and splits it by ```\n```, newline, returning the result as a ```[]string```.  

To allow enough space after the line so that it actually scrolls off screeen, we simply add a handful of empty paragraphs at the end of the list. 

### Section 3 - Start the application 

The last section of ```main``` starts the Gui in a normal manner:
```go
  // ... continuing inside main()
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

### Section 4 - Variables to control behaviour

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

 The state variables in play here are:

 |Variable        |Description                                       | Changed with                              |
 |---             |---                                               |---                                        |
 |```scrollY```   | Scroll the text                                  | Mouse/Trackpad scroll, Arrow Up/Down, J/K |
 |```focusBarY``` | How high up is the red focus bar                 | U (up) and D (down)                       |
 |```textWidth``` | How wide is the area in which we display text    | W (wider) and N (narrower)                |
 |```fontSize```  | How large is the text                            | + (larger) and - (smaller)                |
 |```autoscroll```| Start/stop automatic scrolling                   | Space                                     |
 |```autospeed``` | How fast / slow the text should scroll           | F (faster) or S (slower)                  |
 
For keypresses, ```Shift``` increases the rate of change when making adjustments

### Section 5 - Listen for events

Finally, we get to listen for events. As mentioned above, there are quite a few inputs here, with the various keys and also the use of the mouse. In this application, these can mutually impact each other. For example, if ```textWdith``` increases, more words can be shown per line since there is now space. But if ```fontSize``` increases, each word requires more space and fewer words can be shown. Luckily for us Gio takes care of all of the underlying mechanics, our job is the keep track of the required state variables used to define the visualisation. 

As before the switch statement uses type assertion, ```e.(type)``` to deterimine what just happened:

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
    // Break out and end

  }
}
```

The two new events here are:
 - ```key.Event``` - Was a key just pressed? 
 - ```pointer.Event``` - Was a mouse or trackpad just used?

So let's go through those in more detail:

#### key.Event
If a key is pressed, Gio receives it as a [key.Event](https://pkg.go.dev/gioui.org/io/key#Event). As we see from the docs, the Event is a struct with three variables, ```Name```, ```Modifiers``` and ```State```:

```go
type Event struct {
  // Name of the key. For letters, the upper case form is used, via
  // unicode.ToUpper. The shift modifier is taken into account, all other
  // modifiers are ignored. For example, the "shift-1" and "ctrl-shift-1"
  // combinations both give the Name "!" with the US keyboard layout.
  Name string
  // Modifiers is the set of active modifiers when the key was pressed.
  Modifiers Modifiers
  // State is the state of the key when the event was fired.
  State State
}
```
- ```Name``` is simply the letter pressed, or [special keys](https://pkg.go.dev/gioui.org/io/key#pkg-constants) such as ```key.NameUpArrow``` and ```key.NameSpace```
- ```Modifiers``` are keys like ```key.ModShift``` or ```key.ModCommand```, listede [here](https://pkg.go.dev/gioui.org/io/key#Modifiers). Note the comment on how Shift is taken into account, but not others, which can be worth knowing about. 
- ```State``` can be either Press or Release, if the distinction is needed

Ok, that gives us something to work with. Once a key is pressed, this will help us detect which key it was, and weither a modifier like Shift is pressed. Here's the code for this section:

```go
// A keypress
case key.Event:
  if e.State == key.Press {
    // To set increment
    var stepSize int = 1
    if e.Modifiers == key.ModShift {
      stepSize = 10
    }
    // To scroll text down
    if e.Name == key.NameDownArrow || e.Name == "J" {
      scrollY = scrollY + stepSize*4
    }
    // To scroll text up
    if e.Name == key.NameUpArrow || e.Name == "K" {
      scrollY = scrollY - stepSize*4
      if scrollY < 0 {
        scrollY = 0
      }
    }
    // To turn on/off autoscroll, and set the scrollspeed
    if e.Name == key.NameSpace {
      autoscroll = !autoscroll
      if autospeed == 0 {
        autoscroll = true
        autospeed++
      }
    }
    // Faster scrollspeed
    if e.Name == "F" {
      autoscroll = true
      autospeed++
    }
    // Slower scrollspeed
    if e.Name == "S" {
      if autospeed > 0 {
        autospeed--
      }
    }
    // Set Wider space for text to be displayed
    if e.Name == "W" {
      textWidth = textWidth + stepSize*10
    }
    // Set Narrower space for text to be displayed
    if e.Name == "N" {
      textWidth = textWidth - stepSize*10
    }
    // To increase the fontsize
    if e.Name == "+" {
      fontSize = fontSize + stepSize
    }
    // To decrease the fontsize
    if e.Name == "-" {
      fontSize = fontSize - stepSize
    }
    // Move the focusBar Up
    if e.Name == "U" {
      focusBarY = focusBarY - stepSize
    }
    // Move the focusBar Down
    if e.Name == "D" {
      focusBarY = focusBarY + stepSize
    }
    // Force re-rendering to use the new states set above
    w.Invalidate()
  }
```

With the expception of ```stepSize``` all these variables are explained earlier. The role of ```stepSize``` is to control how large the change to the other parameters will be. Should a scroll be long or short? Should the focus bar move by lot or a little? Should text size adjustments be considerable or minor? Should ... you get it. 

The point is that for a user it can sometimes be important to quickly navigate or adjust quite quickly, and thereafter finetune to perfection. Therefor it's useful to define a variable that controls the rate of change. This defaults to 1, but when ```Shift``` is pressed it increases to 10. Why those value? Well, it worked well in my experimentation. Try it out.

For all the other keypresses, the code adjusts one or two state variables. These are all used later when rendering the actual frame. I went a bit back and forth on the logic around adjusting speed, but conlcuded that if you ask for ```F```aster scrolling, that should start up the autoscroll if it wasn't running already. Conversely, if speed is 0 and the user presses ```Space``` to start the scroll, speed must increase. Negative speed is avoided, although it was fun times before I nerfed it. 

The point is that for interacting behaviour, it makes sense to experiemnt and think through how the various state variables should be tuned in relation to each other. Keeping it all togehter in this input section makes the code easier to grasp than if these states were handled in various other parts of the program. 

Finally, at the end we call ```w.Invalidate()```, forcing the program to re-render so that any new state information is take into account at once. Try commenting this out and re-run. What happens now, and why?

#### pointer.Event
If the mouse is used, Gio receives it as a pointer.Event. That can be any type, such as movement, scrolling or clicking. Once we detect with ```case pointer.Event:``` it is up to us to define what to do with it. Here's the code inside that case:

```go
//TODO DESCRIBE POINTER EVENT CODE HERE
```

#### system.FramveEvent
If a request to re-render is sent, typically from a call to ```invalidate```, program redraws. The layout deserves it's own section though

```go
//TODO DESCRIBE LAYOUT HERE
```

### TODO - Describe the layout