---
layout: default
title: Chapter 7 
nav_order: 2
parent: Egg timer
has_children: false
---

# Chapter 7 - Progressbar

## Goals
The intention of this section is to add a progressbar

![Progressbar](07_progressbar.gif)

## Outline

This chapter is a highlight, we will do so many cool things here.
 - Introduce the **material.Progressbar** component
 - Add some state variables to control behaviour
 - Setup a separate channel that generates ticks to progress the progressbar

Let's look at the new pieces to the puzzle.

## Generate and listen to ticks

For the progressbar, it's good to be able to process it precisely. To do that, a separate go-routine that ticks with a steady beat is a good solution. By doing that we get
1. Precise timing, where we control the execution exactly as we want it
1. Consistent timing, simlar across fast and slow hardware
1. Concurrent timing, the rest of the application continues as before


### Code
```go

```

## Comments

