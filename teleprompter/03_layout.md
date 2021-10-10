---
layout: default
title: Chapter 3 - Layout
nav_order: 2
parent: Teleprompter
has_children: false 
---

# Chapter 3 - Layout

In this chapter we'll make use of the updated state variables and lay out our program.

Now that we have processed all incoming input, both ```key.Event``` and ```pointer.Scroll```, it's time to wait for a request to redraw. Those are sent by call ```w.Invalidate``` at the end of the key and pointer event sections. Also, we'll add an ```op.InvlidateOp{}``` operation will also when autoscroll is turned on as we'll see below. 

## system.FramveEvent
When we receive a ```system.FrameEvent``` it is time to lay out and redraw. As we'll get into, it's a nested structure with three main components. 

### Setup
The setup for rendering provides us the with the building blocks we will need:

```go
// A re-render request?
case system.FrameEvent:
  // ops are the operations from the UI
  var ops op.Ops

  // Graphical context
  gtx := layout.NewContext(&ops, e)

  // Bacground
  paint.Fill(&ops, color.NRGBA{R: 0xff, G: 0xfe, B: 0xe0, A: 0xff})

  // Textscroll
  if autoscroll {
    scrollY = scrollY + autospeed
    op.InvalidateOp{At: gtx.Now.Add(time.Second / 50)}.Add(&ops)
  }
```
We have identified a ```FrameEvent```.  First we define ```ops```, the list of operations, as well as the graphical context we will work within. The background then is filled with a soothing papyrus-like color. Finally we check if ```autoscroll``` is activated. If so, we move the starting point for text by a small amount, ```autospeed```, and request a redraw in 0.02 seconds. This last part is interesting, effectively setting the framerate of our change. The higher the smoother, but also effectively alter the speed. As you remember, [there are some nuances](../egg_timer/11_improved_animation.md), between ```w.Invalidate``` and ```op.InvalidateOp{}.Add```. Maybe most interesting here is the timing functionality though. Feel free to experiment.

Let's continue the coding.

### Three main parts

The three parts of the layout are 
 - Margins: [layout.Inset](https://pkg.go.dev/gioui.org/layout#Inset)
 - A list of paragraphs: [layout.List](https://pkg.go.dev/gioui.org/layout#List)
 - Each single paragraph: [material.Label](https://pkg.go.dev/gioui.org/widget/material#Label)

```go
  // Margins
  marginWidth := (gtx.Constraints.Max.X - textWidth) / 2
  margins := layout.Inset{
    Left:   unit.Dp(float32(marginWidth)),
    Right:  unit.Dp(float32(marginWidth)),
    Top:    unit.Dp(0),
    Bottom: unit.Dp(0),
  }

  // Visualisation of the speech, using a list where each paragraph is a separate item.
  // Offset is the distance from the top of the screen to the first element.
  // I.e. it controls how far we have scrolled.
  var visList = layout.List{
    Axis: layout.Vertical,
    Position: layout.Position{
      Offset: scrollY,
    },
  }

  // Layout the list inside margins
  // 1) First the margins ...
  margins.Layout(gtx,
    func(gtx C) D {
      // 2) ... then the list inside those margins ...
      return visList.Layout(gtx, len(paragraphList),
        // 3) ... where each paragraph is a separate item
        func(gtx C, index int) D {
          // One label per paragraph
          paragraph := material.Label(th, unit.Dp(float32(fontSize)), paragraphList[index])
          // The text is centered
          paragraph.Alignment = 2
          // Return the laid out paragraph
          return paragraph.Layout(gtx)
        },
      )
    },
  )
```

The ```margins``` are on the right and left side of the screen. There role is to grow and shring so that the text in the middle is squeezed together or can flow wide to fill the screen. Since wide text requres narrow margins, the ```marginWidth``` is calculated by subtracting the ```textWidth``` state variable from full screenwidth from ```gtx.Constraints.Max.X```.

The list of visualised paragraphs is defined in ```visList```. As the struct defines it aligns vertically, i.e. elements are above and below each other. Most interesting is the ```Offset: scrollY``` which defines the distance in pixels from the top of the screen to the first element in the list. In other words, by setting the offset to the value of our ```scrollY``` state variable we move the whole list up and down. And voila, we're scrolling.  

The third nested block reads as follows:
- First define the margins
  - Within those margins define a list
    - Within each element of the list, define a paragraph
    - Return the paragraph
  - Once each element in the list is visited, return the list
- Done

By using a list, Gio takes care of only showing the elements currently on screen. Off screen elements are not processed until they appear, reducing the load on the system and allowing for really long lists. In developing this app I played around with some really long ones, like [The Complete Works of William Shakespeare](https://www.gutenberg.org/ebooks/100) for example. No problem.

### Focusbar

Finally we add the focusbar. This is done in the following steps:
- Use ```op.Offset()``` to move to a new Y position, the one defined by our state variable ```focusBarY```.
- From there, create a new rectangle, width = fullscreen and height = 50
- Color it with transparent red. ```A: 0x66``` controls the transparency, where 0 means zero visibility (full transparency) and 0xff means full visibility (no transparency). 
- Add the Paint

At the end we complete the FrameEvent by ```e.Frame()```.

```go
// Draw a transparent red rectangle.
op.Offset(f32.Pt(0, float32(focusBarY))).Add(&ops)
clip.Rect{Max: image.Pt(gtx.Constraints.Max.X, 50)}.Add(&ops)
paint.ColorOp{Color: color.NRGBA{R: 0xff, A: 0x66}}.Add(&ops)
paint.PaintOp{}.Add(&ops)

// Frame completes the FrameEvent by drawing the graphical operations from ops into the window.
e.Frame(&ops)
```

## system.DestroyEvent
To ensure the program terminates gracefully we receive the ```system.DestroyEvent``` which returns an ```Err``` and breaks the ```draw()``` loop.

```go
// Shutdown?
case system.DestroyEvent:
  return e.Err
}
```


## Wrapping it all up

That´s it. We've got yet another Gio project in our belt, great work!. This one was all about processing input, which we did by listening to events, ```key.Event``` and ```pointer.Event``` respectively, and using custom logic to update a set of state variables. Later, in ```system.FrameEvent``` we used those state variables to control our layout. 

Thank you again so much for following the writeup. If you found this useful, share it with a friend, star it on Github or drop me a line. It's really motivating to hear back from you. Good luck with all your projects!

---

[View it on GitHub](https://github.com/jonegil/gui-with-gio/tree/main/teleprompter){: .btn .fs-5 .mb-4 .mb-md-0 }
