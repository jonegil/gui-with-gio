---
layout: default
title: Chapter 3 - Layout
nav_order: 2
parent: Teleprompter
has_children: false
---

# Chapter 3 - Layout

## Goals
In this chapter we'll make use of the updated state variables and draw our program on screen.

## Outline
After processing events, `system.FrameEvent` continues with laying out and drawing on screen. It has multiple pieces, but we'll calmly walk through them one by one.

## Code

### Background 

We start with coloring the bacground:
```go
// ---------- LAYOUT ---------- 
// Let's start with a background color
paint.Fill(&ops, color.NRGBA{R: 0xff, G: 0xfe, B: 0xe0, A: 0xff})
```

### Scrolling text

Then we check if we should autoscroll. If yes, we continously increase `scrollY`, used to set the [Position](https://pkg.go.dev/gioui.org/layout#Position) of the [List](https://pkg.go.dev/gioui.org/layout#List). 

We request a redraw in 2/100th of a second by calling [op.InvalidateOp](https://pkg.go.dev/gioui.org/op#InvalidateOp) at a future point in time.

```go
// ---------- THE SCROLLING TEXT ----------
// First, check if we should autoscroll
// That's done by increasing the value of scrollY
if autoscroll {
  scrollY = scrollY + autospeed
  op.InvalidateOp{At: gtx.Now.Add(time.Second * 2 / 100)}.Add(&ops)
}
// Then we use scrollY to control the distance from the top of the screen to the first element.
// We visualize the text using a list where each paragraph is a separate item.
var visList = layout.List{
  Axis: layout.Vertical,
  Position: layout.Position{
    Offset: int(scrollY),
  },
}
```

### Margins
These are the margins for the text, effectively controlling the empty space on the right- and left side of the column of text in the middle.

```go
// ---------- MARGINS ----------
// Margins
var marginWidth unit.Dp
marginWidth = (unit.Dp(gtx.Constraints.Max.X) - textWidth) / 2
margins := layout.Inset{
  Left:   marginWidth,
  Right:  marginWidth,
  Top:    unit.Dp(0),
  Bottom: unit.Dp(0),
}
```


### Layout the list

Now it's time to lay out the list withing a set of margins. That's a nested structure with 3 pieces. 
- First define the margins
  - Within those margins, define a list
    - For each element of the list, define a paragraph

```go
  // ---------- LIST WITHIN MARGINS ----------
  // 1) First the margins ...
  margins.Layout(gtx,
    func(gtx C) D {
      // 2) ... then the list inside those margins ...
      return visList.Layout(gtx, len(paragraphList),
        // 3) ... where each paragraph is a separate item
        func(gtx C, index int) D {
          // One label per paragraph
          paragraph := material.Label(th, unit.Sp(float32(fontSize)), paragraphList[index])
          // The text is centered
          paragraph.Alignment = text.Middle
          // Return the laid out paragraph
          return paragraph.Layout(gtx)
        },
      )
    },
  )
```

By using the list widget, Gio takes care of only showing the elements currently on screen. Off screen elements are not processed until they are needed, reducing the load on the system and allowing for really long lists. In developing this app I played around with for example [The Complete Works of William Shakespeare](https://www.gutenberg.org/ebooks/100). No problem.

### Focusbar

Finally we add a focusbar. This is the red ribbon that helps the reader keep focus on what needs to be said right now, while at the same time allow more of the speech to be seen around it. 

```go 
// ---------- THE FOCUS BAR ----------
// Draw the transparent red focus bar.
focusBar := clip.Rect{
  Min: image.Pt(0, int(focusBarY)),
  Max: image.Pt(gtx.Constraints.Max.X, int(focusBarY)+int(unit.Dp(50)),
}.Push(&ops)
paint.ColorOp{Color: color.NRGBA{R: 0xff, A: 0x66}}.Add(&ops)
paint.PaintOp{}.Add(&ops)
focusBar.Pop()
```

The `focusBar` is simply a rectangle drawn from left (x=0) to right (x=`Max.X`) and from top at `focusBarY` and height of 50 Dp. As with all variables we run them as `unit.Dp`, for device compatability, but cast them to `int` here for the Point.

Once the `focusBar` is pushed to the stack of operations, we add color to it, full Red but transparent.
`A: 0x66` controls the transparency, where 0 means zero visibility (full transparency) and 0xff means full visibility (no transparency)
`PaintOp` actually paints it and the focusBar can be popped from the stack.


## Comments

We've completed laying out our application on screen. By using state variables, we control the looks and behaviour of the app, and allow the users commands to actually create changes in the app. 

What's missing now is how to actually capture those inputs. So that's exactly what we'll complete in [Chapter 4](04_event_area.md). See you there!


---
[Next chapter](04_event_area.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[View it on GitHub](https://github.com/jonegil/gui-with-gio/tree/main/teleprompter){: .btn .fs-5 .mb-4 .mb-md-0 }
