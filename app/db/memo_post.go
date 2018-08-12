package db

import (
	"bytes"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/jchavannes/jgo/jerr"
	"github.com/memocash/memo/app/bitcoin/script"
	"github.com/memocash/memo/app/bitcoin/wallet"
	"github.com/memocash/memo/app/db/view"
	"html"
	"net/url"
	"time"
)

type MemoPost struct {
	Id           uint        `gorm:"primary_key"`
	TxHash       []byte      `gorm:"unique;size:50"`
	ParentHash   []byte
	PkHash       []byte      `gorm:"index:pk_hash"`
	PkScript     []byte      `gorm:"size:500"`
	Address      string
	ParentTxHash []byte      `gorm:"index:parent_tx_hash"`
	Parent       *MemoPost
	RootTxHash   []byte      `gorm:"index:root_tx_hash"`
	Replies      []*MemoPost `gorm:"foreignkey:ParentTxHash"`
	Topic        string      `gorm:"index:tag;size:500"`
	Message      string      `gorm:"size:500"`
	IsPoll       bool
	IsVote       bool
	BlockId      uint
	Block        *Block
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (m *MemoPost) Save() error {
	result := save(&m)
	if result.Error != nil {
		return jerr.Get("error saving memo test", result.Error)
	}
	return nil
}

func (m MemoPost) GetTransactionHashString() string {
	hash, err := chainhash.NewHash(m.TxHash)
	if err != nil {
		jerr.Get("error getting chainhash from memo post", err).Print()
		return ""
	}
	return hash.String()
}

func (m MemoPost) GetParentTransactionHashString() string {
	hash, err := chainhash.NewHash(m.ParentTxHash)
	if err != nil {
		jerr.Get("error getting chainhash from memo post", err).Print()
		return ""
	}
	return hash.String()
}

func (m MemoPost) GetRootTransactionHashString() string {
	hash, err := chainhash.NewHash(m.RootTxHash)
	if err != nil {
		jerr.Get("error getting chainhash from memo post", err).Print()
		return ""
	}
	return hash.String()
}

func (m MemoPost) GetAddressString() string {
	return m.GetAddress().GetEncoded()
}

func (m MemoPost) GetAddress() wallet.Address {
	return wallet.GetAddressFromPkHash(m.PkHash)
}

func (m MemoPost) GetScriptString() string {
	return html.EscapeString(script.GetScriptString(m.PkScript))
}

func (m MemoPost) GetMessage() string {
	return m.Message
}

func (m MemoPost) GetUrlEncodedTopic() string {
	return url.QueryEscape(m.Topic)
}

func (m MemoPost) GetTimeString() string {
	if m.BlockId != 0 {
		if m.Block != nil {
			return m.Block.Timestamp.Format("2006-01-02 15:04:05")
		} else {
			return "Unknown"
		}
	}
	return "Unconfirmed"
}

func GetMemoPost(txHash []byte) (*MemoPost, error) {
	var memoPost MemoPost
	err := findPreloadColumns([]string{
		BlockTable,
	}, &memoPost, MemoPost{
		TxHash: txHash,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo post", err)
	}
	return &memoPost, nil
}

func GetPostsByTxHashes(txHashes [][]byte) ([]*MemoPost, error) {
	var memoPosts []*MemoPost
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	result := db.
		Preload(BlockTable).
		Where("tx_hash IN (?)", txHashes).
		Find(&memoPosts)
	if result.Error != nil {
		return nil, jerr.Get("error getting memo posts", result.Error)
	}
	return memoPosts, nil
}

func GetMemoPostById(id uint) (*MemoPost, error) {
	var memoPost MemoPost
	err := find(&memoPost, MemoPost{
		Id: id,
	})
	if err != nil {
		return nil, jerr.Get("error getting memo post", err)
	}
	return &memoPost, nil
}

func GetPostReplyCount(txHash []byte) (uint, error) {
	cnt, err := count(MemoPost{
		ParentTxHash: txHash,
	})
	if err != nil {
		return 0, jerr.Get("error running count query", err)
	}
	return cnt, nil
}

type TxHashCount struct {
	TxHash []byte
	Count  uint
}

func GetPostReplyCounts(txHashes [][]byte) ([]TxHashCount, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	query := db.
		Table("memo_posts").
		Select("parent_tx_hash, COUNT(*) AS count").
		Where("parent_tx_hash IN (?)", txHashes).
		Group("parent_tx_hash")
	rows, err := query.Rows()
	if err != nil {
		return nil, jerr.Get("error running query", err)
	}
	defer rows.Close()
	var txHashCounts []TxHashCount
	for rows.Next() {
		var txHash []byte
		var count uint
		err := rows.Scan(&txHash, &count)
		if err != nil {
			return nil, jerr.Get("error scanning rows", err)
		}
		txHashCounts = append(txHashCounts, TxHashCount{
			TxHash: txHash,
			Count:  count,
		})
	}
	return txHashCounts, nil
}

func GetPostReplies(txHash []byte, offset uint) ([]*MemoPost, error) {
	var posts []*MemoPost
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}

	query := db.
		Table("memo_posts").
		Preload(BlockTable).
		Select("memo_posts.*, COUNT(DISTINCT memo_likes.pk_hash) AS count").
		Joins("LEFT OUTER JOIN blocks ON (memo_posts.block_id = blocks.id)").
		Joins("LEFT OUTER JOIN memo_likes ON (memo_posts.tx_hash = memo_likes.like_tx_hash)").
		Group("memo_posts.id").
		Order("count DESC, memo_posts.id DESC").
		Limit(25).
		Offset(offset)

	result := query.Find(&posts, MemoPost{
		ParentTxHash: txHash,
	})
	if result.Error != nil {
		return nil, jerr.Get("error finding post replies", result.Error)
	}
	return posts, nil
}

func GetPostsFeedForPkHash(pkHash []byte, offset uint) ([]*MemoPost, error) {
	var memoPosts []*MemoPost
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	joinSelect := "SELECT " +
		"	follow_pk_hash " +
		"FROM memo_follows " +
		"JOIN (" +
		"	SELECT MAX(id) AS id" +
		"	FROM memo_follows" +
		"	WHERE pk_hash = ?" +
		"	GROUP BY pk_hash, follow_pk_hash" +
		") sq ON (sq.id = memo_follows.id) " +
		"WHERE unfollow = 0"
	result := db.
		Limit(25).
		Offset(offset).
		Preload(BlockTable).
		Joins("JOIN ("+joinSelect+") fsq ON (memo_posts.pk_hash = fsq.follow_pk_hash)", pkHash).
		Order("id DESC").
		Find(&memoPosts)
	if result.Error != nil {
		return nil, jerr.Get("error getting memo posts", result.Error)
	}
	return memoPosts, nil
}

func GetPostsForPkHash(pkHash []byte, offset uint) ([]*MemoPost, error) {
	if len(pkHash) == 0 {
		return nil, nil
	}
	var memoPosts []*MemoPost
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	query := db.
		Preload(BlockTable).
		Order("id DESC").
		Limit(25).
		Offset(offset)
	result := query.Find(&memoPosts, &MemoPost{
		PkHash: pkHash,
	})
	if result.Error != nil {
		return nil, jerr.Get("error getting memo posts", result.Error)
	}
	return memoPosts, nil
}

func GetRecentPosts(offset uint) ([]*MemoPost, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	db = db.Preload(BlockTable)
	var memoPosts []*MemoPost
	result := db.
		Limit(25).
		Offset(offset).
		Order("id DESC").
		Find(&memoPosts)
	if result.Error != nil {
		return nil, jerr.Get("error running query", result.Error)
	}
	return memoPosts, nil
}

func GetPosts(offset uint) ([]*MemoPost, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	var memoPosts []*MemoPost
	result := db.
		Preload(BlockTable).
		Limit(25).
		Offset(offset).
		Order("id ASC").
		Find(&memoPosts)
	if result.Error != nil {
		return nil, jerr.Get("error running query", result.Error)
	}
	return memoPosts, nil
}

func GetRecentReplyPosts(offset uint) ([]*MemoPost, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	var memoPosts []*MemoPost
	result := db.
		Limit(25).
		Offset(offset).
		Order("id DESC").
		Where("parent_tx_hash IS NOT NULL").
		Find(&memoPosts)
	if result.Error != nil {
		return nil, jerr.Get("error running query", result.Error)
	}
	return memoPosts, nil
}

func GetRecentPostsForTopic(topic string, lastPostId uint) ([]*MemoPost, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	var memoPosts []*MemoPost
	result := db.
		Where("id > ?", lastPostId).
		Order("id ASC").
		Find(&memoPosts, MemoPost{
		Topic: topic,
	})
	if result.Error != nil {
		return nil, jerr.Get("error running recent topic post query", result.Error)
	}
	return memoPosts, nil
}

func GetTopPosts(offset uint, timeStart time.Time, timeEnd time.Time) ([]*MemoPost, error) {
	topLikeTxHashes, err := GetRecentTopLikedTxHashes(offset, timeStart, timeEnd)
	if err != nil {
		return nil, jerr.Get("error getting top liked tx hashes", err)
	}
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	db = db.Preload(BlockTable)
	var memoPosts []*MemoPost
	result := db.Where("tx_hash IN (?)", topLikeTxHashes).Find(&memoPosts)
	if result.Error != nil {
		return nil, jerr.Get("error running query", result.Error)
	}
	var sortedPosts []*MemoPost
	for _, txHash := range topLikeTxHashes {
		for _, memoPost := range memoPosts {
			if bytes.Equal(memoPost.TxHash, txHash) {
				sortedPosts = append(sortedPosts, memoPost)
			}
		}
	}
	return sortedPosts, nil
}

const (
	RankCountBoost int     = 60
	RankGravity    float32 = 2
)

func GetRankedPosts(offset uint) ([]*MemoPost, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	var coalescedTimestamp = "IF(COALESCE(blocks.timestamp, memo_posts.created_at) < memo_posts.created_at, blocks.timestamp, memo_posts.created_at)"
	var scoreQuery = fmt.Sprintf("((COUNT(DISTINCT memo_likes.pk_hash)-1)*%d)/POW(TIMESTAMPDIFF(MINUTE, "+coalescedTimestamp+", NOW())+2,%0.2f)", RankCountBoost, RankGravity)

	var memoPosts []*MemoPost
	result := db.
		Joins("LEFT OUTER JOIN memo_likes ON (memo_posts.tx_hash = memo_likes.like_tx_hash)").
		Joins("LEFT OUTER JOIN blocks ON (memo_posts.block_id = blocks.id)").
		Where(coalescedTimestamp + " > DATE_SUB(NOW(), INTERVAL 3 DAY)").
		Group("memo_posts.tx_hash").
		Order(scoreQuery + " DESC").
		Limit(25).
		Offset(offset).
		Preload(BlockTable).
		Find(&memoPosts)
	if result.Error != nil {
		return nil, jerr.Get("error running query", result.Error)
	}
	return memoPosts, nil
}

func GetPollsPosts(offset uint) ([]*MemoPost, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	var coalescedTimestamp = "IF(COALESCE(blocks.timestamp, memo_posts.created_at) < memo_posts.created_at, blocks.timestamp, memo_posts.created_at)"
	var scoreQuery = fmt.Sprintf("((COUNT(DISTINCT memo_poll_votes.tx_hash)-1)*%d)/POW(TIMESTAMPDIFF(MINUTE, "+coalescedTimestamp+", NOW())+2,%0.2f)", RankCountBoost, RankGravity)

	var memoPosts []*MemoPost
	result := db.
		Joins("LEFT OUTER JOIN memo_poll_options ON (memo_posts.tx_hash = memo_poll_options.poll_tx_hash) ").
		Joins("LEFT OUTER JOIN memo_poll_votes ON (memo_poll_options.tx_hash = memo_poll_votes.option_tx_hash)").
		Joins("LEFT OUTER JOIN blocks ON (memo_posts.block_id = blocks.id)").
		Where("is_poll = 1").
		Group("memo_posts.tx_hash").
		Order(scoreQuery + " DESC").
		Limit(25).
		Offset(offset).
		Preload(BlockTable).
		Find(&memoPosts)
	if result.Error != nil {
		return nil, jerr.Get("error running query", result.Error)
	}
	return memoPosts, nil
}

func GetPersonalizedTopPosts(selfPkHash []byte, offset uint, timeStart time.Time, timeEnd time.Time) ([]*MemoPost, error) {
	topLikeTxHashes, err := GetPersonalizedRecentTopLikedTxHashes(selfPkHash, offset, timeStart, timeEnd)
	if err != nil {
		return nil, jerr.Get("error getting top liked tx hashes", err)
	}
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	db = db.Preload(BlockTable)
	var memoPosts []*MemoPost
	result := db.Where("tx_hash IN (?)", topLikeTxHashes).Find(&memoPosts)
	if result.Error != nil {
		return nil, jerr.Get("error running query", result.Error)
	}
	var sortedPosts []*MemoPost
	for _, txHash := range topLikeTxHashes {
		for _, memoPost := range memoPosts {
			if bytes.Equal(memoPost.TxHash, txHash) {
				sortedPosts = append(sortedPosts, memoPost)
			}
		}
	}
	return sortedPosts, nil
}

func GetCountMemoPosts() (uint, uint, uint, uint, error) {
	db, err := getDb()
	if err != nil {
		return 0, 0, 0, 0, jerr.Get("error getting db", err)
	}
	selectStmt := "" +
		"SUM(IF(IFNULL(topic, '') = '' AND IFNULL(parent_tx_hash, b'') = b'', 1, 0)) AS non_topic_posts, " +
		"SUM(IF(IFNULL(is_vote, 0) = 1, 1, 0)) AS vote_posts, " +
		"SUM(IF(IFNULL(topic, '') = '', 0, 1)) AS topic_posts, " +
		"SUM(IF(IFNULL(parent_tx_hash, b'') = b'', 0, 1)) AS reply_posts"
	query := db.
		Model(&MemoPost{}).
		Where("IFNULL(is_poll, 0) = 0").
		Select(selectStmt)
	row := query.Row()
	var postCount uint
	var votePostCount uint
	var topicPostCount uint
	var replyPostCount uint
	err = row.Scan(&postCount, &votePostCount, &topicPostCount, &replyPostCount)
	if err != nil {
		return 0, 0, 0, 0, jerr.Get("error getting distinct topics", err)
	}
	return postCount, votePostCount, topicPostCount, replyPostCount, nil
}

func GetTopicInfoFromPosts(topicNames ...string) ([]*view.Topic, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	joinSelect := "LEFT JOIN (" +
		"	SELECT MAX(id) AS id" +
		"	FROM memo_topic_follows"
	if len(topicNames) > 0 {
		joinSelect +=
		" WHERE topic IN (?)"
	}
	joinSelect +=
		"	GROUP BY pk_hash, topic" +
		") sq ON (sq.id = memo_topic_follows.id) "
	query := db.
		Table("memo_posts").
		Select("" +
		"memo_posts.topic, " +
		"CAST(MAX(IF(COALESCE(blocks.timestamp, memo_posts.created_at) < memo_posts.created_at, blocks.timestamp, memo_posts.created_at)) AS DATETIME) AS max_time, " +
		"COUNT(DISTINCT memo_posts.id) AS post_count, " +
		"COUNT(DISTINCT case memo_topic_follows.unfollow when 0 then memo_topic_follows.id else null end) AS follower_count").
		Joins("LEFT JOIN memo_topic_follows ON (memo_posts.topic = memo_topic_follows.topic)").
		Joins(joinSelect, topicNames).
		Joins("LEFT JOIN blocks ON (memo_posts.block_id = blocks.id)").
		Group("memo_posts.topic")
	if len(topicNames) > 0 {
		query = query.Where("memo_posts.topic IN (?)", topicNames)
	} else {
		query = query.Where("IFNULL(memo_posts.topic, '') != ''")
	}
	rows, err := query.Rows()
	if err != nil {
		return nil, jerr.Get("error getting distinct topics", err)
	}
	defer rows.Close()
	var topics []*view.Topic
	for rows.Next() {
		var topic view.Topic
		err := rows.Scan(&topic.Name, &topic.RecentTime, &topic.CountPosts, &topic.CountFollows)
		if err != nil {
			return nil, jerr.Get("error scanning row with topic", err)
		}
		topics = append(topics, &topic)
	}
	return topics, nil
}

func AttachUnreadToTopics(topics []*view.Topic, userPkHash []byte) error {
	var topicNames []string
	for _, topic := range topics {
		topicNames = append(topicNames, topic.Name)
	}
	lastTopicPostIds, err := GetLastTopicPostIds(userPkHash, topicNames)
	if err != nil {
		return jerr.Get("error getting last topic post ids", err)
	}
	db, err := getDb()
	if err != nil {
		return jerr.Get("error getting db", err)
	}
	query := db.
		Table("memo_posts").
		Select("MAX(id) AS maxId, topic").
		Where("topic IN (?)", topicNames).
		Group("topic")
	rows, err := query.Rows()
	if err != nil {
		return jerr.Get("error getting max topic post ids", err)
	}
	defer rows.Close()
	for rows.Next() {
		var maxId uint
		var topicName string
		err := rows.Scan(&maxId, &topicName)
		if err != nil {
			return jerr.Get("error scanning row for topic max id", err)
		}
		var lastPostId uint
		for _, lastTopicPostId := range lastTopicPostIds {
			if lastTopicPostId.Topic == topicName {
				lastPostId = lastTopicPostId.LastPostId
			}
		}
		for _, topic := range topics {
			if topic.Name == topicName {
				topic.UnreadPosts = lastPostId < maxId
			}
		}
	}
	return nil
}

func GetPostsForTopic(topic string, offset uint) ([]*MemoPost, error) {
	if len(topic) == 0 {
		return nil, jerr.New("empty topic")
	}
	var memoPosts []*MemoPost
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	query := db.
		Preload(BlockTable).
		Order("id DESC").
		Limit(26).
		Offset(offset)
	result := query.Find(&memoPosts, &MemoPost{
		Topic: topic,
	})
	if result.Error != nil {
		return nil, jerr.Get("error getting memo posts", result.Error)
	}
	for i, j := 0, len(memoPosts)-1; i < j; i, j = i+1, j-1 {
		memoPosts[i], memoPosts[j] = memoPosts[j], memoPosts[i]
	}
	return memoPosts, nil
}

func GetOlderPostsForTopic(topic string, firstPostId uint) ([]*MemoPost, error) {
	if len(topic) == 0 {
		return nil, jerr.New("empty topic")
	}
	var memoPosts []*MemoPost
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	query := db.
		Preload(BlockTable).
		Where("id < ?", firstPostId).
		Order("id DESC").
		Limit(26)
	result := query.Find(&memoPosts, &MemoPost{
		Topic: topic,
	})
	if result.Error != nil {
		return nil, jerr.Get("error getting memo posts", result.Error)
	}
	for i, j := 0, len(memoPosts)-1; i < j; i, j = i+1, j-1 {
		memoPosts[i], memoPosts[j] = memoPosts[j], memoPosts[i]
	}
	return memoPosts, nil
}

func GetThreads(offset uint, topic string) ([]*view.Thread, error) {
	db, err := getDb()
	if err != nil {
		return nil, jerr.Get("error getting db", err)
	}
	query := db.
		Table("memo_posts").
		Select("" +
		"topic_posts.topic AS topic, " +
		"topic_posts.message AS message, " +
		"root_tx_hash, " +
		"COUNT(memo_posts.`id`) AS num_replies, " +
		"CAST(MAX(IF(COALESCE(blocks.timestamp, memo_posts.created_at) < memo_posts.created_at, blocks.timestamp, memo_posts.created_at)) AS DATETIME) AS recent_reply").
		Joins("LEFT JOIN blocks ON (memo_posts.block_id = blocks.id)").
		Group("memo_posts.root_tx_hash").
		Order("recent_reply DESC").
		Limit(25).
		Offset(offset)
	if topic != "" {
		joinSelect := "JOIN (" +
			"	SELECT tx_hash, topic, message" +
			"	FROM memo_posts" +
			"   WHERE topic = ?" +
			") topic_posts ON (memo_posts.root_tx_hash = topic_posts.tx_hash)"
		query = query.Joins(joinSelect, topic)
	} else {
		joinSelect := "JOIN (" +
			"	SELECT tx_hash, topic, message" +
			"	FROM memo_posts" +
			") topic_posts ON (memo_posts.root_tx_hash = topic_posts.tx_hash)"
		query = query.Joins(joinSelect)
	}
	var threads []*view.Thread
	result := query.Scan(&threads)
	if result.Error != nil {
		return nil, jerr.Get("error getting threads", result.Error)
	}
	return threads, nil
}
