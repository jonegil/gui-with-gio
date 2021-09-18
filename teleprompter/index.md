---
layout: default
title: Teleprompter
nav_order: 3
has_children: true
has_toc: false
---


# Teleprompter - Animation and interaction

```

Hi mate!

Please look around while I continue to write this chapter. 

Code is done though - and honestly quite nifty.
Pull this repo and play around.

Cheers

```

## Goals

This project continues where the [egg timer](../egg_timer/) leaves off. The timer was a good start and gave us the foundation to build an app. Hence those details won't be repeated here. But we're not done. Especially we should look closer into how to react to user input from mouse and keyboard, and use that do control an animation.  

For learning it's best with a small codebase. That way it's easier to get an overview. Still it should have many novel and interesting pieces. To achieve that it's best to start a new project and investigate more features in Gio: 
 - User input, both **keyboard** and **mouse**
 - Programatic **animation**, using lists
 - Colors and **transparency**
 - **Dynamic text**, both size and paragraph layouts

![Mr_Gorbachev_tear_down_this_wall](teleprompter_Mr_Gorbachev.gif)

This is a simple window with scrolling text. Nothing fancy. But play with the controls and you'll see it's both lively and responsive. And that's worth diving into.

Let's (sc)roll!

(sorry)

## Outline

A [teleprompter](https://en.wikipedia.org/wiki/Teleprompter) is a device that displays text for the presenter to read. From [rolling parchment in a suitcase](https://www.smithsonianmag.com/history/a-brief-history-of-the-teleprompter-88039053/) to modern screens and camera solutions, the core remains the same - display the right text at the right time.

<p align="center">
  <img src="teleprompter_with_text.jpeg" alt="Teleprompter with text" height="250"/>
  <!--img src="teleprompter.jpeg" alt="Teleprompter and camera" height="250"/-->
</p>

## So what will we actually build?
This sounds like something we could build ourselves. Here's the plan:

 1. Read and display content from an external ```txt``` file.
 1. Display the text in an efficient manner. That's done by breaking the full text into smaller **paragraphs**
 1. Present each paragraph as a separate [material.Label](https://pkg.go.dev/gioui.org/widget/material#Label) widget.
 1. Organize the labels into a **list of widgets** using [layout.List](https://pkg.go.dev/gioui.org/layout?utm_source=gopls#List)
 1. Build **auto scrolling**, including start, stop and live speed adjustment
 1. Add manual scrolling, for full control using both **keyboard** and **mouse**
 1. Allow the user to **resize** and **layout** the text exactly as wanted, live, when scrolling
 1. Add custom graphics to create a **transparent** focusbar, that we **move at will**, making it easier to read the right line.

Sounds like we have our work cut out for us? Not to worry, we'll tackle them step by step.

## Source code

Todo: 
```go
Describe the overall structure of the program, with the loop listening for various events

```


```go
The go through each main block in detail
 keystroke
 mouse
 rendering 
 quit

```

All the source-code is in this repo, in the ```teleprompter/code``` folder.
