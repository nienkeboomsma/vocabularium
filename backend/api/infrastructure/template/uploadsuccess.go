package template

import (
	"fmt"
)

func GetSuccessfulWorkUploadTemplate() string {
	template := `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8" />
		<title>Upload successful</title>
		<meta http-equiv="refresh" content="5;url=/" />
		<style>
			%s
		</style>
	</head>
	<body>
		<h1>✅ Upload successful!</h1>
		<p>You will be redirected shortly.</p>
		<p><a href="/">Click here</a> if you’re not redirected automatically.</p>
	</body>
</html>
`

	return fmt.Sprintf(template, baseStyles)
}
