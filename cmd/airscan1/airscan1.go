// Copyright 2020 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Program airscan1 is a utility program that scans documents from one
// AirScan-compatible scanner at a time. This program was written for
// illustration of the airscan package, but might come in handy.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/brutella/dnssd"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/renameio"
	"github.com/stapelberg/airscan"
	"github.com/stapelberg/airscan/preset"
)

func humanDeviceName(srv dnssd.Service) string {
	if ty := srv.Text["ty"]; ty != "" {
		return ty
	}

	// miekg/dns escapes characters in DNS labels: as per RFC1034 and
	// RFC1035, labels do not actually permit whitespace. The purpose of
	// escaping originally appears to be to use these labels in a DNS
	// master file, but for our UI, backslashes look just wrong:
	return strings.ReplaceAll(srv.Name, "\\", "")
}

func airscan1() error {
	var sc airscanner

	flag.BoolVar(
		&sc.debug,
		"debug",
		false,
		"if true, print extra debug output")

	flag.StringVar(
		&sc.host,
		"host",
		"",
		"if specified, locate the scanner to use based on its Hostname")

	flag.StringVar(
		&sc.scanDir,
		"scan_dir",
		"/tmp",
		"Directory in which to store the scanned page(s). Will be created if it does not exist")

	flag.StringVar(
		&sc.source,
		"source",
		"platen",
		"Source of the document. One of platen (flat bed) or adf (Automatic Document Feeder)")

	flag.StringVar(
		&sc.size,
		"size",
		"A4",
		"Page size. One of A4 or letter")

	flag.StringVar(
		&sc.format,
		"format",
		"image/jpeg",
		"File format to request from the scanner")

	flag.StringVar(
		&sc.color,
		"color",
		"Grayscale8",
		"Color mode to request from the scanner (Grayscale8, RGB24)")

	flag.BoolVar(
		&sc.duplex,
		"duplex",
		true,
		"if false, scan only the front side of the page")

	var (
		discover = flag.Duration("discover",
			0,
			"if non-zero, discover (and list) airscan compatible devices for the specified duration, then exit")

		timeout = flag.Duration("timeout",
			5*time.Second,
			"if non-zero, limit time for finding the device")

		conntest = flag.Bool("conntest",
			false,
			"if true, report which ways of connecting to the discovered device work")
	)
	flag.Parse()

	ctx, canc := context.WithCancel(context.Background())
	defer canc()
	if *discover != 0 {
		log.Printf("finding airscan-compatible devices for %v", *discover)
		ctx, canc = context.WithTimeout(ctx, *discover)
		defer canc()
	} else if *timeout != 0 {
		log.Printf("finding device for %v (use -timeout=0 for unlimited)", *timeout)
		ctx, canc = context.WithTimeout(ctx, *timeout)
		defer canc()
	}

	discoverStart := time.Now()

	addFn := func(service dnssd.Service) {
		if sc.debug {
			log.Printf("DNSSD service discovered: %v", spew.Sdump(service))
		}

		humanName := humanDeviceName(service)
		if sc.host != "" && sc.host == service.Host {
			canc()
			log.Printf("device %q found in %v", humanName, time.Since(discoverStart))
			sc.service = &service
			return
		}

		log.Printf("device %q discovered (use -host=%q)", humanName, service.Host)
	}

	rmvFn := func(srv dnssd.Service) {
		log.Printf("device %q vanished", humanDeviceName(srv))
	}

	// addFn and rmvFn are always called (sequentially) from the same goroutine,
	// i.e. no locking is required.
	if err := dnssd.LookupType(ctx, airscan.ServiceName, addFn, rmvFn); err != nil &&
		err != context.Canceled &&
		err != context.DeadlineExceeded {
		return err
	}

	if *discover != 0 {
		return nil // only discovery requested, exit instead of scanning
	}

	if sc.service == nil {
		return fmt.Errorf("no compatible scanner found")
	}

	if *conntest {
		log.Println("testing reachability of all addresses:")
		ctx = context.Background()
		if *timeout != 0 {
			ctx, canc = context.WithTimeout(ctx, *timeout)
			defer canc()
		}
		testConns(ctx, sc.service)
		return nil
	}

	start := time.Now()

	if err := os.MkdirAll(sc.scanDir, 0700); err != nil {
		return err
	}

	if err := sc.scan1(); err != nil {
		return err
	}

	log.Printf("scan done in %v", time.Since(start))

	return nil
}

type airscanner struct {
	debug   bool
	host    string
	scanDir string
	source  string
	size    string
	format  string
	color   string
	duplex  bool
	service *dnssd.Service
}

func (sc *airscanner) scan1() error {
	cl := airscan.NewClientForService(sc.service)

	settings := preset.GrayscaleA4ADF()
	switch sc.source {
	case "platen":
		settings.InputSource = "Platen"
	case "adf":
	default:
		return fmt.Errorf("unexpected source: got %q, want one of platen or adf", sc.source)
	}
	switch sc.size {
	case "A4":
	case "letter":
		settings.ScanRegions.Regions[0].Width = 2550
		settings.ScanRegions.Regions[0].Height = 3300
	default:
		return fmt.Errorf("unexpected page size: got %q, want one of A4 or letter", sc.size)
	}
	suffix := "jpg"
	switch sc.format {
	case "image/jpeg":
	case "application/pdf":
		suffix = "pdf"
		settings.DocumentFormat = "application/pdf"
	}
	switch sc.color {
	case "Grayscale8":
	case "RGB24":
		settings.ColorMode = "RGB24"
	}
	settings.Duplex = sc.duplex

	scan, err := cl.Scan(settings)
	if err != nil {
		return err
	}
	defer scan.Close()

	pagenum := 1
	for scan.ScanPage() {
		if sc.debug {
			log.Printf("receiving page %d", pagenum)
		}
		var fn string
		for {
			fn = filepath.Join(sc.scanDir, fmt.Sprintf("page%d.%s", pagenum, suffix))
			_, err := os.Stat(fn)
			if err == nil /* file exists */ {
				pagenum++
				continue
			}
			if os.IsNotExist(err) {
				break
			}
		}

		o, err := renameio.TempFile("", fn)
		if err != nil {
			return err
		}
		defer o.Cleanup()

		if _, err := io.Copy(o, scan.CurrentPage()); err != nil {
			return err
		}

		size, err := o.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}

		if err := o.CloseAtomicallyReplace(); err != nil {
			return err
		}

		log.Printf("wrote %s (%d bytes)", fn, size)

		pagenum++
	}
	if err := scan.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := airscan1(); err != nil {
		log.Fatal(err)
	}
}
