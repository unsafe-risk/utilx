package sklmapx

import (
	"math/rand"
	"sync/atomic"
	"time"
)

/******************************************************************************

Original C code: https://xorshift.di.unimi.it/splitmix64.c
Written in 2015 by Sebastiano Vigna (vigna@acm.org)

To the extent possible under law, the author has dedicated all copyright
and related and neighboring rights to this software to the public domain
worldwide. This software is distributed without any warranty.

See <http://creativecommons.org/publicdomain/zero/1.0/>.

*******************************************************************************/

const incrementConstant = 0x9e3779b97f4a7c15

func next(x0 uint64) uint64 {
	x0 = (x0 ^ (x0 >> 30)) * 0xbf58476d1ce4e5b9
	x0 = (x0 ^ (x0 >> 27)) * 0x94d049bb133111eb
	return x0 ^ (x0 >> 31)
}

func splitmix64(state *uint64) uint64 {
	*state += incrementConstant
	v := next(*state)
	return v
}

var root_seed uint64

var _ uint64 = func() uint64 {
	s := rand.NewSource(time.Now().UnixNano())
	root_seed = uint64(s.Int63())
	root_seed += incrementConstant
	root_seed = next(root_seed)
	return 0
}()

func root_next() uint64 {
	s := atomic.AddUint64(&root_seed, incrementConstant)
	return next(s)
}
