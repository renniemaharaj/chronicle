package campaign

import (
	"crypto/rand"
	"fmt"
	"math"
	"strings"
	"time"
)

func Apply(value *Campaign, input ActionInput) error {
	if input.Revision != value.Revision {
		return fmt.Errorf("campaign changed; reload before continuing")
	}

	if input.Part < 1 || input.Part > len(value.Parts) {
		return fmt.Errorf("unknown reel part")
	}

	part := &value.Parts[input.Part-1]
	now := time.Now().UTC().Format(time.RFC3339Nano)

	switch input.Action {
	case "destination":
		if err := saveDestination(value, input, now); err != nil {
			return err
		}
	case "begin":
		if part.Status != "planned" {
			return fmt.Errorf("only an unrecorded part can start a take")
		}
		part.Status, part.TakeID = "recording", "take-"+rand.Text()
	case "complete", "interrupt":
		if part.Status != "recording" || part.TakeID != input.TakeID {
			return fmt.Errorf("recording take does not match; inspect the active take before retrying")
		}
		if input.Action == "complete" {
			if strings.TrimSpace(input.File) == "" || len(input.File) > 4096 || input.Duration <= 0 || math.IsInf(input.Duration, 0) || math.IsNaN(input.Duration) {
				return fmt.Errorf("a completed take requires its output file and measured duration")
			}

			part.Status, part.File, part.Duration = "recorded", input.File, input.Duration
			part.Takes = append(part.Takes, Take{
				ID:          part.TakeID,
				File:        input.File,
				Duration:    input.Duration,
				CompletedAt: now,
			})
		} else {
			part.Status, part.TakeID = "planned", ""
		}
	case "retake":
		if part.Status != "recorded" || anyPublication(*part) {
			return fmt.Errorf("only an unpublished completed part can be recorded again")
		}
		part.Status, part.TakeID, part.File, part.Duration = "planned", "", "", 0
	case "publish":
		if err := publishPlatform(value, part, input.URL, now, Platform(input.Platform)); err != nil {
			return err
		}
	case "replace":
		if err := replacePublication(value, part, input, now); err != nil {
			return err
		}
	case "metrics":
		if ReceiptFor(*part, input.Platform) == nil || input.Metrics == nil {
			return fmt.Errorf("performance observations require a published part")
		}
		metric := *input.Metrics
		metric.Platform = Platform(input.Platform)
		if err := validateMetrics(metric); err != nil {
			return err
		}
		metric.Method = "manual"
		part.Metrics = append(part.Metrics, metric)
	default:
		return fmt.Errorf("unknown campaign action")
	}

	value.Revision++
	value.Audit = append(value.Audit, Activity{
		At:     now,
		Action: input.Action,
		Part:   input.Part,
	})

	return nil
}
