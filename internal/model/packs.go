package model

type (
	//Packs is a map of size to Count and packed Quantity
	Packs map[int]int
)

func (p Packs) AddPacks(size, count int) {
	p[size] += count
}

func (p Packs) AddPack(size int) {
	p[size] += 1
}

func (p Packs) RemovePack(size int) {
	record := p[size]
	if record == 0 {
		return
	}
	record -= 1
	if record <= 0 {
		delete(p, size)
		return
	}
	p[size] = record

}

func (p Packs) SetCount(size, count int) {
	if count == 0 {
		delete(p, size)
		return
	}

	p[size] = count
}
