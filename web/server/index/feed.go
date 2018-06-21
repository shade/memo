package index

import (
	"github.com/jchavannes/jgo/jerr"
	"github.com/jchavannes/jgo/web"
	"github.com/memocash/memo/app/auth"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/profile"
	"github.com/memocash/memo/app/res"
	"net/http"
)

var feedRoute = web.Route{
	Pattern: res.UrlFeed,
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
		r.Helper["Key"] = key
		setFeed(r, key.PkHash, user.Id)
		r.RenderTemplate(res.TmplFeed)
	},
}

func setFeed(r *web.Response, selfPkHash []byte, userId uint) error {
	offset := r.Request.GetUrlParameterInt("offset")
	posts, err := profile.GetPostsFeed(selfPkHash, uint(offset))
	if err != nil {
		return jerr.Get("error getting posts for hashes", err)
	}
	err = profile.AttachParentToPosts(posts)
	if err != nil {
		return jerr.Get("error attaching parent to posts", err)
	}
	err = profile.AttachLikesToPosts(posts)
	if err != nil {
		return jerr.Get("error attaching likes to posts", err)
	}
	err = profile.AttachProfilePicsToPosts(posts)
	if err != nil {
		return jerr.Get("error attaching profile pics to posts", err)
	}
	err = profile.AttachPollsToPosts(posts)
	if err != nil {
		return jerr.Get("error attaching polls to posts", err)
	}
	r.Helper["PostCount"] = len(posts)
	err = profile.SetShowMediaForPosts(posts, userId)
	if err != nil {
		return jerr.Get("error setting show media for posts", err)
	}
	r.Helper["Posts"] = posts
	r.Helper["Offset"] = offset

	var prevOffset int
	if offset > 25 {
		prevOffset = offset - 25
	}
	page := offset/25 + 1
	r.Helper["Page"] = page
	r.Helper["PrevOffset"] = prevOffset
	r.Helper["NextOffset"] = offset + 25
	return nil
}
