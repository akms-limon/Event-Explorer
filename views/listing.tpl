{{ template "partials/header.tpl" . }}

<main class="listing-page">
  <div class="container">

    <section class="listing-header">
      <a href="/" class="back-link">← Change City</a>

      <p class="eyebrow">Event Explorer</p>

      <h1>
        Events in {{ .City }}
      </h1>

      {{ if .CountryCode }}
        <p class="listing-location">
          {{ .CountryCode }}
        </p>
      {{ end }}
    </section>

    <section class="event-section">
      <div class="section-heading">
        <div>
          <p class="eyebrow">Music</p>
          <h2>Music Events</h2>
        </div>
      </div>

      {{ if .MusicEvents }}
        <div class="event-grid">
          {{ range .MusicEvents }}
            {{ template "partials/event_card.tpl" . }}
          {{ end }}
        </div>
      {{ else }}
        <div class="empty-state">
          <h3>No music events found</h3>
          <p>There are no music events available for this city right now.</p>
        </div>
      {{ end }}
    </section>

    <section class="event-section">
      <div class="section-heading">
        <div>
          <p class="eyebrow">Sports</p>
          <h2>Sports Events</h2>
        </div>
      </div>

      {{ if .SportsEvents }}
        <div class="event-grid">
          {{ range .SportsEvents }}
            {{ template "partials/event_card.tpl" . }}
          {{ end }}
        </div>
      {{ else }}
        <div class="empty-state">
          <h3>No sports events found</h3>
          <p>There are no sports events available for this city right now.</p>
        </div>
      {{ end }}
    </section>

  </div>
</main>

{{ template "partials/footer.tpl" . }}