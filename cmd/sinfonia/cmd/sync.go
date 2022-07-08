package cmd

import (
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/spf13/cobra"
	"strconv"
)

func GetSyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "module sync",
	}

	cmd.AddCommand(
		GetSyncAccountCmd(),
		GetSyncFantokenCmd(),
		GetSyncMerkledropCmd(),
	)

	return cmd
}

func GetSyncAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account",
		Short:   "sync accounts from latest blocks",
		Example: "sinfonia sync account --mongo-dbname sinfonia-test",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			mongoURI, mongoDBName, mongoRetryWrites, err := parseMongoFlags(cmd)
			if err != nil {
				return err
			}

			/**
			 * Connect to db
			 */
			defaultDB := db.Database{
				DataBaseRefName: "default",
				URL:             mongoURI,
				DataBaseName:    mongoDBName,
				RetryWrites:     strconv.FormatBool(mongoRetryWrites),
			}
			defaultDB.Init()
			defer defaultDB.Disconnect()

			if err := model.SyncAccounts(); err != nil {
				return err
			}

			return nil
		},
	}

	addMongoFlags(cmd)

	return cmd
}

func GetSyncFantokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fantoken",
		Short:   "sync fantokens from latest blocks",
		Example: "sinfonia sync fantoken --mongo-dbname sinfonia-test",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			mongoURI, mongoDBName, mongoRetryWrites, err := parseMongoFlags(cmd)
			if err != nil {
				return err
			}

			/**
			 * Connect to db
			 */
			defaultDB := db.Database{
				DataBaseRefName: "default",
				URL:             mongoURI,
				DataBaseName:    mongoDBName,
				RetryWrites:     strconv.FormatBool(mongoRetryWrites),
			}
			defaultDB.Init()
			defer defaultDB.Disconnect()

			if err := model.SyncFantokens(); err != nil {
				return err
			}

			return nil
		},
	}

	addMongoFlags(cmd)

	return cmd
}

func GetSyncMerkledropCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "merkledrop",
		Short:   "sync merkledrops from latest blocks",
		Example: "sinfonia sync merkledrop --mongo-dbname sinfonia-test",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			mongoURI, mongoDBName, mongoRetryWrites, err := parseMongoFlags(cmd)
			if err != nil {
				return err
			}

			/**
			 * Connect to db
			 */
			defaultDB := db.Database{
				DataBaseRefName: "default",
				URL:             mongoURI,
				DataBaseName:    mongoDBName,
				RetryWrites:     strconv.FormatBool(mongoRetryWrites),
			}
			defaultDB.Init()
			defer defaultDB.Disconnect()

			if err := model.SyncMerkledrops(); err != nil {
				return err
			}

			return nil
		},
	}

	addMongoFlags(cmd)

	return cmd
}

/*func syncFantokens() error {
	// get last available height on db
	lastBlock := model.GetLastHeight()

	// get last block synced from account
	sync := new(model.Sync)
	sync.One()

	if sync.ID.IsZero() {
		sync.ID = primitive.NewObjectID()
		sync.Pools = int64(0)
	}

	txsLogs, err := model.GetTxsAndLogsByMessageType("/bitsong.fantoken.MsgIssueFanToken", sync.Fantokens, lastBlock)
	if err != nil {
		return err
	}

	for _, txLogs := range txsLogs {
		for _, txlog := range txLogs.Tx.Logs {
			for _, evt := range txlog.Events {
				switch evt.Type {
				case "issue_fantoken":
					denom := evt.Attributes[0].Value

					fantoken := new(model.Fantoken)
					data := &model.FantokenCreate{
						ChainID:  &txLogs.ChainID,
						Height:   &txLogs.Height,
						TxID:     &txLogs.TxID,
						Denom:    &denom,
						Owner:    &txLogs.Signer,
						IssuedAt: &txLogs.Time,
					}

					if err := fantoken.Create(data); err != nil {
						return err
					}
				}
			}
		}
	}

	// update sync with last synced height
	sync.Fantokens = lastBlock
	if err := sync.Save(); err != nil {
		return err
	}

	fmt.Printf("%d fantokens synced to block %d ", len(txsLogs), sync.Fantokens)

	return nil
}
*/
