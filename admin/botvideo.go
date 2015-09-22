package admin

import (
	"database/sql"
	"log"
	"math"
	"strconv"
)

type BotVideo struct {
	Db *sql.DB
}

func NewBotVideo(db *sql.DB) *BotVideo {
	b := new(BotVideo)
	b.Db = db
	return b
}

type FormSearchBotUser struct {
	Username     string
	Q            string
	Status       int
	Page         int32
	IsOrderTitle bool
}

type BotUser struct {
	Username    string
	Description string
	IsSelected  bool
}

type BotStatus struct {
	ID         int32
	Name       string
	IsSelected bool
}

type BotVideos struct {
	Videos      []*BotVideoRow
	CountRow    int32
	CurrentPage int32
	MaxPage     int32
}

type BotVideoRow struct {
	ID          int32
	Username    string
	Description string
	ProgramID   int64
	UserType    string
	Title       string
	VideoID     string
	VideoType   string
	Published   string
	Status      int
}

func (b *BotVideo) getBotStatuses(id int) []*BotStatus {
	botStatuses := []*BotStatus{}
	botStatuses = append(botStatuses, &BotStatus{0, "Waiting", (id == 0)})
	botStatuses = append(botStatuses, &BotStatus{1, "Updated", (id == 1)})
	botStatuses = append(botStatuses, &BotStatus{-1, "Rejected", (id == -1)})
	botStatuses = append(botStatuses, &BotStatus{2, "Suspended", (id == 2)})
	return botStatuses
}

func (b *BotVideo) getBotStatusID(status string) int {
	switch status {
	case "Rejected":
		return -1
	case "Updated":
		return 1
	case "Suspended":
		return 2
	default:
		return 0
	}
}

func (b *BotVideo) getBotUsers(selectUsername string) []*BotUser {
	botUsers := []*BotUser{}
	rows, err := b.Db.Query("SELECT DISTINCT v.username, u.description from tv_bot_videos v LEFT JOIN tv_youtube_users u ON (v.username = u.username) WHERE u.description != '' ORDER BY description")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			username    string
			description string
		)
		if err := rows.Scan(&username, &description); err != nil {
			log.Fatal(err)
		}
		isSelected := selectUsername == username
		botUser := &BotUser{username, description, isSelected}
		botUsers = append(botUsers, botUser)
	}

	return botUsers
}

func (b *BotVideo) getBotVideos(f *FormSearchBotUser) *BotVideos {
	var countRow int32
	botVideos := []*BotVideoRow{}

	var rows *sql.Rows
	var err error

	if f.Username == "all" || f.Username == "" {
		err = b.Db.QueryRow("SELECT count(id) from tv_bot_videos WHERE status = ? AND title LIKE ?", f.Status, "%"+f.Q+"%").Scan(&countRow)
	} else {
		err = b.Db.QueryRow("SELECT count(id) from tv_bot_videos WHERE status = ? AND username = ? AND title LIKE ?", f.Status, f.Username, "%"+f.Q+"%").Scan(&countRow)
	}

	extendedOrder := ""
	if f.IsOrderTitle {
		extendedOrder = " v.title ASC,"
	}

	if f.Username == "all" || f.Username == "" {
		rows, err = b.Db.Query("SELECT v.id, v.username, u.description, u.program_id, u.user_type, v.title, video_id, video_type, DATE_ADD(published, INTERVAL 7 HOUR), status from tv_bot_videos v LEFT JOIN tv_youtube_users u ON (v.username = u.username) WHERE status = ? AND title LIKE ? ORDER BY  u.official DESC, v.username ASC,"+extendedOrder+" published DESC LIMIT ?, ?", f.Status, "%"+f.Q+"%", (f.Page * limitRow), limitRow)
	} else {
		rows, err = b.Db.Query("SELECT v.id, v.username, u.description, u.program_id, u.user_type, v.title, video_id, video_type, DATE_ADD(published, INTERVAL 7 HOUR), status from tv_bot_videos v LEFT JOIN tv_youtube_users u ON (v.username = u.username) WHERE status = ? AND v.username = ? AND title LIKE ? ORDER BY"+extendedOrder+" published DESC LIMIT ?, ?", f.Status, f.Username, "%"+f.Q+"%", (f.Page * limitRow), limitRow)
	}

	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			id           int32
			username     string
			sDescription sql.NullString
			description  string
			sProgramID   sql.NullInt64
			programID    int64
			sUserType    sql.NullString
			userType     string
			title        string
			videoID      string
			videoType    string
			published    string
			status       int
		)
		if err := rows.Scan(&id, &username, &sDescription, &sProgramID, &sUserType, &title, &videoID, &videoType, &published, &status); err != nil {
			log.Fatal(err)
		}

		if sDescription.Valid {
			description = sDescription.String
		}

		if sProgramID.Valid {
			programID = sProgramID.Int64
		}

		if sUserType.Valid {
			userType = sUserType.String
		} else {
			userType = "channel"
		}

		botVideo := &BotVideoRow{id, username, description, programID, userType, title, videoID, videoType, published, status}
		botVideos = append(botVideos, botVideo)
	}

	return &BotVideos{botVideos, countRow, f.Page, int32(math.Ceil(float64(countRow / limitRow)))}
	// return botVideos
}

func (b *BotVideo) setBotVideoStatus(id int, status int) {
	_, err := b.Db.Exec("UPDATE tv_bot_videos SET status = ? WHERE id = ?", status, id)
	if err != nil {
		panic(err)
	}
}

func (b *BotVideo) setBotVideosStatus(videoIDs []string, updateStatus string) {
	statusID := b.getBotStatusID(updateStatus)
	for _, videoID := range videoIDs {
		id, _ := strconv.Atoi(videoID)
		b.setBotVideoStatus(id, statusID)
	}
}
