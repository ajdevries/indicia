package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Result struct {
	Query  string
	Photos []*Photo
}

// Start a HTTP server with the storage engine to search for photos
func StartServer(i *Indicia) {
	t, _ := searchTemplate()
	http.HandleFunc("/", searchHandler(i.Storage, t))
	http.HandleFunc("/status", statusHandler(i))
	log.Print("Listening on http://localhost:9000")
	log.Print(http.ListenAndServe(":9000", nil))
}

// Status message
func statusHandler(i *Indicia) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "indexing %d/%d (%.1fs)", i.indexed, i.count, i.elapsed.Seconds())
	}
}

func searchHandler(s Storage, t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.FormValue("q")
		p := []*Photo{}
		if q != "" {
			p = s.Search(q)
		}
		t.Execute(w, Result{Query: q, Photos: p})
	}
}

func searchTemplate() (*template.Template, error) {
	return template.New("search").Parse(
		`
<html>
		<head>
				<title>I N D I C I A</title>

				<!-- Used codepen https://codepen.io/sebastianpopp/pen/myYmmy as starting point -->
				<style type="text/css">
				/*** COLORS ***/
				/*** DEMO ***/
				html,
				body {
					height: 100%;
					margin: 0;
					font: 21px monospace;
				}
				body {
					background: #913d88;
					color: #ffffff;
				}
				p {
					margin-top: 30px;
				}
				.cntr {
					display: table;
					width: 100%;
					height: 100%;
				}
				.cntr .cntr-innr {
					text-align: center;
					padding-top: 20px;
				}
				/*** STYLES ***/
				.search {
					display: inline-block;
					position: relative;
					height: 35px;
					width: 35px;
					box-sizing: border-box;
					margin: 0px 8px 7px 0px;
					padding: 2px 9px 0px 9px;
					border: 3px solid #ffffff;
					border-radius: 25px;
					-webkit-transition: all 200ms ease;
					transition: all 200ms ease;
					cursor: text;
				}
				.search:after {
					content: "";
					position: absolute;
					width: 3px;
					height: 20px;
					right: -5px;
					top: 21px;
					background: #ffffff;
					border-radius: 3px;
					-webkit-transform: rotate(-45deg);
									transform: rotate(-45deg);
					-webkit-transition: all 200ms ease;
					transition: all 200ms ease;
				}
				.search.active,
				.search:hover {
					width: 300px;
					margin-right: 0px;
				}
				.search.active:after,
				.search:hover:after {
					height: 0px;
				}
				.search input {
					width: 100%;
					border: none;
					box-sizing: border-box;
					font: 21px monospace;
					color: inherit;
					background: transparent;
					outline-width: 0px;
				}
				.status {
					float:right;
					font-size:13px;
					margin-right:20px;
				}
				ul {
					list-style-type:none;
				}
				a.photo {
					color: white;
					text-decoration:none;
				}
				.tags {
					background-color: white;
    			color: black;
				}
				.hidden {
					display:none;
				}
</style>
<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.1.1/jquery.js"></script>
<script>
$(function() {
	$("#search").on('focus', function () {
		$(this).parent('label').addClass('active');
	});

	$("#search").on('blur', function () {
		if ($(this).val().length == 0)
			$(this).parent('label').removeClass('active');
	});

	$(".photo").on('click', function () {
		$(this).parent().find('.tags').toggleClass('hidden');
		return false;
	});

	var timeout = setInterval(refreshStatus, 500);
	refreshStatus();
	function refreshStatus() {
		$('#status').load('status');
	}

	if ($("#search").val().length > 0) {
		$("#search").parent('label').addClass('active');
	}
});
</script>
		</head>
		<body>
		<p id="status" class="status"></p>
		<div class="cntr">
			<div class="cntr-innr">
				<form>
					<label class="search" for="search">
							<input id="search" name="q" type="text" value="{{.Query}}"/>
					</label>
				</form>
				<p>Search for '%' to shows all photos</p>
			</div>
			<div class="results">
				<ul>
						{{range .Photos}}
                <li><a class="photo" href="#">{{.Name}}</a>
								<table class="tags hidden">
								{{range $key, $value := .Tags}}
									<tr>
										<td>{{$key}}</td>
										<td>{{$value}}</td>
									</tr>
								{{end}}
								</table>
								</li>
            {{end}}
				</ul>
			</div>
		</div>
		</body>
</html>
`)
}
