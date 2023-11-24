package main

import (
	"github.com/starudream/go-lib/cobra/v2"
	"github.com/starudream/go-lib/core/v2/slog"
	"github.com/starudream/go-lib/core/v2/utils/fmtutil"
	"github.com/starudream/go-lib/core/v2/utils/signalutil"
)

var rootCmd = cobra.NewRootCommand(func(c *cobra.Command) {
	c.Use = "starrail-warp-url"
	c.Run = func(cmd *cobra.Command, args []string) {
		gamePath := detectGamePath()
		if gamePath == "" {
			gamePath = fmtutil.Scan("please input game path (such as: D:\\Star Rail\\Game\\StarRail_Data): ")
		}

		cachePath := detectCachePath(gamePath)
		if cachePath == "" {
			slog.Error("please run game and review warp history first")
			return
		}

		gachaRaw := detectGachaURL(cachePath)
		if gachaRaw == "" {
			slog.Error("please run game and review warp history first")
			return
		}

		// gachaURL, err := url.Parse(gachaRaw)
		// if err != nil {
		// 	slog.Error("warp url is invalid: %v", err)
		// 	return
		// }

		if !checkGachaIsValid(gachaRaw) {
			slog.Error("warp url is expired, please run game and review warp history first")
			return
		}

		startServer(gachaRaw)
	}
	c.PostRun = func(cmd *cobra.Command, args []string) {
		for {
			if fmtutil.Scan("press enter `q` to exit: ") == "q" {
				signalutil.Cancel()
				return
			}
		}
	}
})
