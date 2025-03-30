
package fastprime

import (
	"math/big"
	"crypto/rand"
)

func millerRabin(n *big.Int, k int,stop func()bool) bool {
	if n.Cmp(big.NewInt(2)) < 0 {
		return false
	}
	if n.Cmp(big.NewInt(2)) == 0 || n.Cmp(big.NewInt(3)) == 0 {
		return true
	}
	if big.NewInt(0).Mod(n, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		return false
	}

	s := 0
	d := big.NewInt(0).Sub(n, big.NewInt(1))
	for big.NewInt(0).Mod(d, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 && d.Cmp(big.NewInt(0))>0 &&!stop(){
	
		s++
		d.Div(d, big.NewInt(2))
	}
	
	for i := 0; i < k&&!stop(); i++ {
		a,err := rand.Int(rand.Reader, d)
		if err!=nil {
			continue
		}
		a.Add(a, big.NewInt(2))
		x := big.NewInt(0).Exp(a, d, n)
		if x.Cmp(big.NewInt(1)) == 0 || x.Cmp(big.NewInt(0).Sub(n, big.NewInt(1))) == 0 {
			continue
		}
		for j := 0; j < s-1&&!stop(); j++ {
			x = big.NewInt(0).Exp(x, big.NewInt(2), n)
			if x.Cmp(big.NewInt(0).Sub(n, big.NewInt(1))) == 0 {
				break
			}
		}
		if x.Cmp(big.NewInt(0).Sub(n, big.NewInt(1))) != 0 {
			return false
		}
	}

	return true && !stop()
}

