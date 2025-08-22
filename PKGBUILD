# Maintainer: Dmitry Modder
pkgname=lunafan-control
pkgver=1.0
pkgrel=1
pkgdesc="Fan control daemon"
arch=('x86_64')
url="local"
license=('MIT')
depends=('glibc')
makedepends=('go')
source=('main.go' 'config.json' 'lunafan-control.service')
sha256sums=('bf20687fc3cbe8a461432d10115bbbeb91a39832e7c8135983e70438eb394fbc' '04e12a20a7dddd2e2305520677249624044cecacee8c935dfa35ee6fb106f62b' '2b6dd45c3164e5a7e86665332a2306536c1ab27d0fd3708bf4a3778bd59aa861')

build() {
  go build -o "$pkgname" main.go
}

package() {
  install -Dm755 "$pkgname" "$pkgdir/usr/bin/$pkgname"
  install -Dm644 config.json "$pkgdir/etc/lunafan-control/config.json"
  install -dm755 "$pkgdir/etc/lunafan-control/configs"
  cp config.json "$pkgdir/etc/lunafan-control/configs/default.json"
  install -Dm644 lunafan-control.service "$pkgdir/usr/lib/systemd/system/lunafan-control.service"
}
