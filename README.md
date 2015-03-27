ims
===

ims is a standalone image manipulation and serving service written in Go. ims
provides on-the-fly resizing, cropping, rotation etc. (See full docs for an
explanation of all available filters).

ims uses the excellent [gift](https://github.com/disintegration/gift) library
under the hood to perform all image manipulation.



## Work in progress

ims is very much a work in progress and isn't ready for production use. Please
treat it as such and file issues where appropriate. As always, pull requests or
other offers of help are greatly appreciated.



## Usage

Until a release build of ims is made, you'll need to build it if you want to use it:

```
$ go build
```

Once built, you just need to point the binary at an existing directory full of images.

```
$ ./ims --source ./myimages
```

And then all images within the folder should be available at
http://localhost:5995. Assuming you have an image called `apple.png`, applying
filters is simple.


### Resizing

To resize an image you can use a URL like:

```
http://localhost:5995/apple.png?resize=500,0
```

Where the first parameter is the desired width, and the second is the desired
height. Either value can be set to `0` (not both) to maintain aspect ratio
during resize.


### Cropping

To crop an image you can use a URL like:

```
http://localhost:5995/apple.png?crop=0,0,100,200
```

Where the first 2 parameters are the top left corner of the box to use for
cropping, the second 2 values are the bottom corner of the box.


### Chaining filters

Filters can be chained arbitrarily in the URL like so:

```
http://localhost:5995/apple.png?resize=500,0&crop=0,0,100,200&resize=100,100
```

Filters will be applied in the order they are specified.



## TODO

- Docs
- Images are always written to jpeg format at present, this needs to be format specific.
- Expose more manipulations/filters
- Unit/integration tests
- Release testing builds



## Copyright & License

- Copyright © 2015 Patrick Carey (https://github.com/paddycarey)
- Licensed under the **MIT** license.
