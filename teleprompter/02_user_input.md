---
layout: default
title: Chapter 2 - User input
nav_order: 2
parent: Teleprompter
has_children: false
---

# Chapter 2 - User input

In this chapter we receive and process user input.

## Listen for events from keyboard and mouse

Listening for events is the heart of this application. As mentioned earlier, there are quite many inputs here, with various keys as well as scrolling with a mouse or trackpad. In this application, the various changes can affect each other. For example, if `textWdith` increases, more words can be shown per line since there is now space. But if `fontSize` increases, each word requires more space and therefore fewer words can fit. Luckily for us Gio takes care of the underlying mechanics, but we're in charge of receiving input and telling Gio what to do.

Let's start by walking through the structure of our `draw()` function:
  1. Listen for events in the window using `w.Events()`
  1. We're especially interested in events that require us to redraw a new frame. Those are called `system.FrameEvent`
  1. When those arrive, we open it's event queue on range through all events since last frame.

Simplified in code we see this:

```go
// The main draw function
func draw(w *app.Window) error {

  // Listen for events in the window
  for windowEvent := range w.Events() {
		switch winE := windowEvent.(type) {

		// Should we draw a new frame?
		case system.FrameEvent:
      // Open an new context
			gtx := layout.NewContext(&ops, winE)

      // Procress all events from eventArea 0 (to be exlpained)
	    for _, gtxEvent := range gtx.Events(0) {
				switch gtxE := gtxEvent.(type) {

          // ... process the events depending on their type, such as pointer or key for example

```
**Point no 1**
- OK, we recognice `draw()`, and also how we go through events from the window using `range w.Events()`. 
- The `system.FrameEvent` is one such event from the window - the one Gio sends when a new frame needs to be drawn. Fair enough. 
- But before we actually draw the frame, we investigate if something interesting has happened since last frame. That could be for example a pointer scroll or a keystroke. To do that we 
  - Creating a `gtx := NewContext()` and 
  - Open it's queue of events since last frame with `gtx.Events(0)`.

**Point no 2**
But wait. What´s this zero? 

An application can have many different visual areas. When Gio sends an event that a pointer-click has happened, it's kind of useful to know *where* it happened. Therefore we will later define specifc *areas* on screen, and give each it's unique name, to make sure a click is clearly defined. For our app we will only have one area, so `Tag: 0` works just fine. It could have been any number, or a bool, or even a pointer address, but numbers are easy to work with - hence `0`. We'll revisit this one later. 

**Point no 3**
The code example is, by intention, verbose. To make sure it's very clear when we work with events from the window, `windowEvent`, vs when we're working with a FrameEvent and it's context `gtxEvent`, long verbose names are used. Also, since we use [type switches](https://go.dev/tour/methods/16), which are very handy but also a but compact, it was helpful to be very explicit with `winE` and `gtxE`. When reading mature applications you will often see this simplified to `e := range` followed by `e := e.(type)`. That´s all fine, and we did that when boiling eggs too. Here it simply was useful to separate them into more explicit variables, and hopefully that helps understanding the a little easier.

---
WORK IN PROGRESS 

THE REST OF THE TUTORIAL IS REMNANTS FROM AN EARLIER VERSION OF GIO. CODE RUNS WELL, BUT THESE DOCS ARE LAGGING. 

MOST OF THE EXPLANATION MAKES SENSE THOUGH


DEC 16, 2022

https://lists.sr.ht/~eliasnaur/gio/%3CCAFcc3FQNTp_UXr7oA97SsVPD7D91jSw30ZtALcT9vmopFDTeZQ%40mail.gmail.com%3E#%3CCJ5LZODOOR0F.UO9JC0VWAN9I@themachine%3E
https://go.dev/play/p/VDQg6sxRyA4
https://go.dev/play/p/SDHy1LZRljf

---



        // A certain letter?
				case key.EditEvent:
          // Update and store state for size

        // A certain key?
        case key.Event:
          // Update and store state forwidth and positioning

        // A mouse event?
        case pointer.Event:
          // Update and store positioning state


  }
}
```

The two new events here are:

- `key.Event` - Was a key just pressed?
- `pointer.Event` - Was a mouse or trackpad just used?

So let's go through those in more detail:

## key.Event

If a key is pressed, Gio receives it as a [key.Event](https://pkg.go.dev/gioui.org/io/key#Event). As we see from the docs, the Event is a struct with three variables, `Name`, `Modifiers` and `State`:

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

- `Name` is simply the letter pressed, or [special keys](https://pkg.go.dev/gioui.org/io/key#pkg-constants) such as `key.NameUpArrow` and `key.NameSpace`
- `Modifiers` are keys like `key.ModShift` or `key.ModCommand`, listede [here](https://pkg.go.dev/gioui.org/io/key#Modifiers). Note the comment on how Shift is taken into account, but not others, which can be worth knowing about.
- `State` can be either Press or Release, if the distinction is needed

Ok, that gives us something to work with. Once a key is pressed, this will help us detect which key it was, and weither a modifier like Shift is pressed. Here's the code for this section:

```go
// A keypress?
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

With the expception of `stepSize` all these variables are explained earlier. The role of `stepSize` is to control how large the change to the other parameters will be. Should a scroll be long or short? Should the focus bar move by lot or a little? Should text size adjustments be considerable or minor? Should ... you get it.

The point is that for a user it can sometimes be important to quickly navigate or adjust quite quickly, and thereafter finetune to perfection. Therefor it's useful to define a variable that controls the rate of change. This defaults to 1, but when `Shift` is pressed it increases to 10. Why those value? Well, it worked well in my experimentation. Try it out.

For all the other keypresses, the code adjusts one or two state variables. These are all used later when rendering the actual frame. I went a bit back and forth on the logic around adjusting speed, but conlcuded that if you ask for `F`aster scrolling, that should start up the autoscroll if it wasn't running already. Conversely, if speed is 0 and the user presses `Space` to start the scroll, speed must increase. Negative speed is avoided, although it was fun times before I nerfed it.

The point is that for interacting behaviour, it makes sense to experiemnt and think through how the various state variables should be tuned in relation to each other. Keeping it all togehter in this input section makes the code easier to grasp than if these states were handled in various other parts of the program.

At the end we call `w.Invalidate()`, forcing the program to re-render so that any new state information is take into account at once. Try commenting this out and re-run. What happens now, and why?

With this in place, here's an example of how it looks to change fontsize:

![Size adjustments](teleprompter_fontsize.gif)

## pointer.Event

If the mouse is used, Gio receives it as a pointer.Event. That can be any type, such as movement, scrolling or clicking. Once we detect with `case pointer.Event:` it is up to us to define what to do with it.

From [pointer.Event](https://pkg.go.dev/gioui.org/io/pointer#Event) we learn that the pointer event is quite a rich struct:

```go
type Event struct {
	Type   Type
	Source Source
	// PointerID is the id for the pointer and can be used
	// to track a particular pointer from Press to
	// Release or Cancel.
	PointerID ID
	// Priority is the priority of the receiving handler
	// for this event.
	Priority Priority
	// Time is when the event was received. The
	// timestamp is relative to an undefined base.
	Time time.Duration
	// Buttons are the set of pressed mouse buttons for this event.
	Buttons Buttons
	// Position is the position of the event, relative to
	// the current transformation, as set by op.TransformOp.
	Position image.Point
	// Scroll is the scroll amount, if any.
	Scroll image.Point
	// Modifiers is the set of active modifiers when
	// the mouse button was pressed.
	Modifiers key.Modifiers
}
```

What we need here are the two bottom entries, `Scroll` and `Modifiers`. The former returns a `Point`, which is simply a set of X and Y variables that indicate how far the user scrolled in those directions:

```go
type Point struct {
	X, Y float32
}
```

With a scroll-wheel on a mouse it's always Y only and often in fixed clicking amounts. On a laptop trackpad however it can often be both, and with various amounts.

`Modifiers` are just as for the `key.Event` a helper to indicate if `Shift` or `Alt` or any of those are pressed when the mouse event occurs. We'll continute to listen for the former of those. Like this:

```go
// A mouse event?
case pointer.Event:
  if e.Type == pointer.Scroll {
    var stepSize int = 1
    if e.Modifiers == key.ModShift {
      stepSize = 3
    }
    // By how much should the user scroll this time?
    thisScroll := int(e.Scroll.Y)
    // Increment scrollY with that distance
    scrollY = scrollY + thisScroll*stepSize
    if scrollY < 0 {
      scrollY = 0
    }
    w.Invalidate()
  }
```

As with keys we listen for certain events, in this case only the `pointer.Scroll`. We want to scroll faster if `Shift` is pressed, but the stepSize of 10 from `key.Event` proved excessive. So we're content by increasing it by x3 this time.

After some manipulations, the Y value of a scroll is added to the state variable `scrollY` which indicates how far down into the speech we have reached. To reduce confusion we disallow scrolling to before the start by limiting `scrollY` to minimum 0

And just as for `key.Event` we end by invalidating the frame. Show it to me!

---

[Next chapter](03_layout.md){: .btn .fs-5 .mb-4 .mb-md-0 }
