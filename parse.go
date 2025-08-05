package main

import (
	"regexp"
	"strconv"
	"strings"
)

const directivePrefix = "!"

func slidesFromString(s string) (b slideshow, e error) {
	lines := linesFromString(s)
	blocks := lines.toBlocks()
	for _, block := range blocks {
		if len(block) == 0 {
			continue
		}

		if strings.HasPrefix(block[0].string, directivePrefix+"global") {
			dir, args := parseDirective(block[0].string)
			if args != nil {
				e = parseError{l: block[0].index, r: noArgs, d: dir, t: block[0].string}
				return
			}
			if err := b.parseGlobalBlock(block[1:]); err != nil {
				e = err
				return
			}
			continue
		}

		slide, err := parseSlide(block)

		if err != nil {
			e = err
			return
		}

		b.Slides = append(b.Slides, slide)
	}
	return
}

func (ss *slideshow) parseGlobalBlock(block lines) (e error) {
	for _, line := range block {
		if strings.HasPrefix(line.string, "#") {
			continue
		}
		if !strings.HasPrefix(line.string, directivePrefix) {
			return parseError{l: line.index, t: line.string, r: onlyDirectives}
		}
		dir, args := parseDirective(line.string)

		switch dir {
		case "global":
			return parseError{l: line.index, d: dir, r: redundantDirective}
		case "noautoplay":
			if args != nil {
				e = parseError{l: line.index, d: dir, t: *args, r: noArgs}
				return
			}
			if ss.NoAutoPlay {
				return parseError{l: line.index, d: dir, r: redundantDirective}
			}
			ss.NoAutoPlay = true
		case "pfoot":
			if args == nil {
				return parseError{l: line.index, r: takesArgs, d: dir, t: line.string}
			}
			if ss.PrimaryFooter.Has {
				return parseError{l: line.index, d: dir, r: redundantDirective}
			}
			ss.PrimaryFooter = wrap(*args)
		case "sfoot":
			if args == nil {
				return parseError{l: line.index, r: takesArgs, d: dir, t: line.string}
			}
			if ss.SecondaryFooter.Has {
				return parseError{l: line.index, d: dir, r: redundantDirective}
			}
			ss.SecondaryFooter = wrap(*args)
		case "exec":
			if args == nil {
				return parseError{l: line.index, r: takesArgs, d: dir, t: line.string}
			}
			if err := hasForbiddenPatterns(line, dir); err != nil {
				return err
			}
			if ss.Script.Has {
				return parseError{l: line.index, d: dir, r: redundantDirective}
			}
			ss.Script = wrap(*args)
		case "styles":
			if args == nil {
				return parseError{l: line.index, r: takesArgs, d: dir, t: line.string}
			}
			if err := hasForbiddenPatterns(line, dir); err != nil {
				return err
			}
			if ss.Styles.Has {
				return parseError{l: line.index, d: dir, r: redundantDirective}
			}
			ss.Styles = wrap(*args)
		case "time":
			if args == nil {
				return parseError{l: line.index, r: takesArgs, d: dir, t: line.string}
			}
			if err := hasForbiddenPatterns(line, dir); err != nil {
				return err
			}
			if ss.CustomDuration.Has {
				return parseError{l: line.index, d: dir, r: redundantDirective}
			}
			dur, err := strconv.ParseFloat(*args, 64)
			if err != nil {
				return parseError{l: line.index, r: invalidNum, t: *args}
			}
			ss.CustomDuration = wrap(dur)
		case "id", "image", "title":
			return parseError{l: line.index, r: invalidDirective, d: dir}
		default:
			return parseError{l: line.index, r: uknownDirective, d: dir}
		}
	}
	return nil
}

func parseSlide(block lines) (s slide, e error) {
	hasBullets := false
	hasText := false
	for _, line := range block {
		if strings.HasPrefix(line.string, "#") {
			continue
		}
		if strings.HasPrefix(line.string, directivePrefix) {
			dir, args := parseDirective(line.string)
			switch dir {
			case "title":
				if s.IsTitleSlide {
					e = parseError{l: line.index, t: line.string, r: redundantDirective}
					return
				}
				if args != nil {
					e = parseError{l: line.index, d: dir, t: *args, r: noArgs}
					return
				}
				s.IsTitleSlide = true
			case "id":
				if args == nil {
					e = parseError{l: line.index, d: dir, t: line.string, r: takesArgs}
					return
				}
				if err := hasForbiddenPatterns(line, dir); err != nil {
					e = err
					return
				}
				if s.Id.Has {
					e = parseError{l: line.index, t: line.string, r: redundantDirective}
					return
				}
				s.Id = wrap(*args)
			case "image":
				if args == nil {
					e = parseError{l: line.index, d: dir, t: line.string, r: takesArgs}
					return
				}
				if err := hasForbiddenPatterns(line, dir); err != nil {
					e = err
					return
				}
				if s.Image.Has {
					e = parseError{l: line.index, t: line.string, r: redundantDirective}
					return
				}
				s.Image = wrap(*args)
			case "time":
				if args == nil {
					e = parseError{l: line.index, d: dir, t: line.string, r: takesArgs}
				}

				if err := hasForbiddenPatterns(line, dir); err != nil {
					e = err
					return
				}
				if s.CustomDuration.Has {
					e = parseError{l: line.index, t: line.string, r: redundantDirective}
					return
				}

				num, err := strconv.ParseFloat(*args, 64)

				if err != nil {
					e = parseError{l: line.index, t: *args, r: invalidNum}
					return
				}

				s.CustomDuration = wrap(num)
			case "global", "pfoot", "sfoot", "exec", "styles":
				e = parseError{l: line.index, d: dir, r: invalidDirective}
				return
			default:
				e = parseError{l: line.index, d: dir, r: uknownDirective}
				return
			}
			continue
		}

		if strings.HasPrefix(line.string, "- ") {
			if hasText || s.IsTitleSlide {
				e = parseError{l: line.index, r: cannotMixBulletsAndLines, t: line.string}
				return
			}
			hasBullets = true
			s.Bullets = append(s.Bullets, strings.Replace(line.string, "- ", "", 1))
			continue
		}

		if s.IsTitleSlide || (len(s.HeaderLines) != 1) {
			s.HeaderLines = append(s.HeaderLines, line.string)
			continue
		}

		if hasBullets {
			e = parseError{l: line.index, r: cannotMixBulletsAndLines, t: line.string}
			return
		}
		hasText = true
		s.Lines = append(s.Lines, line.string)
	}
	return
}

func parseDirective(directiveLine string) (directive string, args *string) {
	trim := strings.TrimSpace(directiveLine)
	removeMarker := strings.Replace(trim, directivePrefix, "", 1)
	split := strings.SplitN(removeMarker, " ", 2)
	if len(split) < 2 {
		return split[0], nil
	}
	return split[0], &split[1]
}

var (
	HasDT   = regexp.MustCompile(`[^\\]<<.*?>>`)
	HasBold = regexp.MustCompile(`(?:[^\\]|^)\*(.*?)\*`)
)

func hasForbiddenPatterns(l line, dir string) error {
	if HasDT.MatchString(l.string) {
		return parseError{l: l.index, r: illegalUseOfDT, t: l.string, d: dir}
	}

	if HasBold.MatchString(l.string) {
		return parseError{l: l.index, r: illegalUseofBold, t: l.string, d: dir}
	}
	return nil
}
