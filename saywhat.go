package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type subTitle struct {
	// Index shows which subtitle this subtitle is in the sequence.
	Index int
	// Appear is the time that this subtitle is supposed to appear on screen. The time at the
	// beginning of the video is November 30th, 2001 @ 12:00 AM UTC.
	Appear time.Time
	// Disappear is  time that this subtitle is supposed to disappear from the screen. The time
	// at the beginning of the video is November 30th, 2001 @ 12:00 AM UTC.
	Disappear time.Time
	// Text is content of the subtitle.
	Text string
}

func parseTimeCode(s string) (time.Time, error) {
	milliAndRest := strings.Split(s, ",")
	if len(milliAndRest) != 2 {
		return time.Time{}, fmt.Errorf("parseTimeCode: expected len 2 array when splitting with comma, got len: %d", len(milliAndRest))
	}
	hhMMSS := strings.Split(milliAndRest[0], ":")
	if len(hhMMSS) != 3 {
		return time.Time{}, fmt.Errorf("parseTimeCode: expected len 3 array when splitting with colon , got len: %d", len(hhMMSS))
	}
	milli, err := strconv.Atoi(milliAndRest[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("parseTimeCode: error when parsing millis: %s", err)
	}
	hours, err := strconv.Atoi(hhMMSS[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("parseTimeCode: error when parsing hours: %s", err)
	}
	mins, err := strconv.Atoi(hhMMSS[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("parseTimeCode: error when parsing mins: %s", err)
	}
	secs, err := strconv.Atoi(hhMMSS[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("parseTimeCode: error when parsing secs: %s", err)
	}
	return time.Date(0, 0, 0, hours, mins, secs, milli*1000000, time.UTC), nil
}

func parseSingle(indexField, timeField, textField string) (subTitle, error) {
	index, err := strconv.Atoi(indexField)
	if err != nil {
		return subTitle{}, fmt.Errorf("parseSingle: index parsing failed, err: %s", err)
	}

	times := strings.Split(timeField, "-->")
	if len(times) != 2 {
		return subTitle{}, fmt.Errorf("parseSingle: times parsing failed, expected length 2 array, instead got %d", len(times))
	}

	appear, err := parseTimeCode(strings.TrimSpace(times[0]))
	if err != nil {
		return subTitle{}, fmt.Errorf("parseSingle: error when parsing appear time, err: %s", err)
	}

	disappear, err := parseTimeCode(strings.TrimSpace(times[1]))
	if err != nil {
		return subTitle{}, fmt.Errorf("parseSingle: error when parsing disappear time, err: %s", err)
	}

	return subTitle{
		Index:     index,
		Appear:    appear,
		Disappear: disappear,
		Text:      textField,
	}, nil
}

func parseSubs(src io.Reader) ([]subTitle, error) {
	var subs []subTitle
	var fields []string
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		field := scanner.Text()
		if len(field) == 0 {
			if len(fields) >= 3 {
				sub, err := parseSingle(fields[0], fields[1], strings.Join(fields[2:], " "))
				if err != nil {
					return nil, fmt.Errorf("parse: error when parsing src, err: %s", err)
				}
				subs = append(subs, sub)
				fields = nil
			}
		} else {
			fields = append(fields, field)
		}
	}
	if len(fields) != 0 {
		return nil, errors.New("parse, incomplete subtitle at EOF")
	}
	return subs, nil
}

func main() {
	args := os.Args
	file, err := os.Open(args[1])
	if err != nil {
		log.Fatal(err)
	}
	subs, err := parseSubs(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(subs)
}
