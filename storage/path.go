package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

const AppFolderName = ".notion-wb"
const themeFolderName = "theme"
const settingFolderName = "setting"
const mdFolderName = "file"
const notionImageFolderName = "notion_images"
const ThemeDefaultWebsiteFolderName = "exampleSite"
const staticFolderName = "static"
const publicFolderName = "public"
const contentFolderName = "content"

func GetSavePath(websiteName string) string {
	return fullPath(AppFolderName, mdFolderName, websiteName)
}

func HugoSource(theme string) string {
	return fullPath(AppFolderName, themeFolderName, theme, ThemeDefaultWebsiteFolderName)
}

func HugoDest(theme string) string {
	return fullPath(AppFolderName, themeFolderName, theme, ThemeDefaultWebsiteFolderName, publicFolderName)
}

func HugoPost(theme string) string {
	return fullPath(AppFolderName, themeFolderName, theme, ThemeDefaultWebsiteFolderName, contentFolderName, "blog")
}

func HugoStatic(theme string) string {
	return fullPath(AppFolderName, themeFolderName, theme, ThemeDefaultWebsiteFolderName, staticFolderName, notionImageFolderName)
}

func ImageSavePath(websiteName string) string {
	return fullPath(AppFolderName, mdFolderName, websiteName, notionImageFolderName, "posts")
}

func MdSavePath(websiteName string) string {
	return fullPath(AppFolderName, mdFolderName, websiteName)
}

func MdImageSavePath(websiteName string) string {
	return fullPath(AppFolderName, mdFolderName, websiteName, notionImageFolderName)
}

func HugoTheme(theme ...string) string {
	root, err := os.UserHomeDir()
	if err != nil {
		fmt.Errorf("could not get the default root directory to use for user-specific configuration data: %w", err)
	}

	if !filepath.IsAbs(root) {
		fmt.Errorf("storage root must be an absolute path, got %s", root)
	}

	// Fyne does not allow to customize the root for a storage
	// so we'll use the same
	prePath := filepath.Join(root, AppFolderName, themeFolderName)
	sufPath := filepath.Join(theme...)
	savePath := filepath.Join(prePath, sufPath)

	if err := os.MkdirAll(savePath, 0755); err != nil {
		fmt.Errorf("couldn't create content folder: %s", err)
	}
	return savePath
}

func fullPath(paths ...string) string {
	root, err := os.UserHomeDir()
	if err != nil {
		fmt.Errorf("could not get the default root directory to use for user-specific configuration data: %w", err)
	}

	if !filepath.IsAbs(root) {
		fmt.Errorf("storage root must be an absolute path, got %s", root)
	}

	var pathArr []string
	pathArr = append(pathArr, root)
	if len(paths) > 0 {
		for _, path := range paths {
			pathArr = append(pathArr, path)
		}
	}

	savePath := filepath.Join(pathArr...)

	if err := os.MkdirAll(savePath, 0755); err != nil {
		fmt.Errorf("couldn't create content folder: %s", err)
	}

	return savePath
}
