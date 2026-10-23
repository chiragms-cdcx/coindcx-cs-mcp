# 01 — Market Landscape: Who Already Does This in India?

> Short answer: **many players do pieces of this, but nobody in India does
> "development-linked, explainable, back-tested appreciation forecasts" for a
> specific building or micro-market as their core product.** That gap is the
> opportunity.

Status, pricing and features of the companies below change often. Verify each
one (product page, recent funding news, a test account) before quoting them
in a pitch deck.

## 1. Competitive map

| Category | Examples | What they do well | What they don't do |
|---|---|---|---|
| Listing portals with price trends | 99acres, MagicBricks (PropIndex), Housing.com / PropTiger, NoBroker, Square Yards | Huge traffic, asking-price trends by locality, some "locality insights" | Prices are mostly *asking* prices, backward-looking, not tied to upcoming infrastructure, no forward forecast with confidence |
| Registration / transaction data | Square Yards (IGR-based transaction data), Zapkey, PropStack, CRE Matrix (commercial) | Actual registered transaction values, owner history | Historical only; sold mainly as data to B2B; little forecasting |
| Research & analytics firms | Liases Foras, PropEquity, Anarock Research, Knight Frank, JLL, CBRE, Colliers | Credible city/micro-market reports, supply/absorption data, developer-side relationships | Quarterly PDFs, city-level granularity, expensive, not self-serve, not per-building |
| Official indices | RBI House Price Index, NHB RESIDEX | Trusted, transaction-based | Very coarse (city level), lagging, no forecasts |
| Buyer advisory / curated platforms | Propsoch and similar "buyer-side advisor" startups | Data-led advice to home buyers | Mostly advisory/brokerage revenue; forecasting not the core productised IP |
| Fractional / investment platforms | Strata, hBits, Alt DRX-style platforms | Investor UX, yield framing | Commercial-asset focused, not a forecasting engine |
| Bank/NBFC valuation | Empanelled valuers, in-house risk teams | Regulated collateral valuation | Manual, slow, no forward risk/appreciation view |

### Global reference points (useful for the pitch)

- **Zillow Zestimate / Redfin Estimate (US)** — automated valuation models (AVMs) on current value. They prove the consumer demand for "what's my home worth", but they are about *today's* value, not development-driven *future* value.
- **HouseCanary, CoreLogic, Cotality (US)** — forecasting and AVMs sold to lenders and investors. The B2B playbook in this repo copies this model.
- **Local Logic (Canada)** — location scores built from geospatial data. Close to our "catalyst map" concept.

## 2. The gap we exploit

1. **Forward-looking, not backward-looking.** Everyone shows what prices *did*. Nobody credibly shows what a specific metro line, expressway, airport, IT park or master-plan rezoning *will do* to a specific building.
2. **Explainability.** "₹X/sq ft by 2029 because the Purple Line extension station is 650 m away (62% probability of opening by 2028), the ring road is in tender stage, and 4,200 new units are coming in the same micro-market."
3. **Evidence and trust.** Every catalyst is linked to its government source document, and the model publishes its own back-test accuracy. Portals have a conflict of interest because they earn from developers and brokers; we don't sell leads (see business model).
4. **Granularity.** Building/project level and 500 m grid level, not "Pune West".

## 3. Why now

- **RERA (2016)** made project data (approvals, timelines, carpet area, complaints) public per state.
- **State IGR portals** (Maharashtra IGR, Karnataka Kaveri, Telangana/AP registration, etc.) expose registration and circle-rate data.
- **Infrastructure push**: PM Gati Shakti National Master Plan, metro expansion in 20+ cities, Bharatmala expressways, new airports (Navi Mumbai, Noida/Jewar), and Smart Cities give a large, visible pipeline of catalysts.
- **Cheap geospatial and LLM tooling** makes it possible to read tender PDFs, master plans and gazette notifications at scale. Five years ago this would have needed a large team of analysts.

## 4. Key risks highlighted by the landscape

| Risk | Mitigation |
|---|---|
| Portals copy the feature | Moat = curated catalyst database + back-test track record + neutrality; portals can't credibly claim neutrality |
| Data access (IGR scraping, ToS) | Prefer official bulk data / RTI / licensed partners (e.g. license from a transaction-data vendor at start); legal review per state |
| Forecast wrong → reputational/legal risk | Ranges not point estimates, published accuracy, strong disclaimers, "information not advice" positioning |
| Indian real estate is opaque (cash component, under-reporting) | Blend registration values, listing prices, and developer price sheets; model the gap explicitly |
