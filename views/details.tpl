{{ template "partials/header.tpl" . }}

<main class="details-page">
  <div class="container">

    <div class="details-topbar">
      <a href="javascript:history.back()" class="back-link">
      ← Back to events
    </a>
    </div>

    <section class="event-details">

      <div class="event-main">

        {{ if .Event.ImageURL }}
        <div class="event-image">
          <img
            src="{{ .Event.ImageURL }}"
            alt="{{ .Event.Name }}"
          >
        </div>
        {{ else }}
        <div class="event-image event-image-placeholder">
          No image available
        </div>
        {{ end }}

        <div class="event-content">

          {{ if .Event.Category }}
          <p class="event-category">
            {{ .Event.Category }}
          </p>
          {{ end }}

          <h1>
            {{ .Event.Name }}
          </h1>

          {{ if .Event.Description }}
          <section class="about-event">

            <p class="section-label">
              ABOUT THIS EVENT
            </p>

            <h2>
              About this event
            </h2>

            <p>
              {{ .Event.Description }}
            </p>

          </section>
          {{ end }}

        </div>

      </div>

      <aside class="event-sidebar">

        <div class="event-info-card">

          <p class="section-label">
            MAKE A PLAN
          </p>

          <div class="event-info">

            <span>WHEN</span>

            <strong>
              {{ .Event.LocalDate }}

              {{ if .Event.LocalTime }}
              · {{ .Event.LocalTime }}
              {{ end }}
            </strong>

            {{ if .Event.Timezone }}
            <small>
              {{ .Event.Timezone }}
            </small>
            {{ end }}

          </div>

          <div class="event-info">

            <span>WHERE</span>

            {{ if .Event.Venue }}
            <strong>
              {{ .Event.Venue }}
            </strong>
            {{ end }}

            <small>
              {{ .Event.City }}

              {{ if .Event.State }}
              , {{ .Event.State }}
              {{ end }}

              {{ if .Event.Country }}
              , {{ .Event.Country }}
              {{ end }}
            </small>

            {{ if .Event.Address }}
            <small>
              {{ .Event.Address }}
            </small>
            {{ end }}

          </div>

          {{ if .Event.Category }}
          <div class="event-info">

            <span>CATEGORY</span>

            <strong>
              {{ .Event.Category }}
            </strong>

          </div>
          {{ end }}

          <a
            href="/redirect/{{ .Event.ID }}" target="_blank"
            class="ticket-button"
          >
            Get Tickets
            <span>↗</span>
          </a>

        </div>

      </aside>

    </section>

  </div>
</main>

{{ template "partials/footer.tpl" . }}