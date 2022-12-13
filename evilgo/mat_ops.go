package evilgo

func Transpose[T any](f [][]T) [][]T {
	w := len(f[0])
	h := len(f)
	newForest := make([][]T, w)
	for i := 0; i < w; i++ {
		newForest[i] = make([]T, h)
	}
	for i, l := range f {
		for j, num := range l {
			newForest[j][i] = num
		}
	}

	return newForest
}

func Reverse[T any](f [][]T) [][]T {
	w := len(f[0])
	h := len(f)
	newMat := make([][]T, h)

	for i, l := range f {
		newMat[i] = make([]T, w)
		for j, k := 0, len(l)-1; j <= k; j, k = j+1, k-1 {
			newMat[i][j], newMat[i][k] = f[i][k], f[i][j]
		}
	}
	return newMat
}
