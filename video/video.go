package video

import (
	"io"
	"log"
	"mime"
	"os"
	"strings"

	"github.com/dsnet/try"
	"github.com/kkdai/youtube/v2"
)

// quality:720p
func DonwloadVideo(videoId string, quality string) error {
	var client youtube.Client
	video, err := client.GetVideo(videoId)
	if err != nil {
		return err
	}

	format := video.Formats.FindByQuality(quality)
	log.Printf("found format:%s, size:%.2f, mimetype:%s", format.QualityLabel, float64(format.ContentLength)/1024/1024/1024, format.MimeType)
	stream, _, err := client.GetStream(video, format)
	if err != nil {
		return err
	}

	mediatype, _ := try.E2(mime.ParseMediaType(format.MimeType))
	ext := strings.Split(mediatype, "/")[1]

	file, err := os.Create(videoId + "." + ext)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return err
	}

	return nil
}
