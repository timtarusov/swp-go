package models

type GapsClosure struct {
	//спрос
	Demand float64

	//фактическая численность на этот год, с учетом закрытия разрывов в прошлом году
	Headcount float64

	// доп потребность в этом году
	AdditionalDemand float64

	//количество сотрудников, которые соответствуют профилю
	CurrentGood float64

	// количество сотрудников, которые не соответствуют профилю
	CurrentBad float64

	// количество плохих сотрудников, которые уйдут и мы не будет их удерживать
	AttritionBuyBad float64

	// количество хороших сотрудников, которые могут уйти, но мы их удержим
	AttritionRetainGood float64

	// то же для плохих сотрудников
	AttritionRetainBad float64

	// #стоимость стратегий в этом году
	CostOfStrategies map[string]*CostOfStrategy

	// #определяем порядок применения стратегий
	// #сначала самая дешевая из предпочтительных, потом дорогая из предпочтительных
	// #потом дешевая из альтернативных
	// #в конце - дорогая из альтернативных
	StratPriority []string

	// #считаем ограничения по стратегиям
	// #Buy - не больше, чем доступность на рынке
	// #Build - не больше, чем доступно кандидатов в других ролях
	// #Retain - не больше, чем уходящих добровольно сотрудников, которых мы оставляем
	// #Borrow - TODO
	StratLimits map[string]float64

	// #фактическое применение стратегий
	// #покрытие доп потребности и закрытие оттока
	StratApplication *stratApplication
}

//NewGapsClosure returns a new GapsClosure object
func NewGapsClosure() *GapsClosure {
	return &GapsClosure{
		CostOfStrategies: make(map[string]*CostOfStrategy),
		StratLimits:      make(map[string]float64),
		StratApplication: NewStratApplication(),
	}
}

func (r *Role) CalculateAnnualGapsClosure(year int, suffThreshold float64) *GapsClosure {
	gc := NewGapsClosure()
	gc.Demand = r.ProjectedHeadcount.Headcount[year]
	gc.AdditionalDemand = r.ProjectAdditionalDemand(year)
	gc.CurrentGood, gc.CurrentBad = r.CurrentGoodAndBadEmployees(year, suffThreshold)
	return gc
	// TODO!!! Stopped here
}
