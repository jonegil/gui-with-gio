---
layout: default
title: Chapter 3 
nav_order: 2
parent: Egg timer
has_children: false 
---

# Chapter 3 - Button 

## Goals
The intention of this section is to add a button. Not only can we click it, but it will have a nice hover and click animations.

## Outline
We start by reviewing the new packages that are imported. There are quite many, so let´s spend some time here. Next, we look at how **operations** and **widgets** combine to make a button.
Finally we touch on [Material Design](https://material.io/), the well established framework for user interfaces also available in Gio.

## Code
To make things tidy, let's discuss imports first.

```go
import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)
```

**app** and **unit** we know before, but the rest are new:
- [font/gofont](https://pkg.go.dev/gioui.org/font/gofont) - Did you know Go has it's own dedicated high-quality True Type fonts? Read the [fascinating blog](https://blog.golang.org/go-fonts) and definetly visit [Bigelow & Holmes](https://bigelowandholmes.typepad.com), its creators. True old-school.

- [io/system](https://pkg.go.dev/gioui.org/io/system) - Provides high-level events that are sent from the window. Most important is the **system.FrameEvent**. It's effectively a list of operations that do one of two things: Details how to handle input and describes what to display.
  
- [layuot](https://pkg.go.dev/gioui.org/layout) - Defines useful parts of a layout, such as *dimensions*, *constraints* and *directions*. Also, it includes the layout-concept known as [Flexbox](https://pkg.go.dev/gioui.org/layout#Flex). It's widely used web and user interface development. Among the many introductions, I recommend the one from [Mozilla](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Flexible_Box_Layout/Basic_Concepts_of_Flexbox). 

- [op](https://pkg.go.dev/gioui.org/op) - Operations, or ops, are central in Gio. They are used to update the user interface. There are operations used to draw, handle input, change window properties, scale, rotate and more. Interestingly there are also [macros](https://pkg.go.dev/gioui.org/op#MacroOp), making it possible to record opertions to be executed later. Taken together this means a list of opererations is a *mutable stack*, where you can control the flow.

- [widget](https://pkg.go.dev/gioui.org/widget) - Widget provides the underlying functionality of UI components, such as state tracking and event handling. Is the mouse hovering over a button? Has it been clicked, and if so how many times? 

- [widget/material](https://pkg.go.dev/gioui.org/widget/material) - While the **widget** provides functionality, **widget/material** defines a theme. The default looks good, and is what we'll use, but it's just as easy to adjust by setting propoerties such as color, text size font properties etc.
  - Note: Gio expands the base functionality in a dedicated repo called [gio-x](https://pkg.go.dev/gioui.org/x) where [more material components](https://pkg.go.dev/gioui.org/x/component) are in development, including navigation bars and tooltips.

```go
func main() {
	go func() {
		// create new window
		w := app.NewWindow(
			app.Title("Egg timer"),
			app.Size(unit.Dp(400), unit.Dp(600)),
		)

		// ops are the operations from the UI
		var ops op.Ops

		// startButton is a clickable widget
		var startButton widget.Clickable

		// th defnes the material design style
		th := material.NewTheme(gofont.Collection())

		// listen for events in the window.
		for e := range w.Events() {

			// detect what type of event
			switch e := e.(type) {

			// this is sent when the application should re-render.
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				b := material.Button(th, &startButton, "Start")
				b.Layout(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}


```



## Comments

