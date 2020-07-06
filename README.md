# palettepull
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT) [![GoDoc](https://godoc.org/github.com/zedseven/palettepull?status.svg)](https://godoc.org/github.com/zedseven/palettepull)

A tool to pull a complete colour palette from a directory of sprites.

## Usage
Usage of the tool is very simple. Just call `palettepull "<path-to-sprite/sprite-directory>"`.

The tool will then generate a PNG image named in the format of `<name-of-source-path-base>Palette.png` in the same
location as the directory/sprite that was specified.

## Example
Using `palettepull "media/dialga-sprite.png"` on...

![Dialga Sprite](media/dialga-sprite.png "The sprite of the Pokémon Dialga.")

...generates the following palette (blown up to be visible) in the same location as the source:

![Dialga Sprite Palette](media/dialga-spritePalette.png "The palette generated from the above sprite (blown up).")
