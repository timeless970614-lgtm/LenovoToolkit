<template>
  <div class="test-func-page">

    <!-- Test Tabs (horizontal) -->
    <div class="test-tabs-bar">
      <div 
        v-for="tab in testTabs" :key="tab.id"
        :class="['test-tab', activeTab === tab.id ? 'active' : '', tab.done ? (tab.pass ? 'pass' : 'warn') : '']"
        @click="activeTab = tab.id"
      >
        <svg v-if="tab.id === 'launch'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="5 3 19 12 5 21 5 3"/><circle cx="12" cy="12" r="10"/></svg>
        <svg v-else-if="tab.id === 'brightness'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>
        <svg v-else-if="tab.id === 'fnq'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="6" width="20" height="12" rx="2"/><path d="M6 12h4"/></svg>
        <svg v-else-if="tab.id === 'mode'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/></svg>
        <span class="tab-label">{{ tab.label }}</span>
        <span v-if="tab.done" class="tab-badge">{{ tab.pass ? '✓' : '~' }}</span>
      </div>
    </div>

    <!-- Tab Content Area -->
    <div class="test-tab-content">

      <!-- Launch Speed Tab -->
      <div v-if="activeTab === 'launch'" class="card launch-speed-card">
        <div class="card-header">
          <h3 class="card-title">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="5 3 19 12 5 21 5 3"/><circle cx="12" cy="12" r="10"/></svg>
            应用启动速度
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
          <!-- Boot info -->
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
              <span class="summary-value mono">{{ report.avgLaunch.toFixed(2) }} s</span>
            </div>
            <div class="summary-card total">
              <span class="summary-label">总测试时间</span>
              <span class="summary-value mono">{{ report.totalTime.toFixed(2) }} s</span>
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
            <div v-for="r in report.results" :key="r.appName" :class="['result-row', r.success ? 'success' : 'failed']">
              <span class="col-name"><span :class="['cat-dot', 'cat-' + r.category]"></span>{{ r.appName }}</span>
              <span class="col-cat">{{ categoryLabel(r.category) }}</span>
              <span class="col-process mono">{{ r.success ? r.processMs.toFixed(2) + ' s' : '—' }}</span>
              <span class="col-launch mono" :class="launchSpeedClass(r.launchMs)">{{ r.success ? r.launchMs.toFixed(2) + ' s' : '—' }}</span>
              <span :class="['col-status', r.success ? 'status-ok' : 'status-fail']">{{ r.success ? '✓' : '✗ ' + (r.error || '失败') }}</span>
            </div>
          </div>
          <!-- Custom app -->
          <div class="custom-app-section">
            <div class="custom-title">🔧 自定义应用测试</div>
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
                <span class="result-value mono">{{ customResult.processMs.toFixed(2) }} s</span>
              </div>
              <div class="result-row custom-result-row">
                <span class="result-label">启动耗时</span>
                <span :class="['result-value', 'mono', launchSpeedClass(customResult.launchMs)]">{{ customResult.launchMs.toFixed(2) }} s</span>
              </div>
              <div v-if="customResult.error" class="result-error">{{ customResult.error }}</div>
            </div>
          </div>
          <div v-if="!report && !benchmarkRunning" class="test-hint">
            <p>测试常见应用的启动耗时，支持冷启动和热启动对比。冷启动会先关闭应用再测量，热启动直接测量。</p>
          </div>
        </div>
      </div>

      <!-- Brightness Tab -->
      <div v-if="activeTab === 'brightness'" class="card">
        <div class="card-header">
          <h3 class="card-title">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle></circle>
              <path></path>
            </svg>
            屏幕亮度测试
          </h3>
          <span v-if="brightness.done" :class="['test-badge', brightness.result?.writable ? 'pass' : (brightness.result?.readable ? 'warn' : 'fail')]">
            {{ brightness.result?.writable ? 'PASS' : (brightness.result?.readable ? 'READ-ONLY' : 'FAIL') }}
          </span>
        </div>
        <div class="test-body">
          <div v-if="brightness.running" class="test-loading">
            <div class="spinner-small"></div>
            <span>正在读写亮度...</span>
          </div>
          <div v-else-if="brightness.result" class="test-result-block">
            <div class="result-row">
              <span class="result-label">当前亮度</span>
              <div class="brightness-bar-wrap">
                <div class="brightness-bar" :style="{ width: brightness.result.currentBrightness + '%' }"></div>
              </div>
              <span class="result-value mono">{{ brightness.result.currentBrightness }}%</span>
            </div>
            <div class="result-row" v-if="brightness.result.testBrightness">
              <span class="result-label">测试设置</span>
              <div class="brightness-bar-wrap">
                <div class="brightness-bar test-bar" :style="{ width: brightness.result.testBrightness + '%' }"></div>
              </div>
              <span class="result-value mono">{{ brightness.result.testBrightness }}%</span>
            </div>
            <div class="result-chips">
              <span :class="['chip', brightness.result.readable ? 'chip-ok' : 'chip-fail']">
                {{ brightness.result.readable ? '✓ 可读' : '✗ 不可读' }}
              </span>
              <span :class="['chip', brightness.result.writable ? 'chip-ok' : 'chip-fail']">
                {{ brightness.result.writable ? '✓ 可写' : '✗ 不可写' }}
              </span>
              <span v-if="brightness.result.error" class="chip chip-fail">
                {{ brightness.result.error }}
              </span>
            </div>
          </div>
          <div v-else class="test-hint">
            <p>测试屏幕亮度读取和设置功能（仅支持笔记本内置屏幕）</p>
          </div>
          <div class="test-actions">
            <button class="btn btn-primary btn-sm" @click="testBrightness" :disabled="brightness.running">
              <span v-if="brightness.running" class="spinner-small"></span>
              <span v-else>运行测试</span>
            </button>
          </div>
        </div>
      </div>

      <!-- FN+Q Tab -->
      <div v-if="activeTab === 'fnq'" class="card">
        <div class="card-header">
          <h3 class="card-title">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect></rect>
              <path></path>
            </svg>
            FN+Q 热键测试
          </h3>
          <span v-if="fnq.done" :class="['test-badge', fnq.result?.supported ? 'pass' : 'warn']">
            {{ fnq.result?.supported ? 'PASS' : 'N/A' }}
          </span>
        </div>
        <div class="test-body">
          <div v-if="fnq.running" class="test-loading">
            <div class="spinner-small"></div>
            <span>{{ fnqPhase }}</span>
          </div>
          <div v-else-if="fnq.result" class="test-result-block">
            <!-- Key simulation phases -->
            <div v-if="fnq.result.modeBefore" class="fnq-phases">
              <div class="fnq-phase-step">
                <span class="phase-label">① 按键前模式</span>
                <span :class="['phase-value', 'mono']">{{ fnq.result.modeBefore }}</span>
              </div>
              <div class="fnq-phase-step">
                <span class="phase-label">② 第一次 FN+Q</span>
                <span :class="['phase-value', 'mono', fnq.result.modeChanged1 ? 'text-ok' : 'text-muted']">{{ fnq.result.modeAfter1 }} {{ fnq.result.modeChanged1 ? '✓ 已切换' : '(未切换)' }}</span>
              </div>
              <div class="fnq-phase-step">
                <span class="phase-label">③ 空闲 10 秒</span>
                <span class="phase-value mono">idle</span>
              </div>
              <div class="fnq-phase-step">
                <span class="phase-label">④ 第二次 FN+Q</span>
                <span :class="['phase-value', 'mono', fnq.result.modeChanged2 ? 'text-ok' : 'text-muted']">{{ fnq.result.modeAfter2 }} {{ fnq.result.modeChanged2 ? '✓ 已切换' : '(未切换)' }}</span>
              </div>
              <div class="fnq-phase-step">
                <span class="phase-label">⑤ 恢复原模式</span>
                <span class="phase-value mono">{{ fnq.result.restoredTo }}</span>
              </div>
            </div>
            <div class="test-info-grid">
              <div class="test-info-item">
                <span class="test-info-label">当前模式</span>
                <span class="test-info-value mono">{{ fnq.result.currentMode || 'N/A' }}</span>
              </div>
              <div class="test-info-item">
                <span class="test-info-label">FN+Q 支持</span>
                <span :class="['test-info-value', 'mono', fnq.result.supported ? 'text-ok' : 'text-muted']">
                  {{ fnq.result.supported ? '支持' : '不支持' }}
                </span>
              </div>
              <div class="test-info-item">
                <span class="test-info-label">热键可用</span>
                <span :class="['test-info-value', 'mono', fnq.result.hotkeyAvailable ? 'text-ok' : 'text-muted']">
                  {{ fnq.result.hotkeyAvailable ? '可用' : '不可用' }}
                </span>
              </div>
              <div class="test-info-item">
                <span class="test-info-label">Dispatcher Func</span>
                <span class="test-info-value mono">{{ fnq.result.dispatcherFunction || 'N/A' }}</span>
              </div>
              <div class="test-info-item" v-if="fnq.result.hotkeyStatus">
                <span class="test-info-label">热键状态</span>
                <span class="test-info-value mono">{{ fnq.result.hotkeyStatus }}</span>
              </div>
              <div class="test-info-item">
                <span class="test-info-label">DCC 功能</span>
                <span :class="['test-info-value', 'mono', fnq.result.dccCapability ? 'text-ok' : 'text-muted']">
                  {{ fnq.result.dccCapability ? '支持' : '不支持' }}
                </span>
              </div>
              <div class="test-info-item">
                <span class="test-info-label">GEEK 功能</span>
                <span :class="['test-info-value', 'mono', fnq.result.geekCapability ? 'text-ok' : 'text-muted']">
                  {{ fnq.result.geekCapability ? '支持' : '不支持' }}
                </span>
              </div>
            </div>
            <div v-if="fnq.result.error" class="result-error">
              {{ fnq.result.error }}
            </div>
          </div>
          <div v-else class="test-hint">
            <p>检测 FN+Q 快捷键功能支持状态及热键可用性</p>
          </div>
          <div class="test-actions">
            <button class="btn btn-primary btn-sm" @click="testFNQ" :disabled="fnq.running">
              <span v-if="fnq.running" class="spinner-small"></span>
              <span v-else>运行测试</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Mode Switch Tab -->
      <div v-if="activeTab === 'mode'" class="card">
        <div class="card-header">
          <h3 class="card-title">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline></polyline><polyline></polyline>
            </svg>
            模式切换测试
          </h3>
          <span v-if="mode.done" :class="['test-badge', mode.anyVerified ? 'pass' : 'warn']">
            {{ mode.anyVerified ? 'PASS' : 'N/A' }}
          </span>
        </div>
        <div class="test-body">
          <div v-if="mode.running" class="test-loading">
            <div class="spinner-small"></div>
            <span>正在测试模式切换...</span>
          </div>
          <div v-else-if="mode.result" class="test-result-block">
            <div class="result-row">
              <span class="result-label">当前模式</span>
              <span class="result-value mono">{{ mode.result.currentMode || 'N/A' }}</span>
            </div>
            <div v-if="mode.result.attempts && mode.result.attempts.length > 0" class="mode-attempts">
              <div
                v-for="(attempt, idx) in mode.result.attempts"
                :key="idx"
                :class="['mode-attempt-row', attempt.skipped ? 'skipped' : (attempt.verified ? 'verified' : 'failed')]"
              >
                <span class="attempt-mode mono">{{ attempt.mode }}</span>
                <span class="attempt-status">
                  <span v-if="attempt.skipped" class="status-text skipped">跳过 ({{ attempt.reason }})</span>
                  <span v-else-if="attempt.verified" class="status-text verified">✓ 成功</span>
                  <span v-else class="status-text failed">✗ {{ attempt.error || attempt.note }}</span>
                </span>
                <span v-if="attempt.switchTimeMs" class="attempt-time mono">
                  {{ attempt.switchTimeMs }} ms
                </span>
              </div>
            </div>
            <div v-if="mode.result.error" class="result-error">
              {{ mode.result.error }}
            </div>
          </div>
          <div v-else class="test-hint">
            <p>测试热管理模式切换功能（BSM / STD / AQM）</p>
          </div>
          <div class="test-actions">
            <button class="btn btn-primary btn-sm" @click="testModeSwitch" :disabled="mode.running">
              <span v-if="mode.running" class="spinner-small"></span>
              <span v-else>运行测试</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import {
  TestBrightness,
  TestFNQ,
  TestModeSwitch,
  BenchmarkLaunchSpeed,
  GetBootSpeedInfo,
  BenchmarkCustomAppLaunch,
  GetCommonAppsList,
} from '../../wailsjs/go/main/App'

export default {
  name: 'TestFunction',
  data() {
    return {
      activeTab: 'launch',
      // Launch speed
      bootInfo: null,
      bootLoading: false,
      startupList: [],
      showStartupDetail: false,
      report: null,
      benchmarkRunning: false,
      launchMethod: 'warm',
      doneCount: 0,
      totalCount: 0,
      customAppPath: '',
      customResult: null,
      customRunning: false,
      appsList: [],
      brightness: { running: false, done: false, result: null },
      fnq: { running: false, done: false, result: null },
      fnqPhase: '准备测试 FN+Q...',
      mode: { running: false, done: false, result: null, anyVerified: false },
    }
  },
  computed: {
    testTabs() {
      return [
        { id: 'launch', label: '启动速度', done: this.report !== null, pass: this.report?.fastest },
        { id: 'brightness', label: '亮度', done: this.brightness.done, pass: this.brightness.result?.writable },
        { id: 'fnq', label: 'FN+Q', done: this.fnq.done, pass: this.fnq.result?.supported },
        { id: 'mode', label: '模式切换', done: this.mode.done, pass: this.mode.anyVerified },
      ]
    },
  },
  methods: {
    // ===== Launch Speed =====
    async refreshBootInfo() {
      this.bootLoading = true
      try {
        this.bootInfo = await GetBootSpeedInfo()
        if (this.bootInfo && this.bootInfo.startupItems) {
          try { this.startupList = JSON.parse(this.bootInfo.startupItems) } catch (e) { this.startupList = [] }
        }
      } catch (e) { console.error('Boot info error:', e) }
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
      } catch (e) { console.error('Benchmark error:', e) }
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
      try { this.appsList = await GetCommonAppsList() } catch (e) { console.error('Apps list error:', e) }
    },
    formatUptime(seconds) {
      if (!seconds) return '—'
      const h = Math.floor(seconds / 3600), m = Math.floor((seconds % 3600) / 60)
      if (h > 24) { const d = Math.floor(h / 24); return `${d}天 ${h % 24}时 ${m}分` }
      return `${h}时 ${m}分`
    },
    truncatePath(path) {
      if (!path) return ''
      return path.length > 60 ? path.substring(0, 57) + '...' : path
    },
    categoryLabel(cat) {
      const labels = { system: '系统', browser: '浏览器', office: '办公', media: '媒体', custom: '自定义' }
      return labels[cat] || cat
    },
    launchSpeedClass(ms) {
      if (!ms) return ''
      if (ms < 0.2) return 'speed-fast'
      if (ms < 1.0) return 'speed-normal'
      return 'speed-slow'
    },

    // ===== Brightness =====
    async testBrightness() {
      this.brightness.running = true
      this.brightness.done = false
      try {
        const result = await TestBrightness()
        this.brightness.result = result
      } catch (e) {
        this.brightness.result = { error: String(e), readable: false, writable: false }
      }
      this.brightness.running = false
      this.brightness.done = true
    },

    async testFNQ() {
      this.fnq.running = true
      this.fnq.done = false
      this.fnqPhase = '① 读取初始模式...'
      try {
        // Phase 1: Read initial mode (backend handles delay internally)
        this.fnqPhase = '① 模拟第一次 FN+Q 按键...'
        // The backend TestFNQ now simulates keypresses internally with 10s idle
        // We update phase text on a timer to give visual feedback
        const phaseTimer = setInterval(() => {
          if (this.fnqPhase.includes('①')) this.fnqPhase = '② 等待模式切换 (3s)...'
          else if (this.fnqPhase.includes('②')) this.fnqPhase = '③ 空闲等待 (10s)...'
          else if (this.fnqPhase.includes('③')) this.fnqPhase = '④ 模拟第二次 FN+Q 按键...'
          else if (this.fnqPhase.includes('④')) this.fnqPhase = '⑤ 等待模式切换 (3s)...'
          else if (this.fnqPhase.includes('⑤')) this.fnqPhase = '⑥ 恢复原模式...'
          else clearInterval(phaseTimer)
        }, 3000)
        const result = await TestFNQ()
        clearInterval(phaseTimer)
        this.fnq.result = result
      } catch (e) {
        this.fnq.result = { error: String(e) }
      }
      this.fnq.running = false
      this.fnq.done = true
      this.fnqPhase = '测试完成'
    },

    async testModeSwitch() {
      this.mode.running = true
      this.mode.done = false
      this.mode.anyVerified = false
      try {
        const result = await TestModeSwitch()
        this.mode.result = result
        if (result.attempts) {
          this.mode.anyVerified = result.attempts.some(a => a.verified)
        }
      } catch (e) {
        this.mode.result = { error: String(e) }
      }
      this.mode.running = false
      this.mode.done = true
    },


  },
  mounted() {
    this.refreshBootInfo()
    this.loadAppsList()
  }
}
</script>

<style scoped>
.test-func-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.test-body {
  margin-top: 12px;
}

/* ===== Launch Speed ===== */
.launch-speed-card .test-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.boot-info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 10px;
}
.boot-stat-card {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 10px 14px;
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
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
}
.boot-stat-value.mono { font-family: 'Cascadia Code', 'Fira Code', monospace; font-size: 15px; }
.boot-stat-value.highlight { color: var(--lenovo-red, #E63F32); }
.boot-stat-value.ok { color: #4ade80; }
.boot-stat-value.warn { color: #fbbf24; }
.boot-stat-hint { font-size: 11px; color: var(--text-muted); }
.startup-section { margin-top: 8px; }
.section-header {
  display: flex; align-items: center; gap: 8px;
  padding: 8px 12px; cursor: pointer; border-radius: 6px; transition: background 0.15s;
}
.section-header:hover { background: var(--bg-elevated); }
.section-title { font-size: 13px; font-weight: 600; color: var(--text-primary); }
.section-toggle { font-size: 12px; color: var(--text-muted); }
.startup-list {
  display: flex; flex-direction: column; gap: 4px; margin-top: 8px;
  padding: 8px 12px; background: var(--bg-elevated); border: 1px solid var(--border); border-radius: 8px;
}
.startup-item {
  display: grid; grid-template-columns: 140px 100px 1fr; gap: 8px;
  padding: 6px 8px; border-radius: 4px; font-size: 13px;
}
.startup-item:hover { background: var(--bg); }
.startup-name { font-weight: 600; color: var(--text-primary); }
.startup-source {
  font-size: 11px; color: var(--text-muted); background: var(--bg); padding: 2px 8px;
  border-radius: 4px; text-align: center;
}
.startup-path { font-family: 'Cascadia Code', 'Fira Code', monospace; font-size: 12px; color: var(--text-muted); }
.summary-grid {
  display: grid; grid-template-columns: repeat(auto-fill, minmax(180px, 1fr)); gap: 10px;
}
.summary-card {
  display: flex; flex-direction: column; gap: 4px;
  padding: 10px 14px; border-radius: 8px; border: 1px solid var(--border);
}
.summary-card.best { background: rgba(74, 222, 128, 0.08); border-color: rgba(74, 222, 128, 0.3); }
.summary-card.worst { background: rgba(248, 113, 113, 0.08); border-color: rgba(248, 113, 113, 0.3); }
.summary-card.avg { background: rgba(96, 165, 250, 0.08); border-color: rgba(96, 165, 250, 0.3); }
.summary-card.total { background: var(--bg-elevated); }
.summary-label { font-size: 11px; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.05em; }
.summary-value { font-size: 16px; font-weight: 700; color: var(--text-primary); }
.summary-value.mono { font-family: 'Cascadia Code', 'Fira Code', monospace; }
.results-table {
  display: flex; flex-direction: column; gap: 2px;
  border: 1px solid var(--border); border-radius: 8px; overflow: hidden;
}
.result-header, .results-table .result-row {
  display: grid; grid-template-columns: 2fr 1fr 1.2fr 1.2fr 1.5fr; gap: 8px;
  padding: 8px 12px; font-size: 13px;
}
.results-table .result-header {
  background: var(--bg-elevated); font-weight: 600; color: var(--text-muted);
  text-transform: uppercase; letter-spacing: 0.05em; font-size: 11px;
}
.results-table .result-row { border-bottom: 1px solid var(--border); }
.results-table .result-row:last-child { border-bottom: none; }
.results-table .result-row.success:hover { background: var(--bg-elevated); }
.results-table .result-row.failed { opacity: 0.5; }
.col-name { display: flex; align-items: center; gap: 6px; font-weight: 600; }
.col-cat { font-size: 12px; }
.col-process, .col-launch { font-family: 'Cascadia Code', 'Fira Code', monospace; text-align: right; }
.col-status { font-size: 12px; }
.cat-dot { width: 8px; height: 8px; border-radius: 50%; flex-shrink: 0; }
.cat-system { background: #60a5fa; }
.cat-browser { background: #fbbf24; }
.cat-office { background: #a78bfa; }
.cat-media { background: #f87171; }
.cat-custom { background: #4ade80; }
.speed-fast { color: #4ade80; font-weight: 700; }
.speed-normal { color: var(--text-primary); }
.speed-slow { color: #f87171; font-weight: 700; }
.status-ok { color: #4ade80; }
.status-fail { color: #f87171; font-size: 12px; }
.custom-app-section {
  padding-top: 12px; border-top: 1px solid var(--border);
}
.custom-title { font-size: 14px; font-weight: 600; color: var(--text-primary); margin-bottom: 10px; }
.custom-input-row { display: flex; gap: 8px; margin-bottom: 10px; }
.custom-input {
  flex: 1; padding: 8px 12px; border: 1px solid var(--border); border-radius: 6px;
  background: var(--bg-elevated); color: var(--text-primary);
  font-family: 'Cascadia Code', 'Fira Code', monospace; font-size: 13px;
}
.custom-input:focus { border-color: var(--lenovo-red, #E63F32); outline: none; }
.custom-result {
  display: flex; flex-direction: column; gap: 6px;
  padding: 10px 14px; background: var(--bg-elevated); border: 1px solid var(--border); border-radius: 8px;
}
.custom-result-row { display: flex; align-items: center; gap: 12px; }
.method-select {
  padding: 6px 10px; border: 1px solid var(--border); border-radius: 6px;
  background: var(--bg-elevated); color: var(--text-primary); font-size: 12px; font-weight: 600;
}
.method-select:focus { border-color: var(--lenovo-red, #E63F32); outline: none; }

/* ===== FN+Q Phases ===== */
.fnq-phases {
  display: flex; flex-direction: column; gap: 8px;
  padding: 12px 16px; background: var(--bg-elevated); border: 1px solid var(--border); border-radius: 8px;
  margin-bottom: 12px;
}
.fnq-phase-step {
  display: flex; align-items: center; gap: 12px;
  padding: 6px 10px; border-radius: 6px;
}
.fnq-phase-step:hover { background: var(--bg); }
.phase-label { font-size: 13px; font-weight: 600; color: var(--text-primary); min-width: 120px; }
.phase-value { font-size: 14px; }

/* ===== Tab Bar ===== */
.test-tabs-bar {
  display: flex;
  gap: 2px;
  padding: 4px;
  background: var(--bg-elevated);
  border-radius: 10px;
  border: 1px solid var(--border);
}

.test-tab {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
  background: transparent;
  position: relative;
  user-select: none;
}

.test-tab:hover {
  background: var(--bg);
  color: var(--text-primary);
}

.test-tab.active {
  background: var(--bg);
  color: var(--lenovo-red, #E63F32);
  border: 1px solid rgba(230, 63, 50, 0.3);
  box-shadow: 0 2px 8px rgba(230, 63, 50, 0.1);
}

.test-tab.pass .tab-badge {
  color: #4ade80;
}

.test-tab.warn .tab-badge {
  color: #fbbf24;
}

.tab-label {
  white-space: nowrap;
}

.tab-badge {
  font-size: 11px;
  font-weight: 700;
}

/* ===== Tab Content ===== */
.test-tab-content {
  min-height: 200px;
}

/* Test Info Grid */
.test-info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 8px;
  margin-bottom: 12px;
}

.test-info-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 8px 12px;
  background: var(--bg-elevated);
  border-radius: 6px;
  border: 1px solid var(--border);
}

.test-info-label {
  font-size: 11px;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.test-info-value {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.test-info-value.mono {
  font-family: 'Cascadia Code', 'Fira Code', monospace;
  font-size: 13px;
}

.test-info-value.highlight {
  color: var(--accent);
}

.test-info-value.text-ok {
  color: #4ade80;
}

.test-info-value.text-muted {
  color: var(--text-muted);
}

/* Loading state */
.test-loading {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px;
  color: var(--text-muted);
  font-size: 13px;
}

/* Hint */
.test-hint {
  padding: 12px 16px;
  color: var(--text-muted);
  font-size: 13px;
  background: var(--bg-elevated);
  border-radius: 6px;
}

/* Result */
.test-result-block {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 12px 16px;
  background: var(--bg-elevated);
  border-radius: 6px;
}

.result-row {
  display: flex;
  align-items: center;
  gap: 12px;
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

.brightness-bar-wrap {
  flex: 1;
  height: 8px;
  background: var(--bg);
  border-radius: 4px;
  overflow: hidden;
}

.brightness-bar {
  height: 100%;
  background: linear-gradient(90deg, #60a5fa, #3b82f6);
  border-radius: 4px;
  transition: width 0.5s ease;
}

.brightness-bar.test-bar {
  background: linear-gradient(90deg, #facc15, #f59e0b);
}

/* Chips */
.result-chips {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.chip {
  padding: 3px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.chip-ok {
  background: rgba(74, 222, 128, 0.15);
  color: #4ade80;
  border: 1px solid rgba(74, 222, 128, 0.3);
}

.chip-fail {
  background: rgba(248, 113, 113, 0.15);
  color: #f87171;
  border: 1px solid rgba(248, 113, 113, 0.3);
}

/* Result error */
.result-error {
  padding: 8px 12px;
  background: rgba(248, 113, 113, 0.1);
  border: 1px solid rgba(248, 113, 113, 0.3);
  border-radius: 6px;
  color: #f87171;
  font-size: 12px;
}

/* Mode attempts */
.mode-attempts {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-top: 4px;
}

.mode-attempt-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 13px;
}

.mode-attempt-row.verified {
  background: rgba(74, 222, 128, 0.08);
  border: 1px solid rgba(74, 222, 128, 0.2);
}

.mode-attempt-row.failed {
  background: rgba(248, 113, 113, 0.08);
  border: 1px solid rgba(248, 113, 113, 0.3);
}

.mode-attempt-row.skipped {
  background: var(--bg);
  border: 1px solid var(--border);
  opacity: 0.7;
}

.attempt-mode {
  font-weight: 600;
  width: 60px;
  font-size: 13px;
}

.attempt-status {
  flex: 1;
  font-size: 12px;
}

.attempt-time {
  font-size: 11px;
  color: var(--text-muted);
}

.status-text.verified {
  color: #4ade80;
}

.status-text.failed {
  color: #f87171;
}

.status-text.skipped {
  color: var(--text-muted);
}

/* Test badge */
.test-badge {
  padding: 2px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.05em;
}

.test-badge.pass {
  background: rgba(74, 222, 128, 0.15);
  color: #4ade80;
  border: 1px solid rgba(74, 222, 128, 0.3);
}

.test-badge.warn {
  background: rgba(250, 204, 21, 0.15);
  color: #fbbf24;
  border: 1px solid rgba(250, 204, 21, 0.3);
}

.test-badge.fail {
  background: rgba(248, 113, 113, 0.15);
  color: #f87171;
  border: 1px solid rgba(248, 113, 113, 0.3);
}

/* Actions */
.test-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 10px;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.result-text {
  font-size: 12px;
}

.result-text.success {
  color: #4ade80;
}

.result-text.error {
  color: #f87171;
}

.loading-text {
  color: var(--text-muted);
}

/* Spinner */
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
</style>
