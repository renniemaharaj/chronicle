package campaign

import "testing"

func TestRetakeKeepsCompletedFilesAndProtectsPostedParts(t *testing.T) {
	value, _ := Plan(syntheticProject(), 60)
	apply := func(input ActionInput) {
		t.Helper()

		input.Revision, input.Part = value.Revision, 1
		if err := Apply(&value, input); err != nil {
			t.Fatal(err)
		}
	}
	apply(ActionInput{Action: "begin"})
	apply(ActionInput{
		Action:   "complete",
		TakeID:   value.Parts[0].TakeID,
		File:     "first.mp4",
		Duration: 50,
	})
	apply(ActionInput{Action: "retake"})
	if value.Parts[0].File != "" || value.Parts[0].Status != "planned" || value.Parts[0].Takes[0].File != "first.mp4" {
		t.Fatal("retaking must retain the earlier completed file in history")
	}

	apply(ActionInput{Action: "begin"})
	apply(ActionInput{
		Action:   "complete",
		TakeID:   value.Parts[0].TakeID,
		File:     "second.mp4",
		Duration: 49,
	})
	apply(ActionInput{
		Action: "publish",
		URL:    "https://www.facebook.com/reel/333",
	})
	if len(value.Parts[0].Takes) != 2 || value.Parts[0].File != "second.mp4" {
		t.Fatal("the latest complete take should be selected without losing earlier takes")
	}

	if err := Apply(&value, ActionInput{
		Revision: value.Revision,
		Part:     1,
		Action:   "retake",
	}); err == nil {
		t.Fatal("posted parts cannot be reset for another recording")
	}
}
