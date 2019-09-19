package takanawa

import (
	"fmt"
	"strings"
)

func ParseForwarded(input string) ([]Forwarded, error) {
	s := strings.Split(input, ",")
	fwd := make([]Forwarded, len(s))
	for i, value := range s {
		for _, el := range strings.Split(value, ";") {
			pair := strings.SplitN(el, "=", 2)
			if len(pair) != 2 {
				return nil, fmt.Errorf("parsing pair: %q", el)
			}

			token := strings.ToLower(strings.TrimSpace(pair[0]))
			value := strings.Trim(strings.TrimSpace(pair[1]), "\"")
			switch token {
			case "by":
				fwd[i].By = value
			case "for":
				fwd[i].For = value
			case "host":
				fwd[i].Host = value
			case "proto":
				fwd[i].Proto = value
			}
		}
	}

	return fwd, nil
}

type Forwarded struct {
	By    string
	For   string
	Host  string
	Proto string
}

func (fwd *Forwarded) String() string {
	var str []string
	if len(fwd.By) > 0 {
		str = append(str, "by="+fwd.By)
	}
	if len(fwd.For) > 0 {
		str = append(str, "for="+fwd.For)
	}
	if len(fwd.Host) > 0 {
		str = append(str, "host="+fwd.Host)
	}
	if len(fwd.Proto) > 0 {
		str = append(str, "proto="+fwd.Proto)
	}
	return strings.Join(str, "; ")
}
