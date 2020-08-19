package mlc

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/mattn/go-colorable"
	"github.com/russross/blackfriday/v2"
)

type result struct {
	fileName   string
	rawLink    string
	statusCode int
	err        error
}

type content struct {
	filePath string
	links    []string
}

func (c *content) retrieveLinks(markdown []byte) {
	var links []string
	ast := blackfriday.New().Parse(markdown)
	ast.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if node.Type == blackfriday.Link || node.Type == blackfriday.Image {
			if entering {
				links = append(links, string(node.LinkData.Destination))
			}
		}
		return blackfriday.GoToNext
	})
	c.links = links
}

func Run(cf *config) error {
	out := colorable.NewColorableStdout()
	var contents []*content
	if cf.allMode {
		var err error
		contents, err = cf.fetchFiles(cf.ctx, "")
		if err != nil {
			return err
		}
	} else {
		contents = append(contents, &content{filePath: cf.filePath})
	}

	fmt.Fprint(out, cyan(fmt.Sprintf("%d file found.\n", len(contents))))
	for _, c := range contents {
		fmt.Fprint(out, cyan(fmt.Sprintf("FILE: %s\n", c.filePath)))
	}
	fmt.Fprintln(out, "")

	var totalLinks int
	for _, content := range contents {
		req, err := http.NewRequestWithContext(cf.ctx, "GET", fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", cf.owner, cf.repo, cf.branch, content.filePath), nil)
		if err != nil {
			return err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("response code: %d", resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		content.retrieveLinks(body)
		resp.Body.Close()
		totalLinks += len(content.links)
	}

	DefaultProgressBar.Start(totalLinks)

	resCh := make(chan result)
	var wg sync.WaitGroup

	for _, c := range contents {
		for _, l := range c.links {
			wg.Add(1)
			go func(link string) {
				defer wg.Done()
				if strings.HasPrefix(link, "https") || strings.HasPrefix(link, "http") {
					resp, err := http.Get(link)
					var (
						status = http.StatusOK
						errr   error
					)
					if err != nil {
						status = http.StatusInternalServerError
						errr = err
					} else {
						status = resp.StatusCode
					}
					resp.Body.Close()
					resCh <- result{
						fileName:   c.filePath,
						rawLink:    link,
						statusCode: status,
						err:        errr,
					}
				} else {
					var targetURL string
					root := fmt.Sprintf("https://github.com/%s/%s/blob/%s/", cf.owner, cf.repo, cf.branch)
					path := c.filePath
					if strings.HasPrefix(link, "#") {
						targetURL = root + path + link
					} else if strings.HasPrefix(link, "/") {
						targetURL = root + link
					} else if !strings.HasPrefix(link, "/") {
						targetURL = root + "../" + link
					}
					resp, err := http.Get(targetURL)
					var (
						status = http.StatusOK
						errr   error
					)
					if err != nil {
						status = http.StatusInternalServerError
						errr = err
					} else {
						status = resp.StatusCode
					}
					resp.Body.Close()
					resCh <- result{
						fileName:   c.filePath,
						rawLink:    link,
						statusCode: status,
						err:        errr,
					}
				}
				DefaultProgressBar.Increment()
			}(l)
		}
	}

	go func() {
		wg.Wait()
		close(resCh)
	}()

	var res []result
	for r := range resCh {
		res = append(res, r)
	}

	DefaultProgressBar.Finish()
	print(out, res)

	return nil
}

func print(out io.Writer, res []result) {
	var errRes []result
	for _, r := range res {
		fmt.Fprintf(out, fmt.Sprintf("[%v] %v\n", getStatusLabel(r.statusCode), r.rawLink))
		if r.statusCode >= 400 {
			errRes = append(errRes, r)
		}
	}
	fmt.Fprintf(out, "\n%d links checked.\n", len(res))

	if len(errRes) > 0 {
		fmt.Fprint(out, red(fmt.Sprintf("\nERROR: %d dead links found!\n", len(errRes))))
	}
	for _, r := range errRes {
		fmt.Fprintf(out, fmt.Sprintf("[%v] %v -> Status: %d\n", getStatusLabel(r.statusCode), r.rawLink, r.statusCode))
	}
}
