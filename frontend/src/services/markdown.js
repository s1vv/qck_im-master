// src/utils/markdown.js
import MarkdownIt from 'markdown-it';

const md = new MarkdownIt();

export function formatMarkdown(text) {
  return md.render(text);
}
