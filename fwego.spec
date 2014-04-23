%global import_path     code.google.com/p/go.net
%global commit          2698cf86a6b0b7763f8585240b7b5f5c095277ce
%global shortcommit     %(r=%{commit}; echo ${r:0:7})
%global debug_package   %{nil}
%global __strip         /bin/true

Name:           fwego
Version:        0.1
Release:        1git%{shortcommit}%{?dist}
Summary:        Simple web file browser

License:        GPLv3
URL:            http://github.com/elemc/fwego
Source0:        https://github.com/elemc/%{name}/archive/%{shortcommit}/%{name}-%{shortcommit}.tar.gz

BuildRequires:  golang
BuildRequires:  systemd

Requires(post): systemd
Requires(preun): systemd
Requires(postun): systemd
Requires:       httpd

ExclusiveArch:  %{ix86} x86_64 %{arm}

%description
This is simple executable web service for browsable file system catalog as web page.

%prep
%setup -qn %{name}-%{commit}


%build
go build -o %{name}


%install
rm -rf $RPM_BUILD_ROOT
install -d -m 0755 %{buildroot}%{_bindir}
install -d -m 0755 %{buildroot}%{_unitdir}
install -d -m 0755 %{buildroot}%{_sysconfdir}/httpd/conf.d
install -d -m 0755 %{buildroot}%{_sysconfdir}/sysconfig
install -m 0755 %{name} %{buildroot}%{_bindir}/%{name}
install -m 0644 %{name}.service %{buildroot}%{_unitdir}/%{name}.service
install -m 0644 %{name}.conf %{buildroot}%{_sysconfdir}/httpd/conf.d/%{name}.conf
install -m 0644 %{name}.sysconfig %{buildroot}%{_sysconfdir}/sysconfig/%{name}

%files
%doc README.md
%{_bindir}/%{name}
%{_unitdir}/%{name}.service
%config(noreplace) %{_sysconfdir}/httpd/conf.d/%{name}.conf
%config(noreplace) %{_sysconfdir}/sysconfig/%{name}

%post
%systemd_post fwego.service

%preun
%systemd_preun fwego.service

%postun
%systemd_postun_with_restart fwego.service 

%changelog
* Wed Apr 23 2014 Alexei Panov <me AT elemc DOT name> 0-0.1gitae12030
- Initial build 

