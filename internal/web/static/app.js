// imports
import Prism from "https://cdn.jsdelivr.net/npm/prismjs@1.29.0/+esm";
import "https://cdn.jsdelivr.net/npm/prismjs@1.29.0/components/prism-sql.min.js/+esm";

import mermaid from "https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.esm.min.mjs";

import { marked } from "https://cdn.jsdelivr.net/npm/marked/+esm";
import { gfmHeadingId } from "https://cdn.jsdelivr.net/npm/marked-gfm-heading-id/+esm";

import DOMPurify from "https://cdn.jsdelivr.net/npm/dompurify/+esm";

marked.use(gfmHeadingId());

document.body.addEventListener("htmx:afterSwap", () => {
  Prism.highlightAll();
  mermaid.run({ querySelector: ".mermaid" });

  document.querySelectorAll(".markdown-source").forEach((src) => {
    const target = src.previousElementSibling;

    target.innerHTML = DOMPurify.sanitize(marked.parse(src.innerHTML));
  });
});
