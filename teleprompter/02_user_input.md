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

Simplified, the code looks like this:

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

An application can have many different visual areas. When Gio sends a pointer-event telling us a click has happened, it's useful to know *where* it happened. Therefore we will later define specifc *eventAreas* on screen, and give each one it's unique name, to make sure a click is clearly defined. For our app we will only have one area, so `Tag: 0` works just fine. It could have been any number, or a bool, or even a pointer address, but numbers are easy to work with - hence `0`. We'll revisit this one later. 

**Point no 3**

The code example is, by intention, verbose. To make sure it's very clear when we work with events from the window, `windowEvent`, vs when we're working with a `FrameEvent` and its context `gtxEvent`, longer variable names are used. Also, since we use [type switches](https://go.dev/tour/methods/16), which are very handy but also a bit compact, it was helpful to be very explicit with `winE` and `gtxE` to make it as clear as possible. 

However, this is a bit out of the ordinary. If you read sourcecode from mature applications you will often all event names as `e`, such as `e := range` followed by `e := e.(type)`. That is fine, and we did that when boiling eggs too. However, for this tutorial, it was useful to separate into more explicit variables, which hopefully helps understanding the code the a little easier.


---

Thanks for sticking with it - I know it can be a lot to take in. Finally though, it's time investigate the core of the event handling. In other words, it's time to process the queue of events inside `FrameEvent`. Here's the structure:

```go
    case system.FrameEvent:

      gtx := layout.NewContext(&ops, winE)
      for _, gtxEvent := range gtx.Events(0) {

        switch gtxE := gtxEvent.(type) {

        case key.EditEvent:
  
        case key.Event:
          
        case pointer.Event:
        }
      }

```

What's happening here?
- We range through `gtx.Events()` with `Tag: 0`
- We check their type with `gtxEvent.(type)` and pick up three variants:
  - `key.EditEvent`: These are free edit events, when the user writes freely. All letters are caught.
  - `key.Event`: These are specifically defined events. Only explicitly declared keys are caught.
  - `pointer.Event`: Mouse/trackpad events such as clicks and scrolling.


## key.EditEvent
Let's look at each in turn. First the `key.EditEvent`

```go
  case key.EditEvent:
    // To increase the fontsize
    if gtxE.Text == "+" {
      fontSize = fontSize + unit.Sp(stepSize)
    }
    // To decrease the fontsize
    if gtxE.Text == "-" {
      fontSize = fontSize - unit.Sp(stepSize)
    }
```
The `key.EditEvent` is sent whenever the user types a key, and you get the content of the event with `gtxE.Text`.

While a normal letter (a,b,c...) is the same on all keyboards, the `+` or `-` is not. They are placed differently on US vs non-US keyboards, and also differently on laptop keybaords vs full size keyboards with numerical keys on the right. To ensure all pluses are treated the same, no matter which key was used to produce it, it was useful to catch then through this generic text event. And the same for minuses ofcourse. Fair treatment after all.

## key.Event
Now my favourite, `key.Event`:

When one of the specified keys are pressed, Gio receives it as a [key.Event](https://pkg.go.dev/gioui.org/io/key#Event). As we see from the docs, the Event is a struct with three variables, `Name`, `Modifiers` and `State`:

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
- `Modifiers` are keys like `key.ModShift` or `key.ModCommand`, listed [here](https://pkg.go.dev/gioui.org/io/key#Modifiers). Note the comment on how Shift is taken into account, but not others, which can be worth knowing about.
- `State` can be either Press or Release

Here are the ones we use:
```go
  case key.Event:
    // For better control, we only care about pressing the key down, not releasing it up
    if gtxE.State == key.Press {
      // Inrease the stepSize when pressing Shift
      if gtxE.Modifiers== ModShift {
        stepSize = 5
      }
      // Start/Stop
      if gtxE.Name == "Space" {
        autoscroll = !autoscroll
        if autospeed == 0 {
          autoscroll = true
          autospeed++
        }
      }
      // Scroll up
      if gtxE.Name == "K" {
        scrollY = scrollY - stepSize*4
        if scrollY < 0 {
          scrollY = 0
        }
      }
      // Scroll down
      if gtxE.Name == "J" {
        scrollY = scrollY + stepSize*4
      }
      // Faster scrollspeed
      if gtxE.Name == "F" {
        autoscroll = true
        autospeed++
      }
      // Slower scrollspeed
      if gtxE.Name == "S" {
        if autospeed > 0 {
          autospeed--
        }
        if autospeed == 0 {
          autoscroll = false
        }
      }
      // Wider text to be displayed
      if gtxE.Name == "W" {
        textWidth = textWidth + stepSize*10
      }
      // Narrow text to be displayed
      if gtxE.Name == "N" {
        textWidth = textWidth - stepSize*10
      }
      // Move the focusBar Up
      if gtxE.Name == "U" {
        focusBarY = focusBarY - stepSize
      }
      // Move the focusBar Down
      if gtxE.Name == "D" {
        focusBarY = focusBarY + stepSize
      }
    }

```

The role of `stepSize` is to control how large the change to the other parameters will be. Should a scroll be long or short? Should the focus bar move by lot or a little? Should width-adjustments be considerable or minor? Should ... you get it.

The point is that for a user it can sometimes be important to quickly navigate or adjust quite quickly, and thereafter finetune to perfection. Therefore it's useful to define a variable that controls the rate of change. For simplification this was skipped earlier, but the `stepSize unit.Dp = 1` is actually defined when handling `gtx.Events(0)`. That ensures a change in **D**evice indepentend **p**ixels, making it homogenous across displays. 
```go
  for _, gtxEvent := range gtx.Events(0) {
    // To set how large each change is
    var stepSize unit.Dp = 1
```
When holding `Shift` this increases to 5. Why 5? Well, it worked well in my experimentation. Try it out.

For all the other keypresses, the code adjusts one or two state variables. These are all used later when rendering the actual frame. I went a bit back and forth on the logic, and pondered how it all should interact, but landed on a fairly safe set of behaviours. For exampe, negative speed is avoided, although it was fun times before it was nerfed.

The point is that when defining behaviour, it makes sense to experiment and think through how the various state variables should be tuned in relation to each other. Keeping it all togehter in this input section made the code easier to grasp than if states were handled in various other parts of the program.

There are two reasons why I like the `key.Event`. Firstly, it's because the event is explicit. Only specificly defined keys trigger this event. We'll get to that code later. But the point is that by being strict on which key-strokes are allowed to trigger an event, it is very clear how the event was generated, and also which modifiers (Shift / Ctrl / Command) is allowed with it. 

Secondly, it allowes for fine-tuned control of both press and release. Gio is all about control, and `key.Event` gives you exactly that.

Your preferences may differ, and that's fair, but for me the extra code needed to define which keys generate a `key.Event` is a small investment for the control and precision it yields. 

![Size adjustments](teleprompter_fontsize.gif)


## pointer.Event

If the mouse is used, Gio receives it as a pointer.Event. There are many variants, such as movement, scrolling or clicking. Once we detect with `case pointer.Event:` it is up to us to decide what to do with it.

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
  // Are we scrolling?
  if gtxE.Type == pointer.Scroll {
    if gtxE.Modifiers == key.ModShift {
      stepSize = 3
    }
    // Increment scrollY with gtxE.Scroll.Y
    scrollY = scrollY + unit.Dp(gtxE.Scroll.Y)*stepSize
    if scrollY < 0 {
      scrollY = 0
    }
  }
```

As with keys we listen for certain events, in this case only the `pointer.Scroll`. We want to scroll faster if `Shift` is pressed, and it worked well by increasing it by x3 here.

Since scrollY is a `unit.Dp`, we cast `gtxE.Scroll.Y` and increment the state variable `scrollY` with it. This will control how far down into the text we should present to the user. To reduce confusion we disallow scrolling to before the start by limiting `scrollY` to minimum 0.

## Closing remarks

Phew, that was a long one. We covered a lot of event handling, but althoug it's fairly long I hope it's still clear. Most important now is to understand the overall structure and how the pieces join together. 

With event *handling* under our belt, we still need to cover how events are *created*. That will come on [Chapter 4 - The Event Area](04_event_area.md). However, before we get there, let's have some fun and actually use all the state variables we made to create a visual experience for our users. That's all done in [Chapter 3 - Layout](03_layout.md).

---

[Next chapter](03_layout.md){: .btn .fs-5 .mb-4 .mb-md-0 }
