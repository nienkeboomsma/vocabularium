package template

import (
	"fmt"
)

var uploadFailedStyles = `
.error {
  background-color: #1e1e1e;
  border-radius: 4px;
  color: #d4d4d4;
  font-family: monospace;
  padding: 1rem;
}
`

func GetFailedWorkUploadTemplate(message string, err error) string {
	template := `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8" />
		<title>Upload failed</title>
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
		<p>%s</p>
		<p class="error">%s</p>
	</body>
</html>
`

	return fmt.Sprintf(template, baseStyles, uploadFailedStyles, message, err)
}
