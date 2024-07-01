package main

import "github.com/sniperCore/app/uploader/service"

func main() {
	pService := service.PictureProcessService{}

	// Resize
	input := "/Users/iknow/Works/Go/Project/sniperCore/resources/22"
	output := "/Users/iknow/Works/Go/Project/sniperCore/resources/resize-22-100*80.webp"
	width := 100
	height := 80
	pService.Resize(input, output, width, height)

	// Clip - 剪切

	// Crop
	output1 := "/Users/iknow/Works/Go/Project/sniperCore/resources/crop-22-100*80.png"
	pService.Crop(output, output1, 2000, 1000, 0, 0, "")

	// Compress
	output = "/Users/iknow/Works/Go/Project/sniperCore/resources/compress-imagemagick-22-100*80.png"
	pService.Compress(input, output, "imagemagick", 90)
	output = "/Users/iknow/Works/Go/Project/sniperCore/resources/compress-guetzli-22-100*80.jpg"
	pService.Compress(input, output, "guetzli", 90)

	// Convert
	output = "/Users/iknow/Works/Go/Project/sniperCore/resources/convert-22-100*80.png"
	pService.Convert(input, output, "-channel", "RGB", "-threshold", "100%")

	// format
	output = "/Users/iknow/Works/Go/Project/sniperCore/resources/format-22-100*80.gif"
	pService.Format(input, output, "gif")

}
