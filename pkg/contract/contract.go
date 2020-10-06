package contract

import (
	"fmt"
	"regexp"
	"time"

	"github.com/goodsign/monday"
	"github.com/sirupsen/logrus"
)

// Contract contains the contract template, the raw contract text and the extracted params.
type Contract struct {
	RawContractText string
	Template        *Template
	Params          map[string]string
}

// ParseContractString parses a raw string, finds the contract from the store and extracts the params
// Note: It does not verify the params
func ParseContractString(rawContractText string, contractTemplates TemplateStore) (*Contract, error) {

	// first, find the contract Template
	template, err := contractTemplates.FindFromRawContractText(rawContractText)
	if err != nil {
		return nil, err
	}

	contract := &Contract{
		Template:        template,
		RawContractText: rawContractText,
	}

	if err = contract.initParams(); err != nil {
		return nil, err
	}

	return contract, nil
}

func (sc *Contract) initParams() error {
	// extract the params
	r, _ := regexp.Compile(sc.Template.Regexp)
	matchResult := r.FindSubmatch([]byte(sc.RawContractText))
	if len(matchResult) < 1 {
		return fmt.Errorf("%w: could not match the contract template regex", ErrInvalidContractText)
	}
	matches := matchResult[1:]

	if len(matches) != len(sc.Template.TemplateAttributes) {
		return fmt.Errorf("%w: amount of template attributes does not match the amount of Params: found: %d, expected %d", ErrInvalidContractText, len(matches), len(sc.Template.TemplateAttributes))
	}

	sc.Params = make(map[string]string, len(matches))
	for idx, match := range matches {
		sc.Params[sc.Template.TemplateAttributes[idx]] = string(match)
	}
	return nil
}

// Verify verifies the params with the template
func (sc Contract) Verify() error {
	var (
		err                      error
		ok                       bool
		validFrom, validTo       *time.Time
		validFromStr, validToStr string
	)

	if validFromStr, ok = sc.Params[ValidFromAttr]; !ok {
		return fmt.Errorf("%w: value for ["+ValidFromAttr+"] is missing", ErrInvalidContractText)
	}

	validFrom, err = parseTime(validFromStr, sc.Template.Language)
	if err != nil {
		return fmt.Errorf("%w: unable to parse ["+ValidFromAttr+"]: %s", ErrInvalidContractText, err)

	}

	if validToStr, ok = sc.Params[ValidToAttr]; !ok {
		return fmt.Errorf("%w: value for ["+ValidToAttr+"] is missing", ErrInvalidContractText)
	}

	validTo, err = parseTime(validToStr, sc.Template.Language)
	if err != nil {
		return fmt.Errorf("%w: unable to parse ["+ValidToAttr+"]: %s", ErrInvalidContractText, err)
	}

	// All parsed, check time range
	if validFrom.After(*validTo) {
		return fmt.Errorf("%w: ["+ValidFromAttr+"] must be after ["+ValidToAttr+"]", ErrInvalidContractText)
	}

	amsterdamLocation, _ := time.LoadLocation("Europe/Amsterdam")
	now := NowFunc()
	logrus.Debugf("checking timeframe: now %v, validFrom: %v, validTo: %v", now, *validFrom, *validTo)

	if now.In(amsterdamLocation).Before(*validFrom) {
		return fmt.Errorf("contract is not yet valid. now: %s, validFrom: %s", now, validFrom)

	}
	if now.In(amsterdamLocation).After(*validTo) {
		return fmt.Errorf("contract is expired since: %s", validTo)
	}

	return nil
}

const amsterdamTimeZone = "Europe/Amsterdam"

// parseTime parses the given timeStr in context of the Europe/Amsterdam time zone and uses the given language.
// Note that currently only the language "NL" is supported.
func parseTime(timeStr string, _ Language) (*time.Time, error) {
	contractIssuerTimezone, _ := time.LoadLocation(amsterdamTimeZone)
	// TODO: add support for other languages
	parsedTime, err := monday.ParseInLocation(timeLayout, timeStr, contractIssuerTimezone, monday.LocaleNlNL)
	if err != nil {
		return nil, fmt.Errorf("invalid time string [%v]: %w", timeStr, err)
	}
	return &parsedTime, nil
}
