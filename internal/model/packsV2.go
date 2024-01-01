package model

type (
	//Packs is a map of size to Count and packed Quantity
	PacksV2 map[int]int
)

func (p PacksV2) AddPacks(size, count int) {
	p[size] += count
}

func (p PacksV2) AddPack(size int) {
	p[size] += 1
}

func (p PacksV2) RemovePack(size int) {
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

func (p PacksV2) SetCount(size, count int) {
	if count == 0 {
		delete(p, size)
		return
	}

	p[size] = count
}
