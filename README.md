# Event Explorer

Event Explorer is a web application for searching cities and discovering Music and Sports events. It is built with Go and Beego, with HTML, CSS and vanilla JavaScript on the frontend.

Event data comes from the Ticketmaster Discovery API, and city search is powered by the Google Places API.

## Setup

**1. Clone the repository**

```bash
git clone https://github.com/akms-limon/Event-Explorer
cd Event-Explorer
````

**2. Install dependencies**

```bash
go mod download
```

**3. Configure environment variables**

Copy the example file and add your API keys:

```bash
cp .env.example .env
```

```env
GOOGLE_PLACES_API_KEY=your_google_places_api_key
TICKETMASTER_API_KEY=your_ticketmaster_api_key
```

> Never commit your `.env` file.

**4. Run the application**

```bash
bee run
```

Then open [http://localhost:8080](http://localhost:8080).

## Routes

**Pages**

| Method | Route                                 | Description          |
| ------ | ------------------------------------- | -------------------- |
| GET    | `/`                                   | Home and city search |
| GET    | `/events?city=Toronto&countryCode=CA` | Event listing        |
| GET    | `/events/:eventId`                    | Event details        |
| GET    | `/redirect/:eventId`                  | Ticket redirect      |

**API**

| Method | Route                                                    | Description           |
| ------ | -------------------------------------------------------- | --------------------- |
| GET    | `/api/locations/autocomplete?input=Tor&sessionToken=...` | City suggestions      |
| GET    | `/api/locations/:placeId?sessionToken=...`               | Selected city details |
| GET    | `/api/cache/locations`                                   | Search cached cities  |

**Cache clearing**

| Method | Route                                        | Description                           |
| ------ | -------------------------------------------- | ------------------------------------- |
| POST   | `/cache/invalidate`                          | Clear all cached events               |
| POST   | `/cache/invalidate/:city/:country/:category` | Clear cache for one city and category |

## External APIs

* **Google Places API** provides city autocomplete and selected-city lookup.
* **Ticketmaster Discovery API** provides Music and Sports listings, event details and ticket URLs.

## Testing

Unit tests use mocked API responses, so no live API keys are needed.

Current test coverage: **91.4%**

```bash
go test ./...                            # run all tests
go test -race ./...                      # run with race detection
go test ./... -cover                     # show package coverage

go test ./... -coverprofile=coverage.out # generate coverage profile
go tool cover -func=coverage.out         # show overall coverage
