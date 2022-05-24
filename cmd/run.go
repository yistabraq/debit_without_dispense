/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"time"

	"github.com/istabraq/debit_without_dispense/internal"
	"github.com/istabraq/debit_without_dispense/pkg/config"
	"github.com/istabraq/debit_without_dispense/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	dir         *config.Dir
	database    *config.Database
	sugarLogger *zap.SugaredLogger
	query       string
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long:  ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !config.IsConfigFileExist() {
			err := config.WriteDefaultConfig()
			if err != nil {
				return err
			}
		}
		cfg, db, err := config.ReadConfigFile()
		if err != nil {
			return err
		}
		err = cfg.IsValidConfig()
		if err != nil {
			return err
		}
		dir = cfg
		dir.CheckFolder()
		database = db
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		currentTime := time.Now()
		path := dir.Log + fmt.Sprintf("/log_%s.log", currentTime.Format("2006_01_02"))
		sugarLogger = logger.InitLogger(path)
		defer sugarLogger.Sync()
		internal.Run(*database, query)
	},
}

func init() {
	q := `
	SELECT 
    TO_CHAR(OO.OPER_DATE,'DD/MM/YYYY') DATE_OPER,
    ISS.ACCOUNT_NUMBER,
    OO.OPER_AMOUNT,
    OO.OPER_TYPE,
    AA.EXTERNAL_ORIG_ID UTRNNO ,
    OO.ORIGINATOR_REFNUM REF_NUM,
    OO.ID OPER_ID,
    OO.STATUS,
    OO.STATUS_REASON,
    OO.IS_REVERSAL,
    ISS.INST_ID ISS_INST,
    ACQ.INST_ID ACQ_INST,
    cast(sysdate - OO.OPER_DATE as int) NB
FROM (SELECT * FROM MAIN.OPR_OPERATION@DBLK_TOBO WHERE STATUS IN ('OPST0402','OPST0800')) OO,
     (SELECT * FROM MAIN.OPR_PARTICIPANT@DBLK_TOBO WHERE PARTICIPANT_TYPE = 'PRTYISS') ISS,
     (SELECT * FROM MAIN.OPR_PARTICIPANT@DBLK_TOBO WHERE PARTICIPANT_TYPE = 'PRTYACQ' AND INST_ID IN (5012,9002)) ACQ,
     MAIN.AUT_AUTH@DBLK_TOBO AA 
WHERE OO.ID = ISS.OPER_ID AND 
      OO.ID = ACQ.OPER_ID AND
      OO.ID = AA.ID AND
      OO.MATCH_ID is null AND
      OO.OPER_DATE >= sysdate - 37 AND
      cast(sysdate - OO.OPER_DATE as int) >= 30;
	`
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&query, "query", "q", q, "Query")

}
