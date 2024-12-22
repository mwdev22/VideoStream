package cam

import (
	"flag"
	"log"
	"os"

	"github.com/blackjack/webcam"
)

var (
	dev = flag.String("d", "/dev/video0", "video device to use")
)

func New() (*webcam.Webcam, error) {
	cam, err := webcam.Open(*dev)
	if err != nil {
		log.Fatalf("error opening camera: %v", err)
		os.Exit(1)
	}

	cam.SetFramerate(30)
	format_desc := cam.GetSupportedFormats()
	var format webcam.PixelFormat

	for k, s := range format_desc {
		if s == "Motion-JPEG" {
			format = k
		}
	}
	_, _, _, err = cam.SetImageFormat(format, 640, 480)
	if err != nil {
		log.Fatalf("error setting image format: %v", err)
	}

	err = cam.StartStreaming()
	if err != nil {
		log.Fatalf("error starting camera streaming: %v", err)
	}

	return cam, nil
}
