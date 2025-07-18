package api

import (
	"fmt"

	"github.com/nienkeboomsma/collatinus/domain"
)

type WordListPageData struct {
	Title  string
	Author string
	Words  *[]domain.WordInWork
}

var styles = `
body {
	font-family: system-ui;
	line-height: 1.5;
	padding: 0.8rem 1rem;
}

h1 {
	padding-left: 0.5rem;
}

.subtle {
	font-size: 1.8rem;
	font-style: italic;
	font-weight: 500;
	padding: 0 0.2rem 0 0.25rem;
	opacity: 0.4;
}

a, button {
	background: none;
	border: none;
	border-radius: 4px;
	color: black;
	cursor: pointer;
	padding: 0.25rem 0.35rem;
	text-decoration-line: none;
	transition: 0.2s ease;
}

a:hover, button:hover {
	background-color: rgba(0, 0, 0, 0.07);
	color: black;
	transition: 0.2s ease;
}

a:visited {
	color: black;
}

div.header {
	display: flex;
	justify-content: space-between;
}

.table {
	max-height: calc(100vh - 10.5rem);
	overflow-y: auto;
}

table {
	border-spacing: 0;
	text-align: left;
}

thead {
	position: sticky;
	top: 0;
}

table thead tr,
table tbody tr:nth-child(even) {
	background-color: #f5f5f5;
}

td, th {
	padding: 0.3rem 0.6rem 0.4rem;
}
`

func getWordListTemplate(title string) string {
	template := `
<!DOCTYPE html>
<html>
	<head>
		<title>
			{{if .Title}}
    			{{.Title}} by
    		{{end}}
       		{{.Author}}
        </title>
		<style>
			%s
		</style>
	</head>
	<body>
		<div class="header">
			<a href="http://localhost:4321/works">👈🏻 Back to works</a>
			<a id="toggle-link"></a>
		</div>
		<h1>
			%s <span class="subtle">for</span>
				{{if .Title}}
    				{{.Title}} <span class="subtle">by</span>
    			{{end}}
  			{{.Author}}</h1>
		<div class="table">
			<table>
				<thead>
					<tr>
						<th>Lemma</th>
						<th>Translation</th>
						<th>Count</th>
						<th>Mark</th>
					</tr>
				</thead>
				<tbody>
				{{range .Words}}
					<tr>
						<td>{{.LemmaRich}}</td>
						<td>{{.Translation}}</td>
						<td>{{.Count}}</td>
						<td>
						    {{if .Known}}
								<button title="Mark word as unknown" onclick="toggleKnown(this)" data-id="{{.ID}}" data-known="true">
									❌
								</button>
							{{else}}
								<button title="Mark word as known" onclick="toggleKnown(this)" data-id="{{.ID}}" data-known="false">
									✅
								</button>
							{{end}}
						</td>
					</tr>
				{{else}}
					<tr><td colspan="4">No words to display</td></tr>
				{{end}}
				</tbody>
			</table>
		</div>
	</body>
	<script>
		async function toggleKnown(button) {
			const id = button.getAttribute("data-id");
			const currentKnown = button.getAttribute("data-known") === "true";
			const newKnown = !currentKnown;
			const url = "http://localhost:4321/toggle-known-status/" + id;

			try {
				const response = await fetch(url, { method: "POST" });

				if (!response.ok) {
					console.error("Update failed");
					return;
				}

				button.textContent = newKnown ? "❌" : "✅";
				button.setAttribute("data-known", newKnown.toString());
				button.setAttribute("title", "Mark word as " + newKnown ? "unknown" : "known");
			} catch (error) {
				console.error("Request failed", error);
			}
		}

		const currentURL = new URL(window.location.href);
		const currentPath = currentURL.pathname;

		let text, newPath;

		if (currentPath.includes("true")) {
			newPath = currentPath.replace("true", "false");
			text = "🔄 Show known words";
		} else if (currentPath.includes("false")) {
			newPath = currentPath.replace("false", "true");
			text = "🔄 Hide known words";
		}

		const newURL = new URL(window.location.href);
		newURL.pathname = newPath;

		const toggleLink = document.getElementById("toggle-link");
		toggleLink.href = newURL.href;
		toggleLink.textContent = text;
	</script>
</html>
`
	return fmt.Sprintf(template, styles, title)
}

func getWorkListTemplate() string {
	template := `
<!DOCTYPE html>
<html>
	<head>
		<title>Works</title>
		<style>
			%s
		</style>
	</head>
	<body>
		<h1>Works</h1>
		<div class="table">
			<table>
				<thead>
					<tr>
						<th colspan="2">Author</th>
						<th></td>
						<th colspan="3">Title</th>
					</tr>
				</thead>
				<tbody>
					{{range .}}
						<tr>
							<td>{{.Author.Name}}</td>
							<td style="padding: 0;"><a title="Frequency list" href="http://localhost:4321/frequency-list-author/{{.Author.ID}}/true">📈</a></td>
							<td></td>
							<td style="">{{.Title}}</td>
							<td style="padding: 0;"><a title="Frequency list" href="http://localhost:4321/frequency-list/{{.ID}}/true">📈</a></td>
							<td style="padding: 0;"><a title="Glossary" href="http://localhost:4321/glossary/{{.ID}}/true">📖</a></td>
						</tr>
					{{else}}
						<tr><td colspan="5">No works to display</td></tr>
					{{end}}
				</tbody>
			</table>
		</div>
	</body>
</html>
`

	return fmt.Sprintf(template, styles)
}
