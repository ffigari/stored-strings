package ui

import (
	"fmt"
	"net/http"
	"strings"
)

func LabeledInput(label, input string) string {
	return fmt.Sprintf(`
		<div class="row mb-3">
			<label
				class="col-sm-2 col-form-label"
				for="description-input"
			>%s</label>
			<div
				class="col-sm-10"
			>%s</div>
		</div>
	`, label, input)
}

func Form(submitText string, inputs []string) string {
	return fmt.Sprintf(`
		<form method="POST">
			%s
			<button class="btn btn-primary" type="submit">
				%s
			</button>
		</form>
	`, strings.Join(inputs, ""), submitText)
}

func mapS(xs []string, fn func(string) string) []string {
	ys := []string{}

	for _, x := range xs {
		ys = append(ys, fn(x))
	}

	return ys
}

func Paragraphs(ps []string) string {
	return strings.Join(mapS(ps, func(s string) string {
		return fmt.Sprintf(`<p>%s</p>`, s)
	}), "")
}

func HTMLHeader(w http.ResponseWriter, baseHTML string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
			<head>
				<meta charset="utf-8">
				<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">
				<meta name="viewport" content="width=device-width, initial-scale=1">
			</head>
			<body>
				<div class="container">
					<div class="my-3">
						%s
					</div>
				</div>
				<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/js/bootstrap.bundle.min.js" integrity="sha384-ka7Sk0Gln4gmtz2MlQnikT1wXgYsOg+OMhuP+IlRH9sENBO0LRn5q+8nbTov4+1p" crossorigin="anonymous"></script>
			</body>
		</html>
	`, baseHTML))
}
