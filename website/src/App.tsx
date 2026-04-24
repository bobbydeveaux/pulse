function App() {
  return (
    <>
      <Nav />
      <Hero />
      <Features />
      <TimeSeries />
      <Metrics />
      <Install />
      <GithubAction />
      <Cta />
      <Footer />
    </>
  );
}

/* ─── NAV ─── */
function Nav() {
  return (
    <nav>
      <a href="/" className="nav-logo">
        <svg viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
          <rect width="32" height="32" rx="6" fill="#0d1117" stroke="#22d3ee" strokeWidth="1.5" />
          <polyline
            points="4,16 9,16 12,8 16,24 20,12 23,16 28,16"
            fill="none"
            stroke="#22d3ee"
            strokeWidth="2.5"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
        </svg>
        Pulse
      </a>
      <div className="nav-links">
        <a href="#features">Features</a>
        <a href="#time-series">Trends</a>
        <a href="#metrics">Metrics</a>
        <a href="#install">Install</a>
        <a href="https://github.com/bobbydeveaux/pulse" className="btn-gh">
          <GithubIcon size={16} />
          GitHub
        </a>
      </div>
    </nav>
  );
}

/* ─── HERO ─── */
function Hero() {
  return (
    <section className="hero">
      <div className="hero-badge">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor">
          <circle cx="12" cy="12" r="10" />
        </svg>
        Open Source &middot; Free Forever
      </div>
      <h1>
        Take the <span className="accent">pulse</span> of
        <br />
        your codebase
      </h1>
      <p className="hero-sub">
        Pulse is a multi-language code quality analyzer that tracks complexity, duplication, and
        maintainability — with trends over your entire git history.
      </p>
      <div className="hero-actions">
        <a href="#install" className="btn-primary">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
            <path d="M12 2v13M7 11l5 5 5-5" />
            <path d="M3 19h18" />
          </svg>
          Get Started
        </a>
        <a href="https://github.com/bobbydeveaux/pulse" className="btn-secondary">
          <GithubIcon size={18} />
          View on GitHub
        </a>
      </div>

      <div className="terminal">
        <div className="terminal-bar">
          <div className="dot dot-r" />
          <div className="dot dot-y" />
          <div className="dot dot-g" />
        </div>
        <div className="terminal-body">
          <p>
            <span className="t-prompt">$</span> <span className="t-cmd">pulse check ./src</span>
          </p>
          <p className="t-muted">Analyzing 142 files across 3 languages...</p>
          <p>&nbsp;</p>
          <p>
            <span className="t-muted">  Language      Files    SLOC    Comments    Blanks</span>
          </p>
          <p>
            <span className="t-cmd">  Go              89   12,847     1,203     2,104</span>
          </p>
          <p>
            <span className="t-cmd">  TypeScript      41    6,322       412       891</span>
          </p>
          <p>
            <span className="t-cmd">  Python          12    1,891       203       334</span>
          </p>
          <p>&nbsp;</p>
          <p>
            <span className="t-cyan">Top complexity hotspots:</span>
          </p>
          <p>
            <span className="t-cmd">  engine.go:ProcessBatch       </span>
            <span className="t-muted">CCN:</span>
            <span className="t-err"> 23</span>
            <span className="t-muted">  Cognitive:</span>
            <span className="t-err"> 31</span>
            <span className="t-muted">  Grade:</span>
            <span className="t-grade-d"> D</span>
          </p>
          <p>
            <span className="t-cmd">  handler.go:RouteRequest      </span>
            <span className="t-muted">CCN:</span>
            <span className="t-warn"> 18</span>
            <span className="t-muted">  Cognitive:</span>
            <span className="t-warn"> 22</span>
            <span className="t-muted">  Grade:</span>
            <span className="t-grade-c"> C</span>
          </p>
          <p>
            <span className="t-cmd">  ast.go:WalkTree              </span>
            <span className="t-muted">CCN:</span>
            <span className="t-warn"> 15</span>
            <span className="t-muted">  Cognitive:</span>
            <span className="t-warn"> 19</span>
            <span className="t-muted">  Grade:</span>
            <span className="t-grade-c"> C</span>
          </p>
          <p>&nbsp;</p>
          <p>
            <span className="t-ok">Duplication:</span>
            <span className="t-cmd"> 4.2%</span>
            <span className="t-muted"> (3 clones detected)</span>
          </p>
          <p>
            <span className="t-ok">Maintainability Index:</span>
            <span className="t-cmd"> 71.3</span>
            <span className="t-grade-b"> (B)</span>
          </p>
          <p>
            <span className="t-ok">Estimated effort:</span>
            <span className="t-cmd"> 14.2 person-months</span>
            <span className="t-muted"> (COCOMO)</span>
          </p>
        </div>
      </div>
    </section>
  );
}

/* ─── FEATURES ─── */
const features = [
  {
    icon: '🔄',
    title: 'Cyclomatic Complexity',
    desc: 'Measures the number of independent paths through your code. High CCN means more test cases needed and more places for bugs to hide.',
  },
  {
    icon: '🧠',
    title: 'Cognitive Complexity',
    desc: "How hard is your code to understand? Unlike cyclomatic complexity, cognitive complexity penalises nested logic and broken flow — measuring human readability.",
    badge: { text: 'Unique', color: 'cyan' },
  },
  {
    icon: '📊',
    title: 'Maintainability Index',
    desc: 'A composite score (0-100) combining SLOC, cyclomatic complexity, and Halstead metrics. Every function gets a letter grade from A to F.',
  },
  {
    icon: '📋',
    title: 'Copy-Paste Detection',
    desc: 'Finds duplicated code blocks across your entire codebase. Duplication breeds inconsistency — catch clones before they diverge.',
  },
  {
    icon: '📈',
    title: 'Time-Series Trends',
    desc: 'Track how complexity, duplication, and maintainability change across your git history. See if that refactor actually helped.',
    badge: { text: 'Unique', color: 'cyan' },
  },
  {
    icon: '🔀',
    title: 'PR-Level Diffs',
    desc: '"This PR increased average complexity by +3.2." Quality gates that show the delta, not just absolute numbers — directly in your CI.',
    badge: { text: 'Unique', color: 'cyan' },
  },
  {
    icon: '💰',
    title: 'COCOMO Effort Estimation',
    desc: 'Estimates person-months and cost to develop using the Constructive Cost Model. Useful for project planning and understanding codebase scale.',
  },
  {
    icon: '🌍',
    title: 'Multi-Language',
    desc: 'Go, TypeScript, JavaScript, Python, Java, Rust, C/C++, Ruby, PHP, Kotlin, Swift, and more. One tool for your entire stack.',
  },
  {
    icon: '🚦',
    title: 'Quality Gates',
    desc: 'Set thresholds for complexity, duplication, and maintainability. Fail CI if code quality drops below your standards.',
  },
];

function Features() {
  return (
    <section className="section" id="features">
      <p className="section-label">Features</p>
      <h2 className="section-title">Everything you need to measure code health</h2>
      <p className="section-sub">
        From function-level complexity to project-wide trends — Pulse gives you the full picture.
      </p>
      <div className="features-grid">
        {features.map((f) => (
          <div className="feature-card" key={f.title}>
            <div className="feature-icon">{f.icon}</div>
            <h3>
              {f.title}
              {f.badge && (
                <span className={`badge badge-${f.badge.color}`}>{f.badge.text}</span>
              )}
            </h3>
            <p>{f.desc}</p>
          </div>
        ))}
      </div>
    </section>
  );
}

/* ─── TIME-SERIES EXPLAINER ─── */
function TimeSeries() {
  // Simulated complexity data over 30 commits
  const data = [
    42, 43, 45, 44, 48, 52, 55, 58, 56, 54, 57, 60, 63, 65, 62, 58, 55, 53,
    50, 48, 46, 45, 43, 42, 41, 40, 39, 38, 37, 36,
  ];
  const max = Math.max(...data);

  return (
    <section className="timeseries-section" id="time-series">
      <p className="section-label">Time-Series Tracking</p>
      <h2 className="section-title">See how your code evolves</h2>
      <p className="section-sub">
        Pulse walks your git history and builds a timeline of quality metrics. Every commit tells a
        story.
      </p>

      <div className="timeseries-demo">
        <div className="timeseries-chart">
          <h3>Average Cyclomatic Complexity — last 30 commits</h3>
          <div className="chart-area">
            {data.map((val, i) => {
              const height = (val / max) * 100;
              const improving = i >= 14;
              return (
                <div
                  key={i}
                  className="chart-bar"
                  style={{
                    height: `${height}%`,
                    background: improving
                      ? 'var(--green)'
                      : i >= 10
                        ? 'var(--orange)'
                        : 'var(--cyan)',
                    opacity: 0.8,
                  }}
                  title={`Commit ${i + 1}: CCN ${val}`}
                />
              );
            })}
          </div>
          <div className="chart-labels">
            <span>30 commits ago</span>
            <span>HEAD</span>
          </div>
          <div className="chart-legend">
            <span>
              <span className="legend-dot" style={{ background: 'var(--cyan)' }} />
              Stable
            </span>
            <span>
              <span className="legend-dot" style={{ background: 'var(--orange)' }} />
              Rising
            </span>
            <span>
              <span className="legend-dot" style={{ background: 'var(--green)' }} />
              Improving
            </span>
          </div>
        </div>

        <div className="timeseries-explain">
          <div className="explain-card">
            <h4>What is time-series tracking?</h4>
            <p>
              Instead of just showing you today's metrics, Pulse analyses each commit in your git
              history and records how complexity, duplication, and maintainability have changed over
              time. Think of it like a fitness tracker — but for your code.
            </p>
          </div>
          <div className="explain-card">
            <h4>Spot regressions instantly</h4>
            <p>
              That "quick fix" that doubled the complexity of your auth module? You'll see it as a
              spike in the chart. Track exactly which commit introduced quality regressions and who
              authored it.
            </p>
          </div>
          <div className="explain-card">
            <h4>Prove your refactors work</h4>
            <p>
              Show stakeholders a chart that proves the refactoring sprint actually reduced
              complexity by 40%. Data-driven engineering, not gut feelings.
            </p>
          </div>
          <div className="explain-card">
            <h4>PR-level deltas</h4>
            <p>
              Every pull request gets an automatic comment: "This PR changes average complexity from
              12.3 to 14.1 (+1.8)". Set quality gates to block PRs that exceed your thresholds.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}

/* ─── METRICS TABLE ─── */
function Metrics() {
  return (
    <section className="section" id="metrics">
      <p className="section-label">What it measures</p>
      <h2 className="section-title">Comprehensive metrics</h2>

      <table className="metrics-table">
        <thead>
          <tr>
            <th>Metric</th>
            <th>What it tells you</th>
            <th>Scope</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>
              <strong>Cyclomatic Complexity (CCN)</strong>
            </td>
            <td>Number of independent execution paths. Higher = harder to test exhaustively.</td>
            <td>
              <span className="badge badge-green">Per function</span>
            </td>
          </tr>
          <tr>
            <td>
              <strong>Cognitive Complexity</strong>
            </td>
            <td>How hard code is for a human to read. Penalises nesting, breaks in flow, and boolean logic chains.</td>
            <td>
              <span className="badge badge-green">Per function</span>
            </td>
          </tr>
          <tr>
            <td>
              <strong>Maintainability Index</strong>
            </td>
            <td>Composite score (0-100) with letter grades A-F. Combines SLOC, complexity, and Halstead volume.</td>
            <td>
              <span className="badge badge-cyan">Per file</span>
            </td>
          </tr>
          <tr>
            <td>
              <strong>Halstead Metrics</strong>
            </td>
            <td>Volume, difficulty, and effort based on operators and operands. Measures information content.</td>
            <td>
              <span className="badge badge-cyan">Per file</span>
            </td>
          </tr>
          <tr>
            <td>
              <strong>SLOC / Comments / Blanks</strong>
            </td>
            <td>Source lines of code, comment density, and blank line ratios.</td>
            <td>
              <span className="badge badge-purple">Per project</span>
            </td>
          </tr>
          <tr>
            <td>
              <strong>Duplication %</strong>
            </td>
            <td>Percentage of code that is duplicated. Highlights exact clones and near-duplicates.</td>
            <td>
              <span className="badge badge-purple">Per project</span>
            </td>
          </tr>
          <tr>
            <td>
              <strong>COCOMO Estimate</strong>
            </td>
            <td>Person-months and development cost estimate using the Constructive Cost Model.</td>
            <td>
              <span className="badge badge-purple">Per project</span>
            </td>
          </tr>
        </tbody>
      </table>
    </section>
  );
}

/* ─── INSTALL ─── */
function Install() {
  return (
    <section id="install" style={{ padding: '2rem 0 5rem' }}>
      <div className="install-section">
        <p className="section-label">Installation</p>
        <h2 className="section-title">Up in 30 seconds</h2>

        <pre className="code-block">{`# Install Pulse
go install github.com/bobbydeveaux/pulse/app/cmd/pulse@latest

# Analyse your codebase
pulse check ./src

# See trends over your last 50 commits
pulse trend --last 50

# Check complexity diff for staged changes
pulse diff --staged

# Set quality gates (fail if avg CCN > 15)
pulse gate --max-ccn 15 --max-duplication 5`}</pre>
      </div>
    </section>
  );
}

/* ─── GITHUB ACTION ─── */
function GithubAction() {
  return (
    <section id="github-action" style={{ padding: '2rem 0 5rem' }}>
      <div className="install-section">
        <p className="section-label">CI / CD</p>
        <h2 className="section-title">GitHub Action</h2>
        <p className="section-sub" style={{ marginBottom: '2rem' }}>
          Run Pulse on every pull request. Get complexity diffs and quality gates in your CI pipeline.
        </p>

        <pre className="code-block">{`# .github/workflows/quality.yml
name: Code Quality
on: [pull_request]

jobs:
  pulse:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: bobbydeveaux/pulse@main
        with:
          # Post complexity diff as PR comment
          comment: true
          # Fail if thresholds exceeded
          max_ccn: 15
          max_cognitive: 20
          max_duplication: 5`}</pre>
      </div>
    </section>
  );
}

/* ─── CTA ─── */
function Cta() {
  return (
    <section className="cta-section">
      <h2>Know your code, quantitatively</h2>
      <p>
        Pulse is free, open-source, and built for teams that care about code quality.
        <br />
        Works alongside{' '}
        <a href="https://guardian.stackramp.io" style={{ color: 'var(--green)' }}>
          Guardian
        </a>{' '}
        for security scanning.
      </p>
      <div style={{ display: 'flex', gap: '1rem', justifyContent: 'center', flexWrap: 'wrap' }}>
        <a href="#install" className="btn-primary">
          Get Started
        </a>
        <a href="https://github.com/bobbydeveaux/pulse" className="btn-secondary">
          <GithubIcon size={18} />
          Star on GitHub
        </a>
      </div>
    </section>
  );
}

/* ─── FOOTER ─── */
function Footer() {
  return (
    <footer>
      <p>Pulse is open-source under the MIT licence. Built with Go.</p>
      <p style={{ marginTop: '0.5rem' }}>
        <a href="https://github.com/bobbydeveaux/pulse">GitHub</a> &middot;{' '}
        <a href="https://github.com/bobbydeveaux/pulse/issues">Issues</a> &middot;{' '}
        <a href="https://github.com/bobbydeveaux/pulse/blob/main/LICENSE">MIT Licence</a> &middot;{' '}
        <a href="https://guardian.stackramp.io">Guardian</a>
      </p>
    </footer>
  );
}

/* ─── SHARED ICONS ─── */
function GithubIcon({ size = 16 }: { size?: number }) {
  return (
    <svg width={size} height={size} viewBox="0 0 24 24" fill="currentColor">
      <path d="M12 2C6.48 2 2 6.58 2 12.26c0 4.54 2.87 8.39 6.84 9.75.5.09.68-.22.68-.49v-1.7c-2.78.62-3.37-1.37-3.37-1.37-.45-1.18-1.1-1.5-1.1-1.5-.9-.63.07-.62.07-.62 1 .07 1.52 1.05 1.52 1.05.88 1.55 2.32 1.1 2.88.84.09-.65.35-1.1.63-1.35-2.22-.26-4.56-1.14-4.56-5.07 0-1.12.39-2.03 1.03-2.75-.1-.26-.45-1.3.1-2.71 0 0 .84-.28 2.75 1.05A9.38 9.38 0 0112 7.58c.85 0 1.7.12 2.5.34 1.9-1.33 2.74-1.05 2.74-1.05.55 1.41.2 2.45.1 2.71.64.72 1.03 1.63 1.03 2.75 0 3.94-2.34 4.81-4.57 5.06.36.32.68.94.68 1.9v2.82c0 .27.18.59.69.49C19.13 20.65 22 16.8 22 12.26 22 6.58 17.52 2 12 2z" />
    </svg>
  );
}

export default App;
