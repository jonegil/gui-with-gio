---
layout: default
title: Teleprompter
nav_order: 3
has_children: true
has_toc: false
---


# Teleprompter - Animation and interaction

## Goals

In this project we continue from the [egg timer](../egg_timer/) and investigate more features in Gio: 
 - Programatic **animation**, using lists
 - User input, both **keyboard** and **mouse**
 - Colors and **transparency**
 - **Dynamic text**, both size and paragraph layouts

![Mr_Gorbachev_tear_down_this_wall](teleprompter_Mr_Gorbachev.gif)


As you can see, this is a simple window with scrolling text. But looks deceive, this guy packs a serious punch. Let's (sc)roll!


## Outline

A [teleprompter](https://en.wikipedia.org/wiki/Teleprompter) is a device that displays text for a presenter to read. From [rolling parchment in a suitcase](https://www.smithsonianmag.com/history/a-brief-history-of-the-teleprompter-88039053/) to modern screens and camera solutions, the core remains the same - display the right text at the right time.

<p align="center">
  <img src="teleprompter_with_text.jpeg" alt="Teleprompter with text" height="250"/>
  <img src="teleprompter.jpeg" alt="Teleprompter and camera" height="250"/>
</p>

To build this in Gio we will 
 - convert a long ```.txt``` file to a list of short **paragraphs**
 - use [layout.List](https://pkg.go.dev/gioui.org/layout?utm_source=gopls#List) to display the paragraphs
 - allow full user control to directly **resize**, **layout** and **scroll** the text
 - add logic to **start**, **stop** and **adjust** speed
 - add custom graphics to create a **transparent** focusbar, making it easier to read


## Source code
All the source-code is in this repo, in the ```teleprompter/code``` folder.
