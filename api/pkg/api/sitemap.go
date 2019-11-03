package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/ikeikeikeike/go-sitemap-generator/stm"
	"github.com/joelsaunders/blog-go/api/pkg/config"
	"github.com/joelsaunders/blog-go/api/pkg/repository"
)

type sitemapContent struct {
	slug     string
	modified time.Time
}

func buildSitemap(postDatas []*sitemapContent) *stm.Sitemap {
	sm := stm.NewSitemap(1)
	sm.SetVerbose(false)
	sm.SetDefaultHost("https://www.thebookofjoel.com")
	sm.Create()
	sm.Add(stm.URL{{"loc", "/"}, {"changefreq", "weekly"}, {"lastmod", nil}})
	sm.Add(stm.URL{{"loc", "/team"}, {"changefreq", "yearly"}, {"lastmod", nil}})
	sm.Add(stm.URL{{"loc", "/contact"}, {"changefreq", "yearly"}, {"lastmod", nil}})

	for _, postData := range postDatas {
		sm.Add(stm.URL{
			{"loc", fmt.Sprintf("/%s", postData.slug)},
			{"changefreq", "yearly"},
			{"lastmod", postData.modified.Format("2006-01-02")},
			{"priority", "1"},
		})
	}
	return sm
}

// GetSitemap is a handler that returns an xml response with a sitemap for the blog.
func GetSitemap(postStore repository.PostStore, config *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := postStore.List(r.Context(), map[string][]string{})

		if err != nil {
			log.Println(err)
			render.Render(w, r, ErrDatabase(err))
		}

		var postSlugs []*sitemapContent

		for _, post := range posts {
			content := sitemapContent{post.Slug, post.Modified}
			postSlugs = append(postSlugs, &content)
		}

		sm := buildSitemap(postSlugs)
		w.Write(sm.XMLContent())
		return
	}
}
