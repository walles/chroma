package styles

import (
	"embed"
	"io/fs"
	"sort"
	"sync"

	"github.com/alecthomas/chroma/v2"
)

//go:embed *.xml
var embedded embed.FS

// Registry of Styles.
var Registry = func() map[string]*chroma.Style {
	registry := map[string]*chroma.Style{}
	// Register all embedded styles.
	files, err := fs.ReadDir(embedded, ".")
	if err != nil {
		panic(err)
	}

	mutex := sync.Mutex{}
	var wg sync.WaitGroup

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		wg.Add(1)

		// FIXME: What happens if this goroutine panics?
		go func(file fs.DirEntry) {
			defer wg.Done()

			r, err := embedded.Open(file.Name())
			if err != nil {
				panic(err)
			}
			style, err := chroma.NewXMLStyle(r)
			if err != nil {
				panic(err)
			}
			_ = r.Close()

			mutex.Lock()
			defer mutex.Unlock()
			registry[style.Name] = style
		}(file)
	}

	wg.Wait()
	return registry
}()

// Fallback style. Reassign to change the default fallback style.
var Fallback = Registry["swapoff"]

// Register a chroma.Style.
func Register(style *chroma.Style) *chroma.Style {
	Registry[style.Name] = style
	return style
}

// Names of all available styles.
func Names() []string {
	out := []string{}
	for name := range Registry {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// Get named style, or Fallback.
func Get(name string) *chroma.Style {
	if style, ok := Registry[name]; ok {
		return style
	}
	return Fallback
}
