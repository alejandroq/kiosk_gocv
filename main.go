package main

import (
	"log"
	"gocv.io/x/gocv"
	"image/color"
	"image"

)

var (
    blue          = color.RGBA{0, 0, 255, 0}
    faceAlgorithm = "haarcascade_frontalface_default.xml"
)

 
func main() {
	webcam, err := gocv.VideoCaptureDevice(0)
	 if err != nil {
			 log.Fatalf("error opening web cam: %v", err)
	 }
	 defer webcam.Close()

	 // prepare image matrix
	 img := gocv.NewMat()
	 defer img.Close()
	 
	 // load classifier to recognize faces
classifier := gocv.NewCascadeClassifier()

defer classifier.Close()

if !classifier.Load(faceAlgorithm) {
		log.Printf("Error reading cascade file: %v\n", faceAlgorithm)
		return
	}


	 // open display window
	 window := gocv.NewWindow("")
	 defer window.Close()

	 for {
			 if ok := webcam.Read(&img); !ok || img.Empty() {
					 log.Print("cannot read webcam")
					 continue
			 }
			 
			 
			 rects := classifier.DetectMultiScale(img)
	 for _, r := range rects {
			 // Save each found face into the file
			 imgFace := img.Region(r)
			 imgFace.Close()

			 // draw rectangle for the face
			 size := gocv.GetTextSize("I don't know you", gocv.FontHersheyPlain, 3, 2)
			 pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			 gocv.PutText(&img, "I don't know you", pt, gocv.FontHersheyPlain, 3, blue, 2)
			 gocv.Rectangle(&img, r, blue, 3)

			 // show the image in the window, and wait 100ms
			 window.IMShow(img)
			 window.WaitKey(100)
		}
	}
}
