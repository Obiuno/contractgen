// imports
import Prism from "https://cdn.jsdelivr.net/npm/prismjs@1.29.0/+esm";
import "https://cdn.jsdelivr.net/npm/prismjs@1.29.0/components/prism-sql.min.js/+esm";

import mermaid from "https://cdn.jsdelivr.net/npm/mermaid@11/dist/mermaid.esm.min.mjs";

import { marked } from "https://cdn.jsdelivr.net/npm/marked/+esm";
import { gfmHeadingId } from "https://cdn.jsdelivr.net/npm/marked-gfm-heading-id/+esm";
import "https://cdn.jsdelivr.net/npm/prismjs@1.29.0/components/prism-markdown.min.js/+esm";

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
    
    // Render markdown from each panel's source into its rendered slot
    document.querySelectorAll(".panel .markdown-output").forEach((target) => {
        const panel = target.closest(".panel");
        const source = panel.querySelector(".output-source");
        if (source) {
            target.innerHTML = DOMPurify.sanitize(marked.parse(source.textContent));
        }
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

// copy button event listner
document.body.addEventListener("click", (e) => {
    if (e.target.matches(".btn-copy")) {
        const panel = e.target.closest(".panel");
        // Prefer source, fall back to rendered
        const sourceEl = panel.querySelector(".output-source") || panel.querySelector(".output-rendered");
        copyToClipboard(e.target, sourceEl.textContent);
        return;
    }
    
    if (e.target.matches(".btn-toggle")) {
        const panel = e.target.closest(".panel");
        const showingSource = panel.classList.toggle("showing-source");
        e.target.textContent = showingSource ? "Rendered" : "Source";
        return;
    }
});