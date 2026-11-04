package campaign

import "regexp"

var sentenceEnd = regexp.MustCompile(`[.!?]["')\]]*$`)

// Keep sentences together across reels; long sentences may exceed the duration target.
// Caption offsets still cover the exact frozen passage, including its whitespace.
func readingUnits(items []Passage) [][]Passage {
	units := [][]Passage{}

	for _, item := range items {
		locations := words.FindAllStringIndex(item.Text, -1)
		start, sentenceStart := 0, 0

		for index, location := range locations {
			if index+1 < len(locations) && !sentenceEnd.MatchString(item.Text[location[0]:location[1]]) {
				continue
			}

			chunks := []Passage{}

			for first := sentenceStart; first <= index; {
				count := min(24, index-first+1)
				end := len(item.Text)
				if first+count < len(locations) {
					end = locations[first+count][0]
				}

				chunk := item
				chunk.Text, chunk.Start, chunk.End = item.Text[start:end], start, end
				chunks = append(chunks, chunk)
				start = end
				first += count
			}

			units = append(units, chunks)
			sentenceStart = index + 1
		}
	}

	return units
}
