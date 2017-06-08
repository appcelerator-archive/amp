# smtpd 

An SMTP server package written in Go, in the style of the built-in HTTP server. It meets the minimum requirements specified by RFC 2821 & 5321. 

It is based on [Brad Fitzpatrick's go-smtpd](https://github.com/bradfitz/go-smtpd). The differences can be summarised as:

* A simplified message handler
* Changes made for RFC compliance
* Testing has been added
* Code refactoring
* TLS support

## Features

* A single message handler for simple mail handling with native data types.
* RFC compliance. It implements the minimum command set, responds to commands and adds a valid Received header to messages as specified in RFC 2821 & 5321.
* Customisable listening address and port. It defaults to listening on all addresses on port 25 if unset.
* Customisable host name and application name. It defaults to the system hostname and "smtpd" application name if they are unset.
* Easy to use TLS support that obeys RFC 3207.

## Usage

In general: create the server and pass a handler function to it as for the HTTP server. The server function has the following definition:

```
func ListenAndServe(addr string, handler Handler, appname string, hostname string) error
```

For TLS support, add the paths to the certificate and key files as for the HTTP server.

```
func ListenAndServeTLS(addr string, certFile string, keyFile string, handler Handler, appname string, hostname string) error
```

The handler function must have the following definition:

```
func handler(remoteAddr net.Addr, from string, to []string, data []byte) 
```

The parameters are:

* remoteAddr: remote end of the TCP connection i.e. the mail client's IP address and port.
* from: the email address sent by the client in the MAIL command.
* to: the set of email addresses sent by the client in the RCPT command.
* data: the raw bytes of the mail message.

## TLS Support

SMTP over TLS works slightly differently to how you might expect if you are used to the HTTP protocol. Some helpful links for background information are:

* [SSL vs TLS vs STARTTLS](https://www.fastmail.com/help/technical/ssltlsstarttls.html)
* [Opportunistic TLS](https://en.wikipedia.org/wiki/Opportunistic_TLS)
* [RFC 2487: SMTP Service Extension for Secure SMTP over TLS](https://tools.ietf.org/html/rfc2487)
* [func (*Client) StartTLS](https://golang.org/pkg/net/smtp/#Client.StartTLS)

The TLS support has three server configuration options. The bare minimum requirement to enable TLS is to supply certificate and key files as in the TLS example below.

* TLSConfig

This option allows custom TLS configurations such as [requiring strong ciphers](https://cipherli.st/) or using other certificate creation methods. If a certificate file and a key file are supplied to the ConfigureTLS function, the default TLS configuration for Go will be used. The default value for TLSConfig is nil, which disables TLS support.

* TLSRequired

This option sets whether TLS is optional or required. If set to true, the only allowed commands are NOOP, EHLO, STARTTLS and QUIT (as specified in RFC 3207) until the connection is upgraded to TLS i.e. until STARTTLS is issued. This option is ignored if TLS is not configured i.e. if TLSConfig is nil. The default is false.

* TLSListener

This option sets whether the listening socket requires an immediate TLS handshake after connecting. It is equivalent to using HTTPS in web servers, or the now defunct SMTPS on port 465. This option is ignored if TLS is not configured i.e. if TLSConfig is nil. The default is false.

There is also a related package configuration option.

* Debug

This option determines if the data being read from or written to the client will be logged. This may help with debugging when using encrypted connections. The default is false.

## Example

The following example code creates a new server with the name "MyServerApp" that listens on the localhost address and port 2525. Upon receipt of a new mail message, the handler function parses the mail and prints the subject header.


```
package main

import (
	"bytes"
	"log"
	"net"
	"net/mail"

	"github.com/mhale/smtpd"
)

func mailHandler(origin net.Addr, from string, to []string, data []byte) {
	msg, _ := mail.ReadMessage(bytes.NewReader(data))
	subject := msg.Header.Get("Subject")
	log.Printf("Received mail from %s for %s with subject %s", from, to[0], subject)
}

func main() {
	smtpd.ListenAndServe("127.0.0.1:2525", mailHandler, "MyServerApp", "")
}
```

## TLS Example

Using the example code above, only the main function would be different to add TLS support.

```
func main() {
	smtpd.ListenAndServeTLS("127.0.0.1:2525", "/path/to/server.crt", "/path/to/server.key", mailHandler, "MyServerApp", "")
}
```

This allows STARTTLS to be listed as a supported extension and allows clients to upgrade connections to TLS by sending a STARTTLS command.

As the package level helper functions do not set the TLSRequired or TLSListener options for compatibility reasons, manual creation of a Server struct is necessary in order to use them.

## Testing

The tests cover the supported SMTP command set and line parsing. A single server is created listening on an ephemeral port (52525) for the duration of the tests. Each test creates a new client connection for processing commands.

For the TLS tests, a different server is created with a net.Pipe connection inside each individual test, in order to change the server settings for each test.

The TLS support has also been manually tested with Go client code, Ruby client code, and macOS's Mail.app.

## Licensing

Some of the code in this package was copied or adapted from code found in [Brad Fitzpatrick's go-smtpd](https://github.com/bradfitz/go-smtpd). As such, those sections of code are subject to their original copyright and license. The remaining code is in the public domain.
