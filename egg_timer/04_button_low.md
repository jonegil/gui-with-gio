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

## Overall structure

### Code
You know the overall strucutre of the program like the back of your hand out from last chapter. So instead of repeating it all, so here we´re only going to focus on what´s happening inside the **system.FrameEvent**:

First look at the Layout

```go
case system.FrameEvent:

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

### Comments

This we can work with:
1. First we define a **Flexbox** through Layout.Flex
1. Then we place a new **Layout** onto that Flexbox
1. That 2nd layout contains two **rigids**. 
  a. First one to contain the button
  b. Then one to contain a spacer below it.

**Flex** lays out a list of children. 
- Rigid children are laid out first and Flexed children laid out after.
- Apart from that, children are positioned in the order they are defined.

**Rigid** is simply a Flex child filling out available space. 

#### Constraint and Dimensions
It´s worth mentioning how a Layout is bound together through [Constraints](https://pkg.go.dev/gioui.org/layout?utm_source=gopls#Constraints) and [Dimensions](https://pkg.go.dev/gioui.org/layout?utm_source=gopls#Dimensions). 
 - Constraints are the Minimum and Maxium size of a widget ´´´Min, Max image.point´´´
   - I.e. how large can it be
 - Dimensions are the Actual size of a widget, effectively ´´´Size image.Point´´´
   - I.e. how large is it actually

There are some nuances, like what to do if constraits can't me met, but for the most part this describes the dialoge between parent and child. From the parent, you create a Widget and pass in the context. The widget responds, with it´s own dimensions, effectively laying itself out. 

Note that layout operations are recursive. A child in a layout can itself be a layout. From generic components you can thus create quite involved user interfaces.

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
 - Vertical alignment. Stuff will be be placed or below each other.
 - Leftover space will be at the start
Sounds like Tetris if you ask me.

Now look at the two calls **layout.Rigid( )**:
 - ```Rigid( )``` is a function that accepts a **Widget**
 - A widget is whatever that returns [Dimensions](https://pkg.go.dev/gioui.org/layout?utm_source=gopls#Dimensions)
    a. ```func (b ButtonStyle) Layout(gtx layout.Context) layout.Dimensions {```
    a. ```func (s Spacer) Layout(gtx Context) Dimensions {```
 - If you added a third, fourth, fifth element, guess what they would each return? Yes, Dimensions. It's pretty neat how this simple element binds the Gui together.








