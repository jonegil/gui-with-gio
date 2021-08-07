---
layout: default
title: Bonus - Improved animation
nav_order: 3
parent: Egg timer
has_children: false
---

# Bonus material - Animation

## Goals
The intention of this section is to discuss a slightly more advanced topic related to animation, namely how and when we invalidate a frame, what that actually means, and how to code well with it. 

![The complete egg timer](egg_timer.gif)

## Outline

The outline of this bonus chapter is as follows:
 - First we discuss what it means to invalidate a frame
 - Then we look at two different method calls to do so
 - Finally we discuss an alternative pattern to generate and control animation


### 1. What is Invalidate?

Gio only updates what you see when a [FrameEvent](https://pkg.go.dev/gioui.org/io/system#FrameEvent) is generated. This can be for example when a key is pressed, mouse is clicked, widget receives or loses focus. That makes perfect sense, with refresh rates of up to 120 frames per second for modern devices, chances are that what should be displayed quite often is identical to the last frame.

Quite often. But not always. 

The main exception to this rule is animations. When animating, you want it to run as smooth as possible. 
>*Gliding, flowing, effortlessly, fluently, lighthlessly, painlessly, readily, well.* [Thank you Merriam-Webster](https://www.merriam-webster.com/thesaurus/smoothly). 

To achieve this, we need to ask Gio to redraw continuosly. And without triggering events we need to explicitly tell Gio to do so. That is done by calling ```invalidate```. There are two alternatives, let's look at both:

- [op.Invalidate](https://pkg.go.dev/gioui.org/op#InvalidateOp)
- [window.Invalidate](https://pkg.go.dev/gioui.org/app#Window.Invalidate)


To show the first, we'll quote the animation example in the excellent [architecture document](https://gioui.org/doc/architecture)
```go
var startTime = time.Now()
var duration = 10 * time.Second

func drawProgressBar(ops *op.Ops, now time.Time) {
	// Calculate how much of the progress bar to draw,
	// based on the current time.
	elapsed := now.Sub(startTime)
	progress := elapsed.Seconds() / duration.Seconds()
	if progress < 1 {
		// The progress bar hasn’t yet finished animating.
		op.InvalidateOp{}.Add(ops)
	} else {
		progress = 1
	}

	defer op.Save(ops).Load()
	width := 200 * float32(progress)
	clip.Rect{Max: image.Pt(int(width), 20)}.Add(ops)
	paint.ColorOp{Color: color.NRGBA{R: 0x80, A: 0xFF}}.Add(ops)
	paint.ColorOp{Color: color.NRGBA{G: 0x80, A: 0xFF}}.Add(ops)
	paint.PaintOp{}.Add(ops)
}

```



Depending on the complexity of your layouts, that might cause undue stress on your system. [Recent changes]() introduces partwise caching, so that only what has changes requires recompute. If that works for your application you


```go

}
```

## Comments

What's happening here, and why does it matter

