ims
===

ims is a standalone image manipulation and serving service written in Go. ims
provides on-the-fly resizing, cropping, rotation etc. (See full docs for an
explanation of all available filters).

ims uses the excellent [gift](https://github.com/disintegration/gift) library
under the hood to perform all image manipulation.


### Work in progress

ims is very much a work in progress and isn't ready for production use. Please
treat it as such and file issues where appropriate. As always, pull requests or
other offers of help are greatly appreciated.


### TODO

- Docs
- Images are always written to jpeg format at present, this needs to be format specific.
- Expose more manipulations/filters
- Unit/integration tests
- Release testing builds


### Copyright & License

- Copyright © 2015 Patrick Carey (https://github.com/paddycarey)
- Licensed under the **MIT** license.
