package template

import (
	"fmt"
)

var uploadStyles = `
form {
	display: flex;
	flex-direction: column;
	gap: 1rem;
	width: fit-content;
}

form label {
	display: flex;
	gap: 1rem;
}

form label span {
	display: block;
	min-width: 5rem;
}

form label input {
	width: 10.7rem;
}

form button {
	all: revert;
	align-self: stretch;
}
`

func GetUploadTemplate() string {
	template := `
<!DOCTYPE html>
<html>
	<head>
		<title>Upload</title>
		<style>
			%s
			%s
		</style>
	</head>
	<body>
		<nav>
			<a href="http://localhost:4321">👈🏻 Back to works</a>
		</nav>
		<h1>Upload work</h1>
		<form action="http://localhost:4321/lemmatise" method="POST" enctype="multipart/form-data">
			<label>
				<span>File (.txt)</span>
				<input type="file" id="file" name="file" accept=".txt" required>
			</label>

			<label placeholder="Plautus">
				<span>Author</span>
				<input type="text" id="author" name="author" required>
			</label>

			<label placeholder="Amphitryo">
				<span>Title</span>
				<input type="text" id="title" name="title" required>
			</label>

			<button type="submit">Submit</button>
	</form>
	</body>
</html>
`
	return fmt.Sprintf(template, baseStyles, uploadStyles)
}
