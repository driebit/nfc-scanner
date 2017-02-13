driebit/nfc-scanner
===================

[![Build Status](https://travis-ci.org/driebit/nfc-scanner.svg?branch=master)](https://travis-ci.org/driebit/nfc-scanner)

A Linux binary for scanning NFC/RFID tags and submitting their UIDs to the 
[Ginger Tagger API](https://github.com/driebit/ginger/blob/master/modules/mod_ginger_tagger/README.md)
over HTTP.

Usage
-----

Following the [Twelve-factor principles](https://12factor.net), nfc-scanner is
configured using environment variables:

| Variable        | Explanation                                   | Example             |
| --------------- | --------------------------------------------- | ------------------- |
| `API_URL`       | URL to the Tagger API                         | "https://ginger.nl" |
| `CLIENT_ID`     | your Tagger API client id                     | "you"               |
| `CLIENT_SECRET` | your Tagger API client secret                 | "007secret"         |
| `PANEL_ID`      | the scanned object, e.g. panel or activity id | 516                 |

Download the binary (replace 0.1.0 with the [latest release version](https://github.com/driebit/nfc-scanner/releases)):

```bash
wget https://github.com/driebit/nfc-scanner/releases/download/0.1.0/nfc-scanner
chmod +x nfc-scanner
```

Set the environment variables:

```bash
$ export API_URL="http://..."
$ export CLIENT_ID="..."
$ ...
```

Then start nfc-scanner (don’t forget to plug in the NFC reader):

```bash
$ ./nfc-scanner
```

When started, nfc-scanner will listen in an infinite loop for NFC tags. When a
tag is scanned, its UID is sent to the Tagger API.
