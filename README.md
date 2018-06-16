# kiosk_gocv

Software to Install

1.  GoCV and OpenCV:  Follow instructions here:  https://gocv.io/getting-started/
2.  FaceBox from MachineBox.io:   https://machinebox.io/docs/facebox
3.  Clone or Download this Repo
4.  Install npm: brew intall npm
5.  Clone or Download the Front End React Repo:  https://github.com/jumpinjan/kiosk_gocv_front
 

Startup instructions on Mac OS.

1. In a terminal window  Go to gocv.io/x/gocv directory under src in your gopath.    Run this command:  source env.sh

2.  In a new terminal window : Startup facebox by running  these commands:

MB_KEY="<your personal machinebox.io key from when you registered on the site>"

docker run -p 8080:8080 -e "MB_KEY=$MB_KEY" machinebox/facebox

3.  In the first terminal window where you ran source env.sh -  run this file:  go run main.go
4. Open another new terminal window and navigate to the front end react repo and type npm start.
