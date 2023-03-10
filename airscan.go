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

package airscan

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/brutella/dnssd"
)

// ServiceName is the AirScan DNSSD service name.
const ServiceName = "_uscan._tcp.local."

type ScannerStatus struct {
	Version  string `xml:"Version"`
	State    string `xml:"State"`
	ADFState string `xml:"AdfState"`
}

// ScanSettings instruct the device how to scan.
//
// It is recommended to use the
// https://pkg.go.dev/github.com/stapelberg/airscan/preset package to start with
// a known-good configuration.
type ScanSettings struct {
	XMLName        xml.Name    `xml:"scan:ScanSettings"`
	XmlnsScan      string      `xml:"xmlns:scan,attr"`
	XmlnsPWG       string      `xml:"xmlns:pwg,attr"`
	Version        string      `xml:"pwg:Version"`
	ScanRegions    ScanRegions `xml:"pwg:ScanRegions"`
	DocumentFormat string      `xml:"pwg:DocumentFormat"`
	InputSource    string      `xml:"pwg:InputSource"`
	ColorMode      string      `xml:"scan:ColorMode"`
	XResolution    int         `xml:"scan:XResolution"`
	YResolution    int         `xml:"scan:YResolution"`
	Duplex         bool        `xml:"scan:Duplex"`
}

func (s *ScanSettings) Marshal() (string, error) {
	b, err := xml.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", err
	}
	return `<?xml version="1.0" encoding="UTF-8" standalone="no"?>` + "\n" + string(b), nil
}

type ScanRegions struct {
	MustHonor bool `xml:"pwg:MustHonor,attr"`
	Regions   []*ScanRegion
}

type ScanRegion struct {
	XMLName            xml.Name `xml:"pwg:ScanRegion"`
	ContentRegionUnits string   `xml:"pwg:ContentRegionUnits"`
	Width              int      `xml:"pwg:Width"`
	Height             int      `xml:"pwg:Height"`
	XOffset            int      `xml:"pwg:XOffset"`
	YOffset            int      `xml:"pwg:YOffset"`
}

// A Client allows scanning documents via AirScan, which is also known as eSCL.
type Client struct {
	// HTTPClient is used for all requests made by this Client and can be
	// overridden for testing or to integrate custom behavior. The default
	// amounts to http.DefaultClient.
	HTTPClient interface {
		Do(*http.Request) (*http.Response, error)
	}

	host string
}

func isPrintable(s string) bool {
	// Verify each rune individually. See also https://blog.golang.org/strings
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

// do wraps c.Client.Do, but tries to report a descriptive error, including a
// server-sent error message, if any (and printable!).
func (c *Client) do(req *http.Request, okayStatuses ...int) (resp *http.Response, err error) {
	if debug {
		log.Printf("%s %s", req.Method, req.URL)
		defer func() {
			if err != nil {
				log.Printf("-> error: %v, %v", resp.Status, err)
			} else {
				log.Printf("-> okay: %v", resp.Status)
			}
		}()
	}
	req.Header.Set("User-Agent", "https://github.com/stapelberg/airscan")
	resp, err = c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	for _, okay := range okayStatuses {
		if resp.StatusCode == okay {
			return resp, nil
		}
	}
	want := fmt.Sprintf("one of %v", okayStatuses)
	if len(okayStatuses) == 1 {
		want = fmt.Sprint(okayStatuses[0])
	}
	b, _ := io.ReadAll(resp.Body)
	message := strings.TrimSpace(string(b))
	if !isPrintable(message) {
		message = "<non-printable body>"
	}
	return nil, fmt.Errorf("%v: unexpected HTTP status: got %v (%s), want %v",
		req.URL,
		resp.Status,
		message,
		want)
}

// ScannerStatus queries the device for its status. This can be used for example
// to find out whether a document has been inserted into the Automatic Document
// Feeder (ADF). The Scan method verifies this, too.
func (c *Client) ScannerStatus() (*ScannerStatus, error) {
	req, err := http.NewRequest("GET", c.GetEndpoint("/eSCL/ScannerStatus"), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, http.StatusOK)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var status ScannerStatus
	if err := xml.Unmarshal(b, &status); err != nil {
		return nil, fmt.Errorf("decoding XML: %v (invalid input? %q)", err, string(b))
	}
	return &status, nil
}

func (c *Client) ScannerCapabilities() (*ScannerCapabilities, error) {
	req, err := http.NewRequest("GET", c.GetEndpoint("/eSCL/ScannerCapabilities"), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, http.StatusOK)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var capabilities ScannerCapabilities
	if err := xml.Unmarshal(b, &capabilities); err != nil {
		return nil, fmt.Errorf("decoding XML: %v (invalid input? %q)", err, string(b))
	}
	return &capabilities, nil
}

func (c *Client) createScanJob(settings string) (*url.URL, error) {
	req, err := http.NewRequest("POST", c.GetEndpoint("/eSCL/ScanJobs"), strings.NewReader(settings))
	if err != nil {
		return nil, err
	}
	resp, err := c.do(req, http.StatusCreated)
	if err != nil {
		return nil, err
	}
	loc, err := resp.Location()
	if err != nil {
		return nil, err
	}
	loc.Host = req.URL.Host
	return loc, nil
}

func (c *Client) deleteScanJob(loc *url.URL) error {
	req, err := http.NewRequest("DELETE", loc.String(), nil)
	if err != nil {
		return err
	}
	if _, err := c.do(req, http.StatusNotFound); err != nil {
		return err
	}
	return nil
}

// ScanState represents an in-progress scan job.
type ScanState struct {
	loc     *url.URL
	scanner *Client
	reader  io.Reader
	err     error
}

// ScanPage requests the next page of this scan job. It returns true if a new
// page is available, or false when all pages were exhausted or an error
// occurred. Errors are available via the Err() method.
//
// Note that some scanners return 503 for NextDocument but will eventually
// return an accurate code when given more chances. See also
// https://github.com/alexpevzner/sane-airscan-ipp/blob/master/airscan-escl.c#L11.
func (s *ScanState) ScanPage() bool {
	if s.err != nil {
		return false // avoid clobbering existing errors
	}

	u, err := url.Parse(s.loc.String())
	if err != nil {
		s.err = err
		return false
	}
	u.Path = path.Join(u.Path, "NextDocument")
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		s.err = err
		return false
	}
	const tries = 10
	for try := 0; try < tries; try++ {
		resp, err := s.scanner.do(req, http.StatusOK, http.StatusNotFound, http.StatusServiceUnavailable)
		if resp != nil {
			switch resp.StatusCode {
			case http.StatusNotFound:
				if debug {
					log.Printf("NotFound: all pages received")
				}
				return false // all pages received, no error
			case http.StatusServiceUnavailable:
				if debug {
					log.Printf("ServiceUnavailable: will retry (try %d/%d)", try+1, tries)
				}
				time.Sleep(1 * time.Second)
				continue
			default:
				s.reader = resp.Body
				return true
			}
		}
		if err != nil {
			s.err = err
			return false
		}
	}
	s.err = fmt.Errorf("503 retry limit (%d) reached while calling NextDocument", tries)
	return false
}

// CurrentPage returns an io.Reader containing the scan data.
//
// CurrentPage must only be called after ScanPage() returned true, and will
// return nil otherwise.
//
// Note: package airscan never interprets scan data, the package only provides
// the data as-is. If you want to decode scan data, you will need to import
// e.g. image/jpeg (depending on scan settings) yourself.
func (s *ScanState) CurrentPage() io.Reader {
	return s.reader
}

// Err returns the first error that occurred. If it returns non-nil, the
// ScanState must no longer be used, except for the Close method.
func (s *ScanState) Err() error {
	return s.err
}

// Close deletes the scan job on the device.
//
// Some devices work just fine if you never call Close. To maximize
// compatibility, it is recommended to call Close. This mirrors what Apple’s
// scan program does, which might be required for certain scanners (speculation
// only).
func (s *ScanState) Close() error {
	if debug {
		log.Printf("Deleting ScanJob %s", s.loc)
	}
	return s.scanner.deleteScanJob(s.loc)
}

// Scan starts a new scan job using the specified settings.
//
// When scanning from an Automatic Document Feeder (ADF), the Scan method
// verifies a document is inserted before creating a scan job (which would
// otherwise fail with a less clear error message).
func (c *Client) Scan(settings *ScanSettings) (*ScanState, error) {
	// Ensure settings are valid before doing anything else:
	s, err := settings.Marshal()
	if err != nil {
		return nil, err
	}

	status, err := c.ScannerStatus()
	if err != nil {
		return nil, err
	}
	if debug {
		log.Printf("scanner status: %+v", status)
	}
	if got, want := status.State, "Idle"; got != want {
		return nil, fmt.Errorf("scanner not ready: in state %q, want %q", got, want)
	}
	if settings.InputSource == "Feeder" && status.ADFState != "" {
		if got, want := status.ADFState, "ScannerAdfLoaded"; got != want {
			return nil, fmt.Errorf("scanner feeder contains no documents: status %q, want %q", got, want)
		}
	}

	// Check capabilities
	caps, err := c.ScannerCapabilities()
	if err != nil {
		return nil, err
	}

	if settings.InputSource == "Feeder" {
		if caps.Adf == nil {
			return nil, fmt.Errorf("this scanner doesn't have an ADF")
		}

		if settings.Duplex && caps.Adf.AdfDuplexInputCaps == nil {
			return nil, fmt.Errorf("this scanner doesn't support duplex mode")
		}

		if !settings.Duplex && caps.Adf.AdfSimplexInputCaps == nil {
			// can this ever happen?
			return nil, fmt.Errorf("this scanner doesn't support simplex mode")
		}
	}

	if debug {
		log.Printf("capabilities: %+v", caps)
	}

	loc, err := c.createScanJob(s)
	if err != nil {
		return nil, err
	}
	if debug {
		log.Printf("ScanJob created: %s", loc)
	}

	return &ScanState{
		loc:     loc,
		scanner: c,
	}, nil
}

func (c *Client) GetEndpoint(s string) string {
	return fmt.Sprintf("http://%s%s", c.host, s)
}

// NewClient returns a ready-to-use Client. It is safe to update its struct
// fields before first using the returned Client.
//
// When using DNSSD service discovery to locate the scanner (the most common
// choice), prefer NewClientForService instead.
func NewClient(host string) *Client {
	return &Client{
		host: host,
		HTTPClient: &http.Client{
			Transport: http.DefaultTransport.(*http.Transport).Clone(),
		},
	}
}

// NewClientForService is like NewClient, but constructs a net.Dialer that
// attempts to connect to the specified DNSSD service using its host or IP
// address(es) directly, gracefully falling back between connection methods.
//
// This maximizes the chance of a successful connection, even when local
// networks do not offer DHCP-based DNS, and when Avahi is not available.
func NewClientForService(service *dnssd.BrowseEntry) *Client {
	port := strconv.Itoa(service.Port)
	hostports := []string{
		// DHCP-based DNS in the local network:
		net.JoinHostPort(service.Host, port),

		// Avahi name resolution:
		net.JoinHostPort(service.Host+"."+service.Domain, port),
	}
	for _, ip := range service.IPs {
		// DNSSD-provided IP address:
		hostports = append(hostports, net.JoinHostPort(ip.String(), port))
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = (&fallbackDialer{
		hostports: hostports,
		underlying: &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		},
	}).DialContext
	return &Client{
		host: service.Host,
		HTTPClient: &http.Client{
			Transport: transport,
		},
	}

}

// Flip debug to true during development
const debug = false
