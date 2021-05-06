# gofontrender

This simple program renders text using the go font render. 
It computes the anti-aliasing by computing the exact 
pixel coverage with the algorithm of Sean Barrett 
described [here](http://nothings.org/gamedev/rasterize/).

FreeType, PathFinder and many other font rendering
library are now using this algorithm because it is fast.

This program allows to explore the effect of different
parameters on font rendering. The result is saved 
in a PNG file which is assumed to be encoded in sRGB
by PNG viewers. 

Some ttf fonts are given for testing and stricly for
private use. The fonts are Microsoft fonts and the
font Magnolia has been obtained from the web site
http://dafont.com where you can find many other fonts
to test.

## Installation

To test the program with the given fonts file, the
most simple is to clone the directory.

`git clone https://github.com/chmike/gofontrender.git`

## Usage

Go in the cloned directory and run `go run main.go`. 
This uses the default parameters. The program will
create a png file with the parameter values in the
file name.

There are many parameters that can be set to control
the rendering. 

```
$ go run main.go --help
Usage of /tmp/go-build2083423820/b001/exe/main:
  -dpi float
    	pixel per inch resolution (default 92)
  -dst string
    	path where to store the output png images (default ".")
  -gamma float
    	correction applied for specified gamma value (default 1)
  -pt float
    	pt size of font (default 16)
  -srgb
    	use sRGB gamma correction
  -text string
    	text to render (default "This is a test")
  -ttf string
    	ttf font file to use (default: goregular)
```

Here is a function call example with some parameters.

`go run main.go -dpi 163 -pt 9 -tth fonts/Verdana.ttf`

The output is black text on white background.

## DPI

Different screens have different pixel density. You can 
find an exhaustive list of screen types and their 
corresponding DPI [here](https://www.sven.de/dpi/).

What is relevant with the DPI value is the distance at
which it matches the retina resolution wich is 0.017°. 

| Size | Resolution | DPI | Dist. (cm) |
|------|------------|-----|------------|
|  24" |  1920x1080 |  92 |  94 |
|  24" |  2560x1440 | 122 |  71 |
|  24" |  3840x2160 | 184 |  48 |
|  27" |  1920x1080 |  82 | 107 |
|  27" |  2560x1440 | 109 |  81 |
|  27" |  3840x2160 | 163 |  53 |
|  32" |  1920x1080 |  70 | 124 |
|  32" |  2560x1440 |  93 |  94 |
|  32" |  3840x2160 | 140 |  64 |

On the desktop computer, the distance between the eyes
and the screen is usually around 55cm. A 27" screen with
a 3840x2160 resolution is thus optimal for the retina 
resolution. 

As you may see from your testing, this is also a pixel 
density that is able to render all details of the font.

There is no visible artifact, even when encoded in sRGB.

When the screen resolution is low, e.g 72 or 92, 
rendering artifacts become visible because the screen 
resolution is too low. The optimal fix to this problem 
is to use high resolution screens. 

## sRGB and gamma encoding

By default, the rendered text is not gamma encoded. This
is an incorrect encoding. The correct encoding is to use
sRGB encoding (`go run main.go -srgb`). This yield
visible anti-aliasiog artifacts at low DPI because
the pixels are too big in respect to the glyph size. 

It is also possible to give a different gamma correction
value. It has been reported 
[here](https://www.puredevsoftware.com/blog/2019/01/22/sub-pixel-gamma-correct-font-rendering/) 
that using a gamma value of 1.43 could reduce the artifacts 
at the price of slightly thickning the glyphs. This is a 
simple hack that gives good results with different fonts 
at low DPI (e.g 92).

The operation performed is to compute pow(gray, 1/gamma)
without sRGB correction applied of course. 
