![GitHub Actions CI](https://github.com/stapelberg/airscan/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/stapelberg/airscan)](https://goreportcard.com/report/github.com/stapelberg/airscan)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://pkg.go.dev/github.com/stapelberg/airscan)

# airscan 📄 🖨️ 🕸️

The `airscan` Go package can be used to scan paper documents 📄 from a scanner
🖨️ via the network 🕸️ using the Apple AirScan (eSCL) protocol.

## Getting started: example program

First, install the example program coming with this package:

```
go install github.com/stapelberg/airscan/cmd/airscan1
```

Then, query the local network for AirScan compatible devices:

```
% airscan1 -discover=5s
2020/08/16 08:50:31 finding airscan-compatible devices for 1s
2020/08/16 08:50:31 device "Brother MFC-L2750DW series" discovered (use -host="BRW405BD8AxxDyz")
```

Now, I can scan the contents of the flatbed scanner:
```
% airscan1 -host=BRW405BD8AxxDyz
2020/08/16 08:52:44 finding device for 5s (use -timeout=0 for unlimited)
2020/08/16 08:52:45 device "Brother MFC-L2750DW series" found in 298.151935ms
2020/08/16 08:52:51 scan done in 6.738205326s
```

…or the page(s) inserted into the Automatic Document Feeder (ADF):
```
% airscan1 -host=BRW405BD8A10D7C -source=adf
2020/08/16 11:10:34 finding device for 5s (use -timeout=0 for unlimited)
2020/08/16 11:10:34 device "Brother MFC-L2750DW series" found in 112.127399ms
2020/08/16 11:10:45 wrote /tmp/page12.jpg (211305 bytes)
2020/08/16 11:10:47 wrote /tmp/page13.jpg (139335 bytes)
2020/08/16 11:10:47 scan done in 13.068799513s
```

## Getting started: using the package in your program

See the [package airscan examples in
godoc](https://pkg.go.dev/github.com/stapelberg/airscan?tab=doc#pkg-examples)
for how to use the package to scan.

See
[airscan1.go](https://github.com/stapelberg/airscan/blob/master/cmd/airscan1/airscan1.go#L100)
for a full example scan program, including network service discovery, timeouts,
and writing scan data to files.

## Project status

The package does what I needed: grayscale/color scan of A4 documents from the
flat bed or the automatic document feeder (ADF).

If you have any improvements, I’d be happy to review a pull request. Please see the [contribution guidelines](/docs/contributing.md).

## Tested devices

If you successfully scanned documents from your device using the `airscan1`
example program as described above, please [send a pull
request](https://github.com/stapelberg/airscan/edit/master/README.md) to include
your report in this table for the benefit of other interested users:

| Device Name | Working features | Known issues |
| ----------- | ---------------- | ------------ |
| Brother MFC-L2750DW | flat bed scan, automatic document feeder scan | |
