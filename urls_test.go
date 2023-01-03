package fishfish_test

import (
	"testing"
	"time"

	"github.com/existentiality/fishfish-go"
)

func TestGetURLs(t *testing.T) {
	urls, err := client.GetURLs(fishfish.CategoryPhishing, false)

	mustPanic(err)

	t.Logf("got %d phishing urls", len(*urls))
}

func TestGetURLsFull(t *testing.T) {
	urls, err := client.GetURLsFull()

	mustPanic(err)

	t.Logf("got %d urls with full data", len(*urls))
}

func TestGetURL(t *testing.T) {
	// There are currently no URLs in the databse, skip
	t.SkipNow()
	url, err := client.GetURL("https://fishfish.gg/api.html", true)

	mustPanic(err)

	t.Logf("got url %s (category %s)", url.URL, url.Category)
}

func TestAddURL(t *testing.T) {
	if !client.HasPermission(fishfish.APIPermissionURLs) {
		t.Skip("missing permission")
	}

	added, err := client.AddURL("https://fishfish.gg/api.html", fishfish.CategorySafe)

	mustPanic(err)

	t.Logf("added url %s with category %s", added.URL, added.Category)
}

func TestUpdateURL(t *testing.T) {
	if !client.HasPermission(fishfish.APIPermissionURLs) {
		t.Skip("missing permission")
	}

	updated, err := client.UpdateURL("https://fishfish.gg/api.html", fishfish.CategorySafe)

	mustPanic(err)

	t.Logf("updated url %s", updated.URL)
}

func TestUpdateURLMetadata(t *testing.T) {
	if !client.HasPermission(fishfish.APIPermissionURLs) {
		t.Skip("missing permission")
	}

	updated, err := client.UpdateURLMetadata("https://fishfish.gg/api.html", fishfish.URLMetadata{
		Target: "fishfish",
		Active: time.Now(),
	})

	mustPanic(err)

	t.Logf("updated metadata for https://fishfish.gg/api.html (last active %s, target %s)", updated.Active, updated.Target)
}

func TestDeleteURL(t *testing.T) {
	if !client.HasPermission(fishfish.APIPermissionURLs) {
		t.Skip("missing permission")
	}

	err := client.DeleteURL("https://fishfish.gg/api.html")

	mustPanic(err)

	t.Logf("successfully deleted url")
}