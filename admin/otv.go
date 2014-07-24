package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
)

const OTV_CATEGORY_URL = "http://api.otv.co.th/api/index.php/CategoryList/index/15/1.0/2.0.1.2"
const OTV_SHOW_LIST_URL = "http://api.otv.co.th/api/index.php/Lists/index/15/1.0/2.0.1.2/%s/0/50"

type Otv struct {
	Db *sql.DB
}

type OtvCategory struct {
	Items []*OtvCategoryItem
}

type OtvCategoryItem struct {
	ID      string `json:"id"`
	ApiName string `json:"api_name"`
	NameTh  string `json:"name_th"`
	NameEn  string `json:"name_en"`
}

type OtvShowList struct {
	Items []*OtvShowListItem `json:"items"`
}

type OtvShowListItem struct {
	ContentSeasonId int    `json:"content_season_id"`
	NameTh          string `json:"name_th"`
	NameEn          string `json:"name_en"`
	ModifiedDate    string `json:"modified_date"`
	Thumbnail       string `json:"thumbnail"`
}

func (o *Otv) CheckOtvExisting() []*OtvShowListItem {
	var shows []*OtvShowListItem
	c := o.getOtvCategory()
	for _, cat := range c.Items {
		fmt.Println("#####", cat.ID, cat.ApiName, cat.NameEn, "#####")
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
		fmt.Println("#####", cat.ID, cat.ApiName, cat.NameEn, "#####")
		if cat.ID != "5" {
			s := o.getOtvShowList(cat.ID)
			for _, show := range s.Items {
				rowAffected, err := o.updateModifiedDate(show)
				if err != nil {
					panic(err)
				}
				shows = append(shows, show)
				if rowAffected == 0 {
					fmt.Println("Break Update")
					break
				}
			}
		}
	}
	return shows
}

func (o *Otv) checkExisting(show *OtvShowListItem) bool {
	var title string
	err := o.Db.QueryRow("SELECT program_title from tv_program WHERE otv_id = ?", show.ContentSeasonId).Scan(&title)
	if err != nil {
		fmt.Println("##### Not Found #####")
		fmt.Println(show.ContentSeasonId, show.NameTh)
		fmt.Println("ModifiedDate", show.ModifiedDate)
		fmt.Println("Thumbanil", show.Thumbnail)
		return false
	}
	return true
}

func (o *Otv) updateModifiedDate(show *OtvShowListItem) (int64, error) {
	fmt.Println("##### Update #####")
	fmt.Println(show.ContentSeasonId, show.NameTh)
	fmt.Println("ModifiedDate", show.ModifiedDate)

	result, err := o.Db.Exec("UPDATE tv_program SET update_date = ? WHERE otv_id = ?", show.ModifiedDate, show.ContentSeasonId)
	if err != nil {
		panic(err)
	}
	return result.RowsAffected()
}

func (o *Otv) getOtvCategory() OtvCategory {
	resp, err := http.Get(OTV_CATEGORY_URL)
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
		panic(err)
	}
	return c
}

func (o *Otv) getOtvShowList(catId string) OtvShowList {
	apiUrl := fmt.Sprintf(OTV_SHOW_LIST_URL, catId)
	resp, err := http.Get(apiUrl)
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
		panic(err)
	}
	return s
}
