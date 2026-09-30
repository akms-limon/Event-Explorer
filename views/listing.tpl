{{ template "partials/header.tpl" . }}

<main class="listing-page">
  <section class="listing-hero container">
    <p class="eyebrow">EVENTS IN {{ .City }}</p>
    <h1>Find something worth<br><em>heading out for.</em></h1>
    <p class="lead">
      Music and sports events happening in {{ .City }}.
    </p>
  </section>

  <section class="events-section container">
    <div class="section-heading">
      <p class="eyebrow">MUSIC</p>
      <h2>Live music</h2>
    </div>

    {{ if .MusicEvents }}
      <div class="event-grid">
        {{ range .MusicEvents }}
          <article class="event-card">
            <img
              src="{{ .ImageURL }}"
              alt="{{ .Name }}"
              class="event-image"
            >

            <div class="event-content">
              <span class="event-category">{{ .Category }}</span>

              <h3>{{ .Name }}</h3>

              <p>
                {{ .LocalDate }}
                {{ if .LocalTime }}
                  · {{ .LocalTime }}
                {{ end }}
              </p>

              <p>{{ .Venue }}</p>

              <a href="/events/{{ .ID }}">
                View Details →
              </a>
            </div>
          </article>
        {{ end }}
      </div>
    {{ else }}
      <p class="empty-state">
        No music events found in {{ .City }}.
      </p>
    {{ end }}
  </section>

  <section class="events-section container">
    <div class="section-heading">
      <p class="eyebrow">SPORTS</p>
      <h2>Sports &amp; matchdays</h2>
    </div>

    {{ if .SportsEvents }}
      <div class="event-grid">
        {{ range .SportsEvents }}
          <article class="event-card">
            <img
              src="{{ .ImageURL }}"
              alt="{{ .Name }}"
              class="event-image"
            >

            <div class="event-content">
              <span class="event-category">{{ .Category }}</span>

              <h3>{{ .Name }}</h3>

              <p>
                {{ .LocalDate }}
                {{ if .LocalTime }}
                  · {{ .LocalTime }}
                {{ end }}
              </p>

              <p>{{ .Venue }}</p>

              <a href="/events/{{ .ID }}">
                View Details →
              </a>
            </div>
          </article>
        {{ end }}
      </div>
    {{ else }}
      <p class="empty-state">
        No sports events found in {{ .City }}.
      </p>
    {{ end }}
  </section>
</main>

{{ template "partials/footer.tpl" . }}