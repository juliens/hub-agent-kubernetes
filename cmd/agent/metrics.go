package main

import (
	"context"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/traefik/neo-agent/pkg/metrics"
	"github.com/urfave/cli/v2"
)

const (
	flagScrapeName = "scrape-name"
	flagScrapeKind = "scrape-kind"
	flagScrapeIP   = "scrape-ip"
)

func metricsFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     flagScrapeName,
			Usage:    "The name of the ingress controller",
			EnvVars:  []string{"SCRAPE_NAME"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     flagScrapeKind,
			Usage:    "The ingress controller type to scrape (nginx or traefik)",
			EnvVars:  []string{"SCRAPE_KIND"},
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     flagScrapeIP,
			Usage:    "An IP of a ingress controller to scrape",
			EnvVars:  []string{"SCRAPE_IP"},
			Required: true,
		},
	}
}

func runMetrics(ctx context.Context, cliCtx *cli.Context) error {
	retryableClient := retryablehttp.NewClient()
	retryableClient.RetryWaitMin = time.Second
	retryableClient.RetryWaitMax = 10 * time.Second
	retryableClient.RetryMax = 4

	httpClient := retryableClient.StandardClient()

	client, err := metrics.NewClient(httpClient, cliCtx.String("platform-url"), cliCtx.String("token"))
	if err != nil {
		return err
	}

	store := metrics.NewStore()

	scraper := metrics.NewScraper(httpClient)

	mgr := metrics.NewManager(client, store, scraper)

	return mgr.Run(ctx, time.Minute, cliCtx.String(flagScrapeKind), cliCtx.String(flagScrapeName), cliCtx.StringSlice(flagScrapeIP))
}
