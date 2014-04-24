PROJECT_NAME:=fwego
GOPATH:="/tmp/gobuild-$$(date +%Y-%m-%d)"
SRCPATH:="${GOPATH}/src/${PROJECT_NAME}"
BINPATH:="${GOPATH}/bin"
GO_EXEC:=/usr/bin/go
DESTDIR?=/usr/local
SYSTEMD_UNIT_DIR?=${DESTDIR}/lib/systemd/system
SYSCONFDIR?=/etc

all: fwego

fwego:
	@echo "building ${PROJECT_NAME}"
	@[ -d ${GOPATH} ] || mkdir -p ${GOPATH}
	@[ -d ${SRCPATH} ] || mkdir -p ${SRCPATH}
	@install -m 0644 ${PROJECT_NAME}.go ${SRCPATH}/${PROJECT_NAME}.go
	@GOPATH=${GOPATH} go get ${PROJECT_NAME}
	@GOPATH=${GOPATH} go install ${PROJECT_NAME}
	@install -m 0755 ${BINPATH}/${PROJECT_NAME} ./${PROJECT_NAME}

install:
	@[ -d ${DESTDIR}/bin ] || install -d -m 0755 ${DESTDIR}/bin
	@[ -d ${SYSTEMD_UNIT_DIR} ] || install -d -m 0755 ${SYSTEMD_UNIT_DIR}
	@[ -d $${SYSCONFDIR}/sysconfig ] || install -d -m 0755 ${SYSCONFDIR}/sysconfig
	install -m 0755 ${PROJECT_NAME} ${DESTDIR}/bin/${PROJECT_NAME}
	install -m 0644 ${PROJECT_NAME}.service ${SYSTEMD_UNIT_DIR}/${PROJECT_NAME}.service
	install -m 0644 ${PROJECT_NAME}.sysconfig ${SYSCONFDIR}/sysconfig/${PROJECT_NAME}

install-httpd-conf:
	install -d -m 0755 ${SYSCONFDIR}/httpd/conf.d
	install -m 0644 ${PROJECT_NAME}.conf ${SYSCONFDIR}/httpd/conf.d/${PROJECT_NAME}.conf

clean:
	@rm -rf ${GOPATH}
	@rm -rf ${PROJECT_NAME}

uninstall:
	rm -rf ${SYSCONFDIR}/sysconfig/${PROJECT_NAME}
	rm -rf ${SYSTEMD_UNIT_DIR}/${PROJECT_NAME}.service
	rm -rf ${DESTDIR}/bin/${PROJECT_NAME}

uninstall-httpd-conf:
	rm -rf ${SYSCONFDIR}/httpd/conf.d/${PROJECT_NAME.conf}
