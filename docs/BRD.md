# Business Requirements Document (BRD)

| Field | Value |
|---|---|
| Product | Pragati (working name): development-linked property appreciation intelligence for India |
| Version | 0.1 (initial draft) |
| Status | Draft for founder review |
| Related docs | [PRD](PRD.md), [Market landscape](01-market-landscape.md), [Blueprint](02-product-blueprint.md), [Business model](03-business-model-and-fundraising.md) |

## 1. Executive summary
Indian households keep most of their wealth in real estate, yet purchase
decisions rely on broker opinion and backward-looking asking-price trends.
Public infrastructure (metro, expressways, airports, industrial corridors,
rezoning) is the biggest driver of future appreciation, and information about
it is public but scattered across thousands of unstructured government
documents. Pragati turns that information into explainable, source-linked,
back-tested appreciation forecasts for any building, project or locality.

## 2. Business problem
1. Buyers and investors cannot quantify how upcoming developments affect a specific property.
2. Existing portals show asking prices and are financially tied to developers and brokers, so their numbers aren't trusted as neutral.
3. Research firms publish city-level PDFs that are expensive, infrequent and not per-building.
4. Lenders value collateral manually and have no forward view of collateral risk.

## 3. Business objectives
| # | Objective | Measure | Target (Month 12) |
|---|---|---|---|
| BO1 | Become the trusted neutral source for forward-looking property data in launch cities | MAU, NPS, branded search volume | 250k MAU, NPS > 40 |
| BO2 | Prove willingness to pay | Paid reports + Pro subscribers | 10k paid reports cumulative; 2,000 Pro subscribers |
| BO3 | Establish the B2B revenue line | Paying B2B accounts, ARR | 10 accounts, ₹3 Cr+ ARR run-rate |
| BO4 | Prove forecast credibility | Published back-test accuracy | ≥ 65% of 3-year back-tests within ±10% |
| BO5 | Raise institutional capital | Seed round | Seed closed by Month 15–18 |

## 4. Scope

### In scope (Year 1)
- Residential apartments and projects (RERA-registered) in 6–8 cities; plots on city peripheries from Day 360.
- Infrastructure, zoning and supply catalysts from government and trusted sources.
- Web app (responsive), PDF reports, alerts, B2C subscriptions, B2B console and API (beta).

### Out of scope (Year 1)
- Brokerage, selling leads to developers, or transacting property.
- Regulated investment advice or portfolio management.
- Native mobile apps (a responsive web app/PWA comes first).
- Commercial office/retail forecasting (Year 2).

## 5. Stakeholders
| Stakeholder | Interest |
|---|---|
| Founders / leadership | Vision, fundraising, prioritisation |
| Data & ML team | Forecast quality, pipelines |
| Research analysts | Catalyst curation, QA |
| Engineering | Platform delivery |
| GTM (growth, B2B sales) | Acquisition, revenue |
| Legal / compliance advisor | Data licensing, disclaimers, DPDP, RERA implications |
| Customers: buyers, investors, NRIs, brokers, developers, lenders, funds | See segments in business model doc |
| Data sources: government agencies, data vendors | Access, licensing, attribution |

## 6. Business requirements
| ID | Requirement | Priority |
|---|---|---|
| BR-01 | Provide a current-value estimate with a range for any covered property or locality | Must |
| BR-02 | Provide a 1–5 year appreciation forecast with bear/base/bull scenarios and a confidence level | Must |
| BR-03 | Explain every forecast by listing the drivers (catalysts, market trend, supply, risks) and their contribution | Must |
| BR-04 | Link every catalyst to its original source document(s), with fetch date | Must |
| BR-05 | Track the lifecycle stage of each catalyst and update it as new information arrives | Must |
| BR-06 | Publish forecast accuracy (back-tested and live) per city | Must |
| BR-07 | Monetise through paid reports and subscriptions (UPI/cards, GST-compliant invoices) | Must |
| BR-08 | Offer B2B bulk valuation/forecast and an API with usage metering | Should (Day 90 pilot, Day 360 GA) |
| BR-09 | Notify users when catalysts affecting saved properties change | Should |
| BR-10 | Support white-label reports for brokers and advisors | Should |
| BR-11 | Maintain neutrality: no paid placement or developer influence on forecasts | Must (policy) |
| BR-12 | Comply with DPDP Act 2023, source terms of use, and consumer-protection norms | Must |

## 7. Success metrics (KPIs)
- **North star:** weekly users who view a full forecast report.
- Acquisition: organic traffic, SEO pages indexed, newsletter subscribers.
- Monetisation: free → paid conversion (target 2–4%), ARPU, B2B ARR.
- Trust: accuracy metrics, report-dispute rate, NPS.
- Data: catalysts tracked, % RERA projects covered, median "announcement → in database" lag (target < 72 hours).

## 8. Assumptions
1. Enough historical transaction and price data can be licensed or collected to train uplift models (at least 3 cities × 8+ years).
2. Users will pay for a neutral, explainable report before a ₹50 lakh+ decision.
3. Lenders and developers will run pilots with an early-stage vendor if the data is demonstrably better.
4. Government sources remain publicly accessible.

## 9. Constraints & dependencies
- Data-licensing terms per state IGR and portal; possible restrictions on redistribution.
- Analyst capacity for human-in-the-loop review in the early phase.
- LLM/OCR accuracy on scanned regional-language documents (Marathi, Kannada, Telugu, Hindi, etc.).
- Mapping and geocoding quality in Indian addresses.

## 10. Risks
| Risk | Impact | Likelihood | Mitigation |
|---|---|---|---|
| Inaccurate forecasts damage trust | High | Medium | Ranges, confidence, published accuracy, conservative launch cities |
| Legal claims from users who relied on forecasts | High | Low–Medium | Clear ToS, "information not advice", no guaranteed returns language, legal review |
| Data access blocked or restricted | High | Medium | Multiple sources, licensing, own archive, RTI |
| Incumbent portals launch a similar feature | Medium | High | Neutrality, depth, track record, B2B focus |
| Slow B2B sales cycles | Medium | High | Start with mid-size developers and NBFCs; land-and-expand |
| Defamation risk on project risk flags | Medium | Low | Only factual, sourced flags (RERA complaints, court orders) |

## 11. High-level timeline
| Phase | Milestone |
|---|---|
| Day 0–30 | Launch city chosen, top 100 catalysts curated, MVP locality pages, waitlist |
| Day 30–90 | AVM v1, uplift model v1, paid reports, Pro plan, 2–3 B2B design partners, first back-test |
| Day 90–180 | 3 cities, document-AI pipeline, alerts, pre-seed raise |
| Day 180–360 | 6–8 cities, B2B console + API, completion-probability model, seed-round preparation |

## 12. Approval
| Role | Name | Sign-off |
|---|---|---|
| Founder / CEO | | |
| Co-founder / CTO | | |
| Advisor (real estate) | | |
| Advisor (legal) | | |
