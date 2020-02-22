package tools

import (
	"crypto/sha256"
	"github.com/gobwas/pool/pbytes"
	"math/big"
)

const (
	bitSize     = 2048
	SrpByteSize = bitSize / 8
	pHex        = "0xAC6BDB41324A9A9BF166DE5E1389582FAF72B6651987EE07FC3192943DB56050A37329CBB4A099ED8193E0757767A13DD52312AB4B03310DCD7F48A9DA04FD50E8083969EDB767B0CF6095179A163AB3661A05FBD5FAAAE82918A9962F0B93B855F97993EC975EEAA80D740ADBF4FF747359D041D5C33EA71D281E446B14773BCA97B43A23FB801676BD207A436C6481F1D2B9078717461A5B9D32E688F87748544523B524B0D57D5EA77A2775D2ECFA032CFBDBF52FB3786160279004E57AE6AF874E7303CE53299CCC041C7BC308D82A5698F3A8D0C38271AE35F8E9DBFBB694B5C803D89F7AE435DE236D525F54759B65E372FCD68EF20FA7111F9E4AFF73"
)

var (
	DefaultPrime     *big.Int
	DefaultGenerator *big.Int
)

func init() {
	DefaultGenerator = big.NewInt(2)
	DefaultPrime, _ = big.NewInt(0).SetString(pHex, 0)
}

func H(data ...[]byte) []byte {
	h := sha256.New()
	for _, d := range data {
		h.Write(d)
	}
	return h.Sum(nil)
}

func SH(data, salt []byte) []byte {
	b := pbytes.GetCap(len(data) + 2*len(salt))
	b = append(b, salt...)
	b = append(b, data...)
	b = append(b, salt...)
	return H(b)
}

func PH1(password, salt1, salt2 []byte) []byte {
	return SH(SH(password, salt1), salt2)
}

func PH2(password, salt1, salt2 []byte) []byte {
	return SH(SH(PH1(password, salt1, salt2), salt1), salt2)
}

func K(p, g *big.Int) []byte {
	return H(append(Pad(p), Pad(g)...))
}

func U(ga, gb *big.Int) []byte {
	return H(append(Pad(ga), Pad(gb)...))
}

func M(p, g *big.Int, s1, s2 []byte, ga, gb, sb *big.Int) []byte {
	return H(H(Pad(p)), H(Pad(g)), H(s1), H(s2), H(Pad(ga)), H(Pad(gb)), H(Pad(sb)))
}

// Pad x to n bytes if needed
func Pad(x *big.Int) []byte {
	b := x.Bytes()
	if len(b) < SrpByteSize {
		z := SrpByteSize - len(b)
		p := make([]byte, SrpByteSize)
		for i := 0; i < z; i++ {
			p[i] = 0
		}

		copy(p[z:], b)
		b = p
	}
	return b
}
