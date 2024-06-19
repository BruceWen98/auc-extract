package main

import (
	"auc-video-extract/video"
	"log"
	"os"
	"sync"

	"github.com/dsnet/try"
	"github.com/kkdai/youtube/v2"
	"github.com/zhiminwen/quote"
)

func init() {
	os.Setenv("MAGEFILE_VERBOSE", "true")
}

// arg: youtube video ID, list video info such as size, quality, etc.
func T01_youtube_video_info(videoID string) {
	client := youtube.Client{}
	video := try.E1(client.GetVideo(videoID))

	for _, v := range video.Formats {
		log.Printf("quality:%s,%s, size=%.2fGB, mimetype=%s", v.QualityLabel, v.Quality, float64(v.ContentLength)/1024/1024/1024, v.MimeType)
	}
}

func T02_download_video(videoID string) {
	err := video.DonwloadVideo(videoID, "720p")
	if err != nil {
		log.Fatalf("Failed to download video:%v", err)
	}
}

func T03_download_all() {
	list := quote.Line(`
		_LVdvPzBKTk
		COWCWO8993g
		IEVT8oYrgfM
		aZBChhViKwg
		Ylqf_arQ49k
		drm16gIfReQ
		93VRp_yb46s
		rOiaWTlhuik
		aPczPCkBEhU
		dCmRpdNUUO8
		08itPVCVS7A
		M9g6bh02mrQ
		U05hBP1U_Bw
		tP8WydkkaMM
		pXzQVSKJW8E
		kQXuc7F9TfI
		u4i0tzxuCTc
		bNnCT5JHY_s
		TM-Jhb9wmIY
	`)

	ch := make(chan string)

	go func() {
		var wg sync.WaitGroup
		threads := 5

		for i := 0; i < threads; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for id := range ch {
					log.Printf("downloading video:%s", id)
					err := video.DonwloadVideo(id, "720p")
					if err != nil {
						log.Printf("Failed to download video:%v", err)
					}
					log.Printf("downloaded video:%s", id)
				}
			}()
		}

		wg.Wait()
	}()

	for _, id := range list {
		ch <- id
	}
	close(ch)

}
