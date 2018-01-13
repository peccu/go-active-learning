package classifier

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"sync"

	"github.com/syou6162/go-active-learning/lib/evaluation"
	"github.com/syou6162/go-active-learning/lib/example"
	"github.com/syou6162/go-active-learning/lib/feature"
	"github.com/syou6162/go-active-learning/lib/util"
)

type MIRAClassifier struct {
	weight map[string]float64
	c      float64
}

func newMIRAClassifier(c float64) *MIRAClassifier {
	return &MIRAClassifier{make(map[string]float64), c}
}

func NewMIRAClassifier(examples example.Examples, c float64) *MIRAClassifier {
	train := util.FilterLabeledExamples(examples)
	model := newMIRAClassifier(c)
	for iter := 0; iter < 30; iter++ {
		util.Shuffle(train)
		for _, example := range train {
			model.learn(*example)
		}
	}
	return model
}

type MIRAResult struct {
	mira   MIRAClassifier
	FValue float64
}

type MIRAResultList []MIRAResult

func (l MIRAResultList) Len() int           { return len(l) }
func (l MIRAResultList) Less(i, j int) bool { return l[i].FValue < l[j].FValue }
func (l MIRAResultList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func NewMIRAClassifierByCrossValidation(examples example.Examples) *MIRAClassifier {
	train, dev := util.SplitTrainAndDev(util.FilterLabeledExamples(examples))

	params := []float64{100, 50, 10.0, 5.0, 1.0, 0.5, 0.1, 0.05, 0.01}
	miraResults := MIRAResultList{}

	wg := &sync.WaitGroup{}
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	models := make([]*MIRAClassifier, len(params))
	for idx, c := range params {
		wg.Add(1)
		go func(idx int, c float64) {
			defer wg.Done()
			model := NewMIRAClassifier(train, c)
			models[idx] = model
		}(idx, c)
	}
	wg.Wait()

	for _, model := range models {
		c := model.c
		devPredicts := make([]example.LabelType, len(dev))
		for i, example := range dev {
			devPredicts[i] = model.Predict(example.Fv)
		}
		accuracy := evaluation.GetAccuracy(ExtractGoldLabels(dev), devPredicts)
		precision := evaluation.GetPrecision(ExtractGoldLabels(dev), devPredicts)
		recall := evaluation.GetRecall(ExtractGoldLabels(dev), devPredicts)
		f := (2 * recall * precision) / (recall + precision)
		fmt.Fprintln(os.Stderr, fmt.Sprintf("C:%0.03f\tAccuracy:%0.03f\tPrecision:%0.03f\tRecall:%0.03f\tF-value:%0.03f", c, accuracy, precision, recall, f))
		miraResults = append(miraResults, MIRAResult{*model, f})
	}

	sort.Sort(sort.Reverse(miraResults))
	return &miraResults[0].mira
}

func (model *MIRAClassifier) learn(example example.Example) {
	tmp := float64(example.Label) * model.PredictScore(example.Fv) // y w^T x
	loss := 0.0
	if tmp < 1.0 {
		loss = 1 - tmp
	}

	norm := float64(len(example.Fv) * len(example.Fv))
	// tau := math.Min(model.c, loss/norm) // update by PA-I
	tau := loss / (norm + 1.0/model.c) // update by PA-II

	if tau != 0.0 {
		for _, f := range example.Fv {
			w, _ := model.weight[f]
			model.weight[f] = w + tau*float64(example.Label)
		}
	}
}

func (model MIRAClassifier) PredictScore(features feature.FeatureVector) float64 {
	result := 0.0
	for _, f := range features {
		w, ok := model.weight[f]
		if ok {
			result = result + w*1.0
		}
	}
	return result
}

func (model MIRAClassifier) Predict(features feature.FeatureVector) example.LabelType {
	if model.PredictScore(features) > 0 {
		return example.POSITIVE
	}
	return example.NEGATIVE
}

func (model MIRAClassifier) SortByScore(examples example.Examples) example.Examples {
	return SortByScore(model, examples)
}

func (model MIRAClassifier) GetWeight(f string) float64 {
	w, ok := model.weight[f]
	if ok {
		return w
	}
	return 0.0
}

func (model MIRAClassifier) GetActiveFeatures() []string {
	result := make([]string, 0)
	for f := range model.weight {
		result = append(result, f)
	}
	return result
}