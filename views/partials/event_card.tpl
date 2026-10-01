{{ define "partials/event_card.tpl" }}
<article class="event-card">
  <div class="event-image">
    {{ if .ImageURL }}
      <img src="{{ .ImageURL }}" alt="{{ .Name }}">
    {{ else }}
      <div class="event-image-placeholder">No image available</div>
    {{ end }}
  </div>

  <div class="event-content">
    {{ if .Category }}
      <span class="event-category">{{ .Category }}</span>
    {{ end }}

    <h3>{{ .Name }}</h3>

    {{ if .LocalDate }}
      <p class="event-date">
        {{ .LocalDate }}
        {{ if .LocalTime }} · {{ .LocalTime }}{{ end }}
      </p>
    {{ end }}

    {{ if .Venue }}
      <p class="event-venue">{{ .Venue }}</p>
    {{ end }}

    <a
      href="/events/{{ .ID }}"
      class="event-link"
    >
      View Details
    </a>
  </div>
</article>
{{ end }}