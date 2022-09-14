package dto

import(
	"os"
	"log"

	"image/png"
	"image"
)



// loading image
func Load(filePath string) *image.NRGBA {
	imgFile, err := os.Open(filePath)
	if err != nil {
		log.Println("Cannot read file:", err)
	}
	
	img, _, err := image.Decode(imgFile)
	if err != nil {
		log.Println("Cannot decode file:", err)
	}
	defer imgFile.Close()
	return img.(*image.NRGBA)
}
// saving image
func Save(filePath string, img *image.NRGBA) {
	imgFile, err := os.Create(filePath)
	if err != nil {
		log.Println("Cannot create file:", err)
	}
	png.Encode(imgFile, img.SubImage(img.Rect))
	defer imgFile.Close()
}
