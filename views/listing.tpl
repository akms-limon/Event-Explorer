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

      <div class="cache-tools">

        <button
          id="clear-full-cache"
          class="cache-button"
          type="button"
        >
          Clear full cache
        </button>

        <div class="cache-location-control">

          <select id="cache-city" class="cache-select">
            <option value="">City</option>
            <option value="{{ .City }}">{{ .City }}</option>
          </select>

          <select id="cache-country" class="cache-select">
            <option value="">Country</option>
            <option value="{{ .CountryCode }}">
              {{ .CountryCode }}
            </option>
          </select>

          <select id="cache-category" class="cache-select">
            <option value="">Category</option>
            <option value="Music">Music</option>
            <option value="Sports">Sports</option>
          </select>

          <button
            id="clear-selected-cache"
            class="cache-button"
            type="button"
            disabled
          >
            Clear selected cache
          </button>

        </div>

        <div
          id="cache-message"
          class="cache-message"
          aria-live="polite"
        ></div>

      </div>
    </section>


    <!-- Music -->

    <section class="event-section">

      <div class="section-heading">

        <div>
          <p class="eyebrow">Music</p>

          <h2>
            Music Events
          </h2>
        </div>

        {{ if .MusicCacheHit }}

          <span class="cache-status cache-hit">
            Cache hit
          </span>

        {{ else }}

          <span class="cache-status cache-fresh">
            Fresh sample
          </span>

        {{ end }}

      </div>


      {{ if .MusicEvents }}

        <div class="event-grid">

          {{ range .MusicEvents }}

            {{ template "partials/event_card.tpl" . }}

          {{ end }}

        </div>

      {{ else }}

        <div class="empty-state">

          <h3>
            No music events found
          </h3>

          <p>
            There are no music events available for this city right now.
          </p>

        </div>

      {{ end }}

    </section>


    <!-- Sports -->

    <section class="event-section">

      <div class="section-heading">

        <div>
          <p class="eyebrow">Sports</p>

          <h2>
            Sports Events
          </h2>
        </div>

        {{ if .SportsCacheHit }}

          <span class="cache-status cache-hit">
            Cache hit
          </span>

        {{ else }}

          <span class="cache-status cache-fresh">
            Fresh sample
          </span>

        {{ end }}

      </div>


      {{ if .SportsEvents }}

        <div class="event-grid">

          {{ range .SportsEvents }}

            {{ template "partials/event_card.tpl" . }}

          {{ end }}

        </div>

      {{ else }}

        <div class="empty-state">

          <h3>
            No sports events found
          </h3>

          <p>
            There are no sports events available for this city right now.
          </p>

        </div>

      {{ end }}

    </section>

  </div>
</main>

<script src="/static/js/cache.js"></script>

{{ template "partials/footer.tpl" . }}