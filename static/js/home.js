const cityInput = document.getElementById("city");
const suggestionsList = document.getElementById("list");
const message = document.getElementById("msg");
const exploreButton = document.getElementById("go");

let sessionToken = createSessionToken();
let selectedLocation = null;
let debounceTimer = null;

function createSessionToken() {
  if (window.crypto && crypto.randomUUID) {
    return crypto.randomUUID();
  }

  return "session-" + Date.now();
}

function resetSuggestions() {
  suggestionsList.innerHTML = "";
  selectedLocation = null;
  exploreButton.disabled = true;
}

function showMessage(text) {
  message.textContent = text;
}

async function searchCities() {
  const input = cityInput.value.trim();

  resetSuggestions();

  if (input.length < 3) {
    showMessage("");
    return;
  }

  showMessage("Searching matching cities.");

  try {
    const params = new URLSearchParams({
      input,
      sessionToken
    });

    const response = await fetch(
      "/api/locations/autocomplete?" + params.toString()
    );

    if (!response.ok) {
      throw new Error("Autocomplete request failed");
    }

    const data = await response.json();
    const suggestions = data.suggestions || [];

    if (suggestions.length === 0) {
      showMessage("No sample cities match.");
      return;
    }

    showMessage("");
    renderSuggestions(suggestions);
  } catch (error) {
    showMessage("Unable to search cities.");
  }
}

function renderSuggestions(suggestions) {
  suggestionsList.innerHTML = "";

  suggestions.forEach((suggestion) => {
    const item = document.createElement("li");
    const button = document.createElement("button");

    button.type = "button";
    button.className = "suggestion";
    button.textContent = suggestion.text;

    button.addEventListener("click", () => {
      selectSuggestion(suggestion);
    });

    item.appendChild(button);
    suggestionsList.appendChild(item);
  });
}

async function selectSuggestion(suggestion) {
  showMessage("Loading selected city.");
  suggestionsList.innerHTML = "";

  try {
    const response = await fetch(
      "/api/locations/" +
        encodeURIComponent(suggestion.placeId) +
        "?sessionToken=" +
        encodeURIComponent(sessionToken)
    );

    if (!response.ok) {
      throw new Error("Place lookup failed");
    }

    const location = await response.json();

    selectedLocation = {
      placeId: suggestion.placeId,
      city: location.city,
      countryCode: location.countryCode
    };

    cityInput.value = suggestion.text;
    showMessage(
      selectedLocation.city +
        ", " +
        selectedLocation.countryCode +
        " selected."
    );

    exploreButton.disabled = false;

    // The current token belongs to this autocomplete/selection session.
    // The next search gets a new session token.
    sessionToken = createSessionToken();
  } catch (error) {
    selectedLocation = null;
    exploreButton.disabled = true;
    showMessage("Unable to select this city.");
  }
}

cityInput.addEventListener("input", () => {
  clearTimeout(debounceTimer);

  if (selectedLocation) {
    sessionToken = createSessionToken();
  }

  debounceTimer = setTimeout(searchCities, 300);
});

exploreButton.addEventListener("click", () => {
  if (!selectedLocation) {
    showMessage("Select a valid city suggestion.");
    return;
  }

  const params = new URLSearchParams({
    city: selectedLocation.city,
    countryCode: selectedLocation.countryCode
  });

  window.location.href = "/events?" + params.toString();
});