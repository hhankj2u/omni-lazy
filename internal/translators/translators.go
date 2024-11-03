package translators

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/hhankj2u/translators/pkg/cache"
	"github.com/hhankj2u/translators/pkg/dicts"
	"github.com/hhankj2u/translators/pkg/settings"
)

type Translators struct {
	dbConnections map[string]*sql.DB
}

func NewTranslators() *Translators {
	dbConnections := make(map[string]*sql.DB)
	dbConnections[settings.WEBSTER] = cache.InitDB(settings.WEBSTER)
	dbConnections[settings.CAMBRIDGE] = cache.InitDB(settings.CAMBRIDGE)
	dbConnections[settings.SOHA] = cache.InitDB(settings.SOHA)
	return &Translators{dbConnections: dbConnections}
}

// SearchDictionary performs the dictionary search and returns the result HTML from all dictionaries
func (t *Translators) SearchDictionary(term string) (map[string]string, error) {
	dictionaries := map[string]dicts.Dictionary{
		settings.WEBSTER:   dicts.WebsterDictionary{},
		settings.CAMBRIDGE: dicts.CambridgeDictionary{},
		settings.SOHA:      dicts.SohaDictionary{},
	}

	results := make(map[string]string)
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(dictionaries))

	for name, dictionary := range dictionaries {
		wg.Add(1)
		go func(name string, dictionary dicts.Dictionary) {
			defer wg.Done()
			con := t.dbConnections[name]
			_, soup, err := dictionary.Search(con, term, false)
			if err != nil {
				errChan <- fmt.Errorf("error searching %s dictionary: %w", name, err)
				return
			}
			html, err := soup.Html()
			if err != nil {
				errChan <- fmt.Errorf("error getting HTML from %s dictionary: %w", name, err)
				return
			}
			mu.Lock()
			results[name] = html
			mu.Unlock()
		}(name, dictionary)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	return results, nil
}
