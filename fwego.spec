%global import_path     code.google.com/p/go.net
%global commit          b691927ccb3a11dd198562aaae70f96fecf5e16b
%global shortcommit     %(r=%{commit}; echo ${r:0:7})
%global debug_package   %{nil}
%global __strip         /bin/true

Name:                   fwego
Version:                0.1
Release:                6git%{shortcommit}%{?dist}
Summary:                Simple web file browser

License:                GPLv3
URL:                    http://github.com/elemc/fwego
Source0:                https://github.com/elemc/%{name}/archive/%{shortcommit}/%{name}-%{shortcommit}.tar.gz

BuildRequires:          golang
BuildRequires:          systemd
BuildRequires:          git

Requires(post):         systemd
Requires(preun):        systemd
Requires(postun):       systemd

ExclusiveArch:          %{ix86} x86_64 %{arm}

%package httpd
BuildArch:              noarch
Requires:               httpd
Summary:                fwego httpd configuration file for httpd

%package nginx
BuildArch:              noarch
Requires:               nginx
Summary:                fwego nginx configuration file for nginx

%description
This is simple executable web service for browsable file system catalog as web page.

%description httpd
This package contain configuration file for httpd

%description nginx
This package contain configuration file for nginx

%prep
%setup -qn %{name}-%{commit}

%build
make %{?_smp_mflags}

%install
rm -rf $RPM_BUILD_ROOT
%make_install SYSCONFDIR=%{buildroot}%{_sysconfdir} SYSTEMD_UNIT_DIR=%{buildroot}%{_unitdir} BINDIR=%{buildroot}%{_bindir}
%{__make} install-httpd-conf SYSCONFDIR=%{buildroot}%{_sysconfdir} SYSTEMD_UNIT_DIR=%{buildroot}%{_unitdir} BINDIR=%{buildroot}%{_bindir} DESTDIR=%{?buildroot}
%{__make} install-nginx-conf SYSCONFDIR=%{buildroot}%{_sysconfdir} SYSTEMD_UNIT_DIR=%{buildroot}%{_unitdir} BINDIR=%{buildroot}%{_bindir} DESTDIR=%{?buildroot}

%files
%doc README.md
%{_bindir}/%{name}
%{_unitdir}/%{name}.service
%config(noreplace) %{_sysconfdir}/sysconfig/%{name}

%files httpd
%config(noreplace) %{_sysconfdir}/httpd/conf.d/%{name}.conf

%files nginx
%config(noreplace) %{_sysconfdir}/nginx/conf.d/%{name}.conf

%post
%systemd_post fwego.service

%preun
%systemd_preun fwego.service

%postun
%systemd_postun_with_restart fwego.service

%changelog
* Mon Oct 10 2016 Alexei Panov <me AT elemc DOT name> 0.1-6
- Disable debug logging

* Mon Oct 10 2016 Alexei Panov <me AT elemc DOT name> 0.1-5b691927
- Refactoring for golang rules

* Sun Oct 09 2016 Alexei Panov <me AT elemc DOT name> 0.1-4
- Fix panic if write failed

* Sun May 04 2014 Alexei Panov <me AT elemc DOT name> 0.1-3git23a7b63
- Added module file for nginx conf.d

* Thu Apr 24 2014 Alexei Panov <me AT elemc DOT name> 0.1-2git
- Split to two packages and use Makefile

* Wed Apr 23 2014 Alexei Panov <me AT elemc DOT name> 0-0.1git
- Initial build
