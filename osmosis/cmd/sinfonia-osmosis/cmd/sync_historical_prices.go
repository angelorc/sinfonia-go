package cmd

import (
	"fmt"
	"github.com/angelorc/sinfonia-go/config"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/mongo/modelv2"
	"github.com/angelorc/sinfonia-go/mongo/repository"
	"github.com/angelorc/sinfonia-go/utility"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"time"
)

func GetSyncHistoricalPricesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "historical-prices",
		Short:   "sync atom, osmo and btsg prices from coingecko",
		Example: "sinfonia-osmosis historical-prices",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, err := cmd.Flags().GetString(flagConfig)
			if err != nil {
				return err
			}

			cfg, err := config.NewConfig(cfgPath)
			if err != nil {
				return err
			}

			defaultDB := db.Database{
				DataBaseRefName: "default",
				URL:             cfg.Mongo.Uri,
				DataBaseName:    cfg.Mongo.DbName,
				RetryWrites:     strconv.FormatBool(cfg.Mongo.Retry),
			}
			defaultDB.Init()
			defer defaultDB.Disconnect()

			hpr := repository.NewHistoricalPriceRepository()
			hpr.EnsureIndexes()

			// get from sinfonia start
			startTime := time.Date(2022, 07, 11, 0, 0, 0, 0, time.UTC)
			endTime := time.Now()
			assets := []string{
				"osmosis", "cosmos", "bitsong",
			}

			log.Printf("getting historical prices from %s to %s", startTime.Format("02-01-2006"), endTime.Format("02-01-2006"))

			for startTime.Before(endTime) {
				for _, asset := range assets {
					log.Printf("getting price for %s from coingecko, time: %s", asset, startTime.Format("02-01-2006"))

					// get price
					price, err := utility.GetHistoricalCoinPrice(asset, startTime, "usd")
					if err != nil {
						log.Fatal(err)
					}

					var prices []modelv2.Price
					prices = append(prices, modelv2.Price{Usd: fmt.Sprintf("%f", price)})

					_, err = hpr.Create(&modelv2.HistoricalPriceCreateReq{
						Asset: asset,
						Price: prices,
						Time:  startTime,
					})
					if err != nil {
						return err
					}

					log.Printf("stored price for %s from coingecko, time: %s, price: %2f", asset, startTime.Format("02-01-2006"), price)

					time.Sleep(2 * time.Second)
				}

				startTime = startTime.Add(24 * time.Hour)
			}

			return nil
		},
	}

	addConfigFlag(cmd)

	return cmd
}
