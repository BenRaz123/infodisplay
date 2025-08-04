![demo](assets/demo.mp4)

# Info Display

[![Go Report Card](https://goreportcard.com/badge/benraz123/infodisplay)](https://goreportcard.com/report/benraz123/infodisplay) 

## What

the purpose of infodisplay is to generate a looping slideshow that can easily incorporate dynamic data like the current time and other things with javascript that runs in the browser.

Infodisplay uses a custom format that aims for simplicity and terseness.

Here is an example that shows one slide, with a header and 3 bullet points:

```slides
Header
- Bullet One
- Bullet Two
- Bullet Three
```

The first non-bullet point line is inferred to be a header.

This slide would also work, showing one headerline and a text line below it:

```slides
Header
Text Line
```

A system of directives is used for more functionality. For example, to make a slide with the id "hello", a custom duration of 5 seconds, and including a custom image "image.png", write:

```
!id hello
!time 5
!image image.png
Header
Text Line
```

Slides are seperated by whitespace. Thus this makes two slides:

```
Slide One
- Hello
- World

Slide Two
Foo
Bar
```

You can also use `<< >>` to embed a live clock:

```
Today is <<%A>>
- Current hour: <<%H>>
```

And you can use `* *` for bold text:

```
This text is *bold*
```

Most slides have one header, but you can change that by using the `!title` layout

```
!title
Header One
Header Two
```

A slideshow can also have one or more global blocks which can have the following directives:

| Attr | Purpose |
|-|-|
| `!time` | Global slide duration |
| `!exec` | Script to be run when presentation loops (for more see `infodisplay -R!exec` |
| `!styles` | Stylesheet to be used |
| `!pfoot` | Primary Footer |
| `!sfoot` | Secondary Footer |
| `!noautoplay` | Dont run the presentation automatically |

## Running a file

    infodisplay <in> -o <out>

## Getting Help

    infodisplay -r

to get a list of directives

or 
    
    infodisplay -R <specific directive>

