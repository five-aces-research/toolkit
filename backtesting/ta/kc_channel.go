package ta

func KetlerChannels(src, c, h, l Series, len1 int, mult float64, MaFn func(src Series, l int) Series) (upper, lower, ma Series) {
	ma = MaFn(src, len1)
	tr := TrueRange(c, h, l)
	atr := Atr(Rma, tr, 10)
	upper = AddF(atr, ma, mult)
	lower = AddF(atr, ma, -mult)
	return
}
