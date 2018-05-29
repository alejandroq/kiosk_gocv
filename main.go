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
	"image/color"
	"log"
	"net/http"
	"time"

	"github.com/machinebox/sdk-go/facebox"

	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
	//"github.com/leprosus/golang-tts"
)

var (
	blue          = color.RGBA{0, 0, 255, 0}
	faceAlgorithm = "haarcascade_frontalface_default.xml"
	stream        *mjpeg.Stream
	fbox          *facebox.Client
)

func main() {

	//create mjpeg stream and to send to web page
	// create the mjpeg stream
	stream = mjpeg.NewStream()

	fbox = facebox.New("http://localhost:8080")

	go kiosk()

	// start http server
	http.Handle("/camera", stream)

	log.Fatal(http.ListenAndServe("localhost:8090", nil))

}

func face(img []byte) {

	if fbox == nil {
		log.Fatal("no fbox :-(")
	}

	faces, err := fbox.Check(bytes.NewReader(img))

	if err != nil {
		log.Printf("unable to recognize face: %v", err)
	}

	if len(faces) > 0 {
		log.Printf("this photo is  %v ", faces[0].Name)
	}

	time.Sleep(500 * time.Millisecond)

}

func kiosk() {

	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		log.Fatalf("error opening web cam: %v", err)
	}
	defer webcam.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	for {
		if ok := webcam.Read(&img); !ok || img.Empty() {
			log.Print("cannot read webcam")
			continue
		}

		buf, err := gocv.IMEncode(".jpg", img)
		if err != nil {
			log.Printf("unable to encode matrix: %v", err)
			continue
		}

		go face(buf)

		stream.UpdateJPEG(buf)

	}
}
