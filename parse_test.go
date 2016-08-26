package main

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestParseSingle(t *testing.T) {
	input := `1
00:00:02,153 --> 00:00:03,984
Mac, you gotta see this movie, dude.

`
	actual, err := parseSubs(strings.NewReader(input))
	if err != nil {
		t.Errorf("TestParseSingle: got err: %s", err)
	}

	baseTime := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)

	expected := []subTitle{
		subTitle{
			Index:     1,
			Appear:    baseTime.Add(2*time.Second + 153*time.Millisecond),
			Disappear: baseTime.Add(3*time.Second + 984*time.Millisecond),
			Text:      "Mac, you gotta see this movie, dude.",
		}}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("TestParseSingle:\n expected: %v,\n got: %v", expected, actual)
	}
}

func TestParseMultiple(t *testing.T) {
	input := `1
00:00:00,500 --> 00:00:02,000


2
00:00:02,153 --> 00:00:03,984
Mac, you gotta see this movie, dude.

3
00:00:04,055 --> 00:00:06,114
Really? I thought it was boring.


`
	actual, err := parseSubs(strings.NewReader(input))
	if err != nil {
		t.Errorf("TestParseMultiple: got err: %s", err)
	}

	baseTime := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)

	expected := []subTitle{
		subTitle{
			Index:     1,
			Appear:    baseTime.Add(500 * time.Millisecond),
			Disappear: baseTime.Add(2 * time.Second),
			Text:      "",
		},
		subTitle{
			Index:     2,
			Appear:    baseTime.Add(2*time.Second + 153*time.Millisecond),
			Disappear: baseTime.Add(3*time.Second + 984*time.Millisecond),
			Text:      "Mac, you gotta see this movie, dude.",
		},
		subTitle{
			Index:     3,
			Appear:    baseTime.Add(4*time.Second + 55*time.Millisecond),
			Disappear: baseTime.Add(6*time.Second + 114*time.Millisecond),
			Text:      "Really? I thought it was boring.",
		},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("TestParseMultiple:\n expected: %v, \n got: %v", expected, actual)
	}
}
