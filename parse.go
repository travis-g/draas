package dice

import (
	"regexp"
	"strconv"
)

var (
	// DiceNotationRegex is the compiled RegEx for parsing supported dice
	// notations.
	DiceNotationRegex = regexp.MustCompile(`(?P<count>\d*)d(?P<size>(?:\d{1,}|F))`)

	// DropKeepNotationRegex is the compiled RegEx for parsing drop/keep dice
	// notations (unimplemented).
	DropKeepNotationRegex = regexp.MustCompile(`((?P<group>\{.*\})|(?P<count>\d+)?d(?P<size>\d{1,}|F))(?P<dropkeep>(?P<op>[dk][lh]?)(?P<num>\d{1,}))?`)
)

// ParseGroup parses a notation into a group of dice. It returns the group
// unrolled.
func ParseGroup(notation string) (Group, error) {
	matches := DropKeepNotationRegex.FindStringSubmatch(notation)
	if len(matches) < 3 {
		return nil, &ErrParseError{notation, notation, "", ": failed to identify dice components"}
	}

	// extract named capture groups to map
	components := make(map[string]string)
	for i, name := range DropKeepNotationRegex.SubexpNames() {
		if i != 0 && name != "" {
			components[name] = matches[i]
		}
	}

	group := components["group"]
	if group != "" {
		return nil, &ErrNotImplemented{"group rolls not yet implemented"}
	}

	// Parse and cast dice properties from regex capture values
	count, err := strconv.ParseInt(components["count"], 10, 0)
	if err != nil {
		// either there was an implied count, ex 'd20', or count was invalid
		count = 1
	}

	var dropkeep int
	op := components["op"]
	num, _ := strconv.ParseInt(components["num"], 10, 0)
	switch op {
	case "d", "dl":
		dropkeep = int(num)
	case "k", "kh":
		dropkeep = int(count - num)
	case "dh":
		dropkeep = -int(num)
	case "kl":
		dropkeep = -int(count - num)
	}

	size, err := strconv.ParseUint(components["size"], 10, 0)
	// err is nil, which means a valid uint size
	if err == nil {
		props := GroupProperties{
			Type:     TypePolyhedron,
			Size:     int(size),
			Count:    int(count),
			DropKeep: dropkeep,
			Unrolled: true,
		}
		set := NewGroup(props)
		if dropkeep != 0 {
			set.Drop(dropkeep)
		}
		return set, nil
	}

	// size was not a uint, check for special dice types
	if components["size"] == "F" {
		props := GroupProperties{
			Type:     TypeFate,
			Count:    int(count),
			DropKeep: dropkeep,
			Unrolled: true,
		}
		set := NewGroup(props)
		if dropkeep != 0 {
			set.Drop(dropkeep)
		}
		return set, nil
	}

	return Group{}, &ErrParseError{notation, components["size"], "size", ": invalid size"}
}