{{ template "partials/header.tpl" . }}

<main class="events-page">

  <section class="events-hero container">
    <p class="eyebrow">EVENTS IN {{ .City }}</p>

    <h1>
      Something worth<br>
      <em>going out for.</em>
    </h1>

    <p class="lead">
      Discover music and sports events happening in {{ .City }}.
    </p>
  </section>


  <section class="events-section container">

    <div class="events-heading">
      <div>
        <p class="eyebrow">01 / MUSIC</p>
        <h2>Live music</h2>
      </div>

      {{ if .MusicEvents }}
        <span class="event-count">
          {{ len .MusicEvents }} events
        </span>
      {{ end }}
    </div>


    {{ if .MusicEvents }}

      <div class="event-grid">

        {{ range .MusicEvents }}

          <article class="event-card">

            <div class="event-image-wrap">
              <img
                src="{{ .ImageURL }}"
                alt="{{ .Name }}"
                class="event-image"
              >
            </div>

            <div class="event-card-body">

              <span class="event-category">
                {{ .Category }}
              </span>

              <h3>{{ .Name }}</h3>

              <p class="event-meta">
                {{ .LocalDate }}
                {{ if .LocalTime }}
                  · {{ .LocalTime }}
                {{ end }}
              </p>

              <p class="event-venue">
                {{ .Venue }}
              </p>

              <a
                href="/events/{{ .ID }}"
                class="event-link"
              >
                View Details
                <span>↗</span>
              </a>

            </div>

          </article>

        {{ end }}

      </div>

    {{ else }}

      <div class="empty-events">
        <p>No music events found in {{ .City }}.</p>
      </div>

    {{ end }}

  </section>


  <section class="events-section container">

    <div class="events-heading">
      <div>
        <p class="eyebrow">02 / SPORTS</p>
        <h2>Sports &amp; matchdays</h2>
      </div>

      {{ if .SportsEvents }}
        <span class="event-count">
          {{ len .SportsEvents }} events
        </span>
      {{ end }}
    </div>


    {{ if .SportsEvents }}

      <div class="event-grid">

        {{ range .SportsEvents }}

          <article class="event-card">

            <div class="event-image-wrap">
              <img
                src="{{ .ImageURL }}"
                alt="{{ .Name }}"
                class="event-image"
              >
            </div>

            <div class="event-card-body">

              <span class="event-category">
                {{ .Category }}
              </span>

              <h3>{{ .Name }}</h3>

              <p class="event-meta">
                {{ .LocalDate }}
                {{ if .LocalTime }}
                  · {{ .LocalTime }}
                {{ end }}
              </p>

              <p class="event-venue">
                {{ .Venue }}
              </p>

              <a
                href="/events/{{ .ID }}"
                class="event-link"
              >
                View Details
                <span>↗</span>
              </a>

            </div>

          </article>

        {{ end }}

      </div>

    {{ else }}

      <div class="empty-events">
        <p>No sports events found in {{ .City }}.</p>
      </div>

    {{ end }}

  </section>

</main>

{{ template "partials/footer.tpl" . }}