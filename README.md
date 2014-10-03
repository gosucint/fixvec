fixvec
======

fixvec is a Go library for a vector representation of values using fixed bits.

fixvec provides a vector representation of value using fixed bits.
Conceptually, fixvec represents a vector V[0...num), and each
value V[i] can represent in [0...2^(blen)).
The total working space is num * blen bits (+ some small overhead).

Usage
=====

```
import "github.com/hillbig/fixvec"

// fv represents V[0...1000), 0 <= V[i] < 2^10
// fv requires 1000 * 10 = 10000bits = 1250 bytes.
fv := fixvec.New(1000, 10)  
                               
// use Set(ind, val) to set value V[ind]= val
fv.Set(10, 777)

// use Get(ind) to get value V[ind]
fmt.Printf("%d\n", fv.Get(10)) // 777

// Encode to binary representation
bytes, err := fv.MarshalBinary() 

// Decode from binary presentation
newfv := fixvec.NewFixVec(0, 0)
err := newfv.UnmarshalBinary(bytes) 
```
