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
	
### eventArea := clip.Rect

```go	
// ---------- COLLECT INPUT ----------
eventArea := clip.Rect(
  image.Rectangle{
    // From top left
    Min: image.Point{0, 0},
    // To bottom right
    Max: image.Point{gtx.Constraints.Max.X, gtx.Constraints.Max.Y},
  },
).Push(gtx.Ops)
```

We start by creating a **clip area** wihtin where we will listen to events. As we discussed in the [Egg Timer](egg_timer/08_egg_as_circle.md), the role of the Clip ares is to define an area we care about, and we ignore what happens outside. 

In this example we care about the full size of the screen, so naturally the `eventArea` extends from top left to bottom right. This is very useful if you want very precise control over areas that can generate input. If you want to experiment, try changig the constraints of image.Rectangle to, for example `Max: image.Point{300, 300}`. Take a guess what happens then. 

Finally, after first creating `clip.Rect()`, we `Push()` the clip to the stack.


### pointer.InputOp

`InputOp` declares an input handler ready for pointer events. In our case we will listen for scrolling of the mousewheel or two fingers on a trackpad.

Since Gio is stateless we must **Tag** events, to make sure we later have enough information to know where they came from. 
Such a tag can anything really, so we simply use `Tag: 0`. Later we retireve these events with `gtx.Events(0)`.


```go
// 1) We first add a pointer.InputOp to catch scrolling:
pointer.InputOp{
  Types: pointer.Scroll,
  Tag:   0, // Use Tag: 0 as the event routing tag, and retireve it through gtx.Events(0)
  // ScrollBounds sets bounds on scrolling, and we want it to be non-zero.
  // In practice it seldom reached 100, so [MinInt8,MaxInt8] or [-128,127] should be enough
  ScrollBounds: image.Rectangle{
    Min: image.Point{
      X: 0,
      Y: math.MinInt8, //-128
    },
    Max: image.Point{
      X: 0,
      Y: math.MaxInt8, //+127
    },
  },
}.Add(gtx.Ops)
```

Since we listen for pointer operations, we can specify bounds on the scrolling to hinder the user from scrolling out of bounds. 
The important thing is to allow both positive and negative range, and on my laptop the trackpad seldom reached 100. For pure 
programmers love I chose from -128 to +127, but it could easily be +/- 100. It can't be zero though, then no scrolls will happen.


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

Hey, well done! High fives, fist bumps, back slaps and cheers all around. Thanks a lot for staying together on this tour of Gio's event handling. We've covered a lot of ground, but also built a really nifty little app that scales, scrolls, zooms and moves, all at the will of our fingertips. That's pretty neat.

Was this useful, please **star** the repo on Github, or even better, drop me a line. I really love hearing what people build with Gio. 


Thanks and all the best!

---

[View it on GitHub](https://github.com/jonegil/gui-with-gio/tree/main/teleprompter){: .btn .fs-5 .mb-4 .mb-md-0 }
