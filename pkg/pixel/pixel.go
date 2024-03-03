package pixel

type Pixel struct {
	R, G, B, A uint32
}

func NewPixel(r, g, b, a uint32) *Pixel {
	return &Pixel{R: r, G: g, B: b, A: a}
}

func (p *Pixel) Avg() float64 {
	return float64((p.R/255 + p.G/255 + p.B/255) / 3)
}
