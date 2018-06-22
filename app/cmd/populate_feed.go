package cmd

import (
	"fmt"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/db"
	"github.com/memocash/memo/app/obj/feed_event"
	"github.com/memocash/memo/app/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"strings"
)

var feedTypes = []string{
	"likes",
	"posts",
	"poll-votes",
	"follows",
	"topic-follows",
	"set-names",
	"set-profiles",
	"set-profile-pics",
}

var populateFeedCmd = &cobra.Command{
	Use:  "populate-feed",
	Long: "populate-feed [" + strings.Join(feedTypes, ", ") + "]",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 || (!util.StringInSlice(args[0], feedTypes) && args[0] != "all") {
			return errors.New(fmt.Sprintf("invalid feed type, must be one of: %s", strings.Join(append(feedTypes, "all"), ", ")))
		}
		if args[0] == "all" {
			for _, feedType := range feedTypes {
				err := populateFeed(feedType)
				if err != nil {
					jerr.Get("error populating feed", err).Print()
					return nil
				}
			}
		} else {
			err := populateFeed(args[0])
			if err != nil {
				jerr.Get("error populating feed", err).Print()
				return nil
			}
		}
		return nil
	},
}

func populateFeed(feedType string) error {
	var printStatus = func(offset uint, itemsAdded int) {
		fmt.Printf("offset: %6d, items-added: %6d\n", offset, itemsAdded)
	}
	fmt.Printf("Populating feed for type: %s\n", feedType)
	var offset uint
	var itemsAdded int
ItemLoop:
	for ; offset < 100000; offset += 25 {
		switch feedType {
		case "likes":
			likes, err := db.GetLikes(offset)
			if err != nil {
				return jerr.Get("error getting likes", err)
			}
			for _, like := range likes {
				err := feed_event.AddLike(like)
				if err != nil {
					return jerr.Get("error adding like feed item", err)
				}
				itemsAdded++
				if itemsAdded%1000 == 0 {
					printStatus(offset, itemsAdded)
				}
			}
			if len(likes) != 25 {
				break ItemLoop
			}
		case "posts":
			posts, err := db.GetPosts(offset)
			if err != nil {
				return jerr.Get("error getting posts", err)
			}
			for _, post := range posts {
				err := feed_event.AddPost(post)
				if err != nil {
					return jerr.Get("error adding post feed item", err)
				}
				itemsAdded++
				if itemsAdded%1000 == 0 {
					printStatus(offset, itemsAdded)
				}
			}
			if len(posts) != 25 {
				break ItemLoop
			}
		case "poll-votes":
			pollVotes, err := db.GetPollVotes(offset)
			if err != nil {
				return jerr.Get("error getting poll votes", err)
			}
			for _, pollVote := range pollVotes {
				err := feed_event.AddPollVote(pollVote)
				if err != nil {
					return jerr.Get("error adding poll vote feed item", err)
				}
				itemsAdded++
				if itemsAdded%1000 == 0 {
					printStatus(offset, itemsAdded)
				}
			}
			if len(pollVotes) != 25 {
				break ItemLoop
			}
		case "follows":
			follows, err := db.GetAllFollows(offset)
			if err != nil {
				return jerr.Get("error getting follows", err)
			}
			for _, post := range follows {
				err := feed_event.AddFollow(post)
				if err != nil {
					return jerr.Get("error adding follow feed item", err)
				}
				itemsAdded++
				if itemsAdded%1000 == 0 {
					printStatus(offset, itemsAdded)
				}
			}
			if len(follows) != 25 {
				break ItemLoop
			}
		case "topic-follows":
			topicFollows, err := db.GetTopicFollows(offset)
			if err != nil {
				return jerr.Get("error getting topic follows", err)
			}
			for _, topicFollow := range topicFollows {
				err := feed_event.AddTopicFollow(topicFollow)
				if err != nil {
					return jerr.Get("error adding topic follow feed item", err)
				}
				itemsAdded++
				if itemsAdded%1000 == 0 {
					printStatus(offset, itemsAdded)
				}
			}
			if len(topicFollows) != 25 {
				break ItemLoop
			}
		case "set-names":
			setNames, err := db.GetSetNames(offset)
			if err != nil {
				return jerr.Get("error getting set names", err)
			}
			for _, setName := range setNames {
				err := feed_event.AddSetName(setName)
				if err != nil {
					return jerr.Get("error adding set name feed item", err)
				}
				itemsAdded++
				if itemsAdded%1000 == 0 {
					printStatus(offset, itemsAdded)
				}
			}
			if len(setNames) != 25 {
				break ItemLoop
			}
		case "set-profiles":
			setProfiles, err := db.GetSetProfiles(offset)
			if err != nil {
				return jerr.Get("error getting set profiles", err)
			}
			for _, setProfile := range setProfiles {
				err := feed_event.AddSetProfile(setProfile)
				if err != nil {
					return jerr.Get("error adding set profile feed item", err)
				}
				itemsAdded++
				if itemsAdded%1000 == 0 {
					printStatus(offset, itemsAdded)
				}
			}
			if len(setProfiles) != 25 {
				break ItemLoop
			}
		case "set-profile-pics":
			setProfilePics, err := db.GetSetProfilePics(offset)
			if err != nil {
				return jerr.Get("error getting set profile pics", err)
			}
			for _, setProfilePic := range setProfilePics {
				err := feed_event.AddSetProfilePic(setProfilePic)
				if err != nil {
					return jerr.Get("error adding set profile pic feed item", err)
				}
				itemsAdded++
				if itemsAdded%1000 == 0 {
					printStatus(offset, itemsAdded)
				}
			}
			if len(setProfilePics) != 25 {
				break ItemLoop
			}
		}
	}
	printStatus(offset, itemsAdded)
	fmt.Println("All done")
	return nil
}
