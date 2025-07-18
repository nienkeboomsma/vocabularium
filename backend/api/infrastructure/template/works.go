package template

import (
	"fmt"
)

var tableStyles = `
h1 {
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
}
`

func GetWorkListTemplate() string {
	template := `
<!DOCTYPE html>
<html>
	<head>
		<title>Works</title>
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

	return fmt.Sprintf(template, baseStyles, tableStyles)
}
