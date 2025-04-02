package serverdebug

import (
	"html/template"

	"github.com/labstack/echo/v4"
)

type page struct {
	Path        string
	Description string
}

type indexPage struct {
	pages    []page
	logLevel string
}

func newIndexPage() *indexPage {
	return &indexPage{}
}

func (i *indexPage) addPage(path string, description string) {
	// НАДО РЕАЛИЗОВАТЬ
	i.pages = append(i.pages, page{
		Path:        path,
		Description: description,
	})
}

// func (i *indexPage) setLogLevel(level string) {
// 	i.logLevel = level
// }

func (i indexPage) handler(eCtx echo.Context) error {
	return template.Must(template.New("index").Parse(`<html>
	<title>Chat Service Debug</title>
<body>
	<h2>Chat Service Debug</h2>
	<ul>
	// НАДО РЕАЛИЗОВАТЬ список страниц
		{{range .Pages}}
			<li><a href="{{.Path}}">{{.Description}}</a></li>
		{{end}}
	</ul>

	<h2>Log Level</h2>
	<form onSubmit="putLogLevel()">
		<select id="log-level-select">
			// НАДО РЕАЛИЗОВАТЬ список уровней логирования
			<option value="debug">debug</option>
			<option value="info">info</option>
			<option value="warn">warn</option>
			<option value="error">error</option>
			// НАДО РЕАЛИЗОВАТЬ по умолчанию выбрана опция, соответствующая текущему уровню
			<option value="debug" selected>debug</option>
		</select>
		<input type="submit" value="Change"></input>
	</form>
	
	<script>
		function putLogLevel() {
			const req = new XMLHttpRequest();
			req.open('PUT', '/log/level', false);
			// НАДО РЕАЛИЗОВАТЬ проставляем нужные заголовки
			req.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
			req.onload = function() { window.location.reload(); };
			req.send('level='+document.getElementById('log-level-select').value);
		};
	</script>
</body>
</html>
`)).Execute(eCtx.Response(), struct {
		Pages    []page
		LogLevel string
	}{
		Pages:    i.pages,
		LogLevel: i.logLevel,
	})
}
