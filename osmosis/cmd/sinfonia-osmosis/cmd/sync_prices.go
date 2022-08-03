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

func GetSyncPricesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "prices",
		Short:   "sync atom, osmo and btsg current prices from coingecko",
		Example: "sinfonia-osmosis prices",
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

			assets := []string{
				"osmosis", "cosmos", "bitsong",
			}

			for _, asset := range assets {
				log.Printf("getting price for %s from coingecko", asset)

				// get price
				price, err := utility.GetCoinPrice(asset, "usd")
				if err != nil {
					log.Fatal(err)
				}

				var prices []modelv2.Price
				prices = append(prices, modelv2.Price{Usd: fmt.Sprintf("%f", price)})

				_, err = hpr.Create(&modelv2.HistoricalPriceCreateReq{
					Asset: asset,
					Price: prices,
					Time:  time.Now(),
				})
				if err != nil {
					return err
				}

				log.Printf("stored price for %s from coingecko, price: %2f", asset, price)

				time.Sleep(2 * time.Second)
			}

			return nil
		},
	}

	addConfigFlag(cmd)

	return cmd
}
