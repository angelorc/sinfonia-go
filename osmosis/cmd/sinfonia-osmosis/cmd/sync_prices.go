package cmd

import (
	"github.com/angelorc/sinfonia-go/config"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/mongo/modelv2"
	"github.com/angelorc/sinfonia-go/mongo/repository"
	"github.com/angelorc/sinfonia-go/utility"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
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

			assets := map[string]string{
				"osmosis": "uosmo",
				"cosmos":  "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2",
				"bitsong": "ibc/4E5444C35610CC76FC94E7F7886B93121175C28262DDFDDE6F84E82BF2425452",
			}

			startTime := time.Now().Add(-4 * time.Hour)
			endTime := time.Now()

			log.Printf("getting historical prices from %s to %s", startTime.Format("02-01-2006"), endTime.Format("02-01-2006"))

			for k, v := range assets {
				log.Printf("getting price for %s from coingecko, time: %s", k, startTime.Format("02-01-2006"))

				// get price
				prices, err := utility.GetHistoricalCoinPrice(k, "usd", startTime, endTime)
				if err != nil {
					log.Fatal(err)
				}

				for _, price := range prices {
					_, err = hpr.Create(&modelv2.HistoricalPriceCreateReq{
						Asset: v,
						Price: price[1],
						Time:  time.Unix(int64(price[0]/1000), 0),
					})

					if err != nil {
						if !strings.Contains(err.Error(), "E11000 duplicate key error") {
							return err
						}
					} else {
						log.Printf("stored price for %s from coingecko, time: %s, price: %2f", k, startTime.Format("02-01-2006"), price)
					}
				}

				time.Sleep(5 * time.Second)
			}

			return nil
		},
	}

	addConfigFlag(cmd)

	return cmd
}
