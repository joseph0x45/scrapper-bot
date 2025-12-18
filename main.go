package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/playwright-community/playwright-go"
)

func parseEuro(s string) (float64, error) {
	s = strings.Trim(s, " €")
	// Remove thousands separator (dot)
	s = strings.ReplaceAll(s, ".", "")
	// Replace decimal comma with dot
	s = strings.ReplaceAll(s, ",", ".")
	return strconv.ParseFloat(s, 64)
}

func main() {
	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not start playwright: %v", err)
	}
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
		Channel:  playwright.String("chrome"),
	})
	if err != nil {
		log.Fatalf("could not launch browser: %v", err)
	}
	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not create page: %v", err)
	}
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	if _, err = page.Goto(
		"https://www.bonusveicolielettrici.mase.gov.it/veicolielettriciBeneficiario/#/plafond",
		playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
		},
	); err != nil {
		log.Fatalf("could not goto: %v", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	for {
		select {
		case <-ticker.C:
			if _, err := page.Reload(playwright.PageReloadOptions{
				WaitUntil: playwright.WaitUntilStateNetworkidle,
			}); err != nil {
				log.Println("reload failed:", err)
			}
			value, err := page.Locator("p.mt-3.mb-0.text-muted > strong:nth-of-type(3)").First().TextContent()
			if err != nil {
				log.Println("[ERROR]: Failed to get value: ", err.Error())
				continue
			}
			parsed, err := parseEuro(value)
			if parsed >= 11000 {
				//send message to discord
			}
			fmt.Println(parsed)

		case <-ctx.Done():
			log.Println("shutting down watcher")
			if err = browser.Close(); err != nil {
				log.Fatalf("could not close browser: %v", err)
			}
			if err = pw.Stop(); err != nil {
				log.Fatalf("could not stop Playwright: %v", err)
			}
			return
		}
	}
}
