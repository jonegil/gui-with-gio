---
layout: default
title: Chapter 4 
nav_order: 2
parent: Egg timer
has_children: false 
---

# Chapter 4 - Low button 

## Goals
The intention of this section is to move the button to the bottom. To do that we start using the [Flexbox layout concept]((https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Flexible_Box_Layout/Basic_Concepts_of_Flexbox))

## Outline
The last chapter was all about the overall structure of the program. Now we zoom into the **system.FrameEvent** and start using the [Flexbox](https://pkg.go.dev/gioui.org/layout#Flex) 

## Code
Here's the whole **system.FrameEvent**:
```go
case system.FrameEvent:
    gtx := layout.NewContext(&ops, e)
    // Let's try out the flexbox layout concept:
    layout.Flex{
        // Vertical alignment, from top to bottom
        Axis: layout.Vertical,
        //Emtpy space is left at the start, i.e. at the top
        Spacing: layout.SpaceStart,
    }.Layout(gtx,
        // We insert two rigid elements:
        // First a button ...
        layout.Rigid(
            func(gtx layout.Context) layout.Dimensions {
                btn := material.Button(th, &startButton, "Start")
                return btn.Layout(gtx)
            },
        ),
        // .. then an empty spacer
        layout.Rigid(
            //The height of the spacer is 25 Device independent pixels
            layout.Spacer{Height: unit.Dp(25)}.Layout,
        ),
    )
    e.Frame(gtx.Ops)
```

Skim the details and look at the major components. Use a foldable editor and hide some of the inner lines.

```go
layout.Flex{
   // ...
}.Layout( //...
    // We insert two rigid elements:
    // First a button ...
    layout.Rigid(),
    // .. then an empty spacer
    layout.Rigid(),
}
```

This we can work with:
1. First we defeine a Flexbox through **Layout.Flex**
1. Then we place a **Layout** onto that Flexbox
1. Within the layout we place two rigids. 
  a. First one to contain the buttion
  b. Then one to contain a spacer below it.


## Comments

