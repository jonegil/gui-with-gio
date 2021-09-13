---
layout: default
title: Teleprompter
nav_order: 3
has_children: true
has_toc: false
---


# Teleprompter - Animation and interaction

## Goals

In this project we continue from the [egg timer](../egg_timer/) and investigate more features in Gio. 
 - User input, both **keyboard** and **mouse**
 - Programatic **animation**, using scrolling lists
 - Colors and **transparency**
 - **Dynamic text**, both size and paragraph layouts

![Mr_Gorbachev_tear_down_this_wall](teleprompter_Mr_Gorbachev.gif)


As you can see, this is a simple window with scrolling text. But looks deceive, this guy packs a serious punch! We will learn so much in this session, and I'm really excited about what we will get out of this: Seriously cool code, great moments in human history and even a trip to the Smithsonian. Let's (sc)roll.

# TODO: Continue writing from here

## Outline

A [teleprompter](https://en.wikipedia.org/wiki/Teleprompter) is a device that displays text for a presenter to read. From [rolled parchment in a suitcase](https://www.smithsonianmag.com/history/a-brief-history-of-the-teleprompter-88039053/) to modern screens and camera solutions, the core remains the same - display the right text at the righ time.

<p align="center">
  <img src="teleprompter_with_text.jpeg" alt="Teleprompter with text" height="250"/>
  <img src="teleprompter.jpeg" alt="Teleprompter and camera" height="250"/>
</p>

That souns like something we can build, right? Yup, tool like that is perfectly suited for [Gio](www.gioui.org) I'd say. Along the way we'll touch on animation/scrolling, text formatting, fontsize and user control of the prompter with both mouse and keyboard. 

All good techniques for any application. 

Let's roll.

Pun intended.


## Source code
All the source-code is in this repo, in the ```teleprompter/code``` folder.
