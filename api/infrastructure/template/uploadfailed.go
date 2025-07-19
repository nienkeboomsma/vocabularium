package template

import (
	"fmt"
)

type UploadFailedData struct {
	Message string
	Error   string
}

var consoleStyles = `
.console {
  background-color: #1e1e1e;
  border-radius: 4px;
  color: #d4d4d4;
  font-family: monospace;
  padding: 1rem;
}
`

func GetFailedWorkUploadTemplate() string {
	template := `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8" />
		<title>Upload failed</title>
		<link rel="icon" href="https://fav.farm/❌" />
		<style>
			%s
			%s
		</style>
	</head>
	<body>
		<nav>
			<a href="http://localhost:4321/upload">👈🏻 Try again</a>
		</nav>
		<h1>❌ Upload failed</h1>
		<p>{{.Message}}</p>
		{{if .Error}}
			<p class="console">{{.Error}}</p>
		{{end}}
	</body>
</html>
`

	return fmt.Sprintf(template, baseStyles, consoleStyles)
}
