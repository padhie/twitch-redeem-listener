package web

import (
	"net/http"
	"twitch-redeem-trigger/config"
	"twitch-redeem-trigger/output"

	"twitch-redeem-trigger/logger"
)

func StartServer(cfgWeb config.Web, output output.Device, l *logger.Logger) error {
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
            <h1>Twitch Redeem Trigger</h1>
            <p>Status: Running</p>
            <form action="/toggle" method="post">
                <button type="submit">Toggle Tapo Manually</button>
            </form>
        `))
	})

	http.HandleFunc("/toggle", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := output.Toggle(); err != nil {
			l.Error("Manual toggle failed: %v", err)
			http.Error(w, "Failed to toggle output", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/status", http.StatusSeeOther)
	})

	l.Info("Starting web server on port %s", cfgWeb.Port)
	return http.ListenAndServe(":"+cfgWeb.Port, nil)
}
