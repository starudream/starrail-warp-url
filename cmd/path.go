package main

import (
	"bufio"
	"cmp"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"

	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/starudream/go-lib/core/v2/slog"
)

const NAME = "崩坏：星穹铁道"

func detectGamePath() string {
	// C:\Users\...\AppData\Roaming
	appDataPath := os.Getenv("APPDATA")
	if appDataPath == "" {
		slog.Error("APPDATA is not defined")
		return ""
	}

	// C:\Users\...\AppData\LocalLow\miHoYo\崩坏：星穹铁道
	gameLogPath := filepath.Join(appDataPath, "..", "LocalLow", "miHoYo", NAME)

	playerLogPath := filepath.Join(gameLogPath, "Player.log")
	slog.Info("detecting game player log file: %s", playerLogPath)

	playerLogFile, err := os.Open(playerLogPath)
	if err != nil {
		slog.Error("cannot open game player log file: %v", err)
		return ""
	}
	defer gh.Close(playerLogFile)

	// Loading player data from .../Game/StarRail_Data/data.unity3d
	scanner := bufio.NewScanner(playerLogFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Loading player data from ") && strings.HasSuffix(line, "/StarRail_Data/data.unity3d") {
			gamePath := strings.TrimPrefix(line, "Loading player data from ")
			gamePath = strings.TrimSuffix(gamePath, "/data.unity3d")
			return gamePath
		}
	}

	slog.Error("cannot detect game path from player log file")
	return ""
}

func detectCachePath(gamePath string) string {
	// ...\Game\StarRail_Data\webCaches\2.18.0.0\Cache\Cache_Data\data_2
	cachePath := filepath.Join(gamePath, "webCaches")

	entries, err := os.ReadDir(cachePath)
	if err != nil {
		slog.Error("cannot read game cache path: %v", err)
		return ""
	}

	if len(entries) == 0 {
		slog.Error("cannot detect game cache path")
		return ""
	}

	slices.SortFunc(entries, func(a, b os.DirEntry) int { return cmp.Compare(a.Name(), b.Name()) })

	current := entries[len(entries)-1].Name()

	return filepath.Join(cachePath, current, "Cache", "Cache_Data")
}

func detectGachaURL(cachePath string) string {
	filePath := filepath.Join(cachePath, "data_2")
	slog.Info("detecting game warp data file: %s", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		slog.Error("cannot read warp data file: %v", err)
		if isWindows() && strings.Contains(err.Error(), "used by another process") {
			data = tryCopyRead(filePath)
			if len(data) > 0 {
				goto next
			}
		}
		return ""
	}

next:

	blocks := strings.Split(string(data), "\u0000")
	for i := len(blocks) - 1; i >= 0; i-- {
		if strings.HasPrefix(blocks[i], "1/0/http") && strings.Contains(blocks[i], "getGachaLog") {
			return blocks[i][4:]
		}
	}

	return ""
}

// only work for windows, just use `copy`
func tryCopyRead(path string) []byte {
	tmp := filepath.Join(os.TempDir(), "data_2_"+time.Now().Format("20060102150405"))

	slog.Info("attempting to copy file to temp: %s", tmp)

	cmd := exec.Command("powershell.exe", "Copy-Item", strconv.Quote(path), "-Destination", strconv.Quote(tmp))
	slog.Debug("copy command: %s", cmd.String())

	output, err := cmd.Output()
	if len(output) > 0 {
		gbk, _ := simplifiedchinese.GBK.NewDecoder().Bytes(output)
		slog.Info("copy output: %s", gbk)
	}
	if err != nil {
		slog.Error("cannot copy file to temp: %v", err)
		return nil
	}

	data, err := os.ReadFile(tmp)
	if err != nil {
		slog.Error("cannot read temp file: %v", err)
		return nil
	}

	return data
}
