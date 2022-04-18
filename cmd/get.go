/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get your r34",
	Long:  `Gets you some high quality (Not guaranteed) rule34 images`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, _ := cmd.Flags().GetString("tags")
		limit, _ := cmd.Flags().GetString("limit")
		page, _ := cmd.Flags().GetString("page")

		get(tags, limit, page)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().String("tags", "", "The tags to search for. (Put a + instead of a space between each tag)")
	getCmd.Flags().String("limit", "1", "How many images you want. (Max 50, Min 1)")
	getCmd.Flags().String("page", "1", "Which page to search in.")
}

type Posts struct {
	XMLName xml.Name `xml:"posts"`
	Text    string   `xml:",chardata"`
	Count   string   `xml:"count,attr"`
	Offset  string   `xml:"offset,attr"`
	Post    []struct {
		Text          string `xml:",chardata"`
		Height        string `xml:"height,attr"`
		Score         string `xml:"score,attr"`
		FileURL       string `xml:"file_url,attr"`
		ParentID      string `xml:"parent_id,attr"`
		SampleURL     string `xml:"sample_url,attr"`
		SampleWidth   string `xml:"sample_width,attr"`
		SampleHeight  string `xml:"sample_height,attr"`
		PreviewURL    string `xml:"preview_url,attr"`
		Rating        string `xml:"rating,attr"`
		Tags          string `xml:"tags,attr"`
		ID            string `xml:"id,attr"`
		Width         string `xml:"width,attr"`
		Change        string `xml:"change,attr"`
		Md5           string `xml:"md5,attr"`
		CreatorID     string `xml:"creator_id,attr"`
		HasChildren   string `xml:"has_children,attr"`
		CreatedAt     string `xml:"created_at,attr"`
		Status        string `xml:"status,attr"`
		Source        string `xml:"source,attr"`
		HasNotes      string `xml:"has_notes,attr"`
		HasComments   string `xml:"has_comments,attr"`
		PreviewWidth  string `xml:"preview_width,attr"`
		PreviewHeight string `xml:"preview_height,attr"`
	} `xml:"post"`
}

func get(tags string, limit string, page string) {
	timeBef := time.Now().UnixMilli()
	var result Posts
	var url string
	intTest, _ := strconv.ParseInt(limit, 10, 0)
	if intTest > 50 || intTest < 1 {
		fmt.Println("Limit either too high or too low. Try: r34-dl help get")
		return
	}
	if tags != "" {
		url = "https://api.rule34.xxx/index.php?page=dapi&s=post&q=index&pid=" + page + "&limit=" + limit + "&tags=" + tags
	} else {
		url = "https://api.rule34.xxx/index.php?page=dapi&s=post&q=index&pid=" + page + "&limit=" + limit
	}
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	xml.NewDecoder(res.Body).Decode(&result)
	filenum := 0
	for i := 0; i < len(result.Post); i++ {
		fileurl := result.Post[i].FileURL
		fileres, err := http.Get(fileurl)
		if err != nil {
			panic(err)
		}
		data, err := io.ReadAll(fileres.Body)
		if err != nil {
			panic(err)
		}
		if fileurl[len(fileurl)-4:] == ".jpg" || fileurl[len(fileurl)-5:] == ".jpeg" {
			os.WriteFile(result.Post[i].ID+".jpg", data, os.ModePerm)
		} else if fileurl[len(fileurl)-4:] == ".png" {
			os.WriteFile(result.Post[i].ID+".png", data, os.ModePerm)
		} else if fileurl[len(fileurl)-4:] == ".mp4" {
			os.WriteFile(result.Post[i].ID+".mp4", data, os.ModePerm)
		} else {
			fileres.Body.Close()
			continue
		}
		filenum++
		fileres.Body.Close()
	}
	if filenum == 0 {
		fmt.Println("No posts with the specified tags could be found.")
	} else {
		fmt.Printf("Sucessfully downloaded %d files in %dms.", filenum, time.Now().UnixMilli()-timeBef)
		fmt.Printf("\nTags: %s", tags)
	}
}
