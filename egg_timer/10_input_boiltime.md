---
layout: default
title: Chapter 10 
nav_order: 2
parent: Egg timer
has_children: false
---

# Chapter 10 - Set the boiltime

## Goals
The intention of this section is to add an input field to set the boiltime.

![The complete egg timer](egg_timer.gif)

## Outline

The code changes in a few ways
 1. Import [gioui.org/text](https://pkg.go.dev/gioui.org/text) from Gio as well as string and number manipulation from the standard library.
 1. Add a new rigid to hold the [widget.Editor()](https://pkg.go.dev/gioui.org/widget#Editor)
 1. Add some logic for the button to make it a little more well behaved. 

That's it. Let's look at the code:

## Code

### 1. New imports
```go
import (
  "fmt"
  "strconv"
  "strings"

  "gioui.org/text"
)
```
The list of imports is looking long now. But we're not importing the whole kitchen sink of Go here. It's more a testament of how useful functionality coming together across multiple parts of the standard library: 

 - ```fmt``` will be used to convert float to string
 - ```strconv```will be used to convert string to float
 - ```strings``` will be used to trim spaces away from the input string

And finally:
 - [gioui.org/text](https://pkg.go.dev/gioui.org/text) provides supporting types for working with text. A lot of is Font support and Caching, but we will use it for Aligmnent.

### 2. The editor widget

**A variable for the editor**

We start by declaring a [widget.Editor](https://pkg.go.dev/gioui.org/widget#Editor) variable inside ```draw()```:
```go
  // boilDurationInput is a textfield to input boil duration
  var boilDurationInput widget.Editor
```

Later, inside our flexbox we now have 4 separate rigids, one for each of the four components. The input box is the second of these:

```go
// The inputbox
layout.Rigid(
  func(gtx C) D {
    // Define characteristics of the input box ...
    boilDurationInput.SingleLine = true
    boilDurationInput.Alignment = text.Middle
    // ... and wrap it in material design
    ed := material.Editor(th, &boilDurationInput, "sec")
```              
We start by defining two of it's characteristics
 - ```SingleLine``` forces the box to always be one line high.
 - ```Àlignment``` centers text inside the box.
 

**Countdown**

Next, since it useful to see remaining time, we count down inside  the inputbox:
```go
if boiling && progress < 1 {
  boilRemain := (1 - progress) * boilDuration
  // Format to 1 decimal.
  inputStr := fmt.Sprintf("%.1f", math.Round(float64(boilRemain)*10)/10)
  // Update the text in the inputbox
  boilDurationInput.SetText(inputStr)
}
```

When we are in the middle of a boil, we here define a new ```boilRemain``` that holds the remaining time until the boil is complete, calculated using ```(1-progress```)

Since [math.Round()](https://pkg.go.dev/math#Round) doesn't allow rounding to a given number of decimals, we must use a trick. 
 - First multiply by 10. 
 - Then round to zero decimals. 
 - Afterwards divide by 10. 
 - Finally convert to text with 1 decimal.
It's a little compact, but hopefully straight forward

At the end, finally some Gio again. We call
[SetText](https://pkg.go.dev/gioui.org/widget#Editor.SetText) which replaces the text in the box with our ```inputStr```


**Layout**

At this point the input box is complete. We then start laying it out:

 1. Define margins. Note that left and right are quite large, this is how we keep the box farily small.

 1. Define a border with custom color and rounded corners.

 1. Wrap the input box in material design. We take the occation to add a [hint](https://pkg.go.dev/gioui.org/widget/material#EditorStyle)

 1. Combine 1+2+3 and return dimensions. 

Do you see how the final layout is now like a Russian doll? Margins contain border and border contains the editor, all returning their own ```layout.Dimensions```:

```go
return margins.Layout(gtx,
  func(gtx C) D {
    return border.Layout(gtx, ed.Layout)
  },
)
```

### 3. Well behaved button
Since I'm conciously lazy, I don't really care about what the user inputs in the input field before they push the start button. Hence, I chose to keep the logic for processing what's in the textbox as part of ```startButtion.Clicked()```. Other applications will have other needs and then it's better to do this in other parts of the code.

```go
if startButton.Clicked() {
  // Start (or stop) the boil
  boiling = !boiling
  // Read from the input box
  if progress >= 0 {
    inputString := boilDurationInput.Text()
    inputString = strings.TrimSpace(inputString)
    inputFloat, _ := strconv.ParseFloat(inputString, 32)
    boilDuration = float32(inputFloat)
  }
  // Resetting the boil
  if progress >= 1 {
    progress = 0
  }
}
```

As before, clicking the button flips the ```boiling``` variable, starting or stopping the boil.

Input is read only at start, i.e. when progress is 0. This is a simplification to make the logic a little neater. (If not we would have to add some code to figure out how progress should be updated if the remaining boiltime is changed mid way. And that's not so interesting.  )

Finally, the progress indicator is reset if the previous boil was completed. This allows us to boil many eggs in a row. Neat!

To present that to the user, we expand the code to paint the button:
```go
func(gtx C) D {
  var text string
  if !boiling {
    text = "Start"
  }
  if boiling && progress < 1 {
    text = "Stop"
  }
  if boiling && progress >= 1 {
    text = "Finished"
  }
  btn := material.Button(th, &startButton, text)
  return btn.Layout(gtx)
},
```
 - "Start" if not boiling
 - "Stop" if boiling but not finished
 - "Finished" if boil has completed

More bells and whistles could be added here. Might I for example challenge you to set a custom [background color](https://pkg.go.dev/gioui.org/widget/material#ButtonStyle) when the boil is done? 

## Final comments

And that's it. Thank you for coming along and I hope you've been tempted to try your hand at GUI development. 

We've only scratched the surface, and there's much more capability in the framework than we've exposed here. But, now that we have ve boiled that egg together, we've come quite far as well, right? 

Hopefully this get's you started on your own project. And when you do, big or small, please drop me a line, I would love to hear from you.

And make sure to follow Gio on the website, newsletter and community call. 

Now, it's finally time for breakfast. Guess what I'm having.