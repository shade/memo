package cmd

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/jchavannes/btcd/txscript"
	"github.com/jchavannes/btcd/wire"
	"github.com/jchavannes/jgo/jerr"
	"github.com/spf13/cobra"
)

var decodeCmd = &cobra.Command{
	Use:   "decode",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) != 1 {
			return jerr.Newf("wrong number of arguments, expected 1, got %d", len(args))
		}
		txHex := args[0]
		bt, err := hex.DecodeString(txHex)
		if err != nil {
			jerr.Get("error decoding tx hex", err).Print()
			return nil
		}
		msgTx := wire.NewMsgTx(1)
		reader := bytes.NewReader(bt)
		err = msgTx.Deserialize(reader)
		if err != nil {
			jerr.Get("error decerializing tx", err).Print()
			return nil
		}
		fmt.Printf("msgTx: %#v\n", msgTx)
		for _, out := range msgTx.TxOut {
			fmt.Printf("out: %x\n", out.PkScript)
			if bytes.Equal(out.PkScript[0:1], []byte{
				txscript.OP_RETURN,
			}) {
				pushData, err := txscript.PushedData(out.PkScript)
				if err != nil {
					jerr.Get("error parsing push data from message", err).Print()
					return nil
				}
				for i, data := range pushData {
					fmt.Printf("i: %d, string: %s, hex: %x\n", i, data, data)
				}
			}
		}
		return nil
	},
}
