module streamdeckd-modules

go 1.25.1

require (
	github.com/Endg4meZer0/go-mpris v1.0.5
	github.com/godbus/dbus/v5 v5.0.4-0.20200513180336-df5ef3eb7cca
	github.com/the-jonsey/pulseaudio v0.0.1
	github.com/unix-streamdeck/api/v2 v2.0.1-0.20250915204217-05040f967038
	github.com/unix-streamdeck/streamdeckd v1.0.3-0.20260218203403-61f7659db6c3
	golang.org/x/sync v0.1.0
)

require (
	github.com/bearsh/hid v1.4.2-0.20220627123055-35af594cb5a7 // indirect
	github.com/bendahl/uinput v1.7.0 // indirect
	github.com/christopher-dG/go-obs-websocket v0.0.0-20200720193653-c4fed10356a5 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/mitchellh/mapstructure v1.1.2 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/unix-streamdeck/driver v0.0.0-20260218200545-edb2f975e560 // indirect
	golang.org/x/image v0.0.0-20201208152932-35266b937fa6 // indirect
)

replace github.com/the-jonsey/pulseaudio v0.0.1 => ../pulseaudio

replace github.com/unix-streamdeck/streamdeckd => ../streamdeckd

//replace github.com/Endg4meZer0/go-mpris => ../go-mpris
