package fastprime

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRandPrime(t *testing.T) {
	rnd:=NewRand()
	for i:=0; i<10; i++ {
		for j:=uint16(64); j<96; j++{
			p:=rnd.RandPrime(j,64)
			assert.Equal(t,int(j),p.BitLen())
			assert.True(t,p.ProbablyPrime(64))
		}
	}
}
