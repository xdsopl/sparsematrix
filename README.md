# sparse - GF(2) sparse matrix fun
Written in 2017 by <Ahmet Inan> <xdsopl@gmail.com>
To the extent possible under law, the author(s) have dedicated all copyright and related and neighboring rights to this software to the public domain worldwide. This software is distributed without any warranty.
You should have received a copy of the CC0 Public Domain Dedication along with this software. If not, see <http://creativecommons.org/publicdomain/zero/1.0/>.

First example is playing with generator and parity check matrices.

Second example is playing with huge random sparse matrices and their inverses.

We load this image as a sparse vector:
![data.png](https://github.com/xdsopl/sparse/raw/master/data.png)
Do a matrix and column vector multiplication:
![encoded.png](https://github.com/xdsopl/sparse/raw/master/encoded.png)
To get the original vector, we do another matrix and column vector multiplication but using the inverse of the same matrix:
![decoded.png](https://github.com/xdsopl/sparse/raw/master/decoded.png)

```
# go run sparse.go
HammingWeight of P = 497
(Min, Max) of HammingWeightsOfRows of P = 0 8
(Min, Max) of HammingWeightsOfCols of P = 0 5
Wrote GT.png
Wrote H.png
HammingWeight of H*GT = 0
Load data.png
Wrote A.png
Wrote BT.png
AB IsIdentity = true
Wrote AB.png
Wrote encoded.png
Wrote decoded.png

# feh GT.png H.png
# feh A.png BT.png AB.png
# feh data.png encoded.png decoded.png
```

