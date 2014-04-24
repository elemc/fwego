%global import_path     code.google.com/p/go.net
%global commit          727db67eaf52f608e67acdd3bf69ef82216bdc0c
%global shortcommit     %(r=%{commit}; echo ${r:0:7})
%global debug_package   %{nil}
%global __strip         /bin/true

Name:           fwego
Version:        0.1
Release:        2git%{shortcommit}%{?dist}
Summary:        Simple web file browser

License:        GPLv3
URL:            http://github.com/elemc/fwego
Source0:        https://github.com/elemc/%{name}/archive/%{shortcommit}/%{name}-%{shortcommit}.tar.gz

BuildRequires:  golang
BuildRequires:  systemd

Requires(post): systemd
Requires(preun): systemd
Requires(postun): systemd

ExclusiveArch:  %{ix86} x86_64 %{arm}

%package -n httpd
BuildArch:      noarch
Requires:       httpd

%description
This is simple executable web service for browsable file system catalog as web page.

%description -n httpd
This package contain configuration file for httpd

%prep
%setup -qn %{name}-%{commit}

%build
make %{?_smp_mflags}

%install
rm -rf $RPM_BUILD_ROOT
%make_install SYSCONFDIR=%{_sysconfdir} SYSTEMD_UNIT_DIR=%{_unitdir}

%files
%doc README.md
%{_bindir}/%{name}
%{_unitdir}/%{name}.service
%config(noreplace) %{_sysconfdir}/sysconfig/%{name}

%files httpd
%config(noreplace) %{_sysconfdir}/httpd/conf.d/%{name}.conf

%post
%systemd_post fwego.service

%preun
%systemd_preun fwego.service

%postun
%systemd_postun_with_restart fwego.service 

%changelog
* Thu Apr 24 2014 Alexei Panov <me AT elemc DOT name> 0.1-2git
- Split to two packages and use Makefile

* Wed Apr 23 2014 Alexei Panov <me AT elemc DOT name> 0-0.1git
- Initial build 

