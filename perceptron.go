package main

type Model struct {
	weight    map[string]float64
	cumWeight map[string]float64
	count     int
}

func NewModel() *Model {
	return &Model{make(map[string]float64), make(map[string]float64), 1}
}

func GetAccuracy(gold []LabelType, predict []LabelType) float64 {
	if len(gold) != len(predict) {
		return 0.0
	}
	sum := 0.0
	for i, v := range gold {
		if v == predict[i] {
			sum += 1.0
		}
	}
	return sum / float64(len(gold))
}

func (model *Model) Learn(example Example) {
	predict := model.predictForTraining(example.fv)
	if example.label != predict {
		for _, f := range example.fv {
			w, _ := model.weight[f]
			cumW, _ := model.cumWeight[f]
			model.weight[f] = w + float64(example.label)*1.0
			model.cumWeight[f] = cumW + float64(model.count)*float64(example.label)*1.0
		}
		model.count += 1
	}
}

func (model *Model) predictForTraining(features FeatureVector) LabelType {
	result := 0.0
	for _, f := range features {
		w, ok := model.weight[f]
		if ok {
			result = result + w*1.0
		}
	}
	if result > 0 {
		return POSITIVE
	}
	return NEGATIVE
}

func (model Model) Predict(features FeatureVector) LabelType {
	result := 0.0
	for _, f := range features {
		w, ok := model.weight[f]
		if ok {
			result = result + w*1.0
		}

		w, ok = model.cumWeight[f]
		if ok {
			result = result - w*1.0/float64(model.count)
		}

	}
	if result > 0 {
		return POSITIVE
	}
	return NEGATIVE
}

func ExtractGoldLabels(examples Examples) []LabelType {
	golds := make([]LabelType, 0, 0)
	for _, e := range examples {
		golds = append(golds, e.label)
	}
	return golds
}
