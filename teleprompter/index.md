---
layout: default
title: Teleprompter
nav_order: 3
has_children: true
has_toc: false
---


# Teleprompter

## Goals

In this project we continue from the [egg timer](../egg_timer/) and investigate other features of Gio. In particular how to listen to inputs from the user, both keyboard and mouse. 

## Outline

A [teleprompter](https://en.wikipedia.org/wiki/Teleprompter) is a device that displays text for a presenter to read. From their [humble beginings](https://www.smithsonianmag.com/history/a-brief-history-of-the-teleprompter-88039053/) as rolled parchment in a suitcase, teleprompters have come a long way, where professional versions costs [tens of thoursans](https://www.bhphotovideo.com/c/product/1576296-REG/acebil_pro_s21thbkit2_21_studio_wuxga_prompter.html) of dollars. 

<img src="teleprompter.jpeg" alt="Teleprompter" width="300"/>

Low cost versons exists though, like the one above, that replace the screed with a tablet or phone. And even pure software versions that use our devices for all parts of the job. Which is exactly what we're gioing to build, using Go and [Gio](www.gioui.org) to build our own, totally free, teleprompter.


## Source code
All the source-code is in this repo, in the ```teleprompter/code``` folder.
