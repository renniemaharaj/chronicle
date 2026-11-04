package campaign

import "testing"

func TestSeparatePlatformOrderAndPartialPublication(t *testing.T) {
	value := Campaign{
		Revision:  1,
		Platforms: []string{"facebook", "youtube"},
		Parts: []Part{
			{
				Number: 1,
				Status: "recorded",
				File:   "one.mp4",
			}, {
				Number: 2,
				Status: "recorded",
				File:   "two.mp4",
			},
		},
	}
	apply := func(part int, platform, address string) {
		t.Helper()
		if err := Apply(&value, ActionInput{
			Revision: value.Revision,
			Part:     part,
			Action:   "publish",
			Platform: platform,
			URL:      address,
		}); err != nil {
			t.Fatal(err)
		}
	}
	apply(1, "facebook", "https://facebook.com/reel/111")
	if value.Parts[0].Status != "recorded" || ReceiptFor(value.Parts[0], "youtube") != nil {
		t.Fatal("one receipt must not mark every selected platform published")
	}

	if Apply(&value, ActionInput{
		Revision: value.Revision,
		Part:     1,
		Action:   "retake",
	}) == nil {
		t.Fatal("a partially published episode cannot be reset")
	}

	apply(2, "facebook", "https://facebook.com/reel/222")
	if Apply(&value, ActionInput{
		Revision: value.Revision,
		Part:     2,
		Action:   "publish",
		Platform: "youtube",
		URL:      "https://youtube.com/shorts/abcdefghijk",
	}) == nil {
		t.Fatal("YouTube order is independent of Facebook receipts")
	}

	apply(1, "youtube", "https://youtu.be/abcdefghijk")
	if value.Parts[0].Status != "published" || value.Parts[0].Publication.URL != "https://www.facebook.com/reel/111" {
		t.Fatal("both receipts are required, and the legacy Facebook receipt must survive")
	}

	apply(2, "youtube", "https://www.youtube.com/watch?v=lmnopqrstuv")
	if ReceiptFor(value.Parts[1], "youtube").URL != "https://www.youtube.com/shorts/lmnopqrstuv" {
		t.Fatal("YouTube receipt aliases should normalize to the Short URL")
	}
}

func TestPublicationTargetsRejectLookalikes(t *testing.T) {
	for _, address := range []string{"https://youtube.com.evil.test/shorts/abcdefghijk", "https://youtube.com/shorts/not-valid", "https://youtube.com/redirect?v=abcdefghijk"} {
		if _, _, err := publicationURL(address, "youtube"); err == nil {
			t.Fatal("invalid YouTube receipt accepted")
		}
	}

	if ValidatePlatforms([]string{"facebook", "facebook"}) == nil || ValidatePlatforms([]string{"instagram"}) == nil {
		t.Fatal("unsupported or duplicate targets must fail")
	}
}
