package preinstall_test

import (
	efi "github.com/canonical/go-efilib"
	secboot_efi "github.com/snapcore/secboot/efi"
	"github.com/snapcore/secboot/internal/testutil"
)

var (
	// Default digests (SHA256)
	shimDigestDefault   []byte = testutil.MustDecodeHexString("25e1b08db2f31ff5f5d2ea53e1a1e8fda6e1d81af4f26a7908071f1dec8611b7")
	grubDigestDefault   []byte = testutil.MustDecodeHexString("d5a9780e9f6a43c2e53fe9fda547be77f7783f31aea8013783242b040ff21dc0")
	kernelDigestDefault []byte = testutil.MustDecodeHexString("2ddfbd91fa1698b0d133c38ba90dbba76c9e08371ff83d03b5fb4c2e56d7e81f")
)

// efiImagesDefault returns the default mock images used for testing.
// It cannot be defined as a static var as some parts get defined in init() sections
// (therefore after static vars get initialized).
func efiImagesDefault() []secboot_efi.Image {
	return []secboot_efi.Image{
		&mockImage{
			contents: []byte("mock shim executable"),
			digest:   shimDigestDefault,
			signatures: []*efi.WinCertificateAuthenticode{
				shimUbuntuSig4WinCert,
			},
		},
		&mockImage{contents: []byte("mock grub executable"), digest: grubDigestDefault},
		&mockImage{contents: []byte("mock kernel executable"), digest: kernelDigestDefault},
	}
}
