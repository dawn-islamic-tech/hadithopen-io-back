package empty

func Coalesce[V string](a, b V) V {
	var k V
	if a == k {
		return b
	}

	return a
}
