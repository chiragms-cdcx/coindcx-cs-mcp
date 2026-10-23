# 02 — Product Blueprint

Working name: **Pragati** ("progress"). It's a placeholder; check trademark and
domain availability before committing to it.

**One-liner:** *Type any building, project or locality in India and see how
much it is likely to appreciate over the next 1–5 years, which upcoming
developments drive that, and how confident we are.*

---

## 1. What the platform looks like

### 1.1 Home / search
- A single search bar: "Search a building, project, RERA ID, locality or drop a pin".
- City switcher, plus a "Hot catalysts this month" feed (e.g. "Metro Line 3 Phase 2 tender awarded").

### 1.2 Property / area report page (the core screen)

```
┌──────────────────────────────────────────────────────────────────────┐
│  Lodha Amara, Kolshet Road, Thane | RERA P51700012345 | Residential  │
├──────────────────────────────┬───────────────────────────────────────┤
│  TODAY'S ESTIMATED VALUE     │  5-YEAR APPRECIATION OUTLOOK          │
│  Rs 14,200 / sq ft           │  Base case +38%  (Rs 19,600 / sq ft)  │
│  (range 13,600 - 14,900)     │  Bear +18%        Bull +55%           │
│  Based on 46 registrations   │  Confidence: Medium (3 of 5)          │
├──────────────────────────────┴───────────────────────────────────────┤
│  [Chart] Price history 2015-today, forecast fan chart to 2031        │
├──────────────────────────────────────────────────────────────────────┤
│  WHAT DRIVES THIS FORECAST                 impact   probability      │
│  [METRO] Line 4, Kasarvadavali stn, 900 m  +9.5%    78% by 2027      │
│  [ROAD]  Thane Coastal Road, 2.1 km        +4.0%    45% by 2029      │
│  [JOBS]  IT park (MIDC notification), 3 km +2.5%    60%              │
│  [RISK]  Supply overhang: 6,800 new units  -3.0%    -                │
│  [TREND] City-wide market trend            +22%     -                │
│  Each row links to its source (tender, gazette, DPR, RERA filing)    │
├──────────────────────────────────────────────────────────────────────┤
│  MAP: catalysts, timelines, isochrones, flood/zoning overlays        │
├──────────────────────────────────────────────────────────────────────┤
│  RISKS: project litigation, RERA complaints, flood zone, CRZ         │
│  COMPARABLES: 5 nearby projects, their price and forecast            │
│  ACCURACY: "3-yr forecasts in Thane within +/-7% in 71% of cases"    │
└──────────────────────────────────────────────────────────────────────┘
   [Download PDF report]  [Set alert]  [Compare]  [Add to portfolio]
```

### 1.3 Catalyst explorer (map)
Every known infrastructure and policy project on a map, filterable by type
(metro, road, airport, industrial corridor, rezoning, SEZ, IT park, hospital
or university) and by lifecycle stage. Clicking a catalyst shows its timeline,
budget, source documents, delay history, and the properties it affects.

### 1.4 Compare & portfolio
- Compare up to 4 properties or localities side by side.
- Portfolio: users add properties they own and track estimated value plus forecast changes.
- Alerts: "A catalyst affecting your property changed stage" (e.g. tender awarded, land acquisition stalled).

### 1.5 B2B console (developers, lenders, funds)
- Bulk upload of addresses or loan books → valuation, forecast and risk flags.
- Land/site-selection heatmaps: "Where will appreciation be highest in the next 5 years within 30 km of the CBD?"
- API access and white-label PDF reports.

---

## 2. Modules

### A. Data acquisition layer
| # | Module | Sources | Notes |
|---|---|---|---|
| A1 | Infrastructure catalyst ingestion | Metro rail corporations, NHAI / MoRTH, PM Gati Shakti, state PWDs, airport authorities, Smart Cities Mission, CPPP / state e-procurement tenders, Union & state budgets, cabinet decisions, gazette notifications | Track the **lifecycle**: Proposed → DPR → Approved → Funded → Land acquisition → Tendered → Awarded → Under construction → Commissioned |
| A2 | Master plan & zoning | Development plans / regional plans (e.g. DDA MPD-2041, MMRDA, BDA/BMRDA, HMDA, PMRDA), TP schemes, FSI/TDR policies | Georeference plan PDFs; rezoning is a big catalyst |
| A3 | Transactions & prices | State IGR registration data, circle rates / guidance values, licensed transaction data vendors, listing portals (asking prices), developer price sheets | Normalise to ₹/sq ft carpet area |
| A4 | Project supply | State RERA portals (project, promoter, timeline, units, complaints), developer launches | Supply pipeline and project risk |
| A5 | Macro & demand | RBI (rates, HPI), NHB RESIDEX, census and population projections, employment hubs (SEZ, IT parks), office absorption reports | Market baseline |
| A6 | Geospatial | OpenStreetMap, Bhuvan (ISRO), satellite imagery (construction progress), flood maps, CRZ, eco-sensitive zones | Distance, travel time, risk overlays |
| A7 | News & signals | Trusted news outlets, press releases (PIB), court orders (land-acquisition litigation) | Detects delays and stage changes |

### B. Intelligence layer
| # | Module | What it does |
|---|---|---|
| B1 | Document AI pipeline | LLM + OCR extraction from tenders, DPRs, gazettes and plan PDFs into structured catalyst records (type, location geometry, budget, dates, stage). Human-in-the-loop review queue before anything goes live |
| B2 | Entity resolution | Map "Hiranandani Estate", RERA IDs, survey numbers and addresses to one canonical property / building / micro-market ID |
| B3 | Geo engine | Geocoding, distance and travel-time (isochrone) to each catalyst, 500 m H3 grid |
| B4 | Current-value model (AVM) | Hedonic + gradient-boosted model on registrations and listings; gives today's value with a range |
| B5 | Catalyst impact model | Learns historical uplift from past projects (event-study / difference-in-differences: e.g. what happened to prices within 1 km of Delhi/Bengaluru/Mumbai metro stations from announcement to opening, compared with control areas). Outputs uplift curves by catalyst type × distance × city tier × lifecycle stage |
| B6 | Completion-probability model | Probability and expected delay of each catalyst, based on agency track record, stage, funding, land-acquisition status and litigation |
| B7 | Forecast engine | Combines market baseline + Σ (catalyst uplift × completion probability × time profile × distance decay) − supply overhang − risk penalties → P10/P50/P90 trajectory, 1–5 years |
| B8 | Explainability | Breaks the forecast into per-driver contributions (SHAP-style) and links every driver to its source |
| B9 | Back-testing & accuracy | Re-runs historical forecasts ("what would we have said in 2018?") and publishes error metrics by city |

**Forecast formula (simplified):**

```
Forecast(t) = Value_today
            × (1 + MarketBaseline(city, segment, t))
            × Π_i (1 + Uplift_i(type, distance, stage, t) × P_i(complete by t))
            × (1 − SupplyOverhangPenalty(t))
            × (1 − RiskPenalty)
```
Run as Monte Carlo over catalyst completion dates to get the fan chart.

### C. Product layer
| # | Module |
|---|---|
| C1 | Web app (search, report page, catalyst map, compare, portfolio) |
| C2 | Report generator (PDF, shareable link, white-label) |
| C3 | Alerts & notifications (email, WhatsApp) |
| C4 | Accounts, subscriptions, payments (Razorpay / UPI, GST invoicing) |
| C5 | B2B console + bulk upload |
| C6 | Public API (keys, metering, rate limits) |
| C7 | Admin / analyst console (catalyst review queue, data QA, overrides with audit log) |

### D. Trust, legal & ops layer
| # | Module |
|---|---|
| D1 | Source provenance: every data point stores URL, document hash, fetch date |
| D2 | Disclaimers and "information, not investment advice" framing; terms of use |
| D3 | DPDP Act 2023 compliance (consent, data minimisation for user portfolios) |
| D4 | Data licensing register (what we may redistribute, per source) |
| D5 | Model governance: versioning, drift monitoring, accuracy dashboards |

---

## 3. Reference architecture

```
 Sources (govt portals, RERA, IGR, tenders, news, OSM, satellite)
        │  scheduled crawlers / APIs / licensed feeds / manual upload
        ▼
 Raw lake (S3/GCS: PDFs, HTML, CSV)  ──► Document AI (OCR + LLM extraction)
        │                                        │
        ▼                                        ▼
 Staging DB  ◄──── Analyst review queue (human-in-the-loop) ────┘
        │
        ▼
 Core DB: PostgreSQL + PostGIS (properties, catalysts, transactions, H3 grid)
        │
        ├──► Feature store ──► ML (AVM, uplift, completion-prob, forecast) ──► model registry
        │                                                     │
        ▼                                                     ▼
 API layer (Node.js/TypeScript or Python FastAPI)  ◄── precomputed forecasts cache
        │
        ├──► Web app (Next.js + Mapbox/MapLibre)
        ├──► PDF report service
        ├──► Alerts (email / WhatsApp Business API)
        └──► Public B2B API
```

Suggested stack: Next.js + TypeScript, PostgreSQL/PostGIS, Python for ML
(LightGBM, PyMC/Prophet for baselines), Airflow/Dagster for pipelines, an LLM
for document extraction, MapLibre for maps, Razorpay for payments, hosted on
AWS Mumbai (ap-south-1) for data residency.

---

## 4. Day 1 vs Day 90 vs Day 360

### Day 1 (MVP, "concierge + map")
Goal: prove that people want this and will pay for it.
- **One city, 10–15 micro-markets.** Pick a city with heavy, visible infrastructure activity and good data access (Pune, Hyderabad, Bengaluru or MMR; MMR has the richest IGR data).
- Catalyst database of the **top 50–100 projects in that city**, entered by hand from government sources, each with a stage and source links.
- Catalyst map + locality page: price history (registrations + listings), nearby catalysts, a simple rule-based "appreciation outlook" (low/medium/high plus an indicative % range).
- Paid **PDF report** (₹499–₹999), assembled semi-manually by an analyst within 24 hours.
- Waitlist, WhatsApp channel, weekly "catalyst update" newsletter.
- Team: founder(s) + 1 full-stack + 1 data/analyst.

### Day 90 ("real model, first revenue")
- Building/project-level pages for all RERA-registered projects in the city.
- AVM v1 + catalyst uplift model v1 trained on historical metro and road projects in 3–4 cities; P10/P50/P90 forecasts.
- Explainability panel (driver contributions) and source links on every driver.
- Back-test v1 published for the launch city.
- Self-serve payments: single report + monthly Pro subscription.
- Alerts on catalyst stage changes.
- Document-AI pipeline that pulls in new tenders and notifications automatically, with analyst review.
- 2–3 **B2B design partners** (a mid-size developer for land selection, a broker network for white-label reports, one HFC/NBFC for collateral risk).
- Targets (hypotheses to validate): 20k MAU, 500 paid reports, 3 signed pilots.

### Day 360 ("multi-city platform")
- 6–8 cities (MMR, Pune, Bengaluru, Hyderabad, Delhi-NCR, Chennai, Ahmedabad, plus one tier-2 growth city such as Nagpur or Lucknow).
- Completion-probability model per agency, Monte Carlo fan charts, supply-overhang model.
- Portfolio tracking, compare, NRI-focused flows (USD/AED view).
- B2B console with bulk upload; public API with metered billing.
- Land and plot forecasting (not only apartments), including TP schemes and rezoning.
- Accuracy leaderboard by city, refreshed quarterly.
- Targets (hypotheses): 250k MAU, ₹40–60 lakh MRR, 10+ paying B2B accounts, at least one bank/HFC in production.

---

## 5. Future use cases (Year 2+)

1. **Lender risk products**: forward-looking LTV, collateral-decline early warning, LAP portfolio stress tests for banks, HFCs and NBFCs.
2. **Developer land intelligence**: "where to buy land next", demand forecasting per configuration and price point.
3. **Government / urban planning analytics**: value-capture financing estimates (betterment levy, TOD premium FSI) for metro corporations and municipal bodies.
4. **Insurance**: flood/climate risk combined with value for home insurers.
5. **Fractional and REIT/InvIT investors**: asset-level forecasts and portfolio analytics.
6. **NRI investor product**: remote due diligence, title and RERA checks, appreciation outlook.
7. **Rental yield forecasting** and rent-vs-buy calculators.
8. **Commercial, warehousing and data-centre land** forecasting along industrial corridors (DMIC, freight corridors).
9. **Agricultural → NA land conversion** opportunity scoring on city peripheries.
10. **Embedded data**: forecast widgets inside listing portals, banks' home-loan apps and wealth-management apps.
11. **"Infra delay index"**: a public index of government project delays (good for PR and policy influence).
12. **Marketplace (careful)**: verified, neutral transaction services (legal, title search, home loans) — only where it doesn't compromise forecast neutrality.
