package fastprime

import (
	"math/big"
	"crypto/rand"
	"sync"
)

var smallPrimes [4]int64 = [4]int64{2,   3,   5,   7}


type Rand struct {
	initMods []*big.Int
	skips []*big.Int
	rnge *big.Int
	bigSmallPrimes []*big.Int
}
func NewRand()*Rand{
	initMods:=make([]*big.Int,4)
	n:=big.NewInt(8)
	big1:=big.NewInt(1)
	big0:=big.NewInt(0)
	var bigSmallPrimes []*big.Int
	var rnge *big.Int
	var skips []*big.Int
	for i:=0; i<4; i++ {
		bigSmallPrimes=append(bigSmallPrimes,big.NewInt(smallPrimes[i]))
		initMods[i]=big.NewInt(0).Mod(n,bigSmallPrimes[i])
	}
	n.Add(n,big1)
	rnge=big.NewInt(1)
	skip:=big.NewInt(1)
	for {
		isPseudoPrime:=true
		isEqInitMods:=true
		for i:=0; i<4; i++ {
			m:=big.NewInt(0).Mod(n,bigSmallPrimes[i])

			if m.Cmp(big0)==0 {
				isPseudoPrime=false
			}
			if m.Cmp(initMods[i])!=0 {
				isEqInitMods=false
			}

		}
		if isEqInitMods {
			skips=append(skips,skip)
			break
		}
		if isPseudoPrime {
			skips=append(skips,skip)
			skip=big.NewInt(0)
		}
		n.Add(n,big1)
		skip.Add(skip,big1)
		rnge.Add(rnge,big1)
		
	}

	return &Rand{
		initMods:initMods,
		skips:skips,
		rnge:rnge,
		bigSmallPrimes:bigSmallPrimes,

	}
	
}

func (rnd *Rand) RandPrime(bits uint16)*big.Int{
	if bits<2 {
		panic("cant generate prime less 2 bit")
	}
	
	
	b := uint(bits % 8)
	if b == 0 {
		b = 8
	}

	bytes := make([]byte, (bits+7)/8)
	p:=big.NewInt(0)
	for {
		rand.Read(bytes)
		bytes[0] &= uint8(int(1<<b) - 1)
		if b >= 2 {
			bytes[0] |= 3 << (b - 2)
		} else {
			bytes[0] |= 1
			if len(bytes) > 1 {
				bytes[1] |= 0x80
			}
		}
		bytes[len(bytes)-1] |= 1
		p.SetBytes(bytes)
		if probablyPrime(p,64) && p.BitLen()==int(bits){
			return p
		}
		var limit bool
		q,limit:=rnd.findPrime(p,int(bits))
		if limit {
			continue
		}
		if p.BitLen()==int(bits) {
			return q
		}
	}


}


func (rnd *Rand) findPrime(p *big.Int,limit int)(*big.Int,bool){
	
	rnge2:=big.NewInt(0).Mul(rnd.rnge,big.NewInt(2))
	a:=p.Add(p,big.NewInt(1))
	big0:=big.NewInt(0)
	big1:=big.NewInt(1)
	for {
		isPseudoPrime:=true
		isEqInitMods:=true
		for i:=0; i<len(rnd.initMods); i++{
			m:=big.NewInt(0).Mod(a,rnd.bigSmallPrimes[i])
			if m.Cmp(big0)==0 {
				isPseudoPrime=false
			}
			if m.Cmp(rnd.initMods[i])!=0{
				isEqInitMods=false
				
			}

		}
		if isPseudoPrime {
			if probablyPrime(a,64){
				return a,false
			}
		}
		if isEqInitMods {
			break
		}
		a.Add(a,big1)
	}
	for {
		var wg sync.WaitGroup
		var mut sync.Mutex
		found:=false
		var q *big.Int
		wg.Add(2)
		a1:=big.NewInt(0).SetBytes(a.Bytes())
		a2:=big.NewInt(0).Add(a,rnge2)
		go func(){
			for i:=0; i<len(rnd.skips)&&!found; i++ {
				a1.Add(a1,rnd.skips[i])
				if probablyPrime(a1,64) {
					mut.Lock()
					found=true
					q=a1
					mut.Unlock()
					break
				}
			}
			wg.Done()
		}()
		go func() {
			for i:=0; i<len(rnd.skips)&&!found; i++ {
				a2.Add(a2,rnd.skips[i])
				if probablyPrime(a2,64) {
					mut.Lock()
					found=true
					q=a2
					mut.Unlock()
					break
				}
			}
			wg.Done()
		}()
		wg.Wait()
		if found {
			return q,false
		}
		a=a2
		if a.BitLen()>limit {
			break
		}
		
	}
	return nil,true
}

func probablyPrime(p *big.Int,n int)bool{
	result:=p.ProbablyPrime(1)
	end:=false
	if n<=4 &&n>1{
		return result && millerRabin(p,n-1,func ()bool{return false})
	}
	if n<=1 {
		return result
	}
	if !result {
		return false
	}
	n4:=n/4
	if n4*4<n {
		n4++
	}
	var wg sync.WaitGroup
	var mut sync.Mutex
	wg.Add(4)
	for i:=0; i<4; i++ {
		go func(){
			r:=millerRabin(p,n4,func ()bool {
				return end
			})
			if result&&!r{
				mut.Lock()
				result=false
				end=true
				mut.Unlock()
			}
			wg.Done()
		
		}()
	}
	wg.Wait()

	return result
}


