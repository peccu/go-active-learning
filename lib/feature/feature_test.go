package feature

import (
	"fmt"
	"testing"
)

func TestIsJapanese(t *testing.T) {
	text := "ほげ"
	if !isJapanese(text) {
		t.Error(fmt.Printf("%s should be Japanese", text))
	}
	text = "This is a pen."
	if isJapanese(text) {
		t.Error(fmt.Printf("%s should be not Japanese", text))
	}
}

func TestEngNounFeatures(t *testing.T) {
	text := "Hello World!"
	fv := extractEngNounFeatures(text, "")
	if len(fv) != 2 {
		t.Error(fmt.Printf("Size of feature vector for %s should be 2", text))
	}
}

func TestExtractPath(t *testing.T) {
	url := "http://b.hatena.ne.jp/search/text?safe=on&q=nlp&users=50"
	path := "/search/text"
	if ExtractPath(url) != path {
		t.Error(fmt.Printf("path should be %s", path))
	}
}

func TestExtractHostFeature(t *testing.T) {
	url := "http://b.hatena.ne.jp/search/text?safe=on&q=nlp&users=50"
	hostFeature := "HOST:b.hatena.ne.jp"
	if ExtractHostFeature(url) != hostFeature {
		t.Error(fmt.Printf("Host feature should be %s", hostFeature))
	}
}
