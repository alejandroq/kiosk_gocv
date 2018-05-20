package main

import (
	"log"
	"gocv.io/x/gocv"
	"image/color"
	"image"
	"github.com/machinebox/sdk-go/facebox"
	"fmt"
	"bytes"
	//"github.com/leprosus/golang-tts"


)

var (
    blue          = color.RGBA{0, 0, 255, 0}
    faceAlgorithm = "haarcascade_frontalface_default.xml"
		
)

 
func main() {
	
/*	polly := golang_tts.New("AKIAI5KHZ4W53L55OQGA", "niFon9xIBQe8VNlnFEhOZR0blygpWweLm/9QhQ7S")
	
	polly.Format(golang_tts.MP3)
  polly.Voice(golang_tts.Nicole)
	_, err := polly.Speech("hello janice")
	
	if err!=nil {
		log.Fatalf("speech fails &v",err)
	} */
	

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
    log.Printf("unable to recognize face: %v", err)
}

var caption = "I don't know you"
if len(faces) > 0 {
	log.Printf("more than 0 faces here %v", len(faces))
    caption = fmt.Sprintf("I know you %s", faces[0].Name)
	
}

			 // draw rectangle for the face
			 size := gocv.GetTextSize(caption, gocv.FontHersheyPlain, 3, 2)
			 pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			 gocv.PutText(&img, caption, pt, gocv.FontHersheyPlain, 3, blue, 2)
			 gocv.Rectangle(&img, r, blue, 3)

log.Println("got the rectangle")
			 // show the image in the window
			 window.IMShow(img)
			  window.WaitKey(100)
			 //
			 
			 //read the caption outloud here
			 /*_, err = polly.Speech(caption)
 			
 			if err!=nil {
 				log.Printf("speech fails &v",err)
 			}
			 log.Println("I have spoken") */
			
			 
			 log.Println("I have waited")
		}
	}
}
