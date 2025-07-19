package template

import (
	"fmt"
)

var tableStyles = `
h1:has(~table) {
	padding-left: 0.5rem;
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
	vertical-align: top;
}

td:has(a),
td:has(button) {
	padding-left: 0.2rem;
}
`

func GetWorkListTemplate() string {
	template := `
<!DOCTYPE html>
<html>
	<head>
		<title>Works</title>
		<link rel="icon" href="https://fav.farm/📜" />
		<style>
			%s
			%s
		</style>
	</head>
	<body>
		<nav>
			<a href="http://localhost:4321/upload">📥 Upload work</a>
		</nav>
		<h1>Works</h1>
		{{if .}}
			<div class="table">
				<table>
					<thead>
						<tr>
							<th colspan="2">Author</th>
							<th></th>
							<th colspan="3">Title</th>
							<th></th>
						</tr>
					</thead>
					<tbody>
						{{range .}}
							<tr>
								<td>{{.Author.Name}}</td>
								<td>
									<a title="Frequency list" href="http://localhost:4321/frequency-list-author/{{.Author.ID}}/true">📈</a>
								</td>
								<td></td>
								<td style="">{{.Title}}</td>
								<td>
									<a title="Frequency list" href="http://localhost:4321/frequency-list/{{.ID}}/true">📈</a>
								</td>
								<td>
									<a title="Glossary" href="http://localhost:4321/glossary/{{.ID}}/true">📖</a>
								</td>
								<td>
									<button title="Delete work" onclick="confirmAndDelete(this)" data-id="{{.ID}}">❌</button>
								</td>
							</tr>
						{{end}}
					</tbody>
				</table>
			</div>
		{{else}}
			<p>No works to display</p>
		{{end}}
	</body>
	<script>
		async function confirmAndDelete(button) {
			const id = button.getAttribute("data-id");
			const confirmed = confirm("Are you sure you want to delete this work?");

			if (!confirmed) return;

			const url = "http://localhost:4321/delete/" + id;

			try {
				const response = await fetch(url, { method: "POST" });

				if (!response.ok) {
					button.textContent = "👎🏻";
					button.title = "Failed to delete work; click to try again";
					return;
				}

				const row = button.closest("tr");
				if (row) row.remove();
		    } catch (error) {
  				button.textContent = "👎🏻";
				button.title = "Failed to delete work; click to try again";
		    }
		}
	</script>
</html>
`

	return fmt.Sprintf(template, baseStyles, tableStyles)
}
