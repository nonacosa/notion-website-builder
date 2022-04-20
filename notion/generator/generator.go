package generator

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"github.com/dstotijn/go-notion"
	"github.com/fyne-io/examples/notion/pkg/tomarkdown"
	"github.com/fyne-io/examples/storage"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Run(config Config, website storage.Website, pg tomarkdown.Progress) error {

	savePath := storage.MdSavePath(website.Name)

	client := notion.NewClient(config.Notion.Secret)
	config.Notion.FilterProp = "Status"
	//config.Notion.FilterValue[0] = "✅"
	//config.Notion.FilterValue[1] = "🖨"
	config.Notion.PublishedValue = "Published"
	q, err := QueryDatabase(client, config.Notion)
	if err != nil {
		return fmt.Errorf("❌ Querying Notion database: %s", err)
	}
	fmt.Println("✔ Querying Notion database: Completed")

	// fetch page children
	changed := 0 // number of article status changed
	var skipPage bool
	for i, page := range q.Results {
		skipPage = false
		fmt.Printf("-- Article [%d/%d] --\n", i+1, len(q.Results))

		storage.Scan(website.PageID, func(items []string, badger *badger.DB) {
			log.Printf("scan item size is : %b ", len(items))
			if items != nil {
				for i := 0; i < len(items); i++ {
					var item storage.WebsiteItem
					err := json.Unmarshal([]byte(items[i]), &item)
					if err != nil {
						fmt.Println(err.Error())
					}
					if item.Id == page.ID {
						localLastModTime := item.LastUpdate
						cloudLastModTime := page.LastEditedTime
						if localLastModTime == cloudLastModTime {
							skipPage = true
						} else {
							fmt.Println("new update fetching ......")
						}
					}
				}
			}
			badger.Close()
		})

		if skipPage {
			continue
		}
		// Get page blocks tree
		blocks, err := queryBlockChildren(client, page.ID)
		if err != nil {
			log.Println("❌ Getting blocks tree:", err)
			continue
		}
		// Generate content to file
		config.Markdown.PostSavePath = savePath
		config.Markdown.ShortcodeSyntax = "hugo"
		config.Markdown.ImageSavePath = storage.ImageSavePath(website.Name)
		// todo domain
		config.Markdown.ImagePublicLink = "/notion_images/posts"
		if err := generate(client, page, blocks, config.Markdown, func(this float64, next float64, msg string) {
			currentProgress := float64(i) / float64(len(q.Results))
			oneLen := float64(1) / float64(len(q.Results))
			nextLen := oneLen * next
			fmt.Println(currentProgress + nextLen)
			pg(currentProgress+nextLen, ((float64(1)/float64(len(q.Results)))+next)/float64(len(q.Results)), fmt.Sprintf(" generate page : %s", msg))
		}); err != nil {
			fmt.Println("❌ Generating blog post:", err)
			continue
		}
		fmt.Println("✔ Generating blog post: Completed")

		// Change status of blog post if desired
		if changeStatus(client, page, config.Notion) {
			changed++
		}

	}

	// Set GITHUB_ACTIONS info variables
	// https://docs.github.com/en/actions/learn-github-actions/workflow-commands-for-github-actions
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		fmt.Printf("::set-output name=articles_published::%d\n", changed)
	}

	return nil
}

func generate(client *notion.Client, page notion.Page, blocks []notion.Block, config Markdown, pg tomarkdown.Progress) error {
	// Create file
	pageName := tomarkdown.ConvertRichText(page.Properties.(notion.DatabasePageProperties)["Name"].Title)
	title := tomarkdown.ConvertRichText(page.Properties.(notion.DatabasePageProperties)["Title"].RichText)
	status := page.Properties.(notion.DatabasePageProperties)["Status"].Select.Name
	var date notion.DateTime
	if page.Properties.(notion.DatabasePageProperties)["Date"].Date != nil {
		date = page.Properties.(notion.DatabasePageProperties)["Date"].Date.Start
	}

	f, err := os.Create(filepath.Join(config.PostSavePath, generateArticleFilename(pageName, page.CreatedTime, config)))
	if err != nil {
		return fmt.Errorf("error create file: %s", err)
	}

	// Generate markdown content to the file
	tm := tomarkdown.New()
	tm.ImgSavePath = filepath.Join(config.ImageSavePath, pageName)
	tm.ImgVisitPath = filepath.Join(config.ImagePublicLink, url.PathEscape(pageName))
	tm.ContentTemplate = config.Template
	// todo edit frontMatter
	tm.WithFrontMatter(page)
	if config.ShortcodeSyntax != "" {
		tm.EnableExtendedSyntax(config.ShortcodeSyntax)
	}

	parentId := strings.ReplaceAll(page.Parent.DatabaseID, "-", "")

	var fm = &storage.FrontMatter{}

	blocks, _ = syncMentionBlocks(client, blocks)

	err = tm.GenerateTo(blocks, f, fm, func(this float64, next float64, msg string) {
		pg(float64(0), next, fmt.Sprintf("%s \n %s", pageName, msg))
	})

	websiteItemMeta := storage.WebsiteItem{
		Id:          page.ID,
		Name:        pageName,
		Title:       title,
		URL:         page.URL,
		Status:      status,
		Date:        date.Format("2022-03-31 23:27:30"),
		FrontMatter: *fm,
	}
	// save last update time
	websiteItemMeta.LastUpdate = page.LastEditedTime
	websiteItemJson, _ := json.Marshal(websiteItemMeta)

	storage.Save(fmt.Sprintf("%s_%s", parentId, page.ID), string(websiteItemJson))

	return err
}

func generateArticleFilename(title string, date time.Time, config Markdown) string {
	escapedTitle := strings.ReplaceAll(
		strings.ToValidUTF8(
			strings.ToLower(title),
			"",
		),
		" ", "_",
	)
	escapedFilename := escapedTitle + ".md"

	if config.GroupByMonth {
		return filepath.Join(date.Format("2006-01-02"), escapedFilename)
	}

	return escapedFilename
}

// todo pref
func syncMentionBlocks(client *notion.Client, blocks []notion.Block) (retBlocks []notion.Block, err error) {

	for _, block := range blocks {
		switch block.Type {
		// todo image
		case notion.BlockTypeParagraph:
			richTexts := block.Paragraph.Text
			for _, rich := range richTexts {
				if rich.Type == "mention" {
					pageId := rich.Mention.Page.ID
					return queryBlockChildren(client, pageId)
				}
			}
		default:
			{
			}
		}
	}
	return nil, nil
}
