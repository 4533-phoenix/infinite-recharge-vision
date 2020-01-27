package main

import (
	"fmt"
	"image"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"gocv.io/x/gocv"
)

var config *Config

func main() {
	// Parse flags
	pflag.Parse()

	// Initialize logging

	// Load configuration
	config, err := LoadConfig(configPath)

	if err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	// Open Capture Device
	capture, err := gocv.VideoCaptureFile(config.Capture.Device)

	if err != nil {
		log.Fatalf("error intializing capture device: %v", err)
	}

	capture.Set(gocv.VideoCaptureFrameHeight, float64(config.Capture.Height))
	capture.Set(gocv.VideoCaptureFrameWidth, float64(config.Capture.Width))

	// Capture frame.
	frame := gocv.NewMat()
	defer frame.Close()

	// HSV copy of the capture frame.
	hsv := gocv.NewMat()
	defer frame.Close()

	// Tracking frame.
	tracking := gocv.NewMat()
	defer tracking.Close()

	// The hue channel from the HSV image.
	hueThreshold := gocv.NewMat()
	defer hueThreshold.Close()

	// The saturation channel from the HSV image.
	satThreshold := gocv.NewMat()
	defer satThreshold.Close()

	// The value channel from the HSV image.
	valThreshold := gocv.NewMat()
	defer valThreshold.Close()

	// Collection of vectors that identifies all the circles found in the frame.
	circles := gocv.NewMat()
	defer circles.Close()

	hMin := gocv.Scalar{Val1: config.Threshold.MinHue}
	hMax := gocv.Scalar{Val1: config.Threshold.MaxHue}
	sMin := gocv.Scalar{Val1: config.Threshold.MinSaturation}
	sMax := gocv.Scalar{Val1: config.Threshold.MaxSaturation}
	vMin := gocv.Scalar{Val1: config.Threshold.MinValue}
	vMax := gocv.Scalar{Val1: config.Threshold.MaxValue}

	kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(5, 5))

	ticker := time.NewTicker(config.Capture.Interval)

	for _ = range ticker.C {
		if ok := capture.Read(&frame); ok {
			// Process the captured frame and convert it to HSV format.
			gocv.MedianBlur(frame, &frame, config.Morph.Blur)
			gocv.CvtColor(frame, &hsv, gocv.ColorBGRToHSV)

			// Process the capture frame into a threshold image.
			channels := gocv.Split(hsv)
			gocv.InRangeWithScalar(channels[0], hMin, hMax, &hueThreshold)
			gocv.InRangeWithScalar(channels[1], sMin, sMax, &satThreshold)
			gocv.InRangeWithScalar(channels[2], vMin, vMax, &valThreshold)

			gocv.BitwiseAnd(hueThreshold, satThreshold, &tracking)
			gocv.BitwiseAnd(tracking, valThreshold, &tracking)

			// Morph the tracking image to remove noise an potential for false
			// positives.
			gocv.Erode(tracking, &tracking, kernel)

			// Find all the circles in the image.
			gocv.HoughCirclesWithParams(
				tracking,
				&circles,
				gocv.HoughGradient,
				config.Transform.DP,
				config.Transform.MinDist,
				config.Transform.Param1,
				config.Transform.Param2,
				config.Transform.MinRadius,
				config.Transform.MaxRadius,
			)

			// If no circles were found then continue to the next frame.
			if circles.Cols() < 0 {
				continue
			}

			var closest gocv.Vecf

			for i := 0; i < circles.Cols(); i++ {

				// Get the circles (x, y, r) information.
				vector := circles.GetVecfAt(0, i)

				// If all of the information was returned, then determine which
				// circle is the 'closest'.
				if len(vector) > 2 {
					radius := vector[2]
					if radius > closest[2] {
						closest = vector
					}
				}
			}

			// Publish the closest circle
			fmt.Printf("%v\n", closest)
		}
	}
}
