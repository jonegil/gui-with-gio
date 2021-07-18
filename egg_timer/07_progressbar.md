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

I've looked forward to writing this cpater, we will do many cool things here.
 - Introduce a new widget, the **material.Progressbar**
 - Start using state variables to control behaviour
 - Use two concurrency technques; one to create and share a beating pulse that progresses the progressbar, one one to select among independent communication operations

Let's look at these new pieces.

## Code - A beating pulse

A good progressbar must grow smoothly and precisely. To achieve that, we first create a separate go-routine that beats with a steady pulse. Then, when we listen for events, we pick up on these beats and grow the bar.

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

