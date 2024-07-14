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
First we listen for mouse or trackpad events. This is done by waiting for en event that matches a filter. 

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
The filter defines that this code will wait for a `pointer.Scroll` of a given range. The range can be anything really, but the most important thing is that it containts positive and negative range. Apart from that it could easily be +/- 1 or +/- 100 or whatever, but check your device sensitivity. It can't be zero though, then no scrolls will happen.

### pointer.Press

Here we filter for `pointer.Press`, i.e. clicking down with the mouse.

```go
// Pressed a mouse button?
for {
  ev, ok := gtx.Event(
    pointer.Filter{
      Target: tag,
      Kinds:  pointer.Press,
    },
  )
  if !ok {
    break
  }
  fmt.Printf("PRESS : %+v\n", ev)
  // Start / stop
  autoscroll = !autoscroll
}
```


### key.Filter

This filter explicitly filters the keys relevant for us. Some are normal letters while others are named keys, such as arrows.

```go
for {
  ev, ok := gtx.Event(
    key.Filter{Name: key.NameSpace},
    key.Filter{Optional: key.ModShift, Name: "U"},
    key.Filter{Optional: key.ModShift, Name: "D"},
    key.Filter{Optional: key.ModShift, Name: "J"},
    key.Filter{Optional: key.ModShift, Name: "K"},
    key.Filter{Optional: key.ModShift, Name: key.NameUpArrow},
    key.Filter{Optional: key.ModShift, Name: key.NameDownArrow},
    key.Filter{Optional: key.ModShift, Name: key.NamePageUp},
    key.Filter{Optional: key.ModShift, Name: key.NamePageDown},
    key.Filter{Optional: key.ModShift, Name: "F"},
    key.Filter{Optional: key.ModShift, Name: "S"},
    key.Filter{Optional: key.ModShift, Name: "+"},
    key.Filter{Optional: key.ModShift, Name: "-"},
    key.Filter{Optional: key.ModShift, Name: "W"},
    key.Filter{Optional: key.ModShift, Name: "N"},
    key.Filter{Optional: key.ModShift, Name: "C"},
  )
```

For each key we then define what is to happen when it is pressed. For example
```go
  // Faster scrollspeed
  if name == "F" {
    autoscroll = true
    autospeed += stepSize
  }

  // Slower scrollspeed
  if name == "S" {
    if autospeed > 0 {
      autospeed -= stepSize
    }
    if autospeed <= 0 {
      autospeed = 0
      autoscroll = false
    }
  }
```


### Finalize

We have now completed the full frame and it's time to draw all operations from `&ops` on screen:

```go
// Frame completes the FrameEvent 
winE.Frame(&ops)
```

## Comments

Thank you so much, and well done my friend! High fives, fist bumps, back slaps and cheers all around. Thanks a lot for staying together on this tour of Gio's event handling. We've covered a lot of ground, but also built a really nifty little app that scales, scrolls, zooms and moves, all at the will of our fingertips. That's pretty neat.

If this was useful, please **star** the repo on Github, or even better, drop me a line. I really love hearing what people build with Gio. 


Thanks and all the best!

---

[View it on GitHub](https://github.com/jonegil/gui-with-gio/tree/main/teleprompter){: .btn .fs-5 .mb-4 .mb-md-0 }
