package web

import (
	"fmt"

	"github.com/playwright-community/playwright-go"

	"github.com/scipunch/myfeed/parser"
)

type Parser struct {
	pw      *playwright.Playwright
	browser playwright.Browser
}

func New() (Parser, error) {
	var p Parser
	err := playwright.Install()
	if err != nil {
		return p, err
	}
	pw, err := playwright.Run()
	if err != nil {
		return p, fmt.Errorf("could not start playwright: %w", err)
	}
	browser, err := pw.Chromium.Launch()
	if err != nil {
		return p, fmt.Errorf("could not launch browser: %w", err)
	}
	p.pw = pw
	p.browser = browser
	return p, nil
}

func (p Parser) Close() error {
	if err := p.browser.Close(); err != nil {
		return err
	}
	return p.pw.Stop()
}

type Response struct {
	HTML string
}

func (r Response) String() string {
	return r.HTML
}

func (p Parser) Parse(uri string) (parser.Response, error) {
	var resp Response
	page, err := p.browser.NewPage()
	if err != nil {
		return resp, fmt.Errorf("could not create page: %w", err)
	}
	defer page.Close()
	if _, err = page.Goto(uri); err != nil {
		return resp, fmt.Errorf("could not goto: %w", err)
	}
	html, err := page.Content()
	if err != nil {
		return resp, fmt.Errorf("could not read page's content with %w", err)
	}

	resp.HTML = html
	return resp, nil
}
