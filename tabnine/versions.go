package tabnine

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	rootBinDir = "binaries"
	binName    = "TabNine"
)

type Versions struct {
	configDir string
	Version   string
	System    string
}

func NewVersions(configDir string) (Versions, error) {
	ver, err := getVersion()
	if err != nil {
		return Versions{}, fmt.Errorf("getVersion: %v", err)
	}

	sys, err := getSystem()
	if err != nil {
		return Versions{}, fmt.Errorf("getSystem: %v", err)
	}

	return Versions{
		configDir: configDir,
		Version:   ver,
		System:    sys,
	}, nil
}

func (vs Versions) BinExists(ver, sys string) (bool, error) {
	binPath := filepath.Join(vs.configDir, rootBinDir, vs.Version, vs.System, binName)
	_, err := os.Stat(binPath)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("Stat: %v", err)
	}
	return true, nil
}

func (vs Versions) EnsureLatestBin() (binPath string, err error) {
	latestBinDir := filepath.Join(vs.configDir, rootBinDir, vs.Version, vs.System)
	latestBinPath := filepath.Join(latestBinDir, binName)

	// check if bin exists. Note that we're just checking to avoid downloading
	// data needlessly. This is not a race condition.
	binExists, err := vs.BinExists(vs.Version, vs.System)
	if err != nil {
		return "", fmt.Errorf("BinExists: %v", err)
	}

	if binExists {
		return latestBinPath, nil
	}
	log.Println("downloading TabNine:", latestBinPath)

	if err := os.MkdirAll(latestBinDir, 0755); err != nil {
		return "", fmt.Errorf("MkdirAll: %v", err)
	}

	// https://update.tabnine.com/0.11.2/x86_64-apple-darwin/TabNine
	url := fmt.Sprintf("https://update.tabnine.com/%s/%s/TabNine", vs.Version, vs.System)
	res, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("http Get: %v", err)
	}

	f, err := os.OpenFile(latestBinPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("OpenFile: %v", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, res.Body); err != nil {
		return "", fmt.Errorf("Copy: %v", err)
	}

	if err := f.Sync(); err != nil {
		return "", fmt.Errorf("Sync: %v", err)
	}

	if err := os.Chmod(latestBinPath, 0100); err != nil {
		return "", fmt.Errorf("Chmod: %v", err)
	}

	return latestBinPath, nil
}

func getVersion() (string, error) {
	res, err := http.Get("https://update.tabnine.com/version")
	if err != nil {
		return "", fmt.Errorf("Get: %v", err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("ReadAll: %v", err)
	}

	// seems the API returns a newline after the ver...
	ver := strings.TrimSpace(string(b))
	return ver, nil
}

func getSystem() (string, error) {
	return "x86_64-apple-darwin", nil
}
