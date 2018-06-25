package memo

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/bitcoin/transaction/build"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/mutex"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var setNameRoute = web.Route{
	Pattern:    res.UrlMemoSetName,
	NeedsLogin: true,
	Handler: func(r *web.Response) {
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
		hasSpendableTxOut, err := db.HasSpendable(key.PkHash)
		if err != nil {
			r.Error(jerr.Get("error getting spendable tx out", err), http.StatusInternalServerError)
			return
		}
		if ! hasSpendableTxOut {
			r.SetRedirect(res.UrlNeedFunds)
			return
		}
		r.Render()
	},
}

var setNameSubmitRoute = web.Route{
	Pattern:     res.UrlMemoSetNameSubmit,
	NeedsLogin:  true,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
		name := r.Request.GetFormValue("name")
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

		pkHash := privateKey.GetPublicKey().GetAddress().GetScriptAddress()
		mutex.Lock(pkHash)

		tx, err := build.SetName(name, privateKey)
		if err != nil {
			var statusCode = http.StatusInternalServerError
			if build.IsNotEnoughValueError(err) {
				statusCode = http.StatusPaymentRequired
			}
			mutex.Unlock(pkHash)
			r.Error(jerr.Get("error building set name tx", err), statusCode)
			return
		}

		transaction.GetTxInfo(tx).Print()
		transaction.QueueTx(tx.MsgTx)
		r.Write(tx.MsgTx.TxHash().String())
	},
}
