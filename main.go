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
	"encoding/json"
	"image/color"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/machinebox/sdk-go/facebox"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
)

var (
	blue          = color.RGBA{0, 0, 255, 0}
	faceAlgorithm = "haarcascade_frontalface_default.xml"
	stream        *mjpeg.Stream
	fbox          *facebox.Client
	c1            = make(chan bool)
)

func main() {

	//create mjpeg stream and to send to web page
	// create the mjpeg stream
	stream = mjpeg.NewStream()

	router := mux.NewRouter()

	fbox = facebox.New("http://localhost:8080")

	go kiosk()

	// start http server
	router.Handle("/camera", stream)
	log.Println("camera routed")

	router.HandleFunc("/face", face)
	router.HandleFunc("/audio/name/{name}", audioGreeting)

	log.Fatal(http.ListenAndServe("localhost:8090", router))

}

func kiosk() {

	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		log.Fatalln("can't find camera")
	}

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

		stream.UpdateJPEG(buf)

	}
}

type jsonface struct {
	StudentName    string `json:studentname`
	CounselorName  string `json:counselorname`
	CounselorImage string `json:counselorimage`
}

func face(w http.ResponseWriter, r *http.Request) {

	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		log.Fatalln("can't find camera")
	}

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	if ok := webcam.Read(&img); !ok || img.Empty() {
		log.Print("cannot read webcam")
		return
	}

	buf, err := gocv.IMEncode(".jpg", img)
	if fbox == nil {
		log.Fatal("no fbox :-(")
	}

	faces, err := fbox.Check(bytes.NewReader(buf))

	if err != nil {
		log.Printf("unable to recognize face: %v", err)
	}

	if len(faces) > 0 {
		log.Printf("this photo is  %v ", faces[0].Name)

	}

	faceName := ""
	image := ""
	counselorName := ""

	if len(faces[0].Name) == 0 {
		faceName = "Who are you?"
		image = "none.jpg"
		counselorName = "Nope"
	} else {
		faceName = faces[0].Name
		if strings.ToLower(string(faceName[0])) < "k" {
			image = "wink.jpg"
			counselorName = "Wink"
		} else {
			image = "lizzie.jpg"
			counselorName = "Lizzie"
		}
	}

	faceJSON := jsonface{StudentName: faceName, CounselorImage: image, CounselorName: counselorName}

	log.Println("faceJSON has  ", faceJSON)
	jData, err := json.Marshal(faceJSON)
	if err != nil {
		log.Fatalln("problem marshalling json", err)
		return
	}
	log.Println("face json is:", string(jData))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)

}

func audioGreeting(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["name"]
	log.Println("Generating text-to-speech for the name " + name)

	// Initialize a session that the SDK uses to load credentials from the shared credentials file. (~/.aws/credentials).
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	pollyService := polly.New(sess)
	textToSpeak := "welcome " + name
	input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"), Text: aws.String(textToSpeak), VoiceId: aws.String("Nicole")}

	output, err := pollyService.SynthesizeSpeech(input)
	if err != nil {
		log.Println("Error calling SynthesizeSpeech: ")
		log.Print(err.Error())
		w.WriteHeader(500)
		w.Write([]byte("Error synthesizing text " + http.StatusText(500)))
	}

	if _, err := io.Copy(w, output.AudioStream); err != nil {
		log.Println("Error reading mp3: ")
		log.Print(err.Error())
		w.WriteHeader(500)
		w.Write([]byte("Error reading mp3 " + http.StatusText(500)))
	}
}
