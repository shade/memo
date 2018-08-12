package format

import (
	"regexp"
	"strings"
)

func AddYoutubeVideos(msg string) string {
	var re = regexp.MustCompile(`(http[s]?://youtu\.be/)([A-Za-z0-9_\-\?=]+)`)
	msg = re.ReplaceAllString(msg, `<div class="video-wrapper"><div class="video-container"><iframe frameborder="0" src="https://www.youtube.com/embed/$2"></iframe></div></div>`)
	re = regexp.MustCompile(`(http[s]?://y2u\.be/)([A-Za-z0-9_\-\?=]+)`)
	msg = re.ReplaceAllString(msg, `<div class="video-wrapper"><div class="video-container"><iframe frameborder="0" src="https://www.youtube.com/embed/$2"></iframe></div></div>`)
	re = regexp.MustCompile(`(http[s]?://(www\.)?youtube\.com/watch\?v=)([A-Za-z0-9_\-\?=]+)`)
	msg = re.ReplaceAllString(msg, `<div class="video-wrapper"><div class="video-container"><iframe frameborder="0" src="https://www.youtube.com/embed/$3"></iframe></div></div>`)
	return msg
}

func AddImgurImages(msg string) string {
	// Album link
	if strings.Contains(msg, "imgur.com/a/") || strings.Contains(msg, "imgur.com/gallery/") {
		return msg
	}
	containsRex := regexp.MustCompile(`\.jpg|\.jpeg|\.png|\.gif|\.gifv`)
	if strings.Contains(msg, ".mp4") {
		var re = regexp.MustCompile(`(http[s]?://([a-z]+\.)?imgur\.com/)([^\s]*)`)
		msg = re.ReplaceAllString(msg, `<div class="video-wrapper"><div class="video-container"><video controls><source src="https://i.imgur.com/$3" type="video/mp4"></video></div></div>`)
	} else if !containsRex.MatchString(msg) {
		var re = regexp.MustCompile(`(http[s]?://([a-z]+\.)?imgur\.com/)([^\s]*)`)
		msg = re.ReplaceAllString(msg, `<a href="https://i.imgur.com/$3.jpg" target="_blank" class="imgur"><img src="https://i.imgur.com/$3.jpg"/></a>`)
	} else {
		var re = regexp.MustCompile(`(http[s]?://([a-z]+\.)?imgur\.com/)([^\s]*)`)
		msg = re.ReplaceAllString(msg, `<a href="https://i.imgur.com/$3" target="_blank" class="imgur"><img src="https://i.imgur.com/$3"/></a>`)
	}
	return msg
}

func AddGiphyImages(msg string) string {
	if strings.Contains(msg, "giphy.com/gifs/") {
		var re = regexp.MustCompile(`(http[s]?://([a-z]+\.)?giphy.com/gifs/[a-z-]*-([A-Za-z0-9]+))`)
		msg = re.ReplaceAllString(msg, `<a href="https://i.giphy.com/$3.gif" target="_blank" class="imgur"><img src="https://i.giphy.com/$3.gif"/></a>`)
	} else {
		var re = regexp.MustCompile(`(http[s]?://([a-z]+\.)?giphy\.com/)([^\s]*)`)
		msg = re.ReplaceAllString(msg, `<a href="https://i.giphy.com/$3" target="_blank" class="imgur"><img src="https://i.giphy.com/$3"/></a>`)
	}
	return msg
}

func AddTwitterImages(msg string) string {
	var re = regexp.MustCompile(`(http[s]?://pbs.twimg.com/media/([A-Za-z0-9_-]+)[A-Za-z0-9?&=.;]*)`)
	msg = re.ReplaceAllString(msg, `<a href="https://pbs.twimg.com/media/$2.jpg" target="_blank" class="imgur"><img src="https://pbs.twimg.com/media/$2.jpg"/></a>`)
	return msg
}

func AddRedditImages(msg string) string {
	var re = regexp.MustCompile(`(http[s]?://i.redd.it/([^\s]*))`)
	msg = re.ReplaceAllString(msg, `<a href="https://i.redd.it/$2" target="_blank" class="imgur"><img src="https://i.redd.it/$2"/></a>`)
	return msg
}

func AddTweets(msg string) string {
	var re = regexp.MustCompile(`(http[s]?://([a-z]+\.)?twitter.com/([A-Za-z0-9_-]+)/status/([0-9]+)[A-Za-z0-9?=.]*)`)
	msg = re.ReplaceAllString(msg, `<blockquote class="twitter-tweet" data-cards="hidden" data-lang="en"><a href="https://twitter.com/$3/status/$4"></a></blockquote>`)
	return msg
}

func AddLinks(msg string) string {
	// Explanation: https://github.com/jchavannes/memo/pull/57
	var re = regexp.MustCompile(`(^|[\s(])(http[s]?://[^\s]*[^.?!,)\s])`)
	s := re.ReplaceAllString(msg, `$1<a href="$2" target="_blank">$2</a>`)
	return strings.Replace(s, "\n", "<br/>", -1)
}

func RemoveTrailingWhiteSpace(msg string) string {
	var re = regexp.MustCompile(`(<br\/>\s*)+$`)
	return re.ReplaceAllString(msg, ``)
}
