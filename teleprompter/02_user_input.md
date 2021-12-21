---
layout: default
title: Chapter 2 - User input
nav_order: 2
parent: Teleprompter
has_children: false
---

# Chapter 2 - User input

In this chapter we look at how we can receive and process user input to control our prompter.

## Listen for events from keyboard and mouse

Now it's time get to listen for events. This is the heart of the application. As mentioned above, there are quite a few inputs here, with the various keys and also the use of the mouse. In this application, these can mutually impact each other. For example, if `textWdith` increases, more words can be shown per line since there is now space. But if `fontSize` increases, each word requires more space and fewer words can be shown. Luckily for us Gio takes care of all of the underlying mechanics, our job is the keep track of the required state variables used to define the visualisation.

As before the switch statement uses type assertion, `e.(type)` to deterimine what just happened:

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
	Position f32.Point
	// Scroll is the scroll amount, if any.
	Scroll f32.Point
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
