---
layout: default
title: Chapter 9 
nav_order: 2
parent: Egg timer
has_children: false
---

# Chapter 9 - Egg as egg

## Goals
The intention of this section is to draw an actual egg.

![An actual egg](09_egg_as_egg.gif)

## Outline

Here we utilize basic Gio functionality to draw totally custom graphics, in the shape of an egg.

## Code
All the new code is within the rigid that previously displayed the circle. 

```go
layout.Rigid(
  func(gtx C) D {
    // Draw a custom path, shaped like an egg
    var egg clip.Path
    op.Offset(f32.Pt(200, 150)).Add(gtx.Ops)
    egg.Begin(gtx.Ops)
    // Rotate from 0 to 360 degrees
    for deg := 0.0; deg <= 360; deg++ {
      // Egg math (really) at this brilliant site. Thanks!
      // https://observablehq.com/@toja/egg-curve
      // Convert degrees to radians
      rad := deg / 360 * 2 * math.Pi
      // Trig gives the distance in X and Y direction
      cosT := math.Cos(rad)
      sinT := math.Sin(rad)
      // Constants to define the eggshape
      a := 110.0
      b := 150.0
      d := 20.0
      // The x/y coordinates
      x := a * cosT
      y := -(math.Sqrt(b*b-d*d*cosT*cosT) + d*sinT) * sinT
      // Finally the point on the outline
      p := f32.Pt(float32(x), float32(y))
      // Draw the line to this point
      egg.LineTo(p)
    }
    //Close the path
    egg.Close()
    clip.Outline{Path: egg.End()}.Op().Add(gtx.Ops)

    // Fill the shape
    //color := color.NRGBA{R: 255, G: 239, B: 174, A: 255}
    color := color.NRGBA{R: 255, G: uint8(239 * (1 - progress)), B: uint8(174 * (1 - progress)), A: 255}
    paint.Fill(gtx.Ops, color)

    d := image.Point{Y: 375}
    return layout.Dimensions{Size: d}
  },
),

```

The main idea is to define a custom egg shaped **clip.Path**. We draw a line to define it, fill the inside and any drawing on the outside is ignored. 

## Comments

First the new path is defined, ```var eggPath clip.Path````

Then an operation is created to move 200 Points right, 150 Points down, ```op.Offset( )```. As before, this is from the top-left corner inside this widget.

We're now at the center of our egg. This is where the path begins, ```eggPath.Begin( )```

From here we rotate a full 360 degrees and continue drawing the outline of the Egg. We use the math for a [Hügelschäffer Egg](https://mathcurve.com/courbes2d.gb/oeuf/oeuf.shtml), as presented on Torben Jansen's [excellent interactive blog](https://observablehq.com/@toja/egg-curve). The formula receives an angle, from 0 to 360, calculates an approporiate distance from center and returns the outline as a point. Math is fun!

![Hügelschäffer egg](09_torben_jansen.gif)


Regarding Gio, the important line is the last in the for-loop, ```eggPath.LineTo(p)```. At this point, math has found the next point ´´´p´´´ on the 360-degree roundtrip around the egg, and we use [eggPath.LineTo()](https://pkg.go.dev/gioui.org/op/clip#Path.LineTo) to move the pen to this specific coordinate point.

After completing the for-loop the egg-shape is almost done. We finalize it explicitly by calling ```eggPath.Close()``` which closes the path.

With the path complete, we want to get the area inside that path. [clip.Outline{ }.Op( )](https://pkg.go.dev/gioui.org/op/clip#Outline.Op) gives gives us the clip operation representing this area. 

Coloring can either be static, as shown in the uncommented colors, and they work fine. However, wouldn't it be cool if the changed color from cold to warm? I think so at least. Remember how progress is a variable from 0 to 1. This state variable can now be used to slowly alter the color as well. ``` * (1 - progress) ``` is just another way of saying *gradually turn off Green and Blue please*. When progress is complete, both are 0, and we're left with red only. Nifty. 

We end by returning **layout.Dimensions**, the height of this widget.