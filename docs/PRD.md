# Product Requirements Document (PRD)

| Field | Value |
|---|---|
| Product | Pragati (working name) |
| Version | 0.1 (MVP → Day 360) |
| Owner | Founder / Product |
| Related | [BRD](BRD.md), [Blueprint](02-product-blueprint.md) |

## 1. Problem statement
"I'm considering buying in [project/locality]. What is it worth today, how
much will it appreciate over the next few years because of what's being built
around it, and how sure can I be?" No product in India answers this credibly
at building level.

## 2. Goals & non-goals

**Goals**
- G1: Answer "value today + 1–5 year outlook + why" for any covered property in under 10 seconds.
- G2: Make every number traceable to data and sources.
- G3: Convert intent into revenue (reports, Pro, B2B).

**Non-goals (v1)**
- Brokerage, listings marketplace, home-loan origination.
- Point-precise "guaranteed" predictions.
- Native mobile apps.

## 3. Personas
| Persona | Description | Top job-to-be-done |
|---|---|---|
| Ananya, first-time buyer (32, Pune, IT) | Choosing between 3 projects, ₹80 L–1.2 Cr budget | "Which option holds and grows its value best?" |
| Rajesh, investor (45, business owner) | Owns 3 properties, wants the next growth corridor | "Where should I buy land/flats for 5-year appreciation?" |
| Farah, NRI (38, Dubai) | Buying remotely in Hyderabad | "Can I trust this project and this area from abroad?" |
| Vikram, channel partner | Advises 200+ clients per year | "Give me credible, branded reports to share" |
| Meera, land-acquisition head at a developer | Evaluates 20 land parcels per quarter | "Which parcel will see the highest demand/price growth?" |
| Arjun, risk head at an HFC | 40k home loans on the book | "Which collateral is at risk of value decline? What's a fair forward LTV?" |

## 4. User stories & requirements

Priority: **P0** = Day 1/MVP, **P1** = Day 90, **P2** = Day 360.

### 4.1 Search & discovery
| ID | User story | Acceptance criteria | P |
|---|---|---|---|
| S-1 | As a user, I can search by project name, RERA ID, locality, or address | Autocomplete < 300 ms; fuzzy matching; top-5 results include the right entity for ≥ 95% of test queries | P0 |
| S-2 | I can drop a pin on the map and get a locality/grid report | Pin → H3 cell → report in < 3 s | P0 |
| S-3 | I see a "hot catalysts" feed for my city | Shows the 10 latest stage changes with dates and sources | P1 |

### 4.2 Property / locality report
| ID | User story | Acceptance criteria | P |
|---|---|---|---|
| R-1 | I see the estimated current price/sq ft with a range | Shows P50 and P10–P90 range; number of comparable transactions; data recency | P0 (locality), P1 (building) |
| R-2 | I see a 1/3/5-year appreciation outlook with bear/base/bull | Fan chart; % and ₹/sq ft; confidence level (5-point) with explanation | P0 (rule-based range), P1 (model) |
| R-3 | I see the drivers behind the forecast | List of catalysts + market trend + supply + risk, each with signed contribution %, probability and expected date; contributions sum to total forecast | P1 |
| R-4 | I can open the source of each catalyst | Each driver links to ≥ 1 source document (URL + archived copy + fetch date) | P0 |
| R-5 | I see risk flags | RERA complaints count, project delay vs RERA timeline, litigation (where sourced), flood/CRZ overlays | P1 |
| R-6 | I see comparable nearby projects | 5 comparables with price and forecast | P1 |
| R-7 | I see historical accuracy for this city | Back-test metric displayed with methodology link | P1 |
| R-8 | Free users see a teaser; paid users see full detail | Free: today's range + outlook band + top 2 drivers. Paid: everything + PDF | P0 |

### 4.3 Catalyst explorer
| ID | User story | Acceptance criteria | P |
|---|---|---|---|
| C-1 | I can view all catalysts on a map, filterable by type and stage | Geometry (point/line/polygon) rendered; filter by type, stage, expected completion year | P0 |
| C-2 | I can open a catalyst detail page | Timeline of stage changes, budget, executing agency, sources, delay history, affected localities | P0 |
| C-3 | I can follow a catalyst | Alerts on stage change | P1 |

### 4.4 Reports, compare, portfolio, alerts
| ID | User story | Acceptance criteria | P |
|---|---|---|---|
| P-1 | I can buy and download a PDF report | Razorpay checkout (UPI/cards); PDF generated < 60 s (MVP: analyst-assisted within 24 h); GST invoice emailed | P0 |
| P-2 | I can subscribe to Pro | Monthly/annual plans; unlimited reports within fair use; cancel any time | P1 |
| P-3 | I can compare up to 4 properties | Side-by-side table + overlay chart | P1 |
| P-4 | I can add owned properties to a portfolio | Value + forecast tracking; monthly email summary | P2 |
| P-5 | I get alerts via email/WhatsApp | Opt-in per channel (DPDP consent); delivered within 24 h of the change | P1 |
| P-6 | Broker can white-label reports | Logo, contact details, client name; usage counted per seat | P1 |

### 4.5 B2B console & API
| ID | User story | Acceptance criteria | P |
|---|---|---|---|
| B-1 | Lender uploads a CSV of addresses and gets valuations + forecasts + risk flags | Up to 50k rows; match rate reported; results in < 1 h | P1 (pilot), P2 (GA) |
| B-2 | Developer views an appreciation heatmap for a region | H3 grid heatmap, filters by horizon and segment, export | P2 |
| B-3 | Partner integrates via REST API | API keys, rate limits, usage metering, OpenAPI docs, 99.5% uptime | P2 |

### 4.6 Internal: data & analyst tools
| ID | Requirement | Acceptance criteria | P |
|---|---|---|---|
| D-1 | Crawlers fetch sources on schedule and archive raw files | Per-source schedule; hashing; diff detection; failure alerts | P0 (top sources), P1 (all) |
| D-2 | Document-AI extracts catalyst fields from PDFs/HTML | Fields: type, name, agency, geometry hint, budget, stage, dates; ≥ 85% field accuracy on a labelled set; confidence score per field | P1 |
| D-3 | Analyst review queue | Approve/edit/reject; side-by-side source view; audit log | P0 |
| D-4 | Entity resolution of projects/addresses | ≥ 95% precision on a labelled set | P1 |
| D-5 | Manual override with reason | Overrides logged and shown in model governance reports | P1 |

## 5. Forecasting model requirements
| ID | Requirement |
|---|---|
| M-1 | **AVM**: hedonic + gradient boosting on registrations (primary) and listings (secondary, discounted for asking-price premium). Output P10/P50/P90. Target median absolute % error ≤ 10% in launch city |
| M-2 | **Catalyst uplift model**: event-study / difference-in-differences on past projects (treatment = within X km of catalyst, control = matched areas). Uplift curve by catalyst type × distance band × city tier × lifecycle stage. Historical uplift is typically partly priced in at announcement; model the timing (announcement, construction start, opening) separately |
| M-3 | **Completion probability**: per catalyst, by agency track record, stage, funding status, land acquisition %, litigation. Output a distribution of completion dates |
| M-4 | **Market baseline**: city/segment trend from macro (rates, absorption, inventory months, income growth) |
| M-5 | **Supply overhang**: penalty from RERA pipeline vs absorption rate in the micro-market |
| M-6 | **Simulation**: Monte Carlo over completion dates and uplift uncertainty → fan chart |
| M-7 | **Explainability**: additive decomposition of the forecast by driver |
| M-8 | **Back-testing**: rolling-origin back-tests (e.g. forecast from 2016, 2018, 2020 and compare with realised prices); publish hit-rate within ±10% and MAPE |
| M-9 | **Governance**: versioned models, reproducible training, drift monitoring, monthly re-scoring |

## 6. Data model (core entities)
- `Property` (id, type, RERA id, developer, geometry, H3 cell, attributes)
- `Locality` / `MicroMarket` (polygon, city)
- `Transaction` (property/locality, date, value, area, ₹/sq ft, source)
- `Listing` (asking price, date, source)
- `Catalyst` (type, name, agency, geometry, budget, current stage, expected completion distribution)
- `CatalystStageEvent` (catalyst, stage, date, source document)
- `SourceDocument` (url, archived file, hash, fetched_at, license tag)
- `Forecast` (entity, model version, horizon, P10/P50/P90, drivers[])
- `User`, `Subscription`, `Report`, `Alert`, `Portfolio`, `Organisation`, `ApiKey`

## 7. Non-functional requirements
| Area | Requirement |
|---|---|
| Performance | Report page TTFB < 1 s (precomputed forecasts); search < 300 ms |
| Availability | 99.5% (B2C), 99.9% (API, Day 360) |
| Security | OWASP top-10, encryption at rest/in transit, SSO for B2B, role-based access |
| Privacy | DPDP Act 2023: consent, purpose limitation, data-principal rights, India data residency |
| Accessibility | WCAG 2.1 AA; English first, then Hindi and regional languages |
| SEO | Server-rendered pages per project/locality/catalyst; structured data |
| Observability | Logs, metrics, tracing; data-freshness dashboards per source |

## 8. Legal & content requirements
- Every forecast shows: "Estimates are for information only and are not investment advice or a guarantee of future value."
- Never use "guaranteed returns" language; show ranges and confidence.
- Respect source terms; maintain a licensing register; attribute government sources.
- Risk flags must be factual and sourced (no subjective allegations about developers).
- Get legal review on whether any feature falls under SEBI investment-adviser rules or RERA agent-registration requirements (e.g. if brokerage is added later).

## 9. Analytics & experiment plan
- Funnel: search → report view → paywall → purchase → repeat.
- Experiments: free-tier depth, price points (₹499 vs ₹999), report-vs-subscription framing, WhatsApp alerts vs email.
- Track trust signals: source-link clicks, accuracy-page views, dispute submissions.

## 10. Release plan
| Release | Contents |
|---|---|
| **MVP (Day 1–30)** | S-1, S-2, R-1 (locality), R-2 (rule-based), R-4, R-8, C-1, C-2, P-1 (analyst-assisted), D-1 (top sources), D-3 |
| **v1 (Day 90)** | Building-level R-1/R-2 with models M-1, M-2, M-4, M-7, M-8; R-3, R-5, R-6, R-7; S-3; C-3; P-2, P-3, P-5, P-6; D-2, D-4, D-5; B-1 pilot |
| **v2 (Day 360)** | Multi-city; M-3, M-5, M-6, M-9; P-4; B-1 GA, B-2, B-3; plots/land |

## 11. Open questions
1. Launch city: MMR (best data, most competition) vs Pune/Hyderabad (strong infrastructure pipeline, less noise)?
2. License transaction data from a vendor at the start, or build our own from IGR?
3. How deep should the free tier be to drive SEO and virality without killing conversion?
4. Should B2B (lenders) be prioritised earlier for revenue, at the cost of B2C brand building?
5. Which regional languages first?
