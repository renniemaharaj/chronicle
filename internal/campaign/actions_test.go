package campaign

import (
	"testing"
	"time"
)

func TestTakeRecoveryAndOrderedPublication(t *testing.T) {
	value, _ := Plan(syntheticProject(), 15)
	apply := func(part int, action ActionInput) {
		t.Helper()

		action.Part, action.Revision = part, value.Revision
		if err := Apply(&value, action); err != nil {
			t.Fatal(err)
		}
	}
	apply(1, ActionInput{Action: "begin"})

	oldTake := value.Parts[0].TakeID
	apply(1, ActionInput{
		Action: "interrupt",
		TakeID: oldTake,
	})
	apply(1, ActionInput{Action: "begin"})
	if err := Apply(&value, ActionInput{
		Part:     1,
		Revision: value.Revision,
		Action:   "complete",
		TakeID:   oldTake,
		File:     "old.mp4",
		Duration: 10,
	}); err == nil {
		t.Fatal("a superseded take cannot complete")
	}

	apply(1, ActionInput{
		Action:   "complete",
		TakeID:   value.Parts[0].TakeID,
		File:     "one.mp4",
		Duration: 15,
	})
	apply(2, ActionInput{Action: "begin"})
	apply(2, ActionInput{
		Action:   "complete",
		TakeID:   value.Parts[1].TakeID,
		File:     "two.mp4",
		Duration: 15,
	})
	if err := Apply(&value, ActionInput{
		Part:     2,
		Revision: value.Revision,
		Action:   "publish",
		URL:      "https://www.facebook.com/reel/222",
	}); err == nil {
		t.Fatal("part two cannot publish before part one")
	}

	apply(1, ActionInput{
		Action: "publish",
		URL:    "https://m.facebook.com/reel/111/?tracking=yes",
	})
	if err := Apply(&value, ActionInput{
		Part:     2,
		Revision: value.Revision,
		Action:   "publish",
		URL:      "https://facebook.com/synthetic/videos/111",
	}); err == nil {
		t.Fatal("duplicate publication receipts must fail")
	}

	apply(2, ActionInput{
		Action: "publish",
		URL:    "https://facebook.com/reel/222",
	})
	if err := Apply(&value, ActionInput{
		Part:     1,
		Revision: value.Revision - 1,
		Action:   "metrics",
	}); err == nil {
		t.Fatal("stale writers cannot overwrite newer checkpoints")
	}

	views := int64(0)
	apply(1, ActionInput{
		Action: "metrics",
		Metrics: &Metrics{
			ObservedAt: time.Now().UTC().Format(time.RFC3339),
			Views:      &views,
		},
	})
	if value.Parts[0].Metrics[0].Views == nil || value.Parts[0].Metrics[0].Revenue != nil {
		t.Fatal("observed zero and unknown performance are different")
	}
}

func TestInvalidPublicationAndMetrics(t *testing.T) {
	for _, address := range []string{"javascript:alert(1)", "https://facebook.com.evil.test/reel/1", "https://www.facebook.com/", "https://user@facebook.com/reel/1"} {
		value := Campaign{Parts: []Part{{
			Number: 1,
			Status: "recorded",
		}}}
		if err := publish(&value, &value.Parts[0], address, "now"); err == nil {
			t.Fatalf("accepted invalid receipt %s", address)
		}
	}

	negative := int64(-1)

	for _, metric := range []Metrics{
		{
			ObservedAt: "not-a-date",
			Views:      &negative,
		},
		{
			ObservedAt: time.Now().Format(time.RFC3339),
			Views:      &negative,
		},
		{ObservedAt: time.Now().Format(time.RFC3339)},
	} {
		if validateMetrics(metric) == nil {
			t.Fatal("invalid performance observation was accepted")
		}
	}
}
