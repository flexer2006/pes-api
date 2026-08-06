# Person Enrichment Service

An API for storing and enriching person data (age, gender, and nationality).

**Note!** This is an educational project intended to practice Go backend development, REST API design, database integration, and working with external services. It is not a production-ready service.

## Running

```sh
git clone https://github.com/flexer2006/pes-api.git
cd simple-person-enrichment-go
cp .env.example .env
cd deploy && docker-compose up -d
```

Available at: `http://localhost`. Swagger UI: `http://localhost/swagger/swagger.html`.

## API (base `/api/v1`)

* `GET /persons` – list persons (filters: `name`, `surname`, `patronymic`, `gender`, `nationality`, `age`; `limit`/`offset`)
* `GET /persons/:id` – get a person
* `POST /persons` – create a person
* `PUT|PATCH /persons/:id` – update a person
* `DELETE /persons/:id` – delete a person
* `POST /persons/:id/enrich` – enrich with age, gender, and nationality

### API Usage Examples

#### 1. List persons

```bash
curl -X GET "http://localhost/api/v1/persons"
```

![alt text](docs/photo/persons-get.png)

#### 1.1. With filters

```bash
curl -X GET "http://localhost/api/v1/persons?limit=5&offset=0&name=Ivan&gender=male"
```

![alt text](docs/photo/persons-filters.png)

#### 2. Create

```bash
curl -X POST "http://localhost/api/v1/persons" \
  -H "Content-Type: application/json" \
  -d '{"name":"Dmitry","surname":"Ushakov","patronymic":"Vasilievich"}'
```

![alt text](docs/photo/persons-post.png)

#### 3. Get by ID

```bash
curl -X GET "http://localhost/api/v1/persons/{UID}"
```

![alt text](docs/photo/persons-id-get.png)

#### 4. Enrich

```bash
curl -X POST "http://localhost/api/v1/persons/{UID}/enrich"
```

![alt text](docs/photo/persons-id-enrich-post.png)

#### 5. Update

```bash
curl -X PUT "http://localhost/api/v1/persons/{UID}" \
  -H "Content-Type: application/json" \
  -d '{"name":"Dmitry","surname":"Ushakov","patronymic":"Alexeevich"}'
```

![alt text](docs/photo/persons-id-put.png)

#### 5.1. Partial update

```bash
curl -X PATCH "http://localhost/api/v1/persons/{UID}" \
  -H "Content-Type: application/json" \
  -d '{"surname":"Ushakov"}'
```

![alt text](docs/photo/persons-id-patch.png)

#### 6. Delete

```bash
curl -X DELETE "http://localhost/api/v1/persons/{UID}"
```

![alt text](docs/photo/persons-id-delete.png)

## External Services

* agify.io – predicts age by name
* genderize.io – predicts gender by name
* nationalize.io – predicts nationality by name

## Database

Table `people`:

| Column                    | Type                     | Description                           |
| ------------------------- | ------------------------ | ------------------------------------- |
| `id`                      | UUID                     | Primary key, unique person identifier |
| `name`                    | VARCHAR(100)             | First name (required)                 |
| `surname`                 | VARCHAR(100)             | Last name (required)                  |
| `patronymic`              | VARCHAR(100)             | Patronymic (optional)                 |
| `age`                     | INTEGER                  | Age                                   |
| `gender`                  | VARCHAR(10)              | Gender                                |
| `gender_probability`      | DECIMAL(5,4)             | Gender prediction probability         |
| `nationality`             | VARCHAR(2)               | Country code (nationality)            |
| `nationality_probability` | DECIMAL(5,4)             | Nationality prediction probability    |
| `created_at`              | TIMESTAMP WITH TIME ZONE | Record creation timestamp             |
| `updated_at`              | TIMESTAMP WITH TIME ZONE | Record last update timestamp          |

Migrations in the `migrations` directory are applied automatically on startup.
