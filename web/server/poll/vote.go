package poll

import (
	"github.com/jchavannes/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/bitcoin/transaction/build"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/html-parser"
	"github.com/memocash/memo/app/mutex"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var voteSubmitRoute = web.Route{
	Pattern:     res.UrlPollVoteSubmit,
	NeedsLogin:  true,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
		txHashString := r.Request.GetFormValue("txHash")
		txHash, err := chainhash.NewHashFromStr(txHashString)
		if err != nil {
			r.Error(jerr.Get("error getting transaction hash", err), http.StatusInternalServerError)
			return
		}
		option := html_parser.EscapeWithEmojis(r.Request.GetFormValue("option"))
		message := r.Request.GetFormValue("message")

		memoPollOption, err := db.GetMemoPollOptionByOption(txHash.CloneBytes(), option)
		if err != nil {
			r.Error(jerr.Get("error getting memo_post", err), http.StatusInternalServerError)
			return
		}

		password := r.Request.GetFormValue("password")
		user, err := auth.GetSessionUser(r.Session.CookieId)
		if err != nil {
			r.Error(jerr.Get("error getting session user", err), http.StatusInternalServerError)
			return
		}
		key, err := db.GetKeyForUser(user.Id)
		if err != nil {
			r.Error(jerr.Get("error getting key for user", err), http.StatusInternalServerError)
			return
		}

		privateKey, err := key.GetPrivateKey(password)
		if err != nil {
			r.Error(jerr.Get("error getting private key", err), http.StatusUnauthorized)
			return
		}

		var tip = int64(r.Request.GetFormValueInt("tip"))

		pkHash := privateKey.GetPublicKey().GetAddress().GetScriptAddress()
		mutex.Lock(pkHash)

		tx, err := build.Vote(memoPollOption.TxHash, message, tip, privateKey)
		if err != nil {
			mutex.Unlock(pkHash)
			r.Error(jerr.Get("error building vote tx", err), http.StatusInternalServerError)
			return
		}

		transaction.GetTxInfo(tx).Print()
		transaction.QueueTx(tx)
		r.Write(tx.MsgTx.TxHash().String())
	},
}
