<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">

    <title>Event Explorer</title>

    <style>
        :root {
            --bg: #1e1f1d;
            --bar: #111111;
            --strip: #292928;
            --panel: #121212;
            --line: #3a3b39;
            --green: #0a6b57;
            --mint: #6bbba6;
            --ice: #d8eef2;
            --muted: #8f9a98;
            --cream: #e5ee9f;
            --orange: #e8834f;
            --amber: #e9a24a;
        }

        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }

        html {
            scroll-behavior: smooth;
        }

        body {
            background: var(--bg);
            color: var(--ice);
            font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
            -webkit-font-smoothing: antialiased;
        }

        a {
            color: inherit;
            text-decoration: none;
        }

        button,
        input {
            font: inherit;
        }

        :focus-visible {
            outline: 2px solid var(--mint);
            outline-offset: 3px;
        }

        .wrap {
            max-width: 1104px;
            margin: 0 auto;
            padding: 0 24px;
        }

        /* Header */

        header {
            background: var(--bar);
            height: 85px;
            border-bottom: 1px solid #222;
        }

        header .wrap {
            height: 100%;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }

        .logo {
            display: flex;
            align-items: center;
            gap: 14px;
            font-size: 24px;
            letter-spacing: -0.04em;
            color: #fff;
        }

        .logo b {
            font-weight: 700;
        }

        .logo span {
            font-weight: 300;
        }

        .mark {
            width: 36px;
            height: 38px;
            background: var(--green);
            border-radius: 9px;
            display: grid;
            place-items: center;
            color: #fff;
            font-weight: 700;
            font-size: 30px;
            letter-spacing: -0.06em;
            position: relative;
        }

        .mark::after {
            content: "";
            position: absolute;
            right: 6px;
            bottom: 8px;
            width: 5px;
            height: 5px;
            background: var(--cream);
            border-radius: 50%;
        }

        nav {
            display: flex;
            align-items: center;
            gap: 32px;
            font-size: 14px;
            font-weight: 700;
        }

        nav a.on {
            color: var(--mint);
            padding: 8px 0;
            border-bottom: 2px solid var(--mint);
        }

        .pill {
            background: #33291a;
            color: var(--amber);
            padding: 9px 16px 9px 14px;
            border-radius: 99px;
            font-size: 12px;
            font-weight: 700;
            display: flex;
            align-items: center;
            gap: 8px;
        }

        .pill::before {
            content: "";
            width: 5px;
            height: 5px;
            border-radius: 50%;
            background: #a86a1c;
        }

        .strip {
            background: var(--strip);
            text-align: center;
            font-size: 11px;
            color: var(--muted);
            padding: 9px 0;
        }

        .strip b {
            color: #b7c4c1;
            margin-right: 6px;
        }

        /* Hero */

        .hero {
            display: grid;
            grid-template-columns: 1fr 340px;
            gap: 40px;
            align-items: center;
            padding: 70px 0 58px;
        }

        .eyebrow {
            display: flex;
            align-items: center;
            gap: 12px;
            font-size: 11px;
            font-weight: 700;
            letter-spacing: 0.14em;
            color: var(--mint);
            margin-bottom: 26px;
        }

        .eyebrow::before {
            content: "";
            width: 24px;
            height: 2px;
            background: var(--green);
        }

        h1 {
            font-size: clamp(40px, 5.2vw, 56px);
            line-height: 1.08;
            letter-spacing: -0.045em;
            font-weight: 700;
            color: var(--ice);
        }

        h1 em {
            font-style: normal;
            color: var(--mint);
        }

        .lead {
            margin-top: 30px;
            color: var(--muted);
            font-size: 17px;
            line-height: 1.65;
        }

        .tags {
            display: flex;
            flex-wrap: wrap;
            gap: 12px;
            margin-top: 28px;
        }

        .tags span {
            border: 1px solid var(--line);
            border-radius: 99px;
            padding: 9px 13px;
            font-size: 11px;
            color: var(--muted);
            background: #1b1c1a;
        }

        /* Poster */

        .stage {
            position: relative;
            height: 340px;
        }

        .stage::before {
            content: "";
            position: absolute;
            inset: 14px -8px -12px 10px;
            background: #2b2c28;
            border-radius: 6px 44px 6px 6px;
            transform: rotate(3deg);
        }

        .poster {
            position: absolute;
            inset: 0;
            background: var(--green);
            border-radius: 4px 44px 4px 4px;
            transform: rotate(3deg);
            padding: 22px 24px;
            overflow: hidden;
        }

        .poster::before,
        .poster::after {
            content: "";
            position: absolute;
            border-radius: 50%;
            background: rgba(0, 0, 0, 0.09);
        }

        .poster::before {
            width: 270px;
            height: 270px;
            right: -110px;
            top: 70px;
        }

        .poster::after {
            width: 150px;
            height: 150px;
            right: -50px;
            top: 135px;
            background: rgba(255, 255, 255, 0.04);
        }

        .poster small {
            position: absolute;
            font-size: 9px;
            letter-spacing: 0.14em;
            color: var(--cream);
        }

        .poster .t {
            left: 24px;
            top: 20px;
        }

        .poster .n {
            right: 24px;
            top: 24px;
        }

        .poster .big {
            position: absolute;
            left: 24px;
            top: 58px;
            font-size: 84px;
            line-height: 0.9;
            font-weight: 700;
            letter-spacing: -0.06em;
        }

        .poster .big i {
            display: block;
            font-style: normal;
            color: var(--cream);
        }

        .poster .big u {
            display: block;
            text-decoration: none;
            color: #fff;
        }

        .poster .b {
            left: 24px;
            bottom: 28px;
            line-height: 1.6;
        }

        .poster .arr {
            position: absolute;
            right: 26px;
            bottom: 26px;
            color: var(--orange);
            font-size: 34px;
            line-height: 1;
        }

        /* Search */

        .search {
            background: var(--panel);
            border: 1px solid var(--line);
            border-radius: 14px;
            padding: 30px;
        }

        .search h2 {
            display: inline;
            font-size: 24px;
            letter-spacing: -0.03em;
        }

        .search .sub {
            font-size: 13px;
            color: var(--muted);
            margin-left: 16px;
        }

        label {
            display: block;
            margin: 26px 0 10px;
            font-size: 12px;
            font-weight: 700;
        }

        .row {
            display: flex;
            gap: 14px;
        }

        .field {
            flex: 1;
            position: relative;
        }

        .field input {
            width: 100%;
            height: 57px;
            background: #181a18;
            border: 1px solid #55595a;
            border-radius: 8px;
            padding: 0 20px 0 52px;
            color: var(--ice);
            font: inherit;
            font-size: 15px;
        }

        .field input::placeholder {
            color: #6f7a78;
        }

        .field svg {
            position: absolute;
            left: 18px;
            top: 19px;
            width: 19px;
            height: 19px;
            stroke: var(--mint);
            fill: none;
            stroke-width: 1.6;
        }

        .go {
            width: 195px;
            height: 57px;
            background: var(--green);
            color: #fff;
            border: 0;
            border-radius: 8px;
            font: inherit;
            font-size: 14px;
            font-weight: 700;
            display: flex;
            align-items: center;
            justify-content: space-between;
            padding: 0 28px;
            cursor: pointer;
        }

        .go:hover {
            background: #0c7c65;
        }

        .go:disabled {
            opacity: 0.45;
            cursor: not-allowed;
        }

        .hints {
            display: flex;
            justify-content: space-between;
            margin-top: 12px;
            font-size: 11px;
            color: var(--muted);
        }

        .list {
            position: absolute;
            left: 0;
            right: 0;
            top: 62px;
            background: #181a18;
            border: 1px solid var(--line);
            border-radius: 8px;
            list-style: none;
            z-index: 5;
            display: none;
            overflow: hidden;
        }

        .list li {
            padding: 12px 20px;
            font-size: 14px;
            cursor: pointer;
        }

        .list li:hover {
            background: #222523;
        }

        .list li.selected {
            background: #222523;
            color: var(--mint);
        }

        .msg {
            margin-top: 14px;
            min-height: 18px;
            font-size: 13px;
            color: var(--mint);
        }

        /* Journey */

        .journey {
            padding: 96px 0 90px;
        }

        .journey .eyebrow {
            margin-bottom: 14px;
        }

        .journey .eyebrow::before {
            display: none;
        }

        h3 {
            font-size: 30px;
            letter-spacing: -0.04em;
        }

        .journey p.d {
            margin-top: 14px;
            color: var(--muted);
            font-size: 14px;
        }

        .steps {
            display: grid;
            grid-template-columns: repeat(3, 1fr);
            gap: 28px;
            margin-top: 38px;
        }

        .step {
            border-top: 1px solid var(--line);
            padding-top: 26px;
        }

        .step .no {
            font-size: 11px;
            font-weight: 700;
            color: var(--mint);
        }

        .step h4 {
            margin: 14px 0 8px;
            font-size: 17px;
            letter-spacing: -0.02em;
        }

        .step p {
            font-size: 13px;
            color: var(--muted);
        }

        .samples {
            margin-top: 44px;
            border: 1px dashed #5a5d5b;
            border-radius: 12px;
            background: #212320;
            padding: 26px 24px;
            display: flex;
            justify-content: space-between;
            align-items: center;
            gap: 20px;
            flex-wrap: wrap;
        }

        .samples h5 {
            font-size: 13px;
        }

        .samples p {
            margin-top: 10px;
            font-size: 12px;
            color: var(--muted);
        }

        .chips {
            display: flex;
            gap: 10px;
        }

        .chips a {
            background: var(--panel);
            border: 1px solid var(--line);
            border-radius: 6px;
            padding: 11px 14px;
            font-size: 12px;
            display: flex;
            gap: 14px;
            align-items: center;
        }

        .chips a::after {
            content: "\2197";
            color: var(--mint);
            font-size: 11px;
        }

        .chips a:hover {
            border-color: var(--mint);
        }

        @media (max-width: 860px) {
            .hero {
                grid-template-columns: 1fr;
            }

            .stage {
                max-width: 340px;
                width: 100%;
            }

            .steps {
                grid-template-columns: 1fr;
            }

            .row {
                flex-direction: column;
            }

            .go {
                width: 100%;
            }

            nav a:not(.on):not(.pill) {
                display: none;
            }
        }

        @media (prefers-reduced-motion: reduce) {
            html {
                scroll-behavior: auto;
            }
        }
    </style>
</head>

<body>

<header>
    <div class="wrap">
        <a class="logo" href="/">
            <span class="mark">e</span>
            <span><b>event</b><span>explorer</span></span>
        </a>

        <nav>
            <a class="on" href="/">Discover</a>
            <a href="#journey">How it works</a>
            <span class="pill">Live events</span>
        </nav>
    </div>
</header>

<div class="strip">
    <b>Event Explorer</b> / Discover music and sports in your city.
</div>

<main class="wrap">

    <section class="hero">

        <div>
            <div class="eyebrow">LESS SCROLLING. MORE GOING.</div>

            <h1>
                A city of possibilities.<br>
                <em>Find your next one.</em>
            </h1>

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

        <div class="stage" aria-hidden="true">
            <div class="poster">
                <small class="t">THE CITY IS CALLING</small>
                <small class="n">01 / 02</small>

                <div class="big">
                    <i>GO</i>
                    <u>OUT.</u>
                </div>

                <small class="b">
                    GOOD PLANS.<br>
                    GREAT MEMORIES.
                </small>

                <span class="arr">&#8599;</span>
            </div>
        </div>

    </section>

    <section class="search">

        <h2>Where are we going?</h2>
        <span class="sub">Start with a city, then explore what is on.</span>

        <label for="city">Choose a city</label>

        <div class="row">

            <div class="field">

                <svg viewBox="0 0 24 24">
                    <circle cx="12" cy="12" r="8"></circle>
                    <circle cx="12" cy="12" r="3"></circle>
                </svg>

                <input
                    id="city"
                    type="text"
                    placeholder="Search a city, e.g. Toronto"
                    autocomplete="off"
                >

                <ul class="list" id="list"></ul>

            </div>

            <button class="go" id="go" disabled>
                Explore events
                <span>&rarr;</span>
            </button>

        </div>

        <div class="hints">
            <span>Type at least 3 characters and select a suggestion.</span>
            <span>Powered by Google Places</span>
        </div>

        <div class="msg" id="msg" role="status"></div>

    </section>

</main>

<section class="journey" id="journey">

    <div class="wrap">

        <div class="eyebrow">A SMALL PROJECT. THE COMPLETE JOURNEY.</div>

        <h3>From a city to a ticket.</h3>

        <p class="d">
            Choose a city and discover music and sports happening there.
        </p>

        <div class="steps">

            <div class="step">
                <div class="no">01</div>
                <h4>Pick a place</h4>
                <p>Find a city with autocomplete.</p>
            </div>

            <div class="step">
                <div class="no">02</div>
                <h4>Find your event</h4>
                <p>Browse music and sports together.</p>
            </div>

            <div class="step">
                <div class="no">03</div>
                <h4>View the details</h4>
                <p>Follow the event through to tickets.</p>
            </div>

        </div>

    </div>

</section>

<script>
    const input = document.getElementById("city");
    const list = document.getElementById("list");
    const message = document.getElementById("msg");
    const exploreButton = document.getElementById("go");

    let sessionToken = createSessionToken();
    let selectedLocation = null;
    let searchTimeout = null;

    function createSessionToken() {
        if (window.crypto && crypto.randomUUID) {
            return crypto.randomUUID();
        }

        return "session-" + Date.now() + "-" + Math.random().toString(36).slice(2);
    }

    function clearSelection() {
        selectedLocation = null;
        exploreButton.disabled = true;
    }

    function showMessage(text) {
        message.textContent = text;
    }

    function clearResults() {
        list.innerHTML = "";
        list.style.display = "none";
    }

    function renderSuggestions(suggestions) {
        list.innerHTML = "";

        if (suggestions.length === 0) {
            list.style.display = "none";
            showMessage("No sample cities match.");
            return;
        }

        suggestions.forEach((suggestion) => {
            const item = document.createElement("li");

            item.textContent = suggestion.text;
            item.dataset.placeId = suggestion.placeId;

            item.addEventListener("click", () => {
                selectSuggestion(suggestion);
            });

            list.appendChild(item);
        });

        list.style.display = "block";
        showMessage("");
    }

    async function searchCities() {
        const query = input.value.trim();

        clearSelection();
        clearResults();

        if (query.length < 3) {
            showMessage("");
            return;
        }

        showMessage("Searching matching cities.");

        try {
            const response = await fetch(
                "/api/locations/autocomplete?input=" +
                encodeURIComponent(query) +
                "&sessionToken=" +
                encodeURIComponent(sessionToken)
            );

            if (!response.ok) {
                throw new Error("Unable to search cities.");
            }

            const data = await response.json();

            renderSuggestions(data.suggestions || []);
        } catch (error) {
            clearResults();
            showMessage("Unable to search cities. Please try again.");
        }
    }

    async function selectSuggestion(suggestion) {
        input.value = suggestion.text;
        clearResults();

        showMessage("Searching matching cities.");

        try {
            const response = await fetch(
                "/api/locations/" +
                encodeURIComponent(suggestion.placeId) +
                "?sessionToken=" +
                encodeURIComponent(sessionToken)
            );

            if (!response.ok) {
                throw new Error("Unable to select city.");
            }

            const location = await response.json();

            selectedLocation = {
                placeId: suggestion.placeId,
                city: location.city,
                countryCode: location.countryCode
            };

            exploreButton.disabled = false;

            showMessage(
                location.city + ", " + location.countryCode + " selected."
            );
        } catch (error) {
            clearSelection();
            showMessage("Unable to select this city. Please try again.");
        }
    }

    input.addEventListener("input", () => {
        clearSelection();

        clearTimeout(searchTimeout);

        searchTimeout = setTimeout(() => {
            searchCities();
        }, 300);
    });

    exploreButton.addEventListener("click", () => {
        if (!selectedLocation) {
            return;
        }

        const params = new URLSearchParams({
            city: selectedLocation.city,
            countryCode: selectedLocation.countryCode
        });

        window.location.href = "/events?" + params.toString();
    });
</script>

</body>
</html>