package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/leeola/gokakoune/api"
	tnp "github.com/leeola/gokakoune/plugins/tabnine"
	"github.com/leeola/gokakoune/tabnine"
)

func main() {
	var metaCmd string
	if len(os.Args) == 2 {
		metaCmd = os.Args[1]
	}

	switch metaCmd {
	case "http-serve":
		if err := httpServeMain(); err != nil {
			log.Fatalf("httpServeMain: %v", err)
		}
		return
	case "http-serve-background":
		if err := httpServeBackgroundMain(); err != nil {
			log.Fatalf("httpServeMain: %v", err)
		}
		return
	}

	if err := gokakouneMain(); err != nil {
		log.Fatalf("gokakouneMain: %v", err)
		return
	}
}

func gokakouneMain() error {
	k := api.New()
	if err := tnp.Plugin(k); err != nil {
		return fmt.Errorf("gokakoune: %v", err)
	}

	return nil
}

func httpServeBackgroundMain() error {
	// use the kak home as the gokakoune-tabnine home as well.
	//
	// NOTE: i'm not sure if Tabnine is best run per project, or for
	// everywhere. Need to contact the author, but currently it seems
	// like Sublime uses a single tabnine process for everything,
	// so that's what i'm basing this off of.
	kakConfigDir := os.Getenv("kak_config")
	if kakConfigDir == "" {
		return fmt.Errorf("missing kak_config value")
	}

	configDir := filepath.Join(kakConfigDir, "tabnine")

	// the bin path and args is *this* bin. That is to say,
	// we're spinning off a background process of the command
	// `gokakoune-tabnine http-serve`. We can use os.Args[0]
	// to refer to whatever this bin is, no need to configure
	// that.
	binPath := os.Args[0]
	binArgs := []string{"http-serve"}

	bp, err := tabnine.NewProcess(configDir, binPath, binArgs)
	if err != nil {
		return fmt.Errorf("NewProcess: %v", err)
	}

	if err := bp.EnsureStarted(); err != nil {
		return fmt.Errorf("EnsureStarted: %v", err)
	}

	return nil
}

func httpServeMain() error {
	// use the kak home as the gokakoune-tabnine home as well.
	//
	// NOTE: i'm not sure if Tabnine is best run per project, or for
	// everywhere. Need to contact the author, but currently it seems
	// like Sublime uses a single tabnine process for everything,
	// so that's what i'm basing this off of.
	kakConfigDir := os.Getenv("kak_config")
	if kakConfigDir == "" {
		return fmt.Errorf("missing kak_config value")
	}

	configDir := filepath.Join(kakConfigDir, "tabnine")

	// TODO(leeola): create a downloader that ensures we have the
	// latest TabNine binary for the current system. For now i'm
	// hard coding the bin path.
	tabnineBin := filepath.Join(configDir, "binaries", "0.11.2", "x86_64-apple-darwin", "TabNine")

	c := tabnine.HTTPServerConfig{
		TabnineBin: tabnineBin,
		ConfigDir:  configDir,
		LogTabnine: true,
	}
	s, err := tabnine.NewHTTPServer(c)
	if err != nil {
		return fmt.Errorf("NewHTTPServer: %v", err)
	}

	if err := s.ListenAndServe("localhost:"); err != nil {
		return fmt.Errorf("ListenAndServe: %v", err)
	}

	return nil
}
