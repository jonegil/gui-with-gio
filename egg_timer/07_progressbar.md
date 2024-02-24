---
layout: default
title: Chapter 7 - Progressbar
nav_order: 2
parent: Egg timer
has_children: false
---

# Chapter 7 - Progressbar

Updated February 24th 2024

## Goals

The intent of this section is to add a progressbar

![Progressbar](07_progressbar.gif)

## Outline

I've looked forward to this chapter ever since I started writing this series. We will cover quite some ground and introduce multiple new ideas:

- Try out a new widget, the `material.Progressbar`
- Start using state variables to control behaviour
- Use a concurrency technique to create and share a beating pulse that progresses the progressbar

Let's look at these in turn pieces.

## Feature 1 - The progressbar

A progressbar is obviously a bar that displays progress. But which progress? And how to control it? How fast should it grow, can it pause, or even reverse? From the [docs](https://pkg.go.dev/gioui.org/widget/material?utm_source=gopls#ProgressBar) we find `ProgressBar(th *Theme, progress float32)` receives progress as a decimal between 0 and 1.

### Code

We start by definding two variables for progress. The first is simply the progress as a number. We also define a channel used to send progress information through, which we'll look closer at later. Both are defined at root level, outside main, so that they are once and we have access to them throughout the whole program:

```go
// Define the progress variables, a channel and a variable
var progress float32
var progressIncrementer chan float32
```

To lay out the progressbar, we turn to our sturdy Flexbox and insert it through a rigid:

```go
// Inside System.FrameEvent
layout.Flex{
  // ...
}.Layout(gtx,
  layout.Rigid(
    func(gtx C) D {
      bar := material.ProgressBar(th, progress)  // Here progress is used
      return bar.Layout(gtx)
    },
  ),

```

Notice how the widget itself has no state. State is maintained in the rest of the program, the widget only knows how to display the progress we send it. Any logic to increase, pause, reverse or reset we control outside the widget.

## Feature 2 - State variables

We mentioned `progress`, the variable that contains progress state. Another useful state to track is whether or not the start button has been clicked. In our app that means tracking if the egg has started to boil. 

### Code

```go
// is the egg boiling?
var boiling bool
```
We want to flip that boolean when the start button is clicked. Thus we listen for a `app.FrameEvent` from the GUI and check if `startButton.Clicked()` is true.

```go
case system.FrameEvent:
  gtx := layout.NewContext(&ops, e)
  // Let's try out the flexbox layout concept
  if startButton.Clicked(gtx) {
    boiling = !boiling
  }
```

Again, the only job of the button is shout out if it recently was clicked. Beyond that, the rest of the program takes care of any actions that needs to be taken. 

One example is what the text on the button should be. We decide that before calling the `material.Button( )` function by first checking what the state of `boiling` is.

```go
// ...the same function we earlier used to create a button
func(gtx C) D {
  var text string
  if !boiling {
    text = "Start"
  } else {
    text = "Stop"
  }
  btn := material.Button(th, &startButton, text)
  return btn.Layout(gtx)
},
```

## Feature 3 - A beating pulse

A good progressbar must grow smoothly and precisely. To achieve that, we first create a separate go-routine that beats with a steady pulse. Later we listen for events, and pick up on these beats to grow the bar.

### Code

Here's the code, first the tick-generator:

```go
// Define the progress variables, a channel and a variable
var progress float32
var progressIncrementer chan float32

func main() {
  // Setup a separate channel to provide ticks to increment progress
  progressIncrementer = make(chan float32)
  go func() {
    for {
      time.Sleep(time.Second / 25)
      progressIncrementer <- 0.004
    }
  }()

  // ...
```

`progressIncrementer` is the [channel](https://tour.golang.org/concurrency/2) into which we send values, in this case of type `float32`.

Again, this is done in an anonymous function, called at creation, meaning this for-loop spins for the entirety of the program. Every 1/25th of a second the number 0.004 is injected into the channel.

Later we pick up from the channel, with this code inside `draw(w *app.window)`:

```go
// .. inside draw()
for {
    // listen for events
    switch e := w.NextEvent().(type) {
        // ...
    }

    // listen for events in the incrementor channel
    select {
    case p := <-progressIncrementer:
        if boiling && progress < 1 {
            progress += p
            w.Invalidate()
        }
    }
}
```

The first part of the loop is as before, evaluating events in the frame. At the end of the for-loop with use [select](https://tour.golang.org/concurrency/5).  This is a concurrency feature of go, where `select` waits patiently for an event that one of its `case` statement can run.

These combine so that
- Events happen in the window. These we extract with `e := <- w.NextEvent().(type)`.
- Events stem from the progress-pulse, and we get them using `p := <- progressIncrementer `

We add the `p` to the `progress` variable if the control variable `boiling` is true, and progress is less than 1. Since `p` is 0.004, and progress increased 25 times per second, it will take 10 seconds to reach 1. Feel free to adjust either of these two to find a combination of speed and smoothness that works for you.

Finally we force the window to draw, by calling `w.Invalidate()`. What is does is to inform Gio that the old rendering is now, well, invalid, and hence a new drawing must be made. Without such notice, Gio would simply not update until forced to do so by a mouse click or button press or other events. Invalidating at _every_ frame though can be costly, and alternatives exists. It's a bit of an advanced topic though, so for now let's leave it as is, but return to it in the [Bonus chapter on improved animation](11_improved_animation.md).

By using a channel like this we get

1. Precise timing, where we control the execution exactly as we want it
1. Consistent timing, similar across fast and slow hardware
1. Concurrent timing, the rest of the application continues as before

While all of these make sense, the 2nd point deserves some extra attention. If you recompile the app without the `time.Sleep(time.Second / 25)`, your machine will work hard to render as quickly as possible. That can consume a lot of cpu resources, which in turn can drain battery as well. It also ensures consistency across devices, all run at the same pulse. As an example, pprof's from 3 different machines are included in the code folder. These include a 1/25th sleep, ensuring the same end result. Please have a look.

## Comments

By combining all these building blocks we now have a stateful program we can control with ease. The user interface tells us when something happens, and the rest of the program uses that to take care of business. We had to pull a few tricks out of the bag, including both a `channel` and a `select`. Now that we have those tools in our belt, we will be well equipped to add some custom graphics in the next chapter.

---

[Next chapter](08_egg_as_circle.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[View it on GitHub](https://github.com/jonegil/gui-with-gio/tree/main/egg_timer){: .btn .fs-5 .mb-4 .mb-md-0 }
