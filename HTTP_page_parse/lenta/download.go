package main

import (
	"github.com/mgutz/logxi/v1"
	"golang.org/x/net/html"
	"net/http"
)

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func isElem(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}

func isText(node *html.Node) bool {
	return node != nil && node.Type == html.TextNode
}

func isDiv(node *html.Node, class string) bool {
	return isElem(node, "div") && getAttr(node, "class") == class
}

func isTbody(node *html.Node) bool {
	return isElem(node, "tbody")
}

func isTr(node *html.Node) bool {
	return isElem(node, "tr")
}

func isTd(node *html.Node) bool {
	return isElem(node, "td")
}

type Item struct {
	Ref, CourseR, CourseD, Title string
}

func readItem(item *html.Node) *Item {
	var r, cr, ts, cd string
	for j := item.FirstChild; j != nil; j = j.NextSibling {
		if isTd(j) {
			for k := j.FirstChild; k != nil; k = k.NextSibling {
				if isDiv(k, "courses_table_inline") {
					for t := k.FirstChild; t != nil; t = t.NextSibling {
						//Name of cryptovalue
						if isDiv(t, "courses_table_name") {
							for y := t.FirstChild; y != nil; y = y.NextSibling {
								if isText(y) {
									ts = y.Data
								}
							}
						} else {
							//for rubles
							if isDiv(t, "courses_table_cost") {
								for y := t.FirstChild; y != nil; y = y.NextSibling {
									if isText(y) {
										cr = y.Data
									}
								}
							} else {
								//for $$
								if isDiv(t, "courses_table_cost1") {
									for y := t.FirstChild; y != nil; y = y.NextSibling {
										if isText(y) {
											cd = y.Data
										}
									}
								} else {
									//it's ref
									if isElem(t, "a") {
										r = getAttr(t, "href")
									}
								}
							}
						}
					}
				}
			}
		}
	}
	//packing
	return &Item{
		Ref:     r,
		CourseR: cr,
		CourseD: cd,
		Title:   ts,
	}
}

func search(node *html.Node) []*Item {
	if isTbody(node) {
		var items []*Item
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if isTr(c) {
				//for each child
				if item := readItem(c); item != nil {
					items = append(items, item)
				}
			}
		}
		return items
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if items := search(c); items != nil {
			return items
		}
	}
	return nil
}

func cryptoFinder() []*Item {
	log.Info("sending request to kibers.com")
	if response, err := http.Get("http://kibers.com/courses.html"); err != nil {
		log.Error("request to kibers.com failed", "error", err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		log.Info("got response from kibers.com", "status", status)
		if status == http.StatusOK {
			if doc, err := html.Parse(response.Body); err != nil {
				log.Error("invalid HTML from kibers.com", "error", err)
			} else {
				log.Info("HTML from kibers.com parsed successfully")
				return search(doc)
			}
		}
	}
	return nil
}
