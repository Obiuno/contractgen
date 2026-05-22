// imports
import Prism from "https://cdn.jsdelivr.net/npm/prismjs@1.29.0/+esm";
import "https://cdn.jsdelivr.net/npm/prismjs@1.29.0/components/prism-sql.min.js/+esm";

import mermaid from "https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.esm.min.mjs";

import { marked } from "https://cdn.jsdelivr.net/npm/marked/+esm";
import { gfmHeadingId } from "https://cdn.jsdelivr.net/npm/marked-gfm-heading-id/+esm";

import DOMPurify from "https://cdn.jsdelivr.net/npm/dompurify/+esm";

mermaid.initialize({
  startOnLoad: false,
  theme: "base",
  er: {
    useMaxWidth: true,
    entityPadding: 15,
    diagramPadding: 20,
  },
  themeVariables: {
    primaryColor: "#1a2020",
    primaryTextColor: "#e0e6e6",
    primaryBorderColor: "#00d9c0",
    lineColor: "#00a890",
  },
});

console.log('mermaid config:', mermaid.mermaidAPI.getConfig());

marked.use(gfmHeadingId());

document.body.addEventListener("htmx:afterSwap", () => {
  Prism.highlightAll();
  mermaid.run({ querySelector: ".mermaid" });

  document.querySelectorAll(".markdown-source").forEach((src) => {
    const target = src.previousElementSibling;

    target.innerHTML = DOMPurify.sanitize(marked.parse(src.innerHTML));
  });
});

function copyToClipboard(button, content) {
  navigator.clipboard
    .writeText(content)
    .then(() => {
      button.textContent = "Copied!";
      button.classList.add("success");
      setTimeout(() => {
        button.textContent = "Copy";
        button.classList.remove("success");
      }, 1500);
    })
    .catch(() => {
      button.textContent = "Failed";
      setTimeout(() => {
        button.textContent = "Copy";
      }, 1500);
    });
}

document.body.addEventListener("click", (e) => {
    if (e.target.matches(".btn-copy")) {
        const target = document.getElementById(e.target.dataset.target);
        if (target) copyToClipboard(e.target, target.textContent);
        return;
    }
    
    if (e.target.matches(".btn-toggle")) {
        const target = document.getElementById(e.target.dataset.target);
        if (!target) return;
        const showingSource = target.classList.toggle("showing-source");
        e.target.textContent = showingSource ? "Rendered" : "Source";
        return;
    }
});