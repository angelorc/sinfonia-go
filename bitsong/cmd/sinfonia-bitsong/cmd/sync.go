package cmd

import (
	"fmt"
	bitsong "github.com/angelorc/sinfonia-go/bitsong/chain"
	"github.com/angelorc/sinfonia-go/config"
	"github.com/angelorc/sinfonia-go/mongo/db"
	"github.com/angelorc/sinfonia-go/mongo/model"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"strings"
)

func SyncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "module sync",
	}

	cmd.AddCommand(
		GetSyncFantokenCmd(),
		GetSyncMerkledropCmd(),
		GetSyncMerkledropClaimCmd(),
	)

	return cmd
}

func GetSyncFantokenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "fantokens",
		Short:   "sync fantokens from latest blocks",
		Example: "sinfonia-bitsong sync fantokens --mongo-dbname sinfonia-test",
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

			/**
			 * Connect to db
			 */
			defaultDB := db.Database{
				DataBaseRefName: "default",
				URL:             cfg.Mongo.Uri,
				DataBaseName:    cfg.Mongo.DbName,
				RetryWrites:     strconv.FormatBool(cfg.Mongo.Retry),
			}
			defaultDB.Init()
			defer defaultDB.Disconnect()

			if err := syncFantokens(); err != nil {
				return err
			}

			return nil
		},
	}

	addConfigFlag(cmd)

	return cmd
}

func GetSyncMerkledropCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "merkledrops",
		Short:   "sync merkledrops from latest blocks",
		Example: "sinfonia-bitsong sync merkledrops --mongo-dbname sinfonia-test",
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

			/**
			 * Connect to db
			 */
			defaultDB := db.Database{
				DataBaseRefName: "default",
				URL:             cfg.Mongo.Uri,
				DataBaseName:    cfg.Mongo.DbName,
				RetryWrites:     strconv.FormatBool(cfg.Mongo.Retry),
			}
			defaultDB.Init()
			defer defaultDB.Disconnect()

			client, err := bitsong.NewClient(&cfg.Bitsong)
			if err != nil {
				return fmt.Errorf("failed to get RPC endpoints on chain %s. err: %v", "bitsong", err)
			}

			if err := syncMerkledrops(client); err != nil {
				return err
			}

			return nil
		},
	}

	addConfigFlag(cmd)

	return cmd
}

func syncFantokens() error {
	// get last available height on db
	lastBlock := model.GetLastHeight()

	// get last block synced from account
	sync := new(model.Sync)
	sync.One()

	if sync.ID.IsZero() {
		sync.ID = primitive.NewObjectID()
		sync.Fantokens = int64(0)
	}

	txsLogs, err := model.GetTxsAndLogsByMessageType("/bitsong.fantoken.MsgIssue", sync.Fantokens, lastBlock)
	if err != nil {
		return err
	}

	for _, txLogs := range txsLogs {
		for _, txlog := range txLogs.Tx.Logs {
			for _, evt := range txlog.Events {
				switch evt.Type {
				case "bitsong.fantoken.v1beta1.EventIssue":
					denom := strings.ReplaceAll(evt.Attributes[0].Value, "\"", "")

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

func getAttrs(attrKey string, attrs []model.Attribute) (value string) {
	for _, attr := range attrs {
		if attr.Key == attrKey {
			return attr.Value
		}
	}

	return value
}

func syncMerkledrops(client *bitsong.Client) error {
	// get last available height on db
	lastBlock := model.GetLastHeight()

	// get last block synced
	sync := new(model.Sync)
	sync.One()

	if sync.ID.IsZero() {
		sync.ID = primitive.NewObjectID()
		sync.Fantokens = int64(0)
	}

	txsLogs, err := model.GetTxsAndLogsByMessageType("/bitsong.merkledrop.v1beta1.MsgCreate", sync.Merkledrops, lastBlock)
	if err != nil {
		return err
	}

	for _, txLogs := range txsLogs {
		for _, txlog := range txLogs.Tx.Logs {
			for _, evt := range txlog.Events {
				switch evt.Type {
				case "bitsong.merkledrop.v1beta1.EventCreate":
					idStr := strings.ReplaceAll(getAttrs("merkledrop_id", evt.Attributes), "\"", "")
					merkledropId, _ := strconv.ParseInt(idStr, 10, 64)

					mRes, err := client.QueryMerkledropByID(uint64(merkledropId))
					if err != nil {
						continue
						// return fmt.Errorf("error while fetching merkedropID %d, err: %s", merkledropId, err.Error())
					}

					amount := mRes.Merkledrop.Amount.Int64()

					merkledrop := new(model.Merkledrop)
					data := &model.MerkledropCreate{
						ChainID:      &txLogs.ChainID,
						Height:       &txLogs.Height,
						TxID:         &txLogs.TxID,
						MsgIndex:     &txLogs.MsgIndex,
						MerkledropId: &merkledropId,
						Denom:        &mRes.Merkledrop.Denom,
						Amount:       &amount,
						StartHeight:  &mRes.Merkledrop.StartHeight,
						EndHeight:    &mRes.Merkledrop.EndHeight,
						Time:         &txLogs.Time,
					}

					if err := merkledrop.Create(data); err != nil {
						return err
					}
				}
			}
		}
	}

	// TODO: prune expired merkledrop and merkleproofs
	// get current height
	// if current height > merkledrop-end-height
	// then prune

	// update sync with last synced height
	sync.Merkledrops = lastBlock
	if err := sync.Save(); err != nil {
		return err
	}

	fmt.Printf("%d merkledrops synced to block %d ", len(txsLogs), sync.Merkledrops)

	return nil
}

func GetSyncMerkledropClaimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "merkledrops-claim",
		Short:   "sync merkledrops-claim from latest blocks",
		Example: "sinfonia-bitsong sync merkledrops-claim --mongo-dbname sinfonia-test",
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

			/**
			 * Connect to db
			 */
			defaultDB := db.Database{
				DataBaseRefName: "default",
				URL:             cfg.Mongo.Uri,
				DataBaseName:    cfg.Mongo.DbName,
				RetryWrites:     strconv.FormatBool(cfg.Mongo.Retry),
			}
			defaultDB.Init()
			defer defaultDB.Disconnect()

			if err := syncMerkledropClaims(); err != nil {
				return err
			}

			return nil
		},
	}

	addConfigFlag(cmd)

	return cmd
}

func syncMerkledropClaims() error {

	// get last available height on db
	lastBlock := model.GetLastHeight()

	// get last block synced from account
	sync := new(model.Sync)
	err := sync.One()
	if err != nil {
		return err
	}

	if sync.ID.IsZero() {
		sync.ID = primitive.NewObjectID()
		sync.Fantokens = int64(0)
	}

	txsLogs, err := model.GetTxsAndLogsByMessageType("/bitsong.merkledrop.v1beta1.MsgClaim", sync.Fantokens, lastBlock)
	if err != nil {
		return err
	}

	for _, txLogs := range txsLogs {
		for _, txlog := range txLogs.Tx.Logs {
			for _, evt := range txlog.Events {
				switch evt.Type {
				case "bitsong.merkledrop.v1beta1.EventClaim":
					merkledropId, _ := strconv.ParseInt(strings.ReplaceAll(evt.Attributes[1].Value, "\"", ""), 10, 64)
					index, _ := strconv.ParseInt(strings.ReplaceAll(evt.Attributes[2].Value, "\"", ""), 10, 64)
					address := txLogs.Signer

					proof := new(model.MerkledropProof)
					if err := proof.SetClaimed(model.MerkledropProofClaim{
						MerkledropID: merkledropId,
						Address:      address,
						Index:        index,
					}); err != nil {
						return err
					}
				}
			}
		}
	}

	// update sync with last synced height
	sync.MerkledropProofs = lastBlock
	if err := sync.Save(); err != nil {
		return err
	}

	fmt.Printf("%d merkledrop proofs synced to block %d ", len(txsLogs), sync.MerkledropProofs)

	return nil
}
