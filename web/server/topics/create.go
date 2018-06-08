package topics

import (
	"fmt"
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

var createRoute = web.Route{
	Pattern:    res.UrlTopicsCreate,
	NeedsLogin: true,
	Handler: func(r *web.Response) {
		preHandler(r)
		r.Render()
	},
}

var createSubmitRoute = web.Route{
	Pattern:     res.UrlTopicsCreateSubmit,
	NeedsLogin:  true,
	CsrfProtect: true,
	Handler: func(r *web.Response) {
		topicName := r.Request.GetFormValue("topic")
		message := r.Request.GetFormValue("message")
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

		tx, err := build.TopicMessage(topicName, message, privateKey)
		if err != nil {
			mutex.Unlock(pkHash)
			r.Error(jerr.Get("error building topic message tx", err), http.StatusInternalServerError)
			return
		}

		fmt.Println(transaction.GetTxInfo(tx))
		transaction.QueueTx(tx)
		r.Write(tx.TxHash().String())
	},
}
