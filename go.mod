module streamdeckd-modules

go 1.25.1

require (
	github.com/Endg4meZer0/go-mpris v1.0.5
	github.com/godbus/dbus/v5 v5.0.4-0.20200513180336-df5ef3eb7cca
	github.com/the-jonsey/pulseaudio v0.0.2-0.20260222211608-58a869b098fe
	github.com/unix-streamdeck/api/v2 v2.0.1-0.20250915204217-05040f967038
	golang.org/x/sync v0.1.0
)

require (
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	golang.org/x/image v0.0.0-20201208152932-35266b937fa6 // indirect
)

//
//replace github.com/the-jonsey/pulseaudio v0.0.1 => ../pulseaudio
//

replace github.com/unix-streamdeck/api/v2 => ../api

replace github.com/Endg4meZer0/go-mpris => ../go-mpris
