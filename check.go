package markdownlink

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/russross/blackfriday/v2"
)

type checker struct {
	inputURL  *url.URL
	inputData []byte
	targets   []target
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
	DefaultProgressBar.Start(len(c.targets))

	res, err := c.checkLinks()
	if err != nil {
		return err
	}

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
		u, err := replaceURL(path)
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

var pat = regexp.MustCompile(`(http.)://raw\.githubusercontent\.com/(.+?)/(.+?)/(.*)`)

func replaceURL(path string) (*url.URL, error) {
	g := pat.FindSubmatch([]byte(path))
	p := fmt.Sprintf("%s://github.com/%s/%s/blob/%s", g[1], g[2], g[3], g[4])
	u, err := url.Parse(p)
	if err != nil {
		return nil, err
	}
	return u, nil
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
					c.targets = append(c.targets, target{
						rawDestPath: string(node.LinkData.Destination),
						URL:         dest,
					})
				}
			}
		}
		return blackfriday.GoToNext
	})
}

func (c *checker) checkLinks() ([]result, error) {
	var res []result
	resCh := make(chan result)

	wg := &sync.WaitGroup{}
	for _, ts := range c.targets {
		wg.Add(1)
		go func(t target) {
			defer wg.Done()
			if t.URL.Scheme == "https" || t.URL.Scheme == "http" {
				resp, err := c.client.Get(t.URL.String())
				status := http.StatusOK
				if err != nil {
					resCh <- result{
						link:       t.rawDestPath,
						statusCode: http.StatusInternalServerError,
						err:        err,
					}
				}
				status = resp.StatusCode
				resp.Body.Close()
				resCh <- result{
					link:       t.rawDestPath,
					statusCode: status,
				}
			} else if t.URL.Scheme == "file" {
				status := http.StatusOK
				fpath := t.URL.Path
				if t.URL.Fragment != "" {
					base := filepath.Base(c.inputURL.String())
					fpath += fmt.Sprintf("/%s#%s", base, t.URL.Fragment)
				}
				nfpath := fmt.Sprintf("%s://%s%s/../%s", c.inputURL.Scheme, c.inputURL.Host, c.inputURL.Path, fpath)
				resp, err := c.client.Get(nfpath)
				if err != nil {
					resCh <- result{
						link:       t.rawDestPath,
						statusCode: http.StatusInternalServerError,
						err:        err,
					}
				}
				status = resp.StatusCode
				resp.Body.Close()
				resCh <- result{
					link:       t.rawDestPath,
					statusCode: status,
				}
			}
			DefaultProgressBar.Increment()
		}(ts)
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	for r := range resCh {
		res = append(res, r)
	}
	return res, nil
}
