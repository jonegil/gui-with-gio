---
layout: default
title: Chapter 6 
nav_order: 2
parent: Egg timer
has_children: false 
---

# Chapter 6 - Low button with margin

## Goals
The intention of this chapter is to add open space around all sides of the button.

![Button with margin](06_button_low_margin.png)

## Outline

After looking at the whole code when refactoring in the last section, this time we only zoom in on the lines that change. Again, the action is happening within **layout.Flex**

## Code

It can be useful to strip away some of the detail, to highlight the structure. 
There are really only three key lines here:
 1. Define margins
 2. Lay out the margins
 3. Create button within those margins

```go
layout.Flex{
    //...
}.Layout(gtx, 
    layout.Rigid(
        func(gtx C) D {
            // ONE: First define margins around the button using layout.Inset ...
            margin := layout.Inset{
                // ...
            }

            // TWO: ... then we lay out those margins ...
            margins.Layout(
                
                // THREE: ... and finally within the margins, we ddefine and lay out the button
                func(gtx C) D {
                    btn := material.Button(th, &startButton, "Start")
                    return btn.Layout(gtx)
                },
            
            )

            }
        }
    )    
)


```




```go
layout.Flex{
    // Vertical alignment, from top to bottom
    Axis: layout.Vertical,
    //Emtpy space is left at the start, i.e. at the top
    Spacing: layout.SpaceStart,
}.Layout(gtx,
    layout.Rigid(
        func(gtx C) D {
            //We start by defining a set of margins
            margins := layout.Inset{
                Top:    unit.Dp(25),
                Bottom: unit.Dp(25),
                Right:  unit.Dp(35),
                Left:   unit.Dp(35),
            }
            //Then we lay out a layout within those margins ...
            return margins.Layout(gtx,
                // ...the same function we earlier used to create a button
                func(gtx C) D {
                    btn := material.Button(th, &startButton, "Start")
                    return btn.Layout(gtx)
                },
            )
        },
    ),
)
```

## Comments

