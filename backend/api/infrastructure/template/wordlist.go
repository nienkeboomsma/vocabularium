package template

import (
	"fmt"

	"github.com/nienkeboomsma/collatinus/domain"
)

type WordListPageData struct {
	Title  string
	Author string
	Words  *[]domain.WordInWork
}

var wordListStyles = `
.subtle {
	font-size: 1.8rem;
	font-style: italic;
	font-weight: 500;
	padding: 0 0.2rem 0 0.25rem;
	opacity: 0.4;
}
`

func GetWordListTemplate(listType, emoji string) string {
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
      		<link rel="icon" href="https://fav.farm/%s" />
		<style>
			%s
			%s
			%s
		</style>
	</head>
	<body>
		<nav>
			<a href="http://localhost:4321">👈🏻 Back to works</a>
			<a id="toggle-link"></a>
		</nav>
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
						<th></th>
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
	  				button.textContent = "👎🏻";
					button.title = "Failed to mark words as " + newKnown ? "known" : "unknown" + "; click to try again";
				}

				button.textContent = newKnown ? "❌" : "✅";
				button.setAttribute("data-known", newKnown.toString());
				button.setAttribute("title", "Mark word as " + newKnown ? "unknown" : "known");
			} catch (error) {
				button.textContent = "👎🏻";
				button.title = "Failed to mark words as " + newKnown ? "known" : "unknown" + "; click to try again";
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
	return fmt.Sprintf(template, emoji, baseStyles, tableStyles, wordListStyles, listType)
}
