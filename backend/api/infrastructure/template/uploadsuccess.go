package template

import (
	"fmt"
)

type UploadSuccessData struct {
	Logs []string
}

func GetSuccessfulWorkUploadTemplate() string {
	template := `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8" />
		<title>Upload successful</title>
		{{if not .Logs}}
			<meta http-equiv="refresh" content="5;url=/" />
		{{end}}
		<style>
			%s
			%s
		</style>
	</head>
	<body>
		<nav>
			<a href="http://localhost:4321">👈🏻 Back to works</a>
		</nav>
		<h1>✅ Upload successful!</h1>
		{{if .Logs}}
			<p>The following messages were logged:</p>
			<p class="console">
				{{range .Logs}}
					{{.}}<br/>
				{{end}}
			</p>
		{{else}}
			<p>You will be redirected shortly.</p>
			<p><a href="/">Click here</a> if you’re not redirected automatically.</p>
		{{end}}
	</body>
</html>
`

	return fmt.Sprintf(template, baseStyles, consoleStyles)
}
