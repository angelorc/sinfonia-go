module github.com/angelorc/sinfonia-go/bitsong

go 1.18

require (
	github.com/avast/retry-go v3.0.0+incompatible // indirect
	github.com/bitsongofficial/go-bitsong v0.10.1-0.20220508161238-70f5f0e033b3 // indirect
	github.com/cosmos/cosmos-sdk v0.45.4 // indirect
	golang.org/x/sync v0.0.0-20220513210516-0976fa681c29 // indirect
)

replace (
	// use cosmos-compatible protobufs
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	// use grpc compatible with cosmos protobufs
	google.golang.org/grpc => google.golang.org/grpc v1.33.2
)
