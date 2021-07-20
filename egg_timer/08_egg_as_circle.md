---
layout: default
title: Chapter 8 
nav_order: 2
parent: Egg timer
has_children: false
---

# Chapter 8 - Custom graphics - a circle

## Goals
The intention of this section is to draw custom graphics that (vaguely) resembles an egg

![Egg as circle](08_egg_as_circle.gif)

## Outline

The code introduces custom graphics. The circle in the app is drawn by Gio, not displaying a static picture. Tho achieve that we combine 
 - A **clip** do define the area we can draw within
 - A **paint** operation to fill that area
 - Some paramaters to set color
 
## Imports

There are some new imports, namely
 - [image](https://pkg.go.dev/image) and [image/color](https://blog.golang.org/image/color), Go's standard 2D image library.
   - Nigel Tao [writes well](https://blog.golang.org/image) about these packages on the Go blog.

 - [f32](https://pkg.go.dev/gioui.org/f32). Go's image library is based on ```int```, while Gio prefers to work with ```float32```. Hence f32 reimplements floating point versions of the two main types, ```Points``` and ```Rectangles```
 
 - [clip]()

 - [paint]()


## Points and Rectangles
Points and Rectangles are used extensively, so it's worth quoting from Nigel's blog mentioned above. Point's are coordinates and Rectangles are defined by Points:

```go
type Point struct {
    X, Y float32
}

type Rectangle struct {
    Min, Max Point
}
```

> A Point is an X, Y coordinate pair. The axes increase right and down (origin = top left corner). It is neither a pixel nor a grid square. A Point has no intrinsic width, height or color, but the visualizations below use a small colored square.
![Point](08_image_package_point.png)
```go
p := f32.Point{2, 1}
```

> A Rectangle contains the points with Min.X <= X < Max.X, Min.Y <= Y < Max.Y. It has no intrinsic color, but the visualization below outlines it with a thin colored line, and call out their Min and Max Points.

![Rectangle](08_image_package_rectangle.png)
```go
r := f32.Rect(2, 1, 5, 5)
```
For convenience, image.Rect(x0, y0, x1, y1) is equivalent to ```Rectangle{Point{x0, y0}, Point{x1, y1}}```, but is much easier to type. It also swaps Minimum and Maximum to ensure it's well formed.
   
That's it. Let's look at the code:

## Code

```go

```

## Comments

Package clip provides operations for clipping paint operations. Drawing outside the current clip area is ignored.

The current clip is initially the infinite set. An Op sets the clip to the intersection of the current clip and the clip area it represents.