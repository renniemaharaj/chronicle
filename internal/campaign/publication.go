package campaign

import (
	"fmt"
	"math"
	"regexp"
	"time"
)

var currencyCode = regexp.MustCompile(`^[A-Z]{3}$`)

func publish(value *Campaign, part *Part, address, now string) error {
	return publishPlatform(value, part, address, now, "facebook")
}

func publishPlatform(value *Campaign, part *Part, address, now, platform string) error {
	if part.Status != "recorded" && part.Status != "published" {
		return fmt.Errorf("only a completed recording can be marked published")
	}

	selected := false

	for _, target := range Targets(*value) {
		selected = selected || platform == target
	}

	if !selected || ReceiptFor(*part, platform) != nil {
		return fmt.Errorf("target is not selected or this episode already has a publication receipt for it")
	}

	for _, previous := range value.Parts[:part.Number-1] {
		if ReceiptFor(previous, platform) == nil {
			return fmt.Errorf("publish earlier parts to this platform first")
		}
	}

	canonical, id, err := publicationURL(address, platform)
	if err != nil {
		return err
	}

	for _, previous := range value.Parts {
		receipt := ReceiptFor(previous, platform)
		if receipt != nil {
			_, previousID, err := publicationURL(receipt.URL, platform)
			if err == nil && previousID == id {
				return fmt.Errorf("this publication receipt is already assigned to a part")
			}
		}
	}

	receipt := Receipt{
		URL:        canonical,
		RecordedAt: now,
		Method:     "manual",
	}
	if part.Publications == nil {
		part.Publications = map[string]Receipt{}
	}

	part.Publications[platform] = receipt
	if platform == "facebook" {
		part.Publication = &receipt
	}

	part.Status = "published"

	for _, target := range Targets(*value) {
		if ReceiptFor(*part, target) == nil {
			part.Status = "recorded"
		}
	}

	return nil
}

func validateMetrics(metric Metrics) error {
	observed, err := time.Parse(time.RFC3339Nano, metric.ObservedAt)
	if err != nil || observed.After(time.Now().Add(time.Minute)) {
		return fmt.Errorf("performance needs a valid observation time, no later than now")
	}

	values := []*int64{metric.Views, metric.Reactions, metric.Comments, metric.Shares}
	hasValue := metric.Revenue != nil

	for _, count := range values {
		if count != nil {
			hasValue = true
			if *count < 0 {
				return fmt.Errorf("performance counts cannot be negative")
			}
		}
	}

	if !hasValue {
		return fmt.Errorf("provide at least one observed performance value")
	}

	if metric.Revenue != nil && (*metric.Revenue < 0 || math.IsInf(*metric.Revenue, 0) || math.IsNaN(*metric.Revenue) || !currencyCode.MatchString(metric.Currency)) {
		return fmt.Errorf("revenue needs a nonnegative amount and a three-letter currency code")
	}

	return nil
}
