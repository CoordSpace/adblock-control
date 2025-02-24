# Pi-hole Adblock Control

A simple golang webapp for temporarily disabling a self-hosted instance of the [Pi-hole V6.0+](https://pi-hole.net/) adblocker.

Useful for when a family member needs to quickly make an order from a retailer that pins their site functionality on 3rd-party services that threaten your ambient privacy online and/or also spread malware. (I'm looking at you Lowes, Macy's, JCPenney, & CVS)

Intended for use in combination with Docker and a reverse proxy such as [Traefik](https://docs.traefik.io/) or [nginx-proxy](https://github.com/nginx-proxy/nginx-proxy).

![UI Interaction](https://thumbs.gfycat.com/AridEasyBoa-small.gif)

## Installation

```bash
$ go get github.com/CoordSpace/adblock-control
```

## Usage

> [!IMPORTANT]
> Due to [authentication changes](https://docs.pi-hole.net/api/auth/) as part of the Pi-hole V6 update, you will need to generate an app password in the web interface on the settings page.

You can specify the apikey, url, and port with flags.
```bash
$ adblock-control -app_pass=XXYYZZ -url='https://pi.hole' -port=9000
```

Or with environment variables (Useful for configuring a Docker container)
```bash
$ export APP_PASS="XXYYZZ"
$ export URL="https://pi.hole"
$ export PORT="9000"
$ adblock-control
```

Once running, go to localhost:port and check it out!

### Docker

A prebuild image is already available and can be run via

```bash
$ docker run --rm -p 8080:8080 -e APP_PASS=xxxxxyyyyyzzzzz -e URL='https://pi.hole' coordspace/adblock-control:latest
```

## Notes 
* Flags supersede env variables at runtime.
* Port is an optional setting, if not specified via flag or env it will default to 8080.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## Legal
This project is in no way endorsed, sponsored by, or associated with the Pi-hole project and/or Pi-hole LLC

## License
MIT License

Copyright (c) 2020 Christopher Earley

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
