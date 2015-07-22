package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/facebookgo/httpcontrol"
	_ "github.com/go-sql-driver/mysql"
)

const (
	OtvCategoryURL = "http://api.otv.co.th/api/index.php/v202/Category/index/15/1.0/2.0.2"
	OtvShowListURL = "http://api.otv.co.th/api/index.php/v202/Lists/index/15/1.0/2.0.2/%s/0/50"
	DateFMT        = "2006-01-02 15:04:05"
)

type Otv struct {
	Db *sql.DB
}

type OtvCategory struct {
	Items []*OtvCategoryItem
}

type OtvCategoryItem struct {
	ID      string `json:"id"`
	APIName string `json:"api_name"`
	NameTh  string `json:"name_th"`
	NameEn  string `json:"name_en"`
}

type OtvShowList struct {
	Items []*OtvShowListItem `json:"items"`
}

type OtvShowListItem struct {
	ContentSeasonID string `json:"content_season_id"`
	NameTh          string `json:"name_th"`
	NameEn          string `json:"name_en"`
	ContentType     string `json:"content_type"`
	ModifiedDate    string `json:"modified_date"`
	Thumbnail       string `json:"thumbnail"`
}

type ProcessType struct {
	Text    string
	Value   string
	Checked bool
}

func OtvProcessOption() []ProcessType {
	options := []ProcessType{
		ProcessType{"OTV Update Modified Date", "modified", true},
		ProcessType{"OTV Existing Show", "existing", false},
		ProcessType{"Find Embed", "findembed", false},
	}
	return options
}

func (o *Otv) CheckOtvExisting() []*OtvShowListItem {
	var shows []*OtvShowListItem
	c := o.getOtvCategory()
	for _, cat := range c.Items {
		fmt.Println("#####", cat.ID, cat.APIName, cat.NameEn, "#####")
		if cat.ID != "5" {
			s := o.getOtvShowList(cat.ID)
			for _, show := range s.Items {
				if !o.checkExisting(show) {
					shows = append(shows, show)
				}
			}
		}
	}
	return shows
}

func (o *Otv) UpdateModified() []*OtvShowListItem {
	var shows []*OtvShowListItem
	c := o.getOtvCategory()
	for _, cat := range c.Items {
		fmt.Println("#####", cat.ID, cat.APIName, cat.NameEn, "#####")
		if cat.ID != "5" {
			s := o.getOtvShowList(cat.ID)
			for _, show := range s.Items {
				rowAffected, err := o.updateModifiedDate(show)
				if err != nil {
					panic(err)
				}
				if rowAffected == 0 {
					fmt.Println("Not Update")
				} else {
					shows = append(shows, show)
				}
			}
		}
	}
	return shows
}

func (o *Otv) updateModifiedDate(show *OtvShowListItem) (int64, error) {
	fmt.Println(show.ContentSeasonID, show.NameTh)
	modifiedDate, errT := time.Parse(DateFMT, show.ModifiedDate)
	if errT != nil {
		fmt.Println(errT)
	}
	var (
		title   string
		strDate string
	)
	err := o.Db.QueryRow("SELECT program_title, update_date from tv_program WHERE otv_id = ?", show.ContentSeasonID).Scan(&title, &strDate)
	if err != nil {
		fmt.Println("####### Program Not Found ####### ")
		return 0, nil
	}
	updateDate, _ := time.Parse(DateFMT, strDate)
	if modifiedDate.After(updateDate) {
		fmt.Println("ModifiedDate", modifiedDate, "After UpdateDate", updateDate)
		result, err := o.Db.Exec("UPDATE tv_program SET update_date = ? WHERE otv_id = ?", show.ModifiedDate, show.ContentSeasonID)
		if err != nil {
			panic(err)
		}
		return result.RowsAffected()
	}
	return 0, nil
}

func (o *Otv) FindEmbed() []*OtvShowListItem {
	var shows []*OtvShowListItem
	c := o.getOtvCategory()
	for _, cat := range c.Items {
		fmt.Println("#####", cat.ID, cat.APIName, cat.NameEn, "#####")
		if cat.ID != "5" {
			s := o.getOtvShowList(cat.ID)
			for _, show := range s.Items {
				if show.ContentType == "embed" {
					rowAffected, err := o.updateEmbedCh7(show)
					if err != nil {
						panic(err)
					}
					if rowAffected == 0 {
						fmt.Println("No Update")
					}
					fmt.Println("#####", show.ContentSeasonID, show.NameTh, show.ContentType, "#####")
					shows = append(shows, show)
				}
			}
		}
	}
	return shows
}

func (o *Otv) checkExisting(show *OtvShowListItem) bool {
	var title string
	err := o.Db.QueryRow("SELECT program_title from tv_program WHERE otv_id = ?", show.ContentSeasonID).Scan(&title)
	if err != nil {
		fmt.Println("##### Not Found #####")
		fmt.Println(show.ContentSeasonID, show.NameTh)
		fmt.Println("ModifiedDate", show.ModifiedDate)
		fmt.Println("Thumbanil", show.Thumbnail)
		return false
	}
	return true
}

func (o *Otv) updateEmbedCh7(show *OtvShowListItem) (int64, error) {
	fmt.Println("##### Update Embed Ch7 #####")
	fmt.Println(show.ContentSeasonID, show.NameTh)

	result, err := o.Db.Exec("UPDATE tv_program SET otv_api_name = ? WHERE otv_id = ?", "Ch7", show.ContentSeasonID)
	if err != nil {
		panic(err)
	}
	return result.RowsAffected()
}

func (o *Otv) getOtvCategory() OtvCategory {
	client := &http.Client{
		Transport: &httpcontrol.Transport{
			RequestTimeout: time.Minute,
			MaxTries:       3,
		},
	}
	resp, err := client.Get(OtvCategoryURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var c OtvCategory
	err = json.Unmarshal(body, &c)
	if err != nil {
		fmt.Println("JSON Parser Error : ", OtvCategoryURL)
		panic(err)
	}
	return c
}

func (o *Otv) getOtvShowList(catID string) OtvShowList {
	apiURL := fmt.Sprintf(OtvShowListURL, catID)
	client := &http.Client{
		Transport: &httpcontrol.Transport{
			RequestTimeout: time.Minute,
			MaxTries:       3,
		},
	}
	resp, err := client.Get(apiURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var s OtvShowList
	err = json.Unmarshal(body, &s)
	if err != nil {
		fmt.Println("JSON Parser Error : ", apiURL)
		panic(err)
	}
	return s
}
