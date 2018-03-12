package distributor

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strings"
	"time"
)

const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>Cortex Ingester Stats</title>
	</head>
	<body>
		<h1>Cortex Ingester Stats</h1>
		<p>Current time: {{ .Now }}</p>
		<form action="" method="POST">
			<input type="hidden" name="csrf_token" value="$__CSRF_TOKEN_PLACEHOLDER__">
			<table width="100%" border="1">
				<thead>
					<tr>
						<th>User</th>
						<th># Series</th>
						<th>Ingest Rate</th>
					</tr>
				</thead>
				<tbody>
					{{ range .Stats }}
					<tr>
						<td>{{ .UserID }}</td>
						<td>{{ .UserStats.NumSeries }}</td>
						<td>{{ .UserStats.IngestionRate }}</td>
					</tr>
					{{ end }}
				</tbody>
			</table>
		</form>
	</body>
</html>`

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("webpage").Parse(tpl))
}

type userStatsByTimeseries []UserIDStats

func (s userStatsByTimeseries) Len() int           { return len(s) }
func (s userStatsByTimeseries) Less(i, j int) bool { return s[i].NumSeries > s[j].NumSeries }
func (s userStatsByTimeseries) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// AllUserStatsHandler shows stats for all users.
func (d *Distributor) AllUserStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := d.AllUserStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sort.Sort(userStatsByTimeseries(stats))

	if encodings, found := r.Header["Accept"]; found &&
		len(encodings) > 0 && strings.Contains(encodings[0], "json") {
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			http.Error(w, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
		}
		return
	}

	if err := tmpl.Execute(w, struct {
		Now   time.Time
		Stats []UserIDStats
	}{
		Now:   time.Now(),
		Stats: stats,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
