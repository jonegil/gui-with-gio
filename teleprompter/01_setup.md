---
layout: default
title: Chapter 1 - Setup
nav_order: 2
parent: Teleprompter
has_children: false
---

# Chapter 1 - Setup

## Goals
The goal of this chapter is to start up our application and get the basic structure in place.

## Outline

In this chapter we'll walk through the following pieces.

1. Introduce new imports to handle user input
1. Read the `.txt` file into a `[]string` slice
1. Start the application
1. Define state variables to control behaviour

These are relatively straight forward. If you completed the [egg timer](../egg_timer/) you will feel right at home. Let's get started.

## Code

### Section 1 - New imports

Many imports are well known, but these two are new:

```go
import (
  // Many imports we discussed earlier ...
  // ... plus two new interesting Gio imports
  "gioui.org/io/key"
  "gioui.org/io/pointer"
)
```

What can be going on here? Something with key´s and pointer´s maybe? From the docs: 

- Package [io/key](https://pkg.go.dev/gioui.org/io/key) implements key and text events and operations.
- Package [io/pointer](https://pkg.go.dev/gioui.org/io/pointer) implements pointer events and operations. A pointer is either a mouse controlled cursor or a touch object such as a finger.

Notice how pointer supports both mouse gestures on a desktop, trackpad on a laptop and fingers on a screen. Nice, again an example of how learning a cross-platform framework gives skills on multiple devices.

### Section 2 - Read text into a slice

We default to reading the text to be read from `speech.txt`. But, as we know, users are crazy and might want other filenames. So we oblige and prepare a command line flag.

```go
// Command line input variables
var filename *string
```

To work with the speech, it´s helpful to store it not as one massive text-variable, but rather as a list of paragraphs. We´ll get into those details later, for now let´s define a placeholder.

```go
// A []string to hold the speech as a list of paragraphs
var paragraphList []string
```

With these to in place we know where to look for a speech and where to place it. Now let´s fire up `main()` and get cracking:
```go
func main() {
  // Step 1 - Read input from command line
  filename = flag.String("file", "speech.txt", "Which .txt file shall I present?")
  flag.Parse()

  // Step 2 - Read from file
  paragraphList = readText(filename)
```


The `readText()` func does what the name suggests, but let's have a look to be sure:

```go
func readText(filename string) []string {
  f, err := os.ReadFile(filename)
  text := []string{}
  if err != nil {
    log.Fatal("Error when reading file:\n  ", err)
  }
  if err == nil {
    // Convert text to a slice of strings.
    text = strings.Split(string(f), "\n")
    // Add extra empty lines a the end. Simple trick to ensure
    // the last line of the speech scrolls out of the screen
    for i := 1; i <= 10; i++ {
      text = append(text, "")
    }
  }

  return text
}
```

The first line of `readText` reads the file. If all goes well, i.e. `err == nil`, we continue to split it by `\n`, newline. 

We also do a little trick at the end. It felt clunky that a speech didn't full scroll of screen after it was finished. An easy fix was to add more empty paragraphs at the end of the list. Easy peasy. 

*Note:* In the sourcecode there's an alternative implementation that generates a very long speech programatically. That was useful when debugging, so I left it in. Please feel free to play around with it. 

### Section 3 - Start the application

The last part of `main` starts the GUI in a normal manner:

```go
  // Step 3 - Start the GUI
  go func() {
    // create new window
    w := app.NewWindow(
      app.Title("Teleprompter"),
      app.Size(unit.Dp(650), unit.Dp(600)),
    )
    // draw on screen
    if err := draw(w); err != nil {
      log.Fatal(err)
    }
    os.Exit(0)
  }()
  app.Main()
}
```

### Section 4 - State variables to control behaviour

Now we're getting into the meat of things. In order to control the behaviour and looks of the program we need multiple variables that control things should be displayed and what the program should do.

There's no one-right-answer on how to do this. Instead, while working with the code over time, the following set emerged and proved helpful:

| Variable     | Description                                   | Changed with                                  |
| ------------ | --------------------------------------------- | --------------------------------------------- |
| `scrollY`    | Scroll the text up/down                       | Mouse/Trackpad scroll, Arrow Up/Down, **J/K** |
| `focusBarY`  | How high up is the red focus bar              | **U** (up) and **D** (down)                   |
| `textWidth`  | How wide is the area in which we display text | **W** (wider) and **N** (narrower)            |
| `fontSize`   | How large is the text                         | **+** (larger) and **-** (smaller)            |
| `autoscroll` | Start/stop automatic scrolling                | **Space**                                     |
| `autospeed`  | How fast / slow the text should scroll        | **F** (faster) or **S** (slower)              |
| `stepSize`   | How large should each change be?              | **Shift** makes any change larger             |


The user should be able to adjust these while using the program. In other words, we can't have them hard coded as constants, but instead want them as variables so they can, well vary. Although we could spread them across the code, I instead opted to collect them in one place to keep the program tidy. Thus, here's the start of the `draw()`:

```go
// The main draw function
func draw(w *app.Window) error {
  // y-position for text
  var scrollY unit.Dp = 0

  // y-position for red focusBar
  var focusBarY unit.Dp = 170

  // width of text area
  var textWidth unit.Dp = 550

  // fontSize
  var fontSize unit.Sp = 35

  // Are we auto scrolling?
  var autoscroll bool = false
  var autospeed unit.Dp = 1
```

The numerical values are set using `unit.Dp` and `unit.Sp`, the [device independent units](https://pkg.go.dev/gioui.org/unit) that ensure similar size across devices and screen resolutions. Note how we also use it for `autospeed`, even though it controls animation speed and not visualization directly. 

With that, let's look at how to collect and handle inputs.

---

[Next chapter](02_user_input.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[View it on GitHub](https://github.com/jonegil/gui-with-gio/tree/main/teleprompter){: .btn .fs-5 .mb-4 .mb-md-0 }