package util

import (
	"testing"

	"github.com/syou6162/go-active-learning/lib/example"
)

func TestFilterLabeledExamples(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", example.POSITIVE)
	e2 := example.NewExample("http://www.yasuhisay.info", example.NEGATIVE)
	e3 := example.NewExample("http://google.com", example.UNLABELED)

	examples := FilterLabeledExamples(example.Examples{e1, e2, e3})
	if len(examples) != 2 {
		t.Error("Number of labeled examples should be 2")
	}
}

func TestFilterUnlabeledExamples(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", example.POSITIVE)
	e2 := example.NewExample("http://www.yasuhisay.info", example.NEGATIVE)
	e3 := example.NewExample("http://google.com", example.UNLABELED)
	e3.Title = "Google"

	examples := FilterUnlabeledExamples(example.Examples{e1, e2, e3})
	if len(examples) != 1 {
		t.Error("Number of unlabeled examples should be 1")
	}
}

func TestFilterStatusCodeOkExamples(t *testing.T) {
	e1 := example.NewExample("http://b.hatena.ne.jp", example.POSITIVE)
	e1.StatusCode = 200
	e2 := example.NewExample("http://www.yasuhisay.info", example.NEGATIVE)
	e2.StatusCode = 404
	e3 := example.NewExample("http://google.com", example.UNLABELED)
	e3.StatusCode = 304

	examples := FilterStatusCodeOkExamples(example.Examples{e1, e2, e3})
	if len(examples) != 1 {
		t.Error("Number of examples (status code = 200) should be 1")
	}
}

func TestRemoveDuplicate(t *testing.T) {
	args := []string{"hoge", "fuga", "piyo", "hoge"}

	result := RemoveDuplicate(args)
	if len(result) != 3 {
		t.Error("Number of unique string in args should be 3")
	}
}

func TestSplitTrainAndDev(t *testing.T) {
	e1 := example.NewExample("http://a.hatena.ne.jp", example.POSITIVE)
	e2 := example.NewExample("http://www.yasuhisay.info", example.NEGATIVE)
	e3 := example.NewExample("http://google.com", example.UNLABELED)
	e4 := example.NewExample("http://a.hatena.ne.jp", example.POSITIVE)
	e5 := example.NewExample("http://www.yasuhisay.info", example.NEGATIVE)
	e6 := example.NewExample("http://a.hatena.ne.jp", example.POSITIVE)
	e7 := example.NewExample("http://www.yasuhisay.info", example.NEGATIVE)
	e8 := example.NewExample("http://google.com", example.UNLABELED)
	e9 := example.NewExample("http://a.hatena.ne.jp", example.POSITIVE)
	e10 := example.NewExample("http://www.yasuhisay.info", example.NEGATIVE)

	train, dev := SplitTrainAndDev(example.Examples{e1, e2, e3, e4, e5, e6, e7, e8, e9, e10})
	if len(train) != 8 {
		t.Error("Number of training examples should be 8")
	}
	if len(dev) != 2 {
		t.Error("Number of dev examples should be 2")
	}
}
