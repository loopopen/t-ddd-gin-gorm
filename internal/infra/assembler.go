package infra

type Assembler[Te any, Tpo any] interface {
	ToEntity(po Tpo) Te
	ToPO(e Te) Tpo
}

type ConvertiblePO[Te any, Tpo any] interface {
	ToEntity() Te
	FromEntity(Te) Tpo
}

type DefAssembler[Te any, Tpo ConvertiblePO[Te, Tpo]] struct{}

func (DefAssembler[Te, Tpo]) ToEntity(po Tpo) Te {
	return po.ToEntity()
}

func (DefAssembler[Te, Tpo]) ToPO(e Te) Tpo {
	var po Tpo
	return po.FromEntity(e)
}
