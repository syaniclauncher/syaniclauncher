package launcher_cache

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"syaniclauncher/src/helpers"
	"syaniclauncher/src/launcher_config"
)

const playerHeadTTL = 2 * time.Hour

var playerHeadMu sync.Mutex

// some public apis to get minecraft heads
var playerHeadSources = map[string]string{
	"mcheads":  "https://api.mcheads.org/head/%s/8",
	"crafatar": "https://crafatar.com/avatars/%s?size=8",
}

func playerHeadDir() (string, error) {
	appDir := launcher_config.AppDir
	if appDir == "" {
		appData, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		appDir = filepath.Join(appData, ".syaniclauncher")
	}

	dir := filepath.Join(appDir, "cache", "player_heads")
	return dir, os.MkdirAll(dir, 0755)
}

func GetPlayerHead(uuid string) ([]byte, error) {
	uuid = helpers.RemoveHyphensUUID(uuid)
	if uuid == "" {
		return nil, fmt.Errorf("player UUID is required")
	}

	playerHeadMu.Lock()
	defer playerHeadMu.Unlock()

	dir, err := playerHeadDir()
	if err != nil {
		return nil, err
	}
	now := time.Now().Unix()
	prefix := uuid + "."

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasPrefix(name, prefix) || !strings.HasSuffix(name, ".png") {
			continue
		}
		expiresText := strings.TrimSuffix(strings.TrimPrefix(name, prefix), ".png")
		expires, parseErr := strconv.ParseInt(expiresText, 10, 64)
		path := filepath.Join(dir, name)
		if parseErr == nil && expires > now {
			data, readErr := os.ReadFile(path)
			if readErr == nil {
				return data, nil
			}
			continue
		}
		_ = os.Remove(path)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	var downloadErr error
	for _, source := range playerHeadSources {
		url := fmt.Sprintf(source, uuid)
		response, err := client.Get(url)
		if err != nil {
			downloadErr = err
			continue
		}
		if response.StatusCode != http.StatusOK {
			downloadErr = fmt.Errorf("HTTP %d", response.StatusCode)
			response.Body.Close()
			continue
		}
		data, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
		response.Body.Close()
		if err != nil {
			downloadErr = err
			continue
		}
		if len(data) < 8 || string(data[:8]) != "\x89PNG\r\n\x1a\n" {
			downloadErr = fmt.Errorf("invalid PNG response")
			continue
		}

		expires := time.Now().Add(playerHeadTTL).Unix()
		path := filepath.Join(dir, fmt.Sprintf("%s.%d.png", uuid, expires))
		if err := os.WriteFile(path, data, 0644); err != nil {
			return nil, err
		}
		return data, nil
	}

	if downloadErr == nil {
		downloadErr = fmt.Errorf("no player head sources available")
	}
	return nil, downloadErr
}
