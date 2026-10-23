"""Combine README and docs/*.md into a single styled PDF using headless Chrome.

Usage: python3 scripts/build_pdf.py  ->  writes Pragati-Founder-Plan.pdf
Requires: pip install markdown; google-chrome on PATH.
"""

import os
import re
import signal
import subprocess
import tempfile
import time
from datetime import date
from pathlib import Path

import markdown

ROOT = Path(__file__).resolve().parent.parent
OUTPUT = ROOT / "Pragati-Founder-Plan.pdf"

SECTIONS = [
    ("summary", "Founder's Summary", ROOT / "README.md"),
    ("market", "Market Landscape", ROOT / "docs/01-market-landscape.md"),
    ("blueprint", "Product Blueprint", ROOT / "docs/02-product-blueprint.md"),
    ("business", "Business Model & Fundraising", ROOT / "docs/03-business-model-and-fundraising.md"),
    ("brd", "Business Requirements Document (BRD)", ROOT / "docs/BRD.md"),
    ("prd", "Product Requirements Document (PRD)", ROOT / "docs/PRD.md"),
]

LINK_TARGETS = {
    "01-market-landscape.md": "market",
    "02-product-blueprint.md": "blueprint",
    "03-business-model-and-fundraising.md": "business",
    "BRD.md": "brd",
    "PRD.md": "prd",
}

CSS = """
@page { size: A4; margin: 18mm 16mm 20mm 16mm; }
body { font-family: "Helvetica Neue", Arial, sans-serif; font-size: 10.2pt;
       line-height: 1.5; color: #1f2933; }
h1 { font-size: 21pt; color: #0b3c5d; border-bottom: 3px solid #0b3c5d;
     padding-bottom: 6px; margin-top: 0; }
h2 { font-size: 14.5pt; color: #0b3c5d; margin-top: 22px; border-bottom: 1px solid #d9e2ec;
     padding-bottom: 3px; }
h3 { font-size: 12pt; color: #243b53; margin-top: 16px; }
h2, h3 { break-after: avoid; }
table { border-collapse: collapse; width: 100%; margin: 10px 0 14px; font-size: 8.8pt; }
th { background: #0b3c5d; color: #fff; text-align: left; }
th, td { border: 1px solid #bcccdc; padding: 5px 7px; vertical-align: top; }
tr:nth-child(even) td { background: #f5f8fb; }
tr { break-inside: avoid; }
pre { background: #f0f4f8; border: 1px solid #d9e2ec; border-radius: 4px; padding: 9px;
      font-size: 7.6pt; line-height: 1.3; white-space: pre; overflow: hidden; break-inside: avoid; }
code { font-family: "DejaVu Sans Mono", Menlo, monospace; font-size: 8.6pt; }
p code, li code, td code { background: #f0f4f8; padding: 1px 4px; border-radius: 3px; }
blockquote { border-left: 4px solid #2bb0ed; background: #eef8fd; margin: 10px 0;
             padding: 6px 12px; color: #243b53; }
a { color: #0967d2; text-decoration: none; }
.section { break-before: page; }
.cover { height: 250mm; display: flex; flex-direction: column; justify-content: center; }
.cover .brand { font-size: 40pt; font-weight: 700; color: #0b3c5d; }
.cover .tag { font-size: 15pt; color: #486581; margin-top: 8px; }
.cover .meta { margin-top: 40px; font-size: 10.5pt; color: #627d98; }
.cover .bar { width: 90px; height: 6px; background: #2bb0ed; margin: 22px 0; }
.toc h1 { border: none; }
.toc ol { font-size: 12pt; line-height: 2.1; }
"""


LIST_ITEM = re.compile(r"^\s*(?:[-*+]|\d+\.)\s")


def separate_lists(md_text: str) -> str:
    """Python-Markdown needs a blank line before a list that follows a paragraph."""
    out, in_code = [], False
    for line in md_text.split("\n"):
        if line.startswith("```"):
            in_code = not in_code
        prev = out[-1] if out else ""
        if (not in_code and LIST_ITEM.match(line) and prev.strip()
                and not LIST_ITEM.match(prev) and not prev.startswith(("```", "|", ">"))):
            out.append("")
        out.append(line)
    return "\n".join(out)


def render(md_text: str) -> str:
    for filename, anchor in LINK_TARGETS.items():
        md_text = re.sub(rf"\]\((?:docs/)?{re.escape(filename)}\)", f"](#{anchor})", md_text)
    md_text = md_text.replace("`01-market-landscape.md`", "the Market Landscape section")
    md_text = separate_lists(md_text)
    return markdown.markdown(md_text, extensions=["tables", "fenced_code", "sane_lists"])


def build_html() -> str:
    cover = f"""
    <div class="cover">
      <div class="brand">Pragati</div>
      <div class="tag">Property Appreciation Intelligence for India</div>
      <div class="bar"></div>
      <div>Founder's plan: market landscape, product blueprint, modules,<br>
           business model, Day 1 / 90 / 360 roadmap, VC fundraising, BRD and PRD</div>
      <div class="meta">Working draft v0.1 &middot; {date.today().strftime('%d %B %Y')}<br>
           "Pragati" is a placeholder name. Figures are planning hypotheses.</div>
    </div>"""
    toc_items = "".join(f'<li><a href="#{a}">{t}</a></li>' for a, t, _ in SECTIONS)
    toc = f'<div class="section toc"><h1>Contents</h1><ol>{toc_items}</ol></div>'
    body = "".join(
        f'<div class="section" id="{a}">{render(p.read_text(encoding="utf-8"))}</div>'
        for a, _, p in SECTIONS
    )
    return f"<!doctype html><html><head><meta charset='utf-8'><style>{CSS}</style></head><body>{cover}{toc}{body}</body></html>"


def main() -> None:
    with tempfile.TemporaryDirectory(ignore_cleanup_errors=True) as tmp:
        html_path = Path(tmp) / "plan.html"
        html_path.write_text(build_html(), encoding="utf-8")
        OUTPUT.unlink(missing_ok=True)
        cmd = [
            "google-chrome", "--headless", "--no-sandbox", "--disable-gpu",
            "--disable-dev-shm-usage", "--no-first-run", "--no-pdf-header-footer",
            f"--user-data-dir={tmp}/chrome", f"--print-to-pdf={OUTPUT}", html_path.as_uri(),
        ]
        # Headless Chrome can keep running after the PDF is written, so stop it
        # once the file size has settled.
        proc = subprocess.Popen(
            cmd, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, start_new_session=True
        )
        last_size, deadline = -1, time.time() + 180
        while time.time() < deadline and proc.poll() is None:
            time.sleep(2)
            size = OUTPUT.stat().st_size if OUTPUT.exists() else 0
            if size and size == last_size:
                break
            last_size = size
        if proc.poll() is None:
            os.killpg(proc.pid, signal.SIGKILL)
            proc.wait()
        if not OUTPUT.exists() or OUTPUT.stat().st_size == 0:
            raise SystemExit("PDF generation failed")
    print(f"Wrote {OUTPUT}")


if __name__ == "__main__":
    main()
