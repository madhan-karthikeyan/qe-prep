# SprintPlanner — Agile Capacity Planning Tool

## Overview

SprintPlanner is a web-based capacity planning tool that helps engineering teams estimate how much work they can commit to in a sprint. It combines historical velocity analysis, team availability modeling, and dependency tracking to produce realistic sprint forecasts. Teams configure their members, schedules, and past sprint data, and the tool outputs recommended point loads, risk flags, and what-if scenarios for different sprint compositions.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Frontend (React SPA)                      │
│  Sprint Board │ Capacity View │ What-If Planner │ Reports    │
└──────────────────────────┬──────────────────────────────────┘
                           │ REST (OpenAPI)
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    API Gateway (FastAPI)                      │
│  Auth │ Rate Limit │ Request Validation │ Caching            │
└──┬──────────────┬──────────────┬───────────────────────────┘
   │              │              │
   ▼              ▼              ▼
┌──────────┐ ┌──────────┐ ┌──────────────────┐
│ Sprint   │ │ Team     │ │ History          │
│ Service  │ │ Service  │ │ Service          │
└────┬─────┘ └────┬─────┘ └───────┬──────────┘
     │            │               │
     ▼            ▼               ▼
┌─────────────────────────────────────────────────────────────┐
│                     PostgreSQL                                │
│  sprints │ team_members │ velocity_log │ estimations          │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                   ML Estimator (sidecar)                      │
│  Historical Velocity → Feature Engineering → Model Inference │
└─────────────────────────────────────────────────────────────┘
```

## Why This Design

The frontend/backend separation allows the API to be consumed by other tools (Slack bots, CI pipelines, Jira plugins) without going through a web UI. The ML estimator runs as a separate sidecar process because training runs can be CPU-intensive and should not block API requests. PostgreSQL was chosen over a document store because sprint data is highly relational — teams, members, sprints, stories, and velocity history have complex join patterns that document databases handle poorly. The OpenAPI spec ensures the frontend and any third-party integrations have a clear contract.

## Key Technical Decisions

- **RESTful API with OpenAPI spec**: Every endpoint is documented with request/response schemas, examples, and error codes. The frontend generates its TypeScript types from the OpenAPI spec via `openapi-typescript`, eliminating hand-written type definitions and catching contract breaks in CI before they reach production. The spec serves as the single source of truth for the API surface.

- **Historical velocity analysis for estimation**: Rather than asking teams to estimate every story from scratch, SprintPlanner computes a team's historical velocity (points completed per sprint over the last 8-12 sprints) and uses it as a baseline. The estimator accounts for sprint length variations, holidays, and team composition changes. A rolling window with exponential decay weights recent sprints more heavily, so a team that recently improved (or degraded) is reflected in the forecast within 2-3 sprints.

- **Team availability modeling with calendar integration**: A team member at 100% allocation for a 2-week sprint is not actually available for 80 hours — meetings, code reviews, interviews, and context switching reduce effective capacity. SprintPlanner integrates with Google Calendar and Outlook to pull events and compute a real availability percentage per person per day. Users can also set planned PTO and part-time schedules. The model supports what-if scenarios like "what if Alice takes parental leave for 2 sprints" or "what if Bob works 4-day weeks."

## Testing Strategy

- **API contract tests**: Using `pytest` with `schemathesis`, every API endpoint is tested against the OpenAPI spec. Tests verify: all required fields return 400 on missing, all responses conform to the response schema, pagination works correctly, and enum values are validated. The contract tests run against a test PostgreSQL instance populated with fixture data. CI blocks merges if the API violates its own spec.

- **Database migration tests**: Each migration is tested by running it against a copy of production schema (anonymized), executing the migration, running the rollback, and verifying schema state matches the expected version. This catches destructive changes (dropped columns, type changes) before they reach staging. Migration tests run in a transaction that is rolled back, so they leave no side effects.

- **ML model evaluation**: The velocity estimator's accuracy is tracked per sprint. Metrics: Mean Absolute Error (MAE) — average absolute difference between predicted and actual points completed; Root Mean Squared Error (RMSE) — penalizes large errors more heavily; and R² — how much variance the model explains. A dashboard shows these metrics over time. If MAE exceeds 30% for three consecutive sprints, an alert fires and the model retrains on the expanded dataset.

- **A/B testing different estimation algorithms**: Teams can opt into an experiment where 50% of sprints use the default model (simple moving average) and 50% use an alternative (e.g., ARIMA or gradient boosting). The API returns a `model_version` field in the estimation response so that downstream analysis can compare performance. Experiments run for at least 8 sprints before declaring a winner.

## Failures & Lessons

- **ML features too complex**: The first version of the estimator used 40+ features: day of week, sprint number in quarter, team mood survey results, PR merge latency, and more. Model accuracy was worse than a simple moving average, and debugging was impossible. Scaled back to a 5-feature model: rolling average velocity, sprint length, team size, avg story complexity, and holiday factor. The simple model beat the complex one in cross-validation and was comprehensible to the product team. Added the complex features back incrementally, each validated with a standalone A/B test.

- **Timezone handling caused date errors**: Sprint start/end dates were stored as `TIMESTAMP WITH TIME ZONE` but the frontend sent them as local time strings without offsets. A team in Sydney overlapped with a team in San Francisco, causing sprints to appear to start a day early or late. Fix: all dates are stored as UTC in the database, the API accepts and returns only ISO 8601 strings with explicit timezone offsets, and the frontend converts to local time only for display. Added a middleware layer that rejects any datetime string without a timezone.

## Tradeoffs

- **ML accuracy vs interpretability**: A gradient boosting model achieves 15% better MAE than a simple moving average but cannot explain why a particular sprint prediction was high or low. The moving average is transparent: "your last 8 sprints averaged 42 points, so we predict 42 ± 5." SprintPlanner defaults to the interpretable model and allows teams to opt into ML models if they accept the opacity tradeoff.

- **API flexibility vs strictness**: A flexible API (JSON fields, dynamic attributes) makes it easy to add new features without version bumps, but clients cannot rely on field presence or types. SprintPlanner uses strict schemas with required/optional clearly marked. Adding a new field requires a minor version bump and a migration period where the field is optional. After two releases, it can be made required. This process is documented and automated via the OpenAPI spec and CI schema checks.

## Interview Questions

**Q: How do you estimate story points using historical data?**
A: We use a rolling weighted average of completed points per sprint over the last 8-12 sprints, with exponential decay (weight halved every 4 sprints). The estimate is adjusted for: sprint length (normalized to standard length), team composition changes (new members get a 0.5 weight multiplier for their first 2 sprints), and known time off (subtracted from available person-hours before converting back to points). The output is a range (P10, P50, P90) derived from the distribution of historical velocities, not a single point estimate. This gives the team a realistic worst-case, expected, and best-case scenario.

**Q: How would you handle team members in different timezones?**
A: Availability is computed per person in their local timezone. The calendar integration fetches events in the user's timezone and computes busy blocks. The API stores all event data in UTC but preserves the original timezone for display. SprintPlanner's availability overlap calculation maps every team member's working hours into UTC and finds common windows. For teams with no overlap (e.g., distributed across 12 timezones), the model assumes async communication and reduces effective capacity by an async overhead factor (default 20% overhead per timezone boundary crossed). The model also surfaces a "timezone risk score" that flags sprints where low overlap may cause delays.

**Q: What's your API versioning strategy?**
A: URL-based versioning (`/v1/sprints`, `/v2/sprints`). A new API version is created when a breaking change is introduced (field removal, type change, required → optional in the breaking direction). Backward-compatible changes (new fields, new endpoints) are made on the current version. Each API version is supported for at least 6 months after a replacement is released. Deprecated versions return a `Sunset` header with the removal date. The OpenAPI spec for each version is published and independently testable. Internal clients (frontend, Slack bot) are updated within one sprint of a new release; external API consumers get the 6-month deprecation window.

## Related Problems from This Repository

- **REST API Tester**: The API contract test framework was built on top of patterns from the REST API Tester problem. The assertion helpers for validating response structure against OpenAPI schemas are direct adaptations.

- **Database CRUD**: All SprintPlanner database operations (CRUD for teams, sprints, members, velocity logs) follow the repository pattern established in the Database CRUD problem. The migration testing framework was also adapted from it.

- **URL Parser**: The calendar integration module uses the URL Parser's normalized URL handling for constructing OAuth redirect URIs and API endpoints for Google Calendar and Outlook.

- **Expression Evaluator**: The what-if scenario engine uses the Expression Evaluator's tree-walking interpreter to parse user-defined constraints like "Alice.points < 20 OR Bob.availability < 0.5" and evaluate them against proposed sprint compositions.
