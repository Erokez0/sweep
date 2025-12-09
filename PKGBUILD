# Maintainer: Parshin Kirill <parshin.k.07@gmail.com>
pkgname=sweep
pkgver=0.1.0
pkgrel=1
# epoch=
pkgdesc="a flexible minesweeping experience in the terminal"
arch=("x86_64" "aarch")
url="https://github.com/erokez0/$pkgname"
license=('MIT')
groups=()
depends=()
makedepends=("git" "go")
checkdepends=()
optdepends=()
provides=()
# conflicts=()
# replaces=()
# backup=()
# options=()
# install=
# changelog=
# source=("$pkgname::git://github.com/erokez0/sweep")
source=("$pkgname::git://github.com/erokez0/$pkgname")
# noextract=()
sha256sums=("SKIP")
# validpgpkeys=()

prepare() {
	go test 
	cd "$pkgname"
	patch -p1 -i "$srcdir/$pkgname-$pkgver.patch"
}

build() {
	cd "$pkgname"
	go build -o ./bin/$pkgname
}

check() {
	cd "$pkgname"
	go test ./game-engine
}

package() {
	cd "$pkgname-$pkgver"
	install -Dm755 ./$pkgname "$pkgdir/usr/bin/$pkgname" 
}
