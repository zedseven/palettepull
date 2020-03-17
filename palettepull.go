package main

import (
	"fmt"
	"github.com/zedseven/coloursorting"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	bitsPerByte    = 8
	targetChannels = 3
)

type fmtInfo struct {
	Model          color.Model
	ChannelsPerPix int
	BitsPerChannel int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("You have to specify a srcPath to a directory.")
		return
	}

	srcPath := os.Args[1]

	palette := make(map[uint]struct{})

	var id bool
	var err error
	if id, err = isDir(srcPath); err != nil {
		fmt.Println("The provided srcPath does not exist.")
		return
	} else if id {
		if err := filepath.Walk(srcPath, func(p string, i os.FileInfo, e error) error {
			if id, err := isDir(p); err == nil && !id {
				return collectColours(&palette, p)
			} else if err != nil {
				return err
			}
			return nil
		}); err != nil {
			fmt.Println("An error occurred while working on the srcPath.", err.Error())
			return
		}
	} else {
		if err := collectColours(&palette, srcPath); err != nil {
			fmt.Println("An error occurred while working on the srcPath.", err.Error())
			return
		}
	}

	paletteColours := make([][3]int, 0, len(palette))
	for k := range palette {
		paletteColours = append(paletteColours, colourToArr(k))
	}
	//fmt.Println(paletteColours)

	sort.Sort(coloursorting.StepSort(paletteColours))

	//fmt.Println(paletteColours)

	srcBase := filepath.Base(srcPath)
	if !id {
		srcBase = strings.TrimSuffix(srcBase, filepath.Ext(srcBase))
	}
	if strings.Contains(srcBase, " ") {
		srcBase += " "
	}
	paletteImgPath := filepath.Join(filepath.Dir(srcPath), srcBase + "Palette.png")
	newFile, err := os.OpenFile(paletteImgPath, os.O_RDWR | os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("An error occurred while creating \"" + paletteImgPath + "\":", err.Error())
		return
	}
	defer func() {
		if err := newFile.Close(); err != nil {
			fmt.Println("An error occurred while writing \"" + paletteImgPath + "\":", err.Error())
		}
	}()
	if err := png.Encode(newFile, paletteToImg(paletteColours)); err != nil {
		fmt.Println("An error occurred while writing image data to \"" + paletteImgPath + "\":", err.Error())
		return
	}

	fmt.Println("Done! The palette image was written to \"" + paletteImgPath + "\".")
}

func isDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, err
	}
	return info.IsDir(), nil
}

func getImgData(img image.Image) (info fmtInfo, channels []uint8) {
	switch img.(type) {
	case *image.RGBA:
		info = fmtInfo{color.RGBAModel, 4, 8}
		simg := img.(*image.RGBA)
		channels = simg.Pix
	case *image.RGBA64:
		info = fmtInfo{color.RGBA64Model, 4, 16}
		simg := img.(*image.RGBA64)
		channels = simg.Pix
	case *image.NRGBA:
		info = fmtInfo{color.NRGBAModel, 4, 8}
		simg := img.(*image.NRGBA)
		channels = simg.Pix
	case *image.NRGBA64:
		info = fmtInfo{color.NRGBA64Model, 4, 16}
		simg := img.(*image.NRGBA64)
		channels = simg.Pix
	case *image.Gray:
		info = fmtInfo{color.GrayModel, 1, 8}
		simg := img.(*image.Gray)
		channels = simg.Pix
	case *image.Gray16:
		info = fmtInfo{color.Gray16Model, 1, 16}
		simg := img.(*image.Gray16)
		channels = simg.Pix
	default:
		info = fmtInfo{color.AlphaModel, 0, 0}
	}

	return
}

func collectColours(palette *map[uint]struct{}, imgPath string) error {
	imgFile, err := os.Open(imgPath)
	if err != nil {
		return err
	}
	img, _, err := image.Decode(imgFile)
	if err == image.ErrFormat {
		return nil
	} else if err != nil {
		return err
	}

	lclPalette := *palette

	info, channels := getImgData(img)
	if info.ChannelsPerPix <= 0 {
		return nil
	}

	for i := 0; i < len(channels) / info.ChannelsPerPix; i++ {
		// Limits to maximum targetChannels channels, and repeats the last channel of the pixel to meet the target
		pix := uint(0)
		for c := 0; c < targetChannels; c++ {
			pix <<= info.BitsPerChannel
			v := c
			if c >= info.ChannelsPerPix {
				v = info.ChannelsPerPix - 1
			}
			pix |= uint(channels[i * info.ChannelsPerPix + v])
		}
		lclPalette[pix] = struct{}{}
	}

	*palette = lclPalette

	return nil
}

func colourToHex(colour uint) string {
	r := uint8(colour >> (2 * bitsPerByte))
	g := uint8((colour >> bitsPerByte) & 0xff)
	b := uint8(colour & 0xff)

	//fmt.Printf("%024b | %d %d %d | %08b %08b %08b\n", colour, r, g, b, r, g, b)

	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

func colourToArr(colour uint) [3]int {
	r := int(uint8(colour >> (2 * bitsPerByte)))
	g := int(uint8((colour >> bitsPerByte) & 0xff))
	b := int(uint8(colour & 0xff))

	return [3]int{r, g, b}
}

func paletteToImg(palette [][3]int) image.Image {
	paletteLen := len(palette)
	outputWidth := int(math.Sqrt(float64(paletteLen)))
	outputHeight := paletteLen / outputWidth
	if paletteLen % outputWidth != 0 {
		outputHeight += 1
	}
	img := image.NewRGBA(image.Rect(0, 0, outputWidth, outputHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Transparent}, image.Pt(0, 0), draw.Src)

	for i := 0; i < paletteLen; i++ {
		x := i % outputWidth
		y := i / outputWidth
		img.Set(x, y, color.RGBA{uint8(palette[i][0]), uint8(palette[i][1]), uint8(palette[i][2]), 255})
	}

	return img
}