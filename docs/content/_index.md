---
title: nux
type: docs
description: "A modern tmux session manager with batch start, session groups, pattern matching, and zero-config project discovery."
---

<div class="nux-hero">
  <h1 class="nux-hero__title">nux</h1>
  <p class="nux-hero__tagline">A modern tmux session manager built for power users.</p>
  <div class="nux-hero__cta">
    <a href="{{< relref "/docs/getting-started/installation" >}}" class="nux-btn nux-btn--primary">Get Started</a>
    <a href="https://github.com/Drew-Daniels/nux" class="nux-btn nux-btn--secondary">GitHub</a>
  </div>
</div>

<div class="nux-terminal">
  <div class="nux-terminal__bar">
    <span class="dot-red"></span>
    <span class="dot-yellow"></span>
    <span class="dot-green"></span>
  </div>
  <div class="nux-terminal__body">
    <span class="comment"># Start a session for any project directory</span><br>
    <span class="prompt">$ </span><span class="cmd">nux blog</span><br><br>
    <span class="comment"># Batch-start everything for your workday</span><br>
    <span class="prompt">$ </span><span class="cmd">nux @work</span><br>
    <span class="output">  started api, web, workers, docs (4 sessions)</span><br><br>
    <span class="comment"># Stop all sessions matching a pattern</span><br>
    <span class="prompt">$ </span><span class="cmd">nux stop web+</span><br><br>
    <span class="comment"># Preview commands without executing</span><br>
    <span class="prompt">$ </span><span class="cmd">nux --dry-run blog</span><br>
    <span class="output">  tmux new-session -d -s blog -c ~/projects/blog</span><br>
    <span class="output">  tmux send-keys -t blog:editor nvim Enter</span>
  </div>
</div>

<h2 class="nux-section-title">Why nux?</h2>
<p class="nux-section-desc">If you run many tmux sessions simultaneously - a dozen or more across different projects - existing tools fall short. They're designed for one session at a time. nux is built from the ground up for <strong>batch-oriented, declarative session management</strong>.</p>

<div class="nux-features">
  <div class="nux-feature">
    <div class="nux-feature__icon">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="7" width="20" height="14" rx="2" ry="2"/><path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16"/></svg>
    </div>
    <div class="nux-feature__title">Session Groups</div>
    <div class="nux-feature__desc">Define named groups in config. <code>nux @work</code> starts 8 sessions at once.</div>
  </div>
  <div class="nux-feature">
    <div class="nux-feature__icon">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/></svg>
    </div>
    <div class="nux-feature__title">Pattern Matching</div>
    <div class="nux-feature__desc"><code>nux web+</code> starts all projects matching the pattern. No quoting needed.</div>
  </div>
  <div class="nux-feature">
    <div class="nux-feature__icon">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"><path d="M13 2 3 14h9l-1 8 10-12h-9l1-8z"/></svg>
    </div>
    <div class="nux-feature__title">Zero Config</div>
    <div class="nux-feature__desc">Projects in <code>~/projects/</code> work without any YAML. Convention over configuration.</div>
  </div>
  <div class="nux-feature">
    <div class="nux-feature__icon">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>
    </div>
    <div class="nux-feature__title">Selective Windows</div>
    <div class="nux-feature__desc">Restart individual windows without tearing down the session with <code>project:window</code> syntax.</div>
  </div>
  <div class="nux-feature">
    <div class="nux-feature__icon">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/><path d="M11 8v6"/><path d="M8 11h6"/></svg>
    </div>
    <div class="nux-feature__title">Smart Discovery</div>
    <div class="nux-feature__desc">fzf/gum picker, zoxide integration, and auto-detect from the current directory.</div>
  </div>
  <div class="nux-feature">
    <div class="nux-feature__icon">
      <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"><path d="M4 17V7l7-4 7 4v10l-7 4z"/><path d="m11 3 7 4"/><path d="M4 7l7 4"/><path d="M11 11v10"/></svg>
    </div>
    <div class="nux-feature__title">Custom Variables</div>
    <div class="nux-feature__desc">Template configs with <code>&#123;&#123;var&#125;&#125;</code> placeholders. Override at runtime with <code>--var</code>.</div>
  </div>
</div>

<h2 class="nux-section-title">Get started</h2>
<p class="nux-section-desc">Head to the <a href="{{< relref "/docs/getting-started/installation" >}}">installation guide</a> to install nux, then follow the <a href="{{< relref "/docs/getting-started/quickstart" >}}">quickstart</a> to set up your first session.<br><br>Coming from tmuxinator? The <a href="{{< relref "/docs/getting-started/migrating-from-tmuxinator" >}}">migration guide</a> covers the differences.</p>
