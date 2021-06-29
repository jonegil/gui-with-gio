---
layout: default
title: Chapter 3 
nav_order: 2
parent: Egg timer
has_children: false 
---

# Chapter 3 - Button 

## Goals
The intention of this section is to add a button. That we can click. It will be so good.

## Outline

We start by reviewing the handful of new packages that are imported. 
Next, we look at how **operations** and **widgets** combine to amke a button.
Finally we touch on [Material Design](material.io), the well established framework for user interfaces also available in Gio.

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
- [font/gofont](https://pkg.go.dev/gioui.org/font/gofont) - Did you know Go had it's own dedicated set of high-quality True Type fonts? 
  - [Read](https://blog.golang.org/go-fonts) the fascinating blog and definetly visit [https://bigelowandholmes.typepad.com/](Bigelow & Holmes), the creators
- [io/system](https://pkg.go.dev/gioui.org/io/system]) - Provides high-level events that are sent from the window. 
  - Most important now will be the **system.FrameEvent** that requests a new frame that describes what to display.
- [layuot](https://pkg.go.dev/gioui.org/layout]) - Defines various layouts, such as *dimensions*, *constraints* and *directions*, as well as [Flexbox](https://pkg.go.dev/gioui.org/layout#Flex) that is well decribed by [Mozilla](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Flexible_Box_Layout/Basic_Concepts_of_Flexbox)





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

