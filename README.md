fwego
=====

This is simple web service for browse a files written on Go.

Build with go
-------------

### Requires
- golang
- configure you GOPATH environment variable

### Download

$ *go get github.com/elemc/fwego*

### Build

$ *go build github.com/elemc/fwego*
$ *go install github.com/elemc/fwego*

### Run standalone
$ *$GOPATH/bin/fwego -root-path=/path/for/share -address="0.0.0.0" -port=80*
See *$GOPATH/bin/fwego --help* for more information

Build with GNU make
-------------------

### Download
$ *git clone git://github.com/elemc/fwego.git*

### Build
$ *cd fwego*
$ *make*
$ *sudo make install*

if you need httpd configuration file also run this:
$ *sudo make install-httpd-conf*

### Clean
$ *make clean*

### Uninstall
$ *sudo make uninstall*

if you have installed httpd config file run this:
$ *sudo make uninstall-httpd-conf

Fedora package
--------------
### Contents
- sysconfig script /etc/sysconfig/fwego
- systemd unit /usr/lib/systemd/system/fwego.service
- apache configuration file /etc/httpd/conf.d/fwego.conf

### Steps for deploy
* Change sysconfig script for set root path for share.
* Change ServerName in apache configuration file
* run fwego: *systemctl start fwego*
* start or restart httpd: *systemctl restart httpd*

Also you may start systemd unit as standalone service. Change sysconfig variable 'FWEGO_LISTEN' for it.

Bugs
----
Please send me bugs and feature request here https://github.com/elemc/fwego/issues

TODO
----
Nothing yet.
