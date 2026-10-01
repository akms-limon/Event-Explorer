const clearFullCacheButton =
  document.getElementById("clear-full-cache");

const citySearchInput =
  document.getElementById("cache-city-search");

const citySuggestions =
  document.getElementById("cache-city-suggestions");

const cityMessage =
  document.getElementById("cache-city-message");

const countrySelect =
  document.getElementById("cache-country");

const categorySelect =
  document.getElementById("cache-category");

const clearSelectedCacheButton =
  document.getElementById("clear-selected-cache");

const actionMessage =
  document.getElementById("cache-action-message");

let selectedLocation = null;
let searchTimer = null;

function showCityMessage(message) {
  cityMessage.textContent = message;
}

function showActionMessage(message) {
  actionMessage.textContent = message;
}

function clearCitySuggestions() {
  citySuggestions.innerHTML = "";
}

function resetCountrySelect() {
  countrySelect.innerHTML = "";

  const option = document.createElement("option");

  option.value = "";
  option.textContent = "Select country";

  countrySelect.appendChild(option);

  countrySelect.disabled = true;
}

function updateClearButton() {
  clearSelectedCacheButton.disabled =
    !selectedLocation ||
    !countrySelect.value ||
    !categorySelect.value;
}

function resetCitySelection() {
  selectedLocation = null;

  resetCountrySelect();

  updateClearButton();
}

function renderCitySuggestions(locations) {
  clearCitySuggestions();

  locations.forEach((location) => {
    const item = document.createElement("li");

    const button = document.createElement("button");

    button.type = "button";
    button.className = "cache-city-suggestion";

    button.textContent =
      location.city +
      ", " +
      location.countryCode;

    button.addEventListener("click", () => {
      selectCity(location);
    });

    item.appendChild(button);

    citySuggestions.appendChild(item);
  });
}

function selectCity(location) {
  selectedLocation = location;

  citySearchInput.value = location.city;

  clearCitySuggestions();

  showCityMessage(
    location.city +
      " selected."
  );

  resetCountrySelect();

  const option = document.createElement("option");

  option.value = location.countryCode;
  option.textContent = location.countryCode;

  countrySelect.appendChild(option);

  countrySelect.disabled = false;

  updateClearButton();
}

async function searchCachedCities() {
  const search =
    citySearchInput.value.trim();

  resetCitySelection();
  clearCitySuggestions();

  if (search.length === 0) {
    showCityMessage("");
    return;
  }

  if (search.length < 3) {
    showCityMessage(
      "Type at least 3 characters to search."
    );
    return;
  }

  showCityMessage(
    "Searching cached cities..."
  );

  try {
    const params = new URLSearchParams({
      search
    });

    const response = await fetch(
      "/api/cache/locations?" +
        params.toString()
    );

    if (!response.ok) {
      throw new Error(
        "Cached city search failed"
      );
    }

    const data = await response.json();

    const locations =
      data.locations || [];

    if (locations.length === 0) {
      showCityMessage(
        "No matching cached cities found."
      );
      return;
    }

    showCityMessage(
      locations.length +
        " cached city" +
        (locations.length === 1 ? "" : "ies") +
        " found."
    );

    renderCitySuggestions(locations);

  } catch (error) {
    showCityMessage(
      "Unable to search cached cities."
    );
  }
}

citySearchInput.addEventListener(
  "input",
  () => {
    clearTimeout(searchTimer);

    searchTimer = setTimeout(
      searchCachedCities,
      300
    );
  }
);

countrySelect.addEventListener(
  "change",
  () => {
    updateClearButton();
  }
);

categorySelect.addEventListener(
  "change",
  () => {
    updateClearButton();
  }
);

clearFullCacheButton.addEventListener(
  "click",
  async () => {
    clearFullCacheButton.disabled = true;

    showActionMessage(
      "Clearing full cache..."
    );

    try {
      const response = await fetch(
        "/cache/invalidate",
        {
          method: "POST"
        }
      );

      if (!response.ok) {
        throw new Error(
          "Full cache clear failed"
        );
      }

      window.location.reload();

    } catch (error) {
      clearFullCacheButton.disabled = false;

      showActionMessage(
        "Unable to clear full cache."
      );
    }
  }
);

clearSelectedCacheButton.addEventListener(
  "click",
  async () => {
    if (
      !selectedLocation ||
      !countrySelect.value ||
      !categorySelect.value
    ) {
      return;
    }

    clearSelectedCacheButton.disabled = true;

    showActionMessage(
      "Clearing selected cache..."
    );

    const city =
      selectedLocation.city;

    const country =
      countrySelect.value;

    const category =
      categorySelect.value;

    const url =
      "/cache/invalidate/" +
      encodeURIComponent(city) +
      "/" +
      encodeURIComponent(country) +
      "/" +
      encodeURIComponent(category);

    try {
      const response = await fetch(
        url,
        {
          method: "POST"
        }
      );

      if (!response.ok) {
        throw new Error(
          "Selected cache clear failed"
        );
      }

      window.location.reload();

    } catch (error) {
      clearSelectedCacheButton.disabled = false;

      showActionMessage(
        "Unable to clear selected cache."
      );
    }
  }
);