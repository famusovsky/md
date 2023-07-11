// Пакет htmltemplates предоставляет функции для работы с шаблонами страниц заметок.
package htmltemplates

import (
	"html/template"
	"path/filepath"

	"github.com/famusovsky/md/internal/models"
)

// Data - структура, которая хранит в себе данные для шаблона страницы заметки.
type Data struct {
	Note         *models.Note
	RenderedNote string
}

// CreateNewCache - создание нового кэша шаблонов.
// Параметр dir - путь к директории с шаблонами.
// Возвращает кэш шаблонов и ошибку.
func CreateNewCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
