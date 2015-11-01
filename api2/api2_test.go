package api2

import (
	"encoding/json"
	"fmt"
	_ "github.com/code-mobi/tvthailand-api/Godeps/_workspace/src/github.com/go-sql-driver/mysql"
	"github.com/code-mobi/tvthailand-api/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCategories_ServeHTTP(t *testing.T) {
	db, err := utils.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	h := &Api2Handler{Db: db}

	server := httptest.NewServer(h)
	defer server.Close()

	url := fmt.Sprintf("%s/api2/category", server.URL)
	fmt.Println(url)
	resp, err := http.Get(url)
	if resp.StatusCode != 200 {
		t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
	}

	actual, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	var c Categories
	err = json.Unmarshal(actual, &c)
	if err != nil {
		t.Fatal(err)
	}

	if len(c.Categories) == 0 {
		t.Error("Expected Categories > 0")
	}

	for _, cat := range c.Categories {
		if cat.ID == "" {
			t.Error("Expected ID Not Empty")
		}
		if cat.Title == "" {
			t.Error("Expected Title Not Empty")
		}
		if cat.Thumbnail == "" {
			t.Error("Expected Thumbnail Not Empty")
		}
	}
}
