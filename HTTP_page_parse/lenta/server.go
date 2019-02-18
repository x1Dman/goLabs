package main

import (
	"github.com/mgutz/logxi/v1"
	"html/template"
	"net/http"
)

const INDEX_HTML = `
    <!doctype html>
    <html lang="ru">
        <head>
			
            <meta charset="utf-8">
			<style type="text/css">
			a {
				text-align:center;
				color:Yellow;
			}
			.referM{
				color:Yellow;
				text-align:center;
			}
			.titles {
				color:Red;
				text-align:center;
			}
			.maintitle{
				color:White;
				text-align:center;
			}
			body {
				background-color:Black
			}
			table {
				width: 80%;
				margin: auto;
			}
			th {
				height: 50px;
			}
			</style>
            <title>Курсы валют</title>
        </head>
        <body>
			<table border="1">
			<caption class="maintitle"><h1>Cryptocurrency value</h1></caption>
            	{{if .}}
					<tr>
					<th class="titles"><h3>Kind</h3></th>
					<th class="titles"><h3>CostInRubles($)</h3></th>
					</tr>
                	{{range .}}
					<tr>
						<td class="referM">
                    	<a href="http://kibers.com{{.Ref}}">{{.Title}}</a></td> 
						<td class="referM">{{.CourseR}} ({{.CourseD}})</td>
                    	<br/>
						</td>
						</tr>
                	{{end}}
            	{{else}}
                Не удалось загрузить новости!
            {{end}}
        </body>
    </html>
    `

var indexHtml = template.Must(template.New("index").Parse(INDEX_HTML))

func serveClient(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	log.Info("got request", "Method", request.Method, "Path", path)
	if path != "/" && path != "/index.html" {
		log.Error("invalid path", "Path", path)
		response.WriteHeader(http.StatusNotFound)
	} else if err := indexHtml.Execute(response, cryptoFinder()); err != nil {
		log.Error("HTML creation failed", "error", err)
	} else {
		log.Info("response sent to client successfully")
	}
}

func main() {
	http.HandleFunc("/", serveClient)
	log.Info("starting listener")
	log.Error("listener failed", "error", http.ListenAndServe("127.0.0.1:6060", nil))
}
