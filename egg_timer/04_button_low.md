---
layout: default
title: Chapter 4 
nav_order: 2
parent: Egg timer
has_children: false 
---

# Chapter 4 - Low button 

## Goals
The button can't fill the screen, obviously. So let's move the button to the bottom. To do that we start using the layout concept known as [Flexbox](https://pkg.go.dev/gioui.org/layout#Flex). 

![A low button with a spacer below](04_button_low.png)

## Outline
The last chapter was all about the overall structure of the program. Now we zoom into the **system.FrameEvent** and start using Flexbox. If it's new to you, read up on it first, for example [this one from Mozilla](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Flexible_Box_Layout/Basic_Concepts_of_Flexbox).

## Overall structure
We don't repeat the whole program here, but instead zoom in on the **system.FrameEvent**:

### Code
We start by removing a lot of the details to better see the structure:

```go
case system.FrameEvent:

    layout.Flex{
    // ...
    }.Layout( //...
        // We insert two rigid elements:
        // First one to hold a button ...
        layout.Rigid(),
        // .. then one to hold an empty spacer
        layout.Rigid(),
    }
```

### Comments

This we can work with:
1. First we define a **Flexbox** through the struct ```layout.Flex{}```
1. Then we send it list of children to be laid out through ```Layout(gtx, ...)```. The graphical context, *gtx*, contais the constraints that the kids have to keep within. 

The children we list are both created by ```layout.Rigid()```: 
  a. First a placeholder for the button, 
  b. Then another placeholder to contain empty space below the button.

A [Rigid](https://pkg.go.dev/gioui.org/layout?utm_source=gopls#Rigid) fills out the space it is given and is laid out first, while [Flexed](https://pkg.go.dev/gioui.org/layout?utm_source=gopls#Flexed) children share the space left after the Rigids. Apart from that, children are positioned in the order they are defined.

#### Constraint and Dimensions
At this point it´s worth mentioning how a Layout is bound together through **Constraints** and **Dimensions**. 
 - [Constraints](https://pkg.go.dev/gioui.org/layout?utm_source=gopls#Constraints) are the Minimum and Maxium size of a widget. How small or large **can** it be
 - [Dimensions](https://pkg.go.dev/gioui.org/layout?utm_source=gopls#Dimensions) are the actual size of a widget. How small or large is it **actually**

There are some nuances, like what to do if constraits can't me met, but for the most part this describes the dialoge between parent and child. From the parent, you create a Widget and pass in the context with constraints. The widget responds with it´s own dimensions, effectively laying itself out. 

As we saw above layout operations are recursive. A child in a layout can itself be a layout. From generic components you can thus create quite involved user interfaces.

## Code in detail

OK, time to look at the whole **system.FrameEvent**:

### Code

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
        // ... then an empty spacer
        layout.Rigid(
            //The height of the spacer is 25 Device independent pixels
            layout.Spacer{Height: unit.Dp(25)}.Layout,
        ),
    )
    e.Frame(gtx.Ops)
```

### Comments

Inside ```Flex { }``` we define two characteristicts:
 - Axis: Vertical alignment means stuff will be be placed or below each other.
 - Spacing: Leftover space will be at the start
Sounds like Tetris if you ask me.

Now let's look at the two calls to **layout.Rigid( )**:
- Rigid accepts a [Widget](https://pkg.go.dev/gioui.org/layout?utm_source=gopls#Widget)
- A Widget is simply something that returns it's own Dimensions
- Note how this is done in very different ways: 
  a. In the first Rigid we pass in a ```func ( )``` that returns Dimensions from btn.Layout()
  a. In the second Rigid we create a ```Spacer { }``` struct, calls it's Layout method, which in turn gives us Dimensions 
- From the parent perspective, it doesn't really matter. As long as the child returns **Dimensions**, it's good.

This takes care of laying the widgets out. But what are the widgets really?
As the name implies, **material.Button** is a [Button](https://pkg.go.dev/gioui.org/widget/material?utm_source=gopls#Button) based on material design, as we detailed in the last chapter.

[Spacer](https://pkg.go.dev/gioui.org@v0.0.0-20210504193539-82fff0178bed/layout?utm_source=gopls#Spacer) adds empty space, here defined by *Height*. Since we've defined the overall layout to be vertical, and excess space should come at the top, this falls to the bottom and the button lands on top of it. Hence some space is created, lifting the button a little from the bottom of the screen. Handy stuff.
