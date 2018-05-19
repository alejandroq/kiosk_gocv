package main

import (
	"log"
	"gocv.io/x/gocv"
	"image/color"
	"image"
	"github.com/machinebox/sdk-go/facebox"
	"fmt"
	"bytes"


)

var (
    blue          = color.RGBA{0, 0, 255, 0}
    faceAlgorithm = "haarcascade_frontalface_default.xml"
)

 
func main() {
	
	fbox := facebox.New("http://localhost:8080")
	log.Println("got facebox")
	
	
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
log.Println("classifer opened")
if !classifier.Load(faceAlgorithm) {
		log.Printf("Error reading cascade file: %v\n", faceAlgorithm)
		return
	}


	 // open display window
	 window := gocv.NewWindow("find me")
	 log.Println("window opened for display")
	 defer window.Close()

	 for {
			 if ok := webcam.Read(&img); !ok || img.Empty() {
					 log.Print("cannot read webcam")
					 continue
			 }
			 
			 log.Println("read  from webcam")
			 rects := classifier.DetectMultiScale(img)
			 log.Println("have rectangles")
	 for _, r := range rects {
			 // Save each found face into the file
			 
			 imgFace := img.Region(r)
			 log.Println("got region")
			 defer imgFace.Close()
			 
			 buf, err := gocv.IMEncode(".jpg", imgFace)
if err != nil {
    log.Printf("unable to encode matrix: %v", err)
    continue
}

faces, err := fbox.Check(bytes.NewReader(buf))


if err != nil {
    log.Println("unable to recognize face: %v", err)
}

log.Println("have faces %v",faces)

var caption = "I don't know you"
if len(faces) > 0 {
    caption = fmt.Sprintf("I know you %s", faces[0].Name)
}

			 // draw rectangle for the face
			 size := gocv.GetTextSize("I don't know you", gocv.FontHersheyPlain, 3, 2)
			 pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			 gocv.PutText(&img, caption, pt, gocv.FontHersheyPlain, 3, blue, 2)
			 gocv.Rectangle(&img, r, blue, 3)

			 // show the image in the window, and wait 100ms
			 window.IMShow(img)
			 window.WaitKey(100)
		}
	}
}
