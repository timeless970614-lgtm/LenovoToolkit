<template>
  <div class="launch-speed-page">

    <!-- Boot Speed Info -->
    <div class="card">
      <div class="card-header">
        <h3 class="card-title">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
          </svg>
          开机速度
        </h3>
        <button class="btn btn-secondary btn-sm" @click="refreshBootInfo" :disabled="bootLoading">
          <span v-if="bootLoading" class="spinner-small"></span>
          <span v-else>↻ 刷新</span>
        </button>
      </div>
      <div class="test-body">
        <div class="boot-info-grid" v-if="bootInfo">
          <div class="boot-stat-card">
            <span class="boot-stat-label">上次开机</span>
            <span class="boot-stat-value mono">{{ bootInfo.lastBootTime || '—' }}</span>
          </div>
          <div class="boot-stat-card">
            <span class="boot-stat-label">运行时长</span>
            <span class="boot-stat-value mono highlight">{{ formatUptime(bootInfo.uptimeSeconds) }}</span>
          </div>
          <div class="boot-stat-card">
            <span class="boot-stat-label">启动项数量</span>
            <span class="boot-stat-value mono" :class="bootInfo.startupCount > 10 ? 'warn' : 'ok'">{{ bootInfo.startupCount }}</span>
            <span class="boot-stat-hint">{{ bootInfo.startupCount > 10 ? '偏多' : '正常' }}</span>
          </div>
          <div class="boot-stat-card">
            <span class="boot-stat-label">自动服务</span>
            <span class="boot-stat-value mono">{{ bootInfo.serviceCount }}</span>
            <span class="boot-stat-hint">个正在运行</span>
          </div>
        </div>
        <div v-else class="test-loading">
          <div class="spinner-small"></div>
          <span>正在获取开机信息...</span>
        </div>

        <!-- Startup items detail -->
        <div v-if="bootInfo && startupList.length" class="startup-section">
          <div class="section-header" @click="showStartupDetail = !showStartupDetail">
            <span class="section-title">📋 启动项列表</span>
            <span class="section-toggle">{{ showStartupDetail ? '▼' : '▶' }}</span>
          </div>
          <div v-if="showStartupDetail" class="startup-list">
            <div v-for="item in startupList" :key="item.name" class="startup-item">
              <span class="startup-name">{{ item.name }}</span>
              <span class="startup-source">{{ item.source }}</span>
              <span class="startup-path mono">{{ truncatePath(item.path) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- App Launch Benchmark -->
    <div class="card">
      <div class="card-header">
        <h3 class="card-title">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="5 3 19 12 5 21 5 3"/>
          </svg>
          应用启动速度测试
        </h3>
        <div class="header-actions">
          <select class="method-select" v-model="launchMethod">
            <option value="warm">Warm（热启动）</option>
            <option value="cold">Cold（冷启动）</option>
          </select>
          <button class="btn btn-primary btn-sm" @click="runBenchmark" :disabled="benchmarkRunning">
            <span v-if="benchmarkRunning" class="spinner-small"></span>
            <span v-else>▶ 开始测试</span>
          </button>
        </div>
      </div>
      <div class="test-body">
        <!-- Summary stats -->
        <div v-if="report" class="summary-grid">
          <div class="summary-card best">
            <span class="summary-label">最快启动</span>
            <span class="summary-value">{{ report.fastest || '—' }}</span>
          </div>
          <div class="summary-card worst">
            <span class="summary-label">最慢启动</span>
            <span class="summary-value">{{ report.slowest || '—' }}</span>
          </div>
          <div class="summary-card avg">
            <span class="summary-label">平均耗时</span>
            <span class="summary-value mono">{{ report.avgLaunch }} ms</span>
          </div>
          <div class="summary-card total">
            <span class="summary-label">总测试时间</span>
            <span class="summary-value mono">{{ report.totalTime }} ms</span>
          </div>
        </div>

        <!-- Running progress -->
        <div v-if="benchmarkRunning" class="test-loading">
          <div class="spinner-small"></div>
          <span>正在测试应用启动速度... ({{ doneCount }}/{{ totalCount }})</span>
        </div>

        <!-- Results table -->
        <div v-if="report && report.results.length" class="results-table">
          <div class="result-header">
            <span class="col-name">应用</span>
            <span class="col-cat">类别</span>
            <span class="col-process">进程创建</span>
            <span class="col-launch">启动耗时</span>
            <span class="col-status">状态</span>
          </div>
          <div v-for="r in report.results" :key="r.appName" 
               :class="['result-row', r.success ? 'success' : 'failed']">
            <span class="col-name">
              <span :class="['cat-dot', 'cat-' + r.category]"></span>
              {{ r.appName }}
            </span>
            <span class="col-cat">{{ categoryLabel(r.category) }}</span>
            <span class="col-process mono">{{ r.success ? r.processMs + ' ms' : '—' }}</span>
            <span class="col-launch mono" :class="launchSpeedClass(r.launchMs)">
              {{ r.success ? r.launchMs + ' ms' : '—' }}
            </span>
            <span :class="['col-status', r.success ? 'status-ok' : 'status-fail']">
              {{ r.success ? '✓' : '✗ ' + (r.error || '失败') }}
            </span>
          </div>
        </div>

        <div v-if="!report && !benchmarkRunning" class="test-hint">
          <p>测试常见应用的启动耗时，支持冷启动和热启动对比。冷启动会先关闭应用再测量，热启动直接测量。</p>
        </div>
      </div>
    </div>

    <!-- Custom App Benchmark -->
    <div class="card">
      <div class="card-header">
        <h3 class="card-title">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 5v14M5 12h14"/>
          </svg>
          自定义应用测试
        </h3>
      </div>
      <div class="test-body">
        <div class="custom-input-row">
          <input type="text" class="custom-input" v-model="customAppPath" placeholder="输入应用路径，如 C:\Windows\notepad.exe" />
          <button class="btn btn-primary btn-sm" @click="testCustomApp" :disabled="customRunning || !customAppPath">
            <span v-if="customRunning" class="spinner-small"></span>
            <span v-else>测试</span>
          </button>
        </div>
        <div v-if="customResult" class="custom-result">
          <div class="result-row custom-result-row">
            <span class="result-label">应用名称</span>
            <span class="result-value mono">{{ customResult.appName }}</span>
          </div>
          <div class="result-row custom-result-row">
            <span class="result-label">进程创建</span>
            <span class="result-value mono">{{ customResult.processMs }} ms</span>
          </div>
          <div class="result-row custom-result-row">
            <span class="result-label">启动耗时</span>
            <span :class="['result-value', 'mono', launchSpeedClass(customResult.launchMs)]">{{ customResult.launchMs }} ms</span>
          </div>
          <div class="result-row custom-result-row">
            <span class="result-label">测试方式</span>
            <span class="result-value">{{ customResult.method === 'cold' ? '冷启动' : '热启动' }}</span>
          </div>
          <div v-if="customResult.error" class="result-error">{{ customResult.error }}</div>
        </div>
      </div>
    </div>

    <!-- Common Apps List -->
    <div class="card">
      <div class="card-header">
        <h3 class="card-title">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><polyline points="14 2 14 8 20 8"/>
          </svg>
          可测试应用列表
        </h3>
      </div>
      <div class="test-body">
        <div class="apps-grid">
          <div v-for="app in appsList" :key="app.name" :class="['app-item', app.exists === 'true' ? 'available' : 'unavailable']">
            <span class="app-name">{{ app.name }}</span>
            <span :class="['app-badge', app.exists === 'true' ? 'badge-ok' : 'badge-na']">
              {{ app.exists === 'true' ? '可用' : '不可用' }}
            </span>
            <span class="app-category">{{ categoryLabel(app.category) }}</span>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script>
import {
  BenchmarkLaunchSpeed,
  GetBootSpeedInfo,
  BenchmarkCustomAppLaunch,
  GetCommonAppsList,
} from '../../wailsjs/go/main/App'

export default {
  name: 'LaunchSpeed',
  props: ['theme'],
  data() {
    return {
      // Boot info
      bootInfo: null,
      bootLoading: false,
      startupList: [],
      showStartupDetail: false,

      // Benchmark
      report: null,
      benchmarkRunning: false,
      launchMethod: 'warm',
      doneCount: 0,
      totalCount: 0,

      // Custom app
      customAppPath: '',
      customResult: null,
      customRunning: false,

      // Apps list
      appsList: [],
    }
  },
  computed: {
  },
  methods: {
    async refreshBootInfo() {
      this.bootLoading = true
      try {
        this.bootInfo = await GetBootSpeedInfo()
        if (this.bootInfo && this.bootInfo.startupItems) {
          try {
            this.startupList = JSON.parse(this.bootInfo.startupItems)
          } catch (e) {
            this.startupList = []
          }
        }
      } catch (e) {
        console.error('Boot info error:', e)
      }
      this.bootLoading = false
    },

    async runBenchmark() {
      this.benchmarkRunning = true
      this.doneCount = 0
      this.totalCount = this.appsList.filter(a => a.exists === 'true').length
      this.report = null

      try {
        this.report = await BenchmarkLaunchSpeed(this.launchMethod)
        this.doneCount = this.report.results.filter(r => r.success).length
      } catch (e) {
        console.error('Benchmark error:', e)
      }

      this.benchmarkRunning = false
    },

    async testCustomApp() {
      if (!this.customAppPath) return
      this.customRunning = true
      this.customResult = null
      try {
        this.customResult = await BenchmarkCustomAppLaunch(this.customAppPath, this.launchMethod)
      } catch (e) {
        this.customResult = { appName: this.customAppPath, error: String(e), success: false }
      }
      this.customRunning = false
    },

    async loadAppsList() {
      try {
        this.appsList = await GetCommonAppsList()
      } catch (e) {
        console.error('Apps list error:', e)
      }
    },

    formatUptime(seconds) {
      if (!seconds) return '—'
      const h = Math.floor(seconds / 3600)
      const m = Math.floor((seconds % 3600) / 60)
      if (h > 24) {
        const d = Math.floor(h / 24)
        return `${d}天 ${h % 24}时 ${m}分`
      }
      return `${h}时 ${m}分`
    },

    truncatePath(path) {
      if (!path) return ''
      if (path.length > 60) return path.substring(0, 57) + '...'
      return path
    },

    categoryLabel(cat) {
      const labels = {
        system: '系统',
        browser: '浏览器',
        office: '办公',
        media: '媒体',
        custom: '自定义',
      }
      return labels[cat] || cat
    },

    launchSpeedClass(ms) {
      if (!ms) return ''
      if (ms < 200) return 'speed-fast'
      if (ms < 1000) return 'speed-normal'
      return 'speed-slow'
    },
  },
  mounted() {
    this.refreshBootInfo()
    this.loadAppsList()
  }
}
</script>

<style scoped>
.launch-speed-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* Boot info */
.boot-info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 10px;
  margin-bottom: 12px;
}

.boot-stat-card {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 12px 16px;
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  border-radius: 8px;
}

.boot-stat-label {
  font-size: 11px;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.boot-stat-value {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
}

.boot-stat-value.mono {
  font-family: 'Cascadia Code', 'Fira Code', monospace;
  font-size: 16px;
}

.boot-stat-value.highlight {
  color: var(--lenovo-red, #E63F32);
}

.boot-stat-value.ok {
  color: #4ade80;
}

.boot-stat-value.warn {
  color: #fbbf24;
}

.boot-stat-hint {
  font-size: 11px;
  color: var(--text-muted);
}

/* Startup section */
.startup-section {
  margin-top: 12px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  cursor: pointer;
  border-radius: 6px;
  transition: background 0.15s;
}

.section-header:hover {
  background: var(--bg-elevated);
}

.section-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
}

.section-toggle {
  font-size: 12px;
  color: var(--text-muted);
}

.startup-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-top: 8px;
  padding: 8px 12px;
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  border-radius: 8px;
}

.startup-item {
  display: grid;
  grid-template-columns: 140px 100px 1fr;
  gap: 8px;
  padding: 6px 8px;
  border-radius: 4px;
  font-size: 13px;
}

.startup-item:hover {
  background: var(--bg);
}

.startup-name {
  font-weight: 600;
  color: var(--text-primary);
}

.startup-source {
  font-size: 11px;
  color: var(--text-muted);
  background: var(--bg);
  padding: 2px 8px;
  border-radius: 4px;
  text-align: center;
}

.startup-path {
  font-family: 'Cascadia Code', 'Fira Code', monospace;
  font-size: 12px;
  color: var(--text-muted);
}

/* Summary grid */
.summary-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 10px;
  margin-bottom: 16px;
}

.summary-card {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 12px 16px;
  border-radius: 8px;
  border: 1px solid var(--border);
}

.summary-card.best {
  background: rgba(74, 222, 128, 0.08);
  border-color: rgba(74, 222, 128, 0.3);
}

.summary-card.worst {
  background: rgba(248, 113, 113, 0.08);
  border-color: rgba(248, 113, 113, 0.3);
}

.summary-card.avg {
  background: rgba(96, 165, 250, 0.08);
  border-color: rgba(96, 165, 250, 0.3);
}

.summary-card.total {
  background: var(--bg-elevated);
}

.summary-label {
  font-size: 11px;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.summary-value {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
}

.summary-value.mono {
  font-family: 'Cascadia Code', 'Fira Code', monospace;
}

/* Results table */
.results-table {
  display: flex;
  flex-direction: column;
  gap: 2px;
  border: 1px solid var(--border);
  border-radius: 8px;
  overflow: hidden;
}

.result-header, .result-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1.2fr 1.2fr 1.5fr;
  gap: 8px;
  padding: 8px 12px;
  font-size: 13px;
}

.result-header {
  background: var(--bg-elevated);
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  font-size: 11px;
}

.result-row {
  border-bottom: 1px solid var(--border);
}

.result-row:last-child {
  border-bottom: none;
}

.result-row.success:hover {
  background: var(--bg-elevated);
}

.result-row.failed {
  opacity: 0.5;
}

.col-name {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
}

.cat-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.cat-system { background: #60a5fa; }
.cat-browser { background: #fbbf24; }
.cat-office { background: #a78bfa; }
.cat-media { background: #f87171; }
.cat-custom { background: #4ade80; }

.col-process, .col-launch {
  font-family: 'Cascadia Code', 'Fira Code', monospace;
  text-align: right;
}

.speed-fast {
  color: #4ade80;
  font-weight: 700;
}

.speed-normal {
  color: var(--text-primary);
}

.speed-slow {
  color: #f87171;
  font-weight: 700;
}

.status-ok {
  color: #4ade80;
}

.status-fail {
  color: #f87171;
  font-size: 12px;
}

/* Custom app */
.custom-input-row {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.custom-input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid var(--border);
  border-radius: 6px;
  background: var(--bg-elevated);
  color: var(--text-primary);
  font-family: 'Cascadia Code', 'Fira Code', monospace;
  font-size: 13px;
}

.custom-input:focus {
  border-color: var(--lenovo-red, #E63F32);
  outline: none;
}

.custom-result {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 12px 16px;
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  border-radius: 8px;
}

.custom-result-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

/* Apps grid */
.apps-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 6px;
}

.app-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid var(--border);
  background: var(--bg-elevated);
  font-size: 13px;
}

.app-item.unavailable {
  opacity: 0.4;
}

.app-name {
  font-weight: 600;
  flex: 1;
}

.app-badge {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 700;
}

.badge-ok {
  background: rgba(74, 222, 128, 0.15);
  color: #4ade80;
  border: 1px solid rgba(74, 222, 128, 0.3);
}

.badge-na {
  background: rgba(248, 113, 113, 0.15);
  color: #f87171;
  border: 1px solid rgba(248, 113, 113, 0.3);
}

.app-category {
  font-size: 10px;
  color: var(--text-muted);
}

/* Common */
.test-body {
  margin-top: 12px;
}

.test-loading {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px;
  color: var(--text-muted);
  font-size: 13px;
}

.test-hint {
  padding: 12px 16px;
  color: var(--text-muted);
  font-size: 13px;
  background: var(--bg-elevated);
  border-radius: 6px;
}

.spinner-small {
  display: inline-block;
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255,255,255,0.15);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  flex-shrink: 0;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.header-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.method-select {
  padding: 6px 10px;
  border: 1px solid var(--border);
  border-radius: 6px;
  background: var(--bg-elevated);
  color: var(--text-primary);
  font-size: 12px;
  font-weight: 600;
}

.method-select:focus {
  border-color: var(--lenovo-red, #E63F32);
  outline: none;
}

/* Buttons */
.btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s;
  font-family: inherit;
  border: 1px solid transparent;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: var(--lenovo-red, #E63F32);
  color: white;
  border-color: var(--lenovo-red, #E63F32);
}

.btn-primary:hover:not(:disabled) {
  background: #d63018;
}

.btn-secondary {
  background: var(--bg-elevated);
  color: var(--text-secondary);
  border-color: var(--border);
}

.btn-secondary:hover:not(:disabled) {
  border-color: var(--lenovo-red, #E63F32);
  color: var(--lenovo-red, #E63F32);
}

.btn-sm {
  padding: 4px 10px;
  font-size: 11px;
}

.result-error {
  padding: 8px 12px;
  background: rgba(248, 113, 113, 0.1);
  border: 1px solid rgba(248, 113, 113, 0.3);
  border-radius: 6px;
  color: #f87171;
  font-size: 12px;
}

.result-label {
  font-size: 12px;
  color: var(--text-muted);
  width: 80px;
  flex-shrink: 0;
}

.result-value {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
}

.result-value.mono {
  font-family: 'Cascadia Code', 'Fira Code', monospace;
}
</style>
