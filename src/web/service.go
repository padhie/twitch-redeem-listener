package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"twitch-redeem-trigger/src/config"
	"twitch-redeem-trigger/src/logger"
	"twitch-redeem-trigger/src/output"
)

var (
	videoMutex sync.RWMutex
)

type VideoStatus struct {
	VideoURL string `json:"video_url"`
}

func StartWebServer(cfgWeb config.Web, output output.Device, l *logger.Logger) error {
	if cfgWeb.Enabled {
		return nil
	}

	// main page
	http.HandleFunc("/", serveMainPage)

	// video handling
	http.Handle("/videos/", http.StripPrefix("/videos/", http.FileServer(http.Dir("./videos"))))
	http.HandleFunc("/api/video/current", getCurrentVideoAPI)
	http.HandleFunc("/players", servePlayerPage)

	bindAddr := cfgWeb.Address
	if bindAddr == "" {
		bindAddr = "0.0.0.0"
	}
	addr := fmt.Sprintf("%s:%s", bindAddr, cfgWeb.Port)

	log.Printf("Server available at http://%s", addr)
	return http.ListenAndServe(addr, nil)
}

func serveMainPage(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Twitch Reward</title>
</head>
<body>
	Under Construction
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func servePlayerPage(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Twitch Reward Videos</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            background: transparent;
            overflow: hidden;
            width: 100vw;
            height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        #video-container {
            width: 100%;
            height: 100%;
            display: none;
            align-items: center;
            justify-content: center;
        }
        #video-container.active {
            display: flex;
        }
        video {
            max-width: 100%;
            max-height: 100%;
            object-fit: contain;
        }
    </style>
</head>
<body>
    <div id="video-container">
        <video id="player" preload="auto"></video>
    </div>

    <script>
        const container = document.getElementById('video-container');
        const player = document.getElementById('player');
        let currentVideoID = 0;

        // Check for new videos
        async function checkForVideo() {
            try {
                const response = await fetch('/api/video/current');
                const data = await response.json();

                // Neues Video gefunden
                if (data.video_id && data.video_id !== currentVideoID) {
                    currentVideoID = data.video_id;
                    playVideo(data.video_url);
                }
            } catch (error) {
                console.error('Error checking for video:', error);
            }
        }

        function playVideo(url) {
            console.log('Playing video:', url);
            player.src = url;
            container.classList.add('active');
            player.play();
        }

        // Video ended - hide container
        player.addEventListener('ended', () => {
            console.log('Video ended');
            container.classList.remove('active');
            player.src = '';
        });

        // Error handling
        player.addEventListener('error', (e) => {
            console.error('Video error:', e);
            container.classList.remove('active');
        });

        // Poll for new videos every 500ms
        setInterval(checkForVideo, 500);
        
        // Initial check
        checkForVideo();
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func getCurrentVideoAPI(w http.ResponseWriter, r *http.Request) {
	currentMediaFilePath := "current_media.txt"

	// 1. read current media file
	content, err := os.ReadFile(currentMediaFilePath)
	if err != nil {
		fmt.Printf("Error during load current media file: %v\n", err)
		return
	}

	videoMutex.RLock()
	defer videoMutex.RUnlock()

	status := VideoStatus{
		VideoURL: string(content),
	}

	// 2. return file to open
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)

	// 3. clear file again
	datei, err := os.OpenFile(currentMediaFilePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error during open file: %s\n%v", currentMediaFilePath, err)
		return
	}
	defer datei.Close()
}
