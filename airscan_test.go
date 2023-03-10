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

package airscan_test

import (
	crypto_rand "crypto/rand"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"sync"
	"testing"

	"github.com/brutella/dnssd"
	"github.com/google/go-cmp/cmp"
	"github.com/stapelberg/airscan"
	"github.com/stapelberg/airscan/preset"
)

var binaryScanDataStandIn = []byte{0x22, 0x33, 0x44}

func getEsclMockFile(t *testing.T, name string) io.ReadCloser {
	filePath := path.Join("resources/eSCL", name)
	f, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("unable to open file %s", filePath)
	}

	return f
}

func mockScanner(t *testing.T) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/eSCL/ScannerStatus", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, getEsclMockFile(t, "ScannerStatus.xml"))
	})

	mux.HandleFunc("/eSCL/ScannerCapabilities", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.Copy(w, getEsclMockFile(t, "ScannerCapabilities.xml"))
	})

	var (
		scansMu sync.Mutex
		scans   = make(map[string]int)
	)

	mux.HandleFunc("/eSCL/ScanJobs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "bad method", http.StatusBadRequest)
			return
		}
		scansMu.Lock()
		defer scansMu.Unlock()
		random := make([]byte, 16)
		crypto_rand.Read(random)
		key := fmt.Sprintf("%x", random)
		scans[key] = 2
		// We intentionally return a never-working URL (port 9 is the discard
		// protocol) so that we verify the package doesn’t accidentally take the
		// scanner-provided host (if any): these might be buggy, so let’s stick
		// to the transfer method that worked for the request itself.
		http.Redirect(w, r, "http://localhost:9/eSCL/ScanJobs/"+key, http.StatusCreated)
	})

	shouldFail := true // return a 503 on `NextDocument`
	mux.HandleFunc("/eSCL/ScanJobs/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/eSCL/ScanJobs/")
		if id == "" {
			http.Error(w, "no such job", http.StatusNotFound)
			return
		}
		var rest string
		if idx := strings.IndexByte(id, '/'); idx > -1 {
			rest = id[idx+1:]
			id = id[:idx]
		}
		if rest != "NextDocument" {
			http.Error(w, "no such handler", http.StatusNotFound)
			return
		}
		scansMu.Lock()
		pages, ok := scans[id]
		if ok && pages > 0 {
			if !shouldFail {
				scans[id] = pages - 1
				if scans[id] == 0 {
					delete(scans, id)
				}
			}
		}
		scansMu.Unlock()
		if shouldFail {
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			shouldFail = false // failed once, work on the next call to `NextDocument`
			return
		}
		if !ok {
			http.Error(w, "no such job", http.StatusNotFound)
			return
		}
		w.Write(binaryScanDataStandIn)
	})

	return mux
}

func clientForMockScanner(t *testing.T) *airscan.Client {
	t.Helper()
	srv := httptest.NewServer(mockScanner(t))
	t.Cleanup(func() { srv.Close() })
	// round-trip the listener address through net.SplitHostPort and
	// net.JoinHostPort to verify that it is indeed a host:port address:
	host, port, err := net.SplitHostPort(srv.Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	cl := airscan.NewClient(net.JoinHostPort(host, port))
	cl.HTTPClient = srv.Client()
	return cl
}

func TestScannerStatus(t *testing.T) {
	cl := clientForMockScanner(t)
	resp, err := cl.ScannerStatus()
	if err != nil {
		t.Fatal(err)
	}
	want := &airscan.ScannerStatus{
		Version:  "2.63",
		State:    "Idle",
		ADFState: "ScannerAdfEmpty",
	}
	if diff := cmp.Diff(want, resp); diff != "" {
		t.Fatalf("unexpected ScannerStatus: diff (-want +got):\n%s", diff)
	}
}

func TestScan(t *testing.T) {
	cl := clientForMockScanner(t)
	grayscaleA4Platen := preset.GrayscaleA4ADF()
	grayscaleA4Platen.InputSource = "Platen"
	job, err := cl.Scan(grayscaleA4Platen)
	if err != nil {
		t.Fatal(err)
	}
	var pages [][]byte
	for job.ScanPage() {
		b, err := io.ReadAll(job.CurrentPage())
		if err != nil {
			t.Fatal(err)
		}
		pages = append(pages, b)
	}
	if err := job.Err(); err != nil {
		t.Fatal(err)
	}
	if got, want := len(pages), 2; got != want {
		t.Fatalf("unexpected number of pages read: got %d, want %d", got, want)
	}
	if diff := cmp.Diff(binaryScanDataStandIn, pages[0]); diff != "" {
		t.Fatalf("unexpected scan data: diff (-want +got):\n%s", diff)
	}
	if err := job.Close(); err != nil {
		t.Fatalf("unexpected cleanup error: %v", err)
	}
}

var discoveredService *dnssd.BrowseEntry // descriptive name for ExampleClient_Scan

func ExampleClient_Scan() {
	// For a full example using DNSSD service discovery, see:
	// https://github.com/stapelberg/airscan/blob/master/cmd/airscan1/airscan1.go
	cl := airscan.NewClientForService(discoveredService)

	// Set up scan job:
	grayscaleA4Platen := preset.GrayscaleA4ADF()
	grayscaleA4Platen.InputSource = "Platen"
	job, err := cl.Scan(grayscaleA4Platen)
	if err != nil {
		panic(err)
	}
	defer job.Close()

	// Scan one individual page at a time:
	for job.ScanPage() {
		// Read and discard scan data. This is where you would typically save
		// the data to a file, send it via the net, display or process it, etc.:
		if _, err := io.Copy(io.Discard, job.CurrentPage()); err != nil {
			panic(err)
		}
	}
	if err := job.Err(); err != nil {
		panic(err)
	}

	// Scan succeeded!
}
