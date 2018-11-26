package addontester

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/bitrise-team/bitrise-add-on-testing-kit/addonprovisioner"
	"github.com/bitrise-team/bitrise-add-on-testing-kit/utils"
)

// LoginTesterParams ...
type LoginTesterParams struct {
	AppSlug   string
	BuildSlug string
	Timestamp int64
	WithRetry bool
}

// Login ...
func (t *Tester) Login(params LoginTesterParams, remainingRetries int) error {
	if len(params.AppSlug) == 0 {
		params.AppSlug, _ = utils.RandomHex(8)
	}

	if len(params.BuildSlug) == 0 {
		params.BuildSlug, _ = utils.RandomHex(8)
	}

	if params.Timestamp == 0 {
		params.Timestamp = time.Now().Unix()
	}

	t.logger.Printf("\nLogin details:")
	t.logger.Printf("App slug: %s", params.AppSlug)
	t.logger.Printf("Build slug: %s", params.BuildSlug)
	t.logger.Printf("Timestamp: %d", params.Timestamp)
	t.logger.Printf("Should retry: %v", params.WithRetry)
	if params.WithRetry {
		t.logger.Printf("No. of test: %d.", numberOfTestsWithRetry-remainingRetries)
	}

	status, body, err := t.addonClient.Login(addonprovisioner.LoginRequestParams{
		AppSlug:   params.AppSlug,
		BuildSlug: params.BuildSlug,
		Timestamp: fmt.Sprintf("%d", params.Timestamp),
	})

	if err != nil {
		return fmt.Errorf("Login failed: %s", err)
	}

	t.logger.Printf("\nResponse status: %d", status)
	t.logger.Printf("Response body: %v\n", body)

	if status < 200 || status > 299 {
		return fmt.Errorf("Login request resulted in a non-2xx response")
	}

	t.logger.Println("\nLogin success.")

	r := strings.NewReader(body)
	d := xml.NewDecoder(r)
	d.Strict = true
	d.Entity = xml.HTMLEntity
	var nodes []html.Node
	err = d.Decode(&nodes)

	if err != nil {
		return fmt.Errorf("Login request responded with invalid HTML")
	}

	if params.WithRetry && remainingRetries > 0 {
		return t.Login(params, remainingRetries-1)
	}

	return nil
}
