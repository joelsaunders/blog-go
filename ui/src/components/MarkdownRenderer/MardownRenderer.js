import React from 'react';
import marked from 'marked';
import hljs from 'highlight.js';
import renderHTML from 'react-render-html';

const renderer = new marked.Renderer();

renderer.code = (code, language) => {
    // Check whether the given language is valid for highlight.js.
    const validLang = !!(language && hljs.getLanguage(language));
    // Highlight only if the language is valid.
    const highlighted = validLang ? hljs.highlight(language, code).value : code;
    // Render the highlighted code with `hljs` class.
    return `<pre><code class="hljs ${language}">${highlighted}</code></pre>`;
};

marked.setOptions({
    renderer: renderer,
    gfm: true,
    breaks: true,
    highlight: (code, language) => {
        return language? hljs.highlight(language, code).value: hljs.highlightAuto(code).value;
    }
});

const MarkdownRenderer = ({text}) => {
    return text? <div className="markdown">{renderHTML(marked(text))}</div>: null
};

export default MarkdownRenderer;