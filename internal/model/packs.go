package model

type (
	CountAndQuantity struct {
		Count    int
		Quantity int
	}
	//Packs is a map of size to Count and packed Quantity
	Packs map[int]CountAndQuantity
)

func NewCountAndQuantity(c, q int) CountAndQuantity {
	return CountAndQuantity{Count: c, Quantity: q}
}

func (p Packs) AddPacks(size, count int, packedQuantity ...int) {
	record := p[size]
	record.Count += count
	if len(packedQuantity) != 0 {
		record.Quantity += packedQuantity[0]
	}
	p[size] = record
}

func (p Packs) AddPack(size int, packedQuantity ...int) {
	record := p[size]
	record.Count += 1
	if len(packedQuantity) != 0 {
		record.Quantity += packedQuantity[0]
	}
	p[size] = record
}

func (p Packs) RemovePack(size int) {
	record := p[size]
	if record.Count == 0 {
		return
	}
	record.Count -= 1
	if record.Count <= 0 {
		delete(p, size)
		return
	}
	p[size] = record

}

func (p Packs) SetCount(size, count int, packedQuantity ...int) {
	if count == 0 {
		delete(p, size)
		return
	}
	record := p[size]
	record.Count = count
	if len(packedQuantity) != 0 {
		record.Quantity = packedQuantity[0]
	}
	p[size] = record
}
