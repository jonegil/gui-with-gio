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

Then we check if we should autoscroll. If so, that's done by incrementing on the Y-axis offset variable `scrollY` which is used to control the [Position](https://pkg.go.dev/gioui.org/layout#Position) of [List](https://pkg.go.dev/gioui.org/layout#List). Using [op.InvalidateOp](https://pkg.go.dev/gioui.org/op#InvalidateOp) to request a redraw in 2/100th of a second creates a smooth animation.
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

By using a list, Gio takes care of only showing the elements currently on screen. Off screen elements are not processed until they are needed, reducing the load on the system and allowing for really long lists. In developing this app I played around with for example [The Complete Works of William Shakespeare](https://www.gutenberg.org/ebooks/100). No problem.

---

THE TEXT BELOW IS NOT UPDATED
Jan 8th

---


### Focusbar

Finally we add a focusbar. This is done in the following steps:

- Use `op.Offset()` to move to a new Y position, the one defined by our state variable `focusBarY`.
- From there, create a new rectangle, width = fullscreen and height = 50
- Color it with transparent red. `A: 0x66` controls the transparency, where 0 means zero visibility (full transparency) and 0xff means full visibility (no transparency).
- Add the Paint

At the end we complete the FrameEvent by `e.Frame()`.

```go
// Draw a transparent red rectangle.
op.Offset(image.Pt(0, focusBarY)).Add(&ops)
clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, 50)}.Add(&ops)
paint.ColorOp{Color: color.NRGBA{R: 0xff, A: 0x66}}.Add(&ops)
paint.PaintOp{}.Add(&ops)

// Frame completes the FrameEvent by drawing the graphical operations from ops into the window.
e.Frame(&ops)
```

### system.DestroyEvent

Finally, just to complete the picture, it's worth mentioning the final event we listen for, namely the `system.DestroyEvent`. It helps us end the program gracefull, returns an `Err` and breaks the `range w.Events()` loop were in to listen for events.

```go
// Shutdown?
case system.DestroyEvent:
  return e.Err
}
```

## Comments

That´s it. We've got yet another Gio project in our belt, great work!. This one was all about processing input, which we did by listening to events, `key.Event` and `pointer.Event` respectively, and using custom logic to update a set of state variables. Later, in `system.FrameEvent` we used those state variables to control our layout.

Thank you again so much for following the writeup. If you found this useful, share it with a friend, star it on Github or drop me a line. It's really motivating to hear back from you. Good luck with all your projects!

---
[Next chapter](04_event_area.md){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[View it on GitHub](https://github.com/jonegil/gui-with-gio/tree/main/teleprompter){: .btn .fs-5 .mb-4 .mb-md-0 }
