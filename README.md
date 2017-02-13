driebit/nfc-scanner
===================

A Linux binary for scanning NFC/RFID tags and submitting their UID to the  
[Ginger Tagger API](https://github.com/driebit/ginger/blob/master/modules/mod_ginger_tagger/README.md)
over HTTP.

Usage
-----

nfc-scanner is configured using environment variables (following the 
[Twelve-factor principles](https://12factor.net)):

| Variable        | Explanation                                   | Example             |
| --------------- | --------------------------------------------- | ------------------- |
| `API_URL`       | URL to the Tagger API                         | "https://ginger.nl" |
| `CLIENT_ID`     | your Tagger API client id                     | "you"               |
| `CLIENT_SECRET` | your Tagger API client secret                 | "007secret"         |
| `PANEL_ID`      | the scanned object, e.g. panel or activity id | 516                 |

Download the binary, then run it:

```bash
$ wget 
$ nfc-scanner 
```

