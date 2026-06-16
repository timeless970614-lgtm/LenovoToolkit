<template>
  <div class="test-func-page">

    <!-- Header Info -->
    <div class="card">
      <div class="card-header">
        <h3 class="card-title">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polygon points="5 3 19 12 5 21 5 3"/>
          </svg>
          Test Function
        </h3>
        <button class="btn btn-secondary btn-sm" @click="runAllTests" :disabled="running">
          <span v-if="running" class="spinner-small"></span>
          <span v-else>▶ Run All</span>
        </button>
      </div>
      <p class="test-desc">对核心功能模块进行诊断测试，验证各子系统是否正常工作。</p>
    </div>

    <!-- App Launch Speed Test -->
    <div class="card">
      <div class="card-header">
        <h3 class="card-title">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
          </svg>
          应用启动速度
        </h3>
        <span v-if="launch.done" :class="['test-badge', launch.fast ? 'pass' : 'warn']">
          {{ launch.fast ? 'PASS' : 'OK' }}
        </span>
      </div>
      <div class="test-body">
        <div class="test-info-grid">
          <div class="test-info-item">
            <span class="test-info-label">Build 版本</span>
            <span class="test-info-value mono">v1.0.21</span>
          </div>
          <div class="test-info-item">
            <span class="test-info-label">构建日期</span>
            <span class="test-info-value mono">2026-06-07</span>
          </div>
          <div class="test-info-item">
            <span class="test-info-label">Framework</span>
            <span class="test-info-value mono">Wails v2.12.0</span>
          </div>
          <div class="test-info-item">
            <span class="test-info-label">启动耗时</span>
            <span :class="['test-info-value', 'mono', 'highlight', launch.done ? '' : 'loading-text']">
              {{ launch.elapsed ? launch.elapsed + ' ms' : '—' }}
            </span>
          </div>
        </div>
        <div class="test-actions">
          <button class="btn btn-primary btn-sm" @click="testLaunchSpeed" :disabled="launch.running">
            <span v-if="launch.running" class="spinner-small"></span>
            <span v-else>测试启动速度</span>
          </button>
          <span v-if="launch.result" :class="['result-text', launch.result.success ? 'success' : 'error']">
            {{ launch.result.message }}
          </span>
        </div>
      </div>
    </div>

    <!-- Brightness Test -->
    <div class="card">
      <div class="card-header">
        <h3 class="card-title">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="5"/>
            <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
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

    <!-- FN+Q Test -->
    <div class="card">
      <div class="card-header">
        <h3 class="card-title">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="2" y="7" width="20" height="14" rx="2"/>
            <path d="M16 7V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v2"/>
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
          <span>正在检测 FN+Q 功能...</span>
        </div>
        <div v-else-if="fnq.result" class="test-result-block">
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

    <!-- Mode Switch Test -->
    <div class="card">
      <div class="card-header">
        <h3 class="card-title">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/>
          </svg>
          模式切换测试
        </h3>
        <span v-if="mode.done" :class="['test-badge', mode.anyVerified ? 'pass' : 'warn']">
          {{ mode.anyVerified ? 'PASS' : 'N/A' }}
        </span>
      </div>
      <div class="test-body">
        <div v-if="mode.running" class="test-loading">
          <div class="spinner"></div>
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
</template>

<script>
import { TestBrightness, TestFNQ, TestModeSwitch } from '../../wailsjs/go/main/App'

export default {
  name: 'TestFunction',
  data() {
    return {
      launch: { running: false, done: false, elapsed: null, fast: false, result: null },
      brightness: { running: false, done: false, result: null },
      fnq: { running: false, done: false, result: null },
      mode: { running: false, done: false, result: null, anyVerified: false },
    }
  },
  methods: {
    async testLaunchSpeed() {
      this.launch.running = true
      this.launch.done = false
      this.launch.result = null
      const start = performance.now()
      // Simulate a lightweight measurement by checking Wails init time
      // The actual launch time is already elapsed when this UI is shown
      // Measure how fast we can get a simple system info call
      try {
        const { GetSystemInfo } = await import('../../wailsjs/go/main/App')
        await GetSystemInfo()
        const elapsed = Math.round(performance.now() - start)
        this.launch.elapsed = elapsed
        this.launch.fast = elapsed < 500
        this.launch.result = { success: true, message: `系统信息调用耗时 ${elapsed}ms` }
      } catch (e) {
        this.launch.result = { success: false, message: '系统信息调用失败: ' + String(e) }
      }
      this.launch.running = false
      this.launch.done = true
    },

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
      try {
        const result = await TestFNQ()
        this.fnq.result = result
      } catch (e) {
        this.fnq.result = { error: String(e) }
      }
      this.fnq.running = false
      this.fnq.done = true
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

    async runAllTests() {
      await this.testLaunchSpeed()
      await this.testBrightness()
      await this.testFNQ()
      await this.testModeSwitch()
    }
  }
}
</script>

<style scoped>
.test-func-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.test-desc {
  margin: 0;
  color: var(--text-muted);
  font-size: 13px;
}

.test-body {
  margin-top: 12px;
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

.test-loading .spinner {
  width: 20px;
  height: 20px;
  border-width: 2px;
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
  border: 1px solid rgba(248, 113, 113, 0.2);
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

/* Spinner small */
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
</style>
