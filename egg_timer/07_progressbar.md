---
layout: default
title: Chapter 7 
nav_order: 2
parent: Egg timer
has_children: false
---

# Chapter 7 - Progressbar

## Goals
The intention of this section is to add a progressbar

![Progressbar](07_progressbar.gif)

## Outline

I've looked forward to writing this cpater, we will introduce many new features here.
 - Try out a new widget, the **material.Progressbar**
 - Start using state variables to control behaviour
 - Use two concurrency technques; one to create and share a beating pulse that progresses the progressbar, one one to select among independent communication operations

Let's look at these in turn pieces.

## Feature 1 - The progressbar

The progressbar is obviously a bar that displays progress. But which progress? And how to control it? How fast should it grow, can it pause, or even reverse? To answer these questions we need a variable that quantifies the current progress. We set it at root level, outside main, so that its set once and we have access to it throughout the whole program

### Code

```go
  // root level, outside main ()
  var progress float32
```

Later on we will look at the logic used to actually set the progress, but suffice to say it needs to be a decimal number between 0 and 1.

To lay out the progressbar, we turn to our sturdy Flexbox and insert it through a rigid:
```go
// Inside System.FrameEvent
layout.Flex{
  // ...
}.Layout(gtx,
  layout.Rigid(
    func(gtx C) D {
      bar := material.ProgressBar(th, progress)  // Here progress is used for display
      return bar.Layout(gtx)
    },
  ),

```

Again we see how the widget is stateless. It knows how to display a progress, but is itself not responsible for what that progress is, not responsible for updating, incrementing, pausing, or progressing it. Pure display, the mission is controled outside the widget. 


## Feature 2 - State variables

Another very useful variable is whether or not the start button has been clicked. In our case, is the egg boiling? 

### Code

```go
	// is the egg boiling?
	var boiling bool
```

This boolean is flipped whenever the button is clicked, we listen for that **system.FrameEvent**

```go
case system.FrameEvent:
  gtx := layout.NewContext(&ops, e)
  // Let's try out the flexbox layout concept
  if startButton.Clicked() {
    boiling = !boiling
  }
```

The only job of the button is to answer yay or nay to weither or not it was just clicked. Beyond that, the rest of the program takes care of any updates that are needed. Among those are for example what the text on the button should be. Here's how that's done, coded where we create and display the button:

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

A good progressbar must grow smoothly and precisely. To achieve that, we first create a separate go-routine that beats with a steady pulse. Then, when we listen for events, we pick up on these beats and grow the bar.


### Code

Here's the code, first the tick-generator:

```go
// Define the progress variables, a channel and a variable
var progressIncrementer chan float32
var progress float32

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

```progressIncrementer``` is the [channel](https://tour.golang.org/concurrency/2) into which we send values, in this case of type ```float32```. 

Again, this is done in an anonymous function, called at creation, meaning this for-loop spins for the entirety of the program. Every 1/25th of a second the number 0.004 is injected into the channel. 

Later we pick up from the channel, with this code inside **draw(w *app.window)**:

```go
  // .. inside draw()
	for {
    select {
      // listen for events in the window.
      case e := <-w.Events():
        // ...    

      // listen for events from the incrementor channel
      case p := <-progressIncrementer:
        if boiling && progress < 1 {
        progress += p
        w.Invalidate()
			}
		}
	}

```

In previous chapters, we ranged over events using ```for e := range w.Events() { }```. Here we insted use a for-loop with a [select](https://tour.golang.org/concurrency/5) inside. The select waits patiently for an event one of its cases can run, then executes that case. 
 - The event can either stem from the window, and if so we extract it using ```e := <- w.Events()```. 
 - Or, the event comes from the progress-pulse, and we get it from ```p := <- progressIncrementer ```

We add the ```p``` to the ```progress``` variable if the control variable ```boiling``` is true, and progress is less than 1. This caps progress at 1, and since it increases by 0.004 every 1/25th of a second, that will take 10 seconds. 

By using the progressIncrementor channel like this we get
1. Precise timing, where we control the execution exactly as we want it
1. Consistent timing, simlar across fast and slow hardware
1. Concurrent timing, the rest of the application continues as before

While all of these make sense, the 2nd point deserves some mention. If you recompile the app without the ```time.Sleep(time.Second / 25)``` (and probably a much smaller increment), your machine will work it's socks off, spinning the loop at insane speeds. That can max the cpu, drain battery, but will not be consistent across machines. For those interested, pprof's from 3 different machines are included in the code folder. These include a 1/25th sleep though, but still work to show the differenences across architectures.

### Code
```go

```

## Comments

