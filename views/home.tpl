{{ template "partials/header.tpl" . }}

<main>
  <section class="hero container">
    <div class="hero-copy">
      <p class="eyebrow">LESS SCROLLING. MORE GOING.</p>
      <h1>A city of possibilities.<br><em>Find your next one.</em></h1>
      <p class="lead">
        Discover music and sports in one place.<br>
        Choose your city. Find something worth heading out for.
      </p>

      <div class="tags">
        <span>Live music</span>
        <span>Sports &amp; matchdays</span>
        <span>One simple search</span>
      </div>
    </div>

    <div class="poster" aria-hidden="true">
      <span>THE CITY IS CALLING</span>
      <strong>GO<br><em>OUT.</em></strong>
      <small>GOOD PLANS.<br>GREAT MEMORIES.</small>
    </div>
  </section>

  <section class="search-section container">
    <div class="search-card">
      <div class="search-heading">
        <p class="eyebrow">01 / CHOOSE YOUR CITY</p>
        <h2>Where are you going?</h2>
        <p>Search for a city and select one of the matching suggestions.</p>
      </div>

      <div class="search-row">
        <div class="search-field">
          <label for="city">City</label>
          <input
            id="city"
            type="text"
            autocomplete="off"
            placeholder="Type a city..."
          >
          <div id="msg" class="search-message" aria-live="polite"></div>
          <ul id="list" class="suggestions" role="listbox"></ul>
          <p class="attribution">Powered by Google</p>
        </div>

        <button id="go" class="go-button" type="button" disabled>
          Explore Events <span>→</span>
        </button>
      </div>
    </div>
  </section>

  <section id="how-it-works" class="how container">
    <p class="eyebrow">HOW IT WORKS</p>
    <div class="how-grid">
      <div>
        <strong>01</strong>
        <h3>Choose a city</h3>
        <p>Select a city from the Google Places suggestions.</p>
      </div>
      <div>
        <strong>02</strong>
        <h3>Discover events</h3>
        <p>Browse Music and Sports events available in that city.</p>
      </div>
      <div>
        <strong>03</strong>
        <h3>Make your plan</h3>
        <p>Open an event and continue to the official ticket provider.</p>
      </div>
    </div>
  </section>
</main>

{{ template "partials/footer.tpl" . }}