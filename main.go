package main

/***
Startup instructions on Mac OS.

1. In a terminal window  Go to gocv.io/x/gocv directory under src in your gopath.    Run this command:  source env.sh

2.  In a new terminal window : Startup facebox by running  these commands:

MB_KEY="<your personal machinebox.io key from when you registered on the site>"

docker run -p 8080:8080 -e "MB_KEY=$MB_KEY" machinebox/facebox

3.  In the first terminal window where you ran source env.sh -  run this file:  go run main.go

*/

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"net/http"

	"github.com/hybridgroup/mjpeg"
	"github.com/machinebox/sdk-go/facebox"
	"gocv.io/x/gocv"
	//"github.com/leprosus/golang-tts"
)

var (
	blue          = color.RGBA{0, 0, 255, 0}
	faceAlgorithm = "haarcascade_frontalface_default.xml"
	stream        *mjpeg.Stream
)

func main() {

	/*polly := golang_tts.New("AKIAI5KHZ4W53L55OQGA", "niFon9xIBQe8VNlnFEhOZR0blygpWweLm/9QhQ7S")
		log.Printf("polly is %v", polly)

		polly.Format(golang_tts.MP3)
	  polly.Voice(golang_tts.Nicole)
	  log.Printf("polly is %v", polly)
	 data , err := polly.Speech("hello janice")


		if err!=nil {
			log.Fatalf("speech fails &v",err)
		}

		if (1==1)   {return } */

	//create mjpeg stream and to send to web page
	// create the mjpeg stream
	stream = mjpeg.NewStream()

	go kiosk()

	// start http server
	http.Handle("/camera", stream)

	log.Fatal(http.ListenAndServe("localhost:8090", nil))

}

func kiosk() {

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
	//window := gocv.NewWindow("find me")
	//log.Println("window opened for display")
	//defer window.Close()

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
			//window.IMShow(img)
			//window.WaitKey(100)
			//

			buf2, _ := gocv.IMEncode("*.jpg", img)
			stream.UpdateJPEG(buf2)
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
