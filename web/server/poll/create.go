package poll

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/bitcoin/memo"
	"github.com/memocash/memo/app/bitcoin/transaction"
	"github.com/memocash/memo/app/bitcoin/transaction/build"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/mutex"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var createRoute = web.Route{
	Pattern:    res.UrlPollCreate,
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		r.Render()
	},
}

var createSubmitRoute = web.Route{
	Pattern:     res.UrlPollCreateSubmit,
	NeedsLogin:  true,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
		pollType := r.Request.GetFormValue("pollType")
		question := r.Request.GetFormValue("question")
		options := r.Request.GetFormValueSlice("options")
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

		memoTxns, err := build.Poll(memo.PollType(pollType), question, options, privateKey)
		if err != nil {
			mutex.Unlock(pkHash)
			r.Error(jerr.Get("error building memo poll tx", err), http.StatusInternalServerError)
			return
		}

		for _, memoTx := range memoTxns {
			transaction.GetTxInfo(memoTx).Print()
			transaction.QueueTx(memoTx.MsgTx)
		}

		r.Write(memoTxns[len(memoTxns)-1].MsgTx.TxHash().String())
	},
}
