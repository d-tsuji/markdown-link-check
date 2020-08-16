package markdownlink

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/russross/blackfriday/v2"
)

type checker struct {
	inputURL  *url.URL
	inputData []byte
	target    []target
	client    *http.Client
}

type target struct {
	rawDestPath string
	URL         *url.URL
}

type result struct {
	link       string
	statusCode int
	err        error
}

// Check parses the markdown file and checks if the link exists.
// It then reports back to you.
func Check(config Config) error {
	out := colorable.NewColorableStdout()
	if len(config.args) == 0 {
		return nil
	}
	path := config.args[0]
	fmt.Fprint(out, cyan(fmt.Sprintf("FILE: %s\n", path)))
	c, err := build(path)
	if err != nil {
		return err
	}
	c.retrieveLinks()
	DefaultProgressBar.Start(len(c.target))

	res := c.checkLinks()

	DefaultProgressBar.Finish()

	var errRes []result
	for _, r := range res {
		fmt.Fprintf(out, fmt.Sprintf("[%v] %v\n", getStatusLabel(r.statusCode), r.link))
		if r.statusCode >= 400 {
			errRes = append(errRes, r)
		}
	}
	fmt.Fprintf(out, "\n%d links checked.\n", len(res))

	if len(errRes) > 0 {
		fmt.Fprint(out, red(fmt.Sprintf("\nERROR: %d dead links found!\n", len(errRes))))
	}
	for _, r := range errRes {
		fmt.Fprintf(out, fmt.Sprintf("[%v] %v -> Status: %d\n", getStatusLabel(r.statusCode), r.link, r.statusCode))
	}

	return nil
}

func build(path string) (*checker, error) {
	transport := new(http.Transport)
	*transport = *http.DefaultTransport.(*http.Transport) // Clone.
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))
	client := &http.Client{Transport: transport, Timeout: 10 * time.Second}

	if strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "http://") {
		u, err := url.Parse(path)
		if err != nil {
			return nil, err
		}
		resp, err := client.Get(path)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, err
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		c := &checker{
			inputURL:  u,
			inputData: b,
			client:    client,
		}
		return c, nil
	} else {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		abs, err := urlFromFilePath(absPath)
		if err != nil {
			return nil, err
		}
		f, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		c := &checker{
			inputURL:  abs,
			inputData: f,
			client:    client,
		}
		return c, nil
	}
}

func (c *checker) retrieveLinks() {
	ast := blackfriday.New().Parse(c.inputData)
	ast.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if node.Type == blackfriday.Link || node.Type == blackfriday.Image {
			if entering {
				dest, err := url.Parse(string(node.LinkData.Destination))
				if err == nil {
					if dest.Scheme == "" {
						dest.Scheme = "file"
					}
					c.target = append(c.target, target{
						rawDestPath: string(node.LinkData.Destination),
						URL:         dest,
					})
				}
			}
		}
		return blackfriday.GoToNext
	})
}

func (c *checker) checkLinks() []result {
	var res []result

	for _, l := range c.target {
		if l.URL.Scheme == "https" || l.URL.Scheme == "http" {
			resp, err := c.client.Get(l.URL.String())
			status := http.StatusOK
			if err != nil {
				res = append(res, result{
					link:       l.rawDestPath,
					statusCode: http.StatusInternalServerError,
					err:        err,
				})
				continue
			}
			status = resp.StatusCode
			resp.Body.Close()
			res = append(res, result{
				link:       l.rawDestPath,
				statusCode: status,
			})
		} else if l.URL.Scheme == "file" {
			status := http.StatusOK
			fpath := l.URL.Path
			if l.URL.Fragment != "" {
				base := filepath.Base(c.inputURL.String())
				fpath += fmt.Sprintf("/%s#%s", base, l.URL.Fragment)
			}
			nfpath := fmt.Sprintf("%s://%s%s/../%s", c.inputURL.Scheme, c.inputURL.Host, c.inputURL.Path, fpath)
			resp, err := c.client.Get(nfpath)
			if err != nil {
				res = append(res, result{
					link:       l.rawDestPath,
					statusCode: http.StatusInternalServerError,
					err:        err,
				})
				continue
			}
			status = resp.StatusCode
			resp.Body.Close()
			res = append(res, result{
				link:       l.rawDestPath,
				statusCode: status,
			})
		}
		DefaultProgressBar.Increment()
	}
	return res
}
