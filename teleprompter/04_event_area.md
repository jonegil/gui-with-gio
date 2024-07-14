---
layout: default
title: Chapter 4 - Event Area
nav_order: 2
parent: Teleprompter
has_children: false
---

# Chapter 4 - Event Area

## Goals
In this chapter we capture user input using a technique we'll call an `eventArea`.

## Outline
So far we've [set up](01_setup.md) our program, [reacted to user input](02_user_input.md) and [laid out](03_layout.md) our application on screed. The final piece to the puzzle is to actually capture keyboard and mouse input. 

We first create a **clip area** where events are to be collected from, and then specify which events to listen for. As always it's fairly straight forward, but we'll take the time to step through it

## Code

### pointer.Scroll

`InputOp` declares an input handler ready for pointer events. In our case we will listen for scrolling of the mousewheel or two fingers on a trackpad.

Since Gio is stateless we must **Tag** events, to make sure we later have enough information to know where they came from. 
Such a tag can anything really, so we simply use `Tag: 0`. Later we retireve these events with `gtx.Events(0)`.


```go
// Scrolled a mouse wheel?
for {
  ev, ok := gtx.Event(
    pointer.Filter{
      Target:  tag,
      Kinds:   pointer.Scroll,
      ScrollY: pointer.ScrollRange{Min: -1, Max: +1},
    },
  )
  if !ok {
    break
  }
  fmt.Printf("SCROLL: %+v\n", ev)
  scrollY = scrollY + unit.Dp(ev.(pointer.Event).Scroll.Y*float32(fontSize))
  if scrollY < 0 {
    scrollY = 0
  }
}
```

Since we listen for pointer operations, we can specify bounds on the scrolling to hinder the user from scrolling out of bounds. 
The important thing is to allow both positive and negative range, but it could easily be +/- 1 or +/- 100 or whatever. It can't be zero though, then no scrolls will happen.


### key.FocusOp

```go
// 2) Next we add key.FocusOp,
key.FocusOp{
  Tag: 0, // Use Tag: 0 as the event routing tag, and retireve it through gtx.Events(0)
}.Add(gtx.Ops)
```

For general keybaord input, we add `key.FocusOp`. These events are retrieved as `key.EditEvent`. Again we add `Tag: 0`. This is a generic catch-all where you don't want to specify the specific keys to listen for. 


### key.InputOp

```go
// 3) Finally we add key.InputOp to catch specific keys
key.InputOp{
  Keys: key.Set("(Shift)-F|(Shift)-S|(Shift)-U|(Shift)-D|(Shift)-J|(Shift)-K|(Shift)-W|(Shift)-N|Space"),
  Tag:  0, // Use Tag: 0 as the event routing tag, and retireve it through gtx.Events(0)
}.Add(gtx.Ops)
```

Finally we add a specific list of keyboard shortcuts to listen for. All are specified as Large Cap, and `(Shift)` means an optional
Shift can be included. These inputs are retrieved as `key.Event`. Here too a `Tag: 0` is used.

### Closing the event Frame

We finally `Pop()` the eventArea from the stack. 
```go
// Finally Pop() the eventArea from the stack
eventArea.Pop()
```

### Finalize

We have now completed the full frame and it's time to draw all operations from `&ops` on screen:

```go
// Frame completes the FrameEvent 
winE.Frame(&ops)
```

## Comments

Thank you so much, and well done my friend! High fives, fist bumps, back slaps and cheers all around. Thanks a lot for staying together on this tour of Gio's event handling. We've covered a lot of ground, but also built a really nifty little app that scales, scrolls, zooms and moves, all at the will of our fingertips. That's pretty neat.

Was this useful, please **star** the repo on Github, or even better, drop me a line. I really love hearing what people build with Gio. 


Thanks and all the best!

---

[View it on GitHub](https://github.com/jonegil/gui-with-gio/tree/main/teleprompter){: .btn .fs-5 .mb-4 .mb-md-0 }
