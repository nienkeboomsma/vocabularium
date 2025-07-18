package template

var baseStyles = `
body {
	font-family: system-ui;
	line-height: 1.5;
	padding: 0.8rem 1rem;
}

nav {
	display: flex;
	justify-content: space-between;
}

a,
button {
	all: initial;
	background: none;
	border: none;
	border-radius: 4px;
	color: black;
	cursor: pointer;
	font-family: inherit;
	margin: -0.25rem -0.35rem;
	padding: 0.25rem 0.35rem;
	text-decoration-line: none;
	transition: 0.2s ease;
}

a:hover,
button:hover {
	background-color: rgba(0, 0, 0, 0.07);
	color: black;
	transition: 0.2s ease;
}

a:visited {
	color: black;
}
`
