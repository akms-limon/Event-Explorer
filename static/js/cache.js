const clearFullCacheButton =
  document.getElementById("clear-full-cache");

const clearSelectedCacheButton =
  document.getElementById("clear-selected-cache");

const citySelect =
  document.getElementById("cache-city");

const countrySelect =
  document.getElementById("cache-country");

const categorySelect =
  document.getElementById("cache-category");

const cacheMessage =
  document.getElementById("cache-message");

function showCacheMessage(message) {
  cacheMessage.textContent = message;
}

function updateClearButton() {
  clearSelectedCacheButton.disabled =
    !citySelect.value ||
    !countrySelect.value ||
    !categorySelect.value;
}

citySelect.addEventListener("change", updateClearButton);
countrySelect.addEventListener("change", updateClearButton);
categorySelect.addEventListener("change", updateClearButton);

clearFullCacheButton.addEventListener(
  "click",
  async () => {
    clearFullCacheButton.disabled = true;

    try {
      const response = await fetch(
        "/cache/invalidate",
        {
          method: "POST"
        }
      );

      if (!response.ok) {
        throw new Error();
      }

      window.location.reload();
    } catch (error) {
      clearFullCacheButton.disabled = false;

      showCacheMessage(
        "Unable to clear cache."
      );
    }
  }
);

clearSelectedCacheButton.addEventListener(
  "click",
  async () => {
    const city = citySelect.value;
    const country = countrySelect.value;
    const category = categorySelect.value;

    if (!city || !country || !category) {
      return;
    }

    clearSelectedCacheButton.disabled = true;

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
        throw new Error();
      }

      window.location.reload();
    } catch (error) {
      clearSelectedCacheButton.disabled = false;

      showCacheMessage(
        "Unable to clear selected cache."
      );
    }
  }
);