{{ template "partials/header.tpl" . }}

<main class="listing-page">
  <div class="container">

    <section class="listing-header">

      <a href="/" class="back-link">
        ← Change City
      </a>

      <p class="eyebrow">
        Event Explorer
      </p>

      <h1>
        Events in {{ .City }}
      </h1>

      {{ if .CountryCode }}
        <p class="listing-location">
          {{ .CountryCode }}
        </p>
      {{ end }}


      <div class="cache-tools">

        <!-- Clear full cache -->

        <div class="cache-full-row">

          <h3>
            Clear full cache:
          </h3>

          <button
            id="clear-full-cache"
            class="cache-button cache-button-primary"
            type="button"
          >
            Clear full cache
          </button>

        </div>


        <!-- Clear cache by category -->

        <div class="cache-category-row">

          <h3>
            Clear cache by category:
          </h3>

          <div class="cache-location-control">

            <!-- City search -->

            <div class="cache-city-field">

              <label for="cache-city-search">
                City
              </label>

              <input
                id="cache-city-search"
                class="cache-city-input"
                type="text"
                autocomplete="off"
                placeholder="Search cached city..."
              >

              <div
                id="cache-city-message"
                class="cache-message"
                aria-live="polite"
              ></div>

              <ul
                id="cache-city-suggestions"
                class="cache-city-suggestions"
              ></ul>

            </div>


            <!-- Country -->

            <div class="cache-select-field">

              <label for="cache-country">
                Country
              </label>

              <select
                id="cache-country"
                class="cache-select"
                disabled
              >
                <option value="">
                  Select country
                </option>
              </select>

            </div>


            <!-- Category -->

            <div class="cache-select-field">

              <label for="cache-category">
                Category
              </label>

              <select
                id="cache-category"
                class="cache-select"
              >
                <option value="">
                  Select category
                </option>

                <option value="Music">
                  Music
                </option>

                <option value="Sports">
                  Sports
                </option>
              </select>

            </div>


            <!-- Clear selected cache -->

            <button
              id="clear-selected-cache"
              class="cache-button cache-button-secondary"
              type="button"
              disabled
            >
              Clear selected cache
            </button>

          </div>

        </div>


        <div
          id="cache-action-message"
          class="cache-action-message"
          aria-live="polite"
        ></div>

      </div>

    </section>


    <!-- =========================
         Music Events
    ========================== -->

    <section class="event-section">

      <div class="section-heading">

        <div>

          <p class="eyebrow">
            Music
          </p>

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
            There are no music events available
            for this city right now.
          </p>

        </div>

      {{ end }}

    </section>


    <!-- =========================
         Sports Events
    ========================== -->

    <section class="event-section">

      <div class="section-heading">

        <div>

          <p class="eyebrow">
            Sports
          </p>

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
            There are no sports events available
            for this city right now.
          </p>

        </div>

      {{ end }}

    </section>

  </div>
</main>


<script src="/static/js/cache.js"></script>

{{ template "partials/footer.tpl" . }}