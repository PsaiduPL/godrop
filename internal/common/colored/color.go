package colored

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

const (
	RedTag   = "<RED>"
	BlueTag  = "<BLUE>"
	GreenTag = "<GREEN>"
)

var allTags = []string{RedTag, BlueTag, GreenTag}

func PrintColored(stringsList ...string) {
	builder := strings.Builder{}
	builder.Grow(len(stringsList) * 5)
	for _, str := range stringsList {
		builder.WriteString(str)
	}
	fmt.Print(builder.String())
}

// PrintColoredWithTags formatuje string (fmt.Sprintf) a następnie wypisuje go
// na stdout, kolorując fragmenty otoczone tym samym tagiem, np.:
//
//	<RED>tekst<RED>       -> "tekst" na czerwono
//	<GREEN>tekst<GREEN>   -> "tekst" na zielono
//
// Tekst poza tagami wypisywany jest bez zmian (domyślnym kolorem).
func BuildColoredString(format string, args ...any) string {
	formattedString := fmt.Sprintf(format, args...)

	builder := strings.Builder{}
	remaining := formattedString

	for {
		tag, idxStart := findClosestTag(remaining)
		if idxStart == -1 {
			builder.WriteString(remaining)
			break
		}

		builder.WriteString(remaining[:idxStart])

		afterOpen := remaining[idxStart+len(tag):]
		idxEnd := strings.Index(afterOpen, tag)

		if idxEnd == -1 {
			builder.WriteString(typeToColoredString(tag, afterOpen))
			break
		}

		coloredText := afterOpen[:idxEnd]
		builder.WriteString(typeToColoredString(tag, coloredText))

		remaining = afterOpen[idxEnd+len(tag):]
	}

	return builder.String()
}

func PrintColoredWithTags(format string, args ...any) {
	formattedString := fmt.Sprintf(format, args...)

	builder := strings.Builder{}
	remaining := formattedString

	for {
		tag, idxStart := findClosestTag(remaining)
		if idxStart == -1 {
			// nie ma już żadnych tagów - dopisz resztę tekstu bez zmian
			builder.WriteString(remaining)
			break
		}

		// tekst przed tagiem otwierającym - bez koloru
		builder.WriteString(remaining[:idxStart])

		afterOpen := remaining[idxStart+len(tag):]
		idxEnd := strings.Index(afterOpen, tag)

		if idxEnd == -1 {
			// brak tagu zamykającego - koloruj resztę stringa i zakończ
			builder.WriteString(typeToColoredString(tag, afterOpen))
			break
		}

		coloredText := afterOpen[:idxEnd]
		builder.WriteString(typeToColoredString(tag, coloredText))

		remaining = afterOpen[idxEnd+len(tag):]
	}

	fmt.Print(builder.String())
}

func findClosestTag(s string) (string, int) {
	closestTag := ""
	closestIdx := -1

	for _, tag := range allTags {
		idx := strings.Index(s, tag)
		if idx == -1 {
			continue
		}
		if closestIdx == -1 || idx < closestIdx {
			closestIdx = idx
			closestTag = tag
		}
	}

	return closestTag, closestIdx
}

func typeToColoredString(tag, str string) string {
	switch tag {
	case RedTag:
		return color.HiRedString(str)
	case BlueTag:
		return color.BlueString(str)
	case GreenTag:
		return color.GreenString(str)
	default:
		return str
	}
}
