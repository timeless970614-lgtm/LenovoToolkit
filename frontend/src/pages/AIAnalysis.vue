<template>
  <div class="ai-analysis-page">

    <!-- Tab Switcher (always visible) -->
    <div class="ai-tabs">
      <button :class="['ai-tab', { active: activeAIType === 'log' }]" @click="activeAIType = 'log'">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
          <polyline points="14 2 14 8 20 8"/>
        </svg>
        Log Analysis
      </button>
      <button :class="['ai-tab', { active: activeAIType === 'etl' }]" @click="switchToETL">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/>
        </svg>
        ETL Trace
      </button>
      <button :class="['ai-tab', { active: activeAIType === 'salog' }]" @click="activeAIType = 'salog'">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
          <line x1="3" y1="9" x2="21" y2="9"/>
          <line x1="9" y1="21" x2="9" y2="9"/>
        </svg>
        SA Log
      </button>
      <button :class="['ai-tab', { active: activeAIType === 'toolkit' }]" @click="switchToToolkit">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>
        </svg>
        Toolkit
      </button>
      <button :class="['ai-tab', { active: activeAIType === 'msr' }]" @click="switchToMSR">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/><line x1="8" y1="7" x2="16" y2="7"/><line x1="8" y1="11" x2="16" y2="11"/><line x1="8" y1="15" x2="12" y2="15"/>
        </svg>
        MSR
      </button>
      <button :class="['ai-tab', { active: activeAIType === 'wmi' }]" @click="switchToWMI">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <rect x="2" y="3" width="20" height="14" rx="2" ry="2"/>
          <line x1="8" y1="21" x2="16" y2="21"/>
          <line x1="12" y1="17" x2="12" y2="21"/>
        </svg>
        WMI
      </button>
      <button :class="['ai-tab', { active: activeAIType === 'ppm' }]" @click="switchToPPM">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/>
        </svg>
        PPM Driver
      </button>
    </div>

    <!-- ETL Trace Tab -->
    <div v-if="activeAIType === 'etl'" class="etl-section">

      <!-- Elevated Warning -->
      <div v-if="!isElevated" class="elevated-warning">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
        </svg>
        <div>
          <strong>Administrator Required</strong>
          <p>ETL trace capture requires running as Administrator. End this process in Task Manager, then right-click the exe and choose "Run as administrator".</p>
        </div>
      </div>

      <!-- Capture Control Card -->
      <div class="card">
        <div class="card-header">
          <span class="card-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:8px">
              <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
            </svg>
            Trace Capture
          </span>
          <span v-if="captureState.isCapturing" class="capture-badge recording">
            <span class="rec-dot"></span> Recording
          </span>
        </div>

        <!-- Profile Selection -->
        <div class="profile-grid">
          <button
            v-for="p in etlProfiles" :key="p.id"
            :class="['profile-btn', { active: selectedProfile === p.id, disabled: captureState.isCapturing }]"
            @click="selectedProfile = p.id" :title="p.description">
            <span class="profile-name">{{ p.name }}</span>
            <span class="profile-desc">{{ p.description }}</span>
          </button>
        </div>

        <!-- Duration & Status Controls -->
        <div class="capture-controls">
          <div class="control-group">
            <label class="ctrl-label">Duration</label>
            <div class="duration-buttons">
              <button v-for="d in [10, 30, 60, 120]" :key="d"
                :class="['dur-btn', { active: captureDuration === d }]"
                @click="captureDuration = d" :disabled="captureState.isCapturing">
                {{ d < 60 ? d + 's' : (d/60) + 'm' }}
              </button>
            </div>
          </div>
          <div class="control-group">
            <label class="ctrl-label">Status</label>
            <span :class="['status-val', captureState.isCapturing ? 'recording' : 'idle']">
              {{ captureState.isCapturing ? 'Recording: ' + captureState.profile : 'Ready' }}
            </span>
          </div>
          <div class="control-group">
            <label class="ctrl-label">Output</label>
            <span class="status-val mono" v-if="captureState.outputPath">{{ truncatePath(captureState.outputPath) }}</span>
            <span class="status-val muted" v-else>C:\Users\Public\ETL_Traces\</span>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="capture-actions">
          <button class="btn-start" @click="startCapture" :disabled="!isElevated || startingCapture || captureState.isCapturing">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" stroke="none"><circle cx="12" cy="12" r="8"/></svg>
            {{ startingCapture ? 'Starting...' : 'Start Trace' }}
          </button>
          <button class="btn-stop" @click="stopCapture" :disabled="!captureState.isCapturing || stoppingCapture">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" stroke="none"><rect x="6" y="6" width="12" height="12" rx="1"/></svg>
            {{ stoppingCapture ? 'Stopping...' : 'Stop Trace' }}
          </button>
          <button class="btn-outline" @click="refreshTraceList" :disabled="refreshing">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/>
            </svg>
            Refresh List
          </button>
          <button class="btn-outline" @click="openTraceFolder">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
            </svg>
            Open Folder
          </button>
        </div>

        <div v-if="captureError" class="capture-error">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
          </svg>
          {{ captureError }}
        </div>
      </div>

      <!-- Load External ETL -->
      <div class="card load-etl-card">
        <div class="card-header">
          <span class="card-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:8px">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
              <polyline points="17 8 12 3 7 8"/>
              <line x1="12" y1="3" x2="12" y2="15"/>
            </svg>
            Load External ETL
          </span>
        </div>
        <div class="etl-load-body">
          <p class="etl-load-hint">Load an existing .etl trace file from disk and analyze it automatically.</p>
          <div class="etl-load-row">
            <button class="btn-start" @click="loadETLFile" :disabled="loadingETL || analyzingTrace">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
                <polyline points="12 11 12 17"/>
                <polyline points="9 14 12 17 15 14"/>
              </svg>
              {{ loadingETL ? 'Opening...' : 'Choose .etl File & Analyze' }}
            </button>
            <code v-if="loadedETLPath" class="loaded-path">{{ loadedETLPath }}</code>
          </div>
        </div>
      </div>


    <!-- Analysis Results (ETL only) -->
    <div v-if="activeAIType === 'etl'">
      <div v-if="analysisResult" class="analysis-results">

        <div class="card">
          <div class="card-header">
            <span class="card-title">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:8px">
                <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/>
              </svg>
              Analysis Result
            </span>
            <span class="trace-meta mono">{{ truncatePath(analysisResult.traceInfo.path) }}</span>
          </div>
        </div>

        <!-- Step Progress -->
        <div class="card">
          <div class="card-header">
            <span class="card-title">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:8px">
                <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
              </svg>
              Analysis Progress (5 Steps)
            </span>
          </div>
          <div class="step-list">
            <div v-for="s in analysisResult.steps" :key="s.step" class="step-row" :class="'step-' + s.status">
              <span class="step-icon">
                <svg v-if="s.status === 'done'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#4CAF50" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
                <svg v-else-if="s.status === 'error'" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#F44336" stroke-width="3"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
                <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#FF9800" stroke-width="3" class="spinning"><path d="M21 12a9 9 0 1 1-6.22-8.56"/></svg>
              </span>
              <span class="step-name">Step {{ s.step }}: {{ s.name }}</span>
              <span class="step-detail" v-if="s.detail">{{ s.detail }}</span>
            </div>
          </div>
        </div>

        <!-- Event Types -->
        <div v-if="analysisResult.eventTypes?.length" class="card">
          <div class="analysis-card">
            <div class="card-title-sm">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/>
              </svg>
              Event Types ({{ analysisResult.eventTypes.length }})
            </div>
            <div class="event-list">
              <div v-for="ev in analysisResult.eventTypes.slice(0, 10)" :key="ev.eventName" class="event-row">
                <span class="event-cat" :class="'cat-' + ev.category.toLowerCase()">{{ ev.category }}</span>
                <span class="event-name">{{ ev.eventName }}</span>
                <span class="event-count">{{ ev.count }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Issues -->
        <div v-if="analysisResult.issues?.length" class="card">
          <div class="analysis-card card-issues">
            <div class="card-title-sm error-title">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
              Potential Issues ({{ analysisResult.issues.length }})
            </div>
            <div class="issue-list">
              <div v-for="iss in analysisResult.issues" :key="iss.keyword" class="issue-row">
                <span class="issue-kw">{{ iss.keyword }}</span>
                <span class="issue-count">{{ iss.foundCount }} 个问题</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Summary Banner -->
        <div v-if="analysisResult.summary" class="result-summary-banner">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/><path d="M12 8v4M12 16h.01"/>
          </svg>
          {{ analysisResult.summary }}
        </div>

        <!-- WPA Launch -->
        <div v-if="analysisResult.steps?.length >= 5 && analysisResult.steps[4].status === 'done'" class="card wpa-card">
          <div class="card-header">
            <span class="card-title">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:8px">
                <rect x="2" y="3" width="20" height="14" rx="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/>
              </svg>
              Step 5: 打开分析结果 (WPA)
            </span>
            <button class="btn-wpa" @click="openWPA">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" stroke="none"><polygon points="5 3 19 12 5 21 5 3"/></svg>
              Open in WPA
            </button>
          </div>
          <div class="wpa-hint">{{ analysisResult.steps[4].detail }}</div>
        </div>

        <!-- Analysis Grid (original CPU/Disk/Network/Power/DPC) -->
        <div class="analysis-grid">
          <!-- CPU -->
          <div class="analysis-card">
            <div class="card-title-sm">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="4" y="4" width="16" height="16" rx="2"/><rect x="9" y="9" width="6" height="6"/>
                <line x1="9" y1="2" x2="9" y2="4"/><line x1="15" y1="2" x2="15" y2="4"/>
                <line x1="9" y1="20" x2="9" y2="22"/><line x1="15" y1="20" x2="15" y2="22"/>
              </svg>
              CPU
            </div>
            <div class="metric-row" v-if="analysisResult.cpu.cpuUsagePct">
              <span class="metric-label">CPU Usage</span>
              <span class="metric-value highlight">{{ analysisResult.cpu.cpuUsagePct }}</span>
            </div>
            <div class="metric-row" v-if="analysisResult.cpu.contextSwitches">
              <span class="metric-label">Context Switches</span>
              <span class="metric-value">{{ formatNum(analysisResult.cpu.contextSwitches) }}</span>
            </div>
            <div class="metric-row" v-if="analysisResult.cpu.interrupts">
              <span class="metric-label">Interrupts</span>
              <span class="metric-value">{{ formatNum(analysisResult.cpu.interrupts) }}</span>
            </div>
            <div v-if="analysisResult.cpu.busyProcesses?.length" class="proc-list">
              <div class="proc-header">Top Processes</div>
              <div v-for="p in analysisResult.cpu.busyProcesses" :key="p.processName" class="proc-row">
                <span class="proc-name">{{ p.processName }}</span>
                <span class="proc-pct">{{ p.cpuPct.toFixed(1) }}%</span>
              </div>
            </div>
            <div v-else class="no-data-sm">Run with xperf profile for process details</div>
          </div>

          <!-- Disk I/O -->
          <div class="analysis-card">
            <div class="card-title-sm">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="2" y="3" width="20" height="14" rx="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/>
              </svg>
              Disk I/O
            </div>
            <div class="metric-row" v-if="analysisResult.disk.totalReadMB">
              <span class="metric-label">Total Read</span>
              <span class="metric-value">{{ analysisResult.disk.totalReadMB }}</span>
            </div>
            <div class="metric-row" v-if="analysisResult.disk.totalWrittenMB">
              <span class="metric-label">Total Written</span>
              <span class="metric-value">{{ analysisResult.disk.totalWrittenMB }}</span>
            </div>
            <div class="metric-row" v-if="analysisResult.disk.readOpsPerSec">
              <span class="metric-label">Read Ops/s</span>
              <span class="metric-value">{{ analysisResult.disk.readOpsPerSec }}</span>
            </div>
            <div class="metric-row" v-if="analysisResult.disk.avgLatencyMs">
              <span class="metric-label">Avg Latency</span>
              <span class="metric-value">{{ analysisResult.disk.avgLatencyMs }}</span>
            </div>
          </div>

          <!-- Network -->
          <div class="analysis-card">
            <div class="card-title-sm">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M5 12.55a11 11 0 0 1 14.08 0M1.42 9a16 16 0 0 1 21.16 0M8.53 16.11a6 6 0 0 1 6.95 0M12 20h.01"/>
              </svg>
              Network
            </div>
            <div class="metric-row" v-if="analysisResult.network.totalSentMB">
              <span class="metric-label">Sent</span>
              <span class="metric-value">{{ analysisResult.network.totalSentMB }}</span>
            </div>
            <div class="metric-row" v-if="analysisResult.network.totalRecvMB">
              <span class="metric-label">Received</span>
              <span class="metric-value">{{ analysisResult.network.totalRecvMB }}</span>
            </div>
            <div class="metric-row" v-if="analysisResult.network.tcpConnections">
              <span class="metric-label">TCP Connections</span>
              <span class="metric-value">{{ analysisResult.network.tcpConnections }}</span>
            </div>
            <div class="metric-row" v-if="analysisResult.gpu.gpuEngineUtilPct">
              <span class="metric-label">GPU Util</span>
              <span class="metric-value">{{ analysisResult.gpu.gpuEngineUtilPct }}</span>
            </div>
            <div class="metric-row" v-if="analysisResult.gpu.gpuMemoryUsedMB">
              <span class="metric-label">GPU Memory</span>
              <span class="metric-value">{{ analysisResult.gpu.gpuMemoryUsedMB }}</span>
            </div>
          </div>

          <!-- Power -->
          <div class="analysis-card">
            <div class="card-title-sm">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/>
              </svg>
              Power
            </div>
            <div class="metric-row" v-if="analysisResult.power.cpuPower">
              <span class="metric-label">CPU Power</span>
              <span class="metric-value">{{ analysisResult.power.cpuPower }}</span>
            </div>
            <div class="metric-row" v-if="analysisResult.power.gpuPower">
              <span class="metric-label">GPU Power</span>
              <span class="metric-value">{{ analysisResult.power.gpuPower }}</span>
            </div>
            <div class="metric-row" v-if="analysisResult.power.packagePower">
              <span class="metric-label">Package Power</span>
              <span class="metric-value">{{ analysisResult.power.packagePower }}</span>
            </div>
            <div class="metric-row" v-if="analysisResult.power.s0ixDuration">
              <span class="metric-label">S0ix Duration</span>
              <span class="metric-value">{{ analysisResult.power.s0ixDuration }}</span>
            </div>
            <div class="metric-row" v-if="analysisResult.power.s0ixTransitions">
              <span class="metric-label">S0ix Transitions</span>
              <span class="metric-value">{{ analysisResult.power.s0ixTransitions }}</span>
            </div>
            <div class="metric-row" v-if="analysisResult.power.processorFreqMHz">
              <span class="metric-label">CPU Freq</span>
              <span class="metric-value">{{ analysisResult.power.processorFreqMHz }}</span>
            </div>
          </div>

          <!-- DPC/ISR -->
          <div class="analysis-card">
            <div class="card-title-sm">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
              </svg>
              DPC / ISR Latency
            </div>
            <div class="metric-row" v-if="analysisResult.dpcrisr.avgDpcMs">
              <span class="metric-label">Avg DPC</span>
              <span class="metric-value">{{ analysisResult.dpcrisr.avgDpcMs }}</span>
            </div>
            <div class="metric-row" v-if="analysisResult.dpcrisr.maxDpcMs">
              <span class="metric-label">Max DPC</span>
              <span class="metric-value warning">{{ analysisResult.dpcrisr.maxDpcMs }}</span>
            </div>
            <div class="metric-row" v-if="analysisResult.dpcrisr.avgIsrMs">
              <span class="metric-label">Avg ISR</span>
              <span class="metric-value">{{ analysisResult.dpcrisr.avgIsrMs }}</span>
            </div>
            <div class="metric-row" v-if="analysisResult.dpcrisr.maxIsrMs">
              <span class="metric-label">Max ISR</span>
              <span class="metric-value warning">{{ analysisResult.dpcrisr.maxIsrMs }}</span>
            </div>
            <div v-if="analysisResult.dpcrisr.highDpcLatencyProcs?.length" class="proc-list">
              <div class="proc-header">High DPC Latency</div>
              <div v-for="p in analysisResult.dpcrisr.highDpcLatencyProcs" :key="p.processName" class="proc-row">
                <span class="proc-name">{{ p.processName || p.module || 'Unknown' }}</span>
                <span class="proc-pct warning">{{ p.maxLatencyMs }}</span>
              </div>
            </div>
          </div>
        </div>

      </div><!-- /analysis-results -->

      <!-- Analyzing Overlay -->
      <div v-if="analyzingTrace" class="analyzing-overlay">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="spinning">
          <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
        </svg>
        <span>Parsing ETL trace with tracerpt...</span>
      </div>

    </div><!-- /ETL section -->
    </div><!-- /Analysis Results wrapper -->
    <!-- SA Log (System Admin Logs) Tab -->
    <div v-if="activeAIType === 'salog'" class="salog-section">

      <!-- Capture System Event Log -->
      <div class="card event-log-card">
        <div class="card-header">
          <span class="card-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:8px">
              <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
              <line x1="3" y1="9" x2="21" y2="9"/>
              <line x1="9" y1="21" x2="9" y2="9"/>
            </svg>
            Capture System Event Log
          </span>
          <span v-if="eventLogResult" class="capture-badge" style="background: rgba(76,175,80,0.1); color:#4CAF50; border-color:rgba(76,175,80,0.2);">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
            Captured
          </span>
        </div>
        <div class="event-log-body">
          <p class="event-log-hint">Capture Windows System event log (errors, warnings, crashes) for quick diagnostics — no ETL trace needed.</p>
          <div class="event-log-controls">
            <div class="event-log-presets">
              <button v-for="p in eventLogPresets" :key="p.label"
                :class="['preset-btn', { active: selectedEventLogPreset === p.label }]"
                @click="selectedEventLogPreset = p.label">
                {{ p.label }}
              </button>
            </div>
            <button class="btn-capture-eventlog" @click="captureEventLog" :disabled="capturingEventLog">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ spinning: capturingEventLog }">
                <circle v-if="capturingEventLog" cx="12" cy="12" r="10"/>
                <path v-if="capturingEventLog" d="M12 6v6l4 2"/>
                <polyline v-else points="22 12 18 12 15 21 9 3 6 12 2 12"/>
              </svg>
              {{ capturingEventLog ? 'Capturing...' : 'Capture Event Log' }}
            </button>
            <button class="btn-open-eventviewer" @click="openEventViewer">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                <line x1="3" y1="9" x2="21" y2="9"/>
                <line x1="9" y1="21" x2="9" y2="9"/>
              </svg>
              Open Event Log
            </button>
            <button v-if="eventLogResult" class="btn-export-eventlog" @click="exportEventLog">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="7 10 12 15 17 10"/>
                <line x1="12" y1="15" x2="12" y2="3"/>
              </svg>
              Export CSV
            </button>
            <button class="btn-load-evtx" @click="loadEVTX" :disabled="loadingEVTX">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14 2 14 8 20 8"/>
                <line x1="12" y1="18" x2="12" y2="12"/>
                <line x1="9" y1="15" x2="15" y2="15"/>
              </svg>
              {{ loadingEVTX ? 'Loading...' : 'Load EVTX' }}
            </button>
            <button v-if="eventLogResult" class="btn-evtx-to-csv" @click="evtxToCSV" :disabled="!evtxFilePath">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="7 10 12 15 17 10"/>
                <line x1="12" y1="15" x2="12" y2="3"/>
              </svg>
              EVTX → CSV
            </button>
          </div>
          <!-- Level Filter -->
          <div v-if="eventLogResult" class="event-level-filter">
            <span class="filter-label">Filter by Level:</span>
            <button :class="['level-btn', 'level-critical', { active: evtxLevelFilter === 'Critical' }]" @click="evtxLevelFilter = evtxLevelFilter === 'Critical' ? '' : 'Critical'">Critical</button>
            <button :class="['level-btn', 'level-error', { active: evtxLevelFilter === 'Error' }]" @click="evtxLevelFilter = evtxLevelFilter === 'Error' ? '' : 'Error'">Error</button>
            <button :class="['level-btn', 'level-warning', { active: evtxLevelFilter === 'Warning' }]" @click="evtxLevelFilter = evtxLevelFilter === 'Warning' ? '' : 'Warning'">Warning</button>
            <button :class="['level-btn', 'level-info', { active: evtxLevelFilter === 'Information' }]" @click="evtxLevelFilter = evtxLevelFilter === 'Information' ? '' : 'Information'">Info</button>
            <button v-if="evtxLevelFilter" class="level-btn level-clear" @click="evtxLevelFilter = ''">✕ Clear</button>
          </div>
        </div>

        <!-- Event Log Results -->
        <div v-if="eventLogResult" class="event-log-results">
          <!-- Summary Bar -->
          <div class="event-log-summary">
            <div class="event-stat critical">
              <span class="stat-count">{{ eventLogResult.criticalCount }}</span>
              <span class="stat-label">Critical</span>
            </div>
            <div class="event-stat error">
              <span class="stat-count">{{ eventLogResult.errorCount }}</span>
              <span class="stat-label">Errors</span>
            </div>
            <div class="event-stat warning">
              <span class="stat-count">{{ eventLogResult.warningCount }}</span>
              <span class="stat-label">Warnings</span>
            </div>
            <div class="event-stat info">
              <span class="stat-count">{{ eventLogResult.infoCount }}</span>
              <span class="stat-label">Info</span>
            </div>
            <div class="event-stat total">
              <span class="stat-count">{{ eventLogResult.totalEvents }}</span>
              <span class="stat-label">{{ eventLogResult.timeRange }}</span>
            </div>
          </div>

          <!-- Events Table (Level Filtered) -->
          <div v-if="filteredEvents.length" class="event-error-table">
            <div class="table-title">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#F44336" stroke-width="2">
                <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
              {{ evtxLevelFilter ? evtxLevelFilter + ' Events' : 'Recent Critical / Errors' }}
              <span class="table-count">{{ filteredEvents.length }}</span>
            </div>
            <div class="table-scroll">
              <table class="event-table">
                <thead>
                  <tr>
                    <th>Time</th>
                    <th>Level</th>
                    <th>Source</th>
                    <th>ID</th>
                    <th>Message</th>
                    <th>EventData</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(ev, i) in filteredEvents" :key="i" :class="{ 'row-critical': ev.level === 'Critical' }">
                      <td class="cell-time">{{ formatEventTime(ev.time) }}</td>
                      <td><span :class="'level-badge level-' + ev.level.toLowerCase()">{{ ev.level }}</span></td>
                      <td class="cell-source">{{ ev.providerName }}</td>
                      <td class="cell-id">{{ ev.eventId }}</td>
                      <td class="cell-msg">{{ ev.message }}</td>
                      <td class="cell-eventdata"><span v-if="ev.eventData" style="white-space:pre-line">{{ ev.eventData }}</span><span v-else>-</span></td>
                    </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- Top Providers -->
          <div v-if="eventLogResult.topProviders && eventLogResult.topProviders.length" class="event-providers">
            <div class="table-title">Top Event Sources</div>
            <div class="provider-list">
              <div v-for="(p, i) in eventLogResult.topProviders" :key="i" class="provider-row">
                <span class="provider-name">{{ p.name }}</span>
                <span class="provider-count">{{ p.count }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>


      <!-- Event Log Analysis Card -->
      <div v-if="eventLogResult" class="card eventlog-analysis-card">
        <div class="card-header">
          <span class="card-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:8px">
              <path d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/>
            </svg>
            Event Log Analysis
          </span>
          <button class="btn-analyze-eventlog" @click="analyzeEventLog" :disabled="analyzingEventLog">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ spinning: analyzingEventLog }">
              <circle v-if="analyzingEventLog" cx="12" cy="12" r="10"/>
              <path v-if="analyzingEventLog" d="M12 6v6l4 2"/>
              <path v-else d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
              <polyline v-if="!analyzingEventLog" points="22 4 12 14.01 9 11.01"/>
            </svg>
            {{ analyzingEventLog ? 'Analyzing...' : 'Analyze' }}
          </button>
        </div>
        <div v-if="eventLogAnalysis" class="eventlog-analysis-body">

          <!-- Health Badge -->
          <div :class="['health-badge', eventLogAnalysis.overallHealth.toLowerCase()]">
            <svg v-if="eventLogAnalysis.overallHealth === 'Good'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#4CAF50" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
            <svg v-else-if="eventLogAnalysis.overallHealth === 'Warning'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#FF9800" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
            <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#F44336" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
            <span class="health-label">{{ eventLogAnalysis.overallHealth }}</span>
            <span class="health-summary">{{ eventLogAnalysis.summary }}</span>
          </div>

          <!-- Root Causes -->
          <div v-if="eventLogAnalysis.rootCauses && eventLogAnalysis.rootCauses.length" class="analysis-section">
            <div class="sub-title">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
              Root Causes
            </div>
            <div class="root-cause-list">
              <div v-for="(rc, i) in eventLogAnalysis.rootCauses" :key="i" :class="['root-cause-item', rc.severity.toLowerCase()]">
                <span class="rc-category">{{ rc.category }}</span>
                <span class="rc-detail">{{ rc.detail }}</span>
                <span :class="['rc-severity', rc.severity.toLowerCase()]">{{ rc.severity }}</span>
              </div>
            </div>
          </div>

          <!-- Repeat Errors -->
          <div v-if="eventLogAnalysis.repeatErrors && eventLogAnalysis.repeatErrors.length" class="analysis-section">
            <div class="sub-title">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="23 6 13.5 15.5 8.5 10.5 1 18"/><polyline points="17 6 23 6 23 12"/></svg>
              Repeat Error Patterns
            </div>
            <div class="repeat-list">
              <div v-for="(rp, i) in eventLogAnalysis.repeatErrors" :key="i" class="repeat-item">
                <div class="repeat-header">
                  <code class="repeat-source">{{ rp.providerName }} #{{ rp.eventId }}</code>
                  <span class="repeat-count">&times;{{ rp.count }}</span>
                  <span class="repeat-time">Last: {{ formatEventTime(rp.lastSeen) }}</span>
                </div>
                <div class="repeat-msg">{{ rp.sampleMsg }}</div>
              </div>
            </div>
          </div>

          <!-- Known Issues -->
          <div v-if="eventLogAnalysis.knownIssues && eventLogAnalysis.knownIssues.length" class="analysis-section">
            <div class="sub-title">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
              Known Issues Detected
            </div>
            <div class="known-issue-list">
              <div v-for="(ki, i) in eventLogAnalysis.knownIssues" :key="i" class="known-issue-item">
                <div class="ki-header">
                  <span class="ki-name">{{ ki.name }}</span>
                  <span :class="['ki-confidence', ki.confidence >= 80 ? 'high' : ki.confidence >= 60 ? 'mid' : 'low']">{{ ki.confidence }}%</span>
                </div>
                <div class="ki-detail">{{ ki.detail }}</div>
              </div>
            </div>
          </div>

          <!-- Correlations -->
          <div v-if="eventLogAnalysis.correlations && eventLogAnalysis.correlations.length" class="analysis-section">
            <div class="sub-title">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/><path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/></svg>
              Correlations
            </div>
            <div class="correlation-list">
              <div v-for="(c, i) in eventLogAnalysis.correlations" :key="i" class="correlation-item">{{ c.description }}</div>
            </div>
          </div>

          <!-- Timeline -->
          <div v-if="eventLogAnalysis.timeline && eventLogAnalysis.timeline.length" class="analysis-section">
            <div class="sub-title">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
              Error Timeline
            </div>
            <div class="timeline-list">
              <div v-for="(te, i) in eventLogAnalysis.timeline" :key="i" :class="['timeline-item', te.level.toLowerCase()]">
                <span class="tl-time">{{ formatEventTime(te.time) }}</span>
                <span :class="['tl-level', te.level.toLowerCase()]">{{ te.level }}</span>
                <code class="tl-source">{{ te.source }} #{{ te.id }}</code>
                <span class="tl-msg">{{ te.message }}</span>
              </div>
            </div>
          </div>

          <!-- Recommendations -->
          <div v-if="eventLogAnalysis.recommendations && eventLogAnalysis.recommendations.length" class="analysis-section">
            <div class="sub-title">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>
              Recommendations
            </div>
            <ul class="recommendation-list">
              <li v-for="(r, i) in eventLogAnalysis.recommendations" :key="i">{{ r }}</li>
            </ul>
          </div>

        </div>
        <div v-else class="analysis-hint">
          <p>Click <strong>Analyze</strong> to identify error patterns, root causes, and get recommendations.</p>
        </div>
      </div>

      <!-- Capture Dispdiag Log -->
      <div class="card dispdiag-card">
        <div class="card-header">
          <span class="card-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:8px">
              <rect x="2" y="3" width="20" height="14" rx="2" ry="2"/>
              <line x1="8" y1="21" x2="16" y2="21"/>
              <line x1="12" y1="17" x2="12" y2="21"/>
            </svg>
            Capture Dispdiag Log
          </span>
          <span v-if="dispdiagResult" class="capture-badge" style="background: rgba(76,175,80,0.1); color:#4CAF50; border-color:rgba(76,175,80,0.2);">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
            Captured
          </span>
        </div>
        <div class="dispdiag-body">
          <p class="dispdiag-hint">Capture Windows display diagnostic data — EDID, link training status, brightness info, and driver details via <code>dispdiag.exe</code>.</p>
          <div class="dispdiag-controls">
            <div class="dispdiag-options">
              <label class="dispdiag-opt">
                <input type="checkbox" v-model="dispdiagDump" /> Dump mode (-d)
              </label>
              <label class="dispdiag-opt">
                Delay: <input type="number" v-model.number="dispdiagDelay" min="0" max="30" class="delay-input" />s
              </label>
            </div>
            <button class="btn-capture-dispdiag" @click="captureDispdiag" :disabled="capturingDispdiag">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :class="{ spinning: capturingDispdiag }">
                <circle v-if="capturingDispdiag" cx="12" cy="12" r="10"/>
                <path v-if="capturingDispdiag" d="M12 6v6l4 2"/>
                <polyline v-else points="22 12 18 12 15 21 9 3 6 12 2 12"/>
              </svg>
              {{ capturingDispdiag ? 'Running dispdiag...' : 'Run Dispdiag' }}
            </button>
            <button class="btn-open-dispdiag" @click="openDispdiagLog">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14 2 14 8 20 8"/>
                <line x1="16" y1="13" x2="8" y2="13"/>
                <line x1="16" y1="17" x2="8" y2="17"/>
                <polyline points="10 9 9 9 8 9"/>
              </svg>
              Open Dispdiag Log
            </button>
          </div>
        </div>

        <!-- Dispdiag Results -->
        <div v-if="dispdiagResult" class="dispdiag-results">
          <!-- Info Bar -->
          <div class="dispdiag-info-bar">
            <div class="dispdiag-info-item">
              <span class="info-label">File:</span>
              <code class="info-value">{{ dispdiagResult.fileName }}</code>
            </div>
            <div class="dispdiag-info-item">
              <span class="info-label">Size:</span>
              <span class="info-value">{{ dispdiagResult.outputSize }}</span>
            </div>
            <div class="dispdiag-info-item">
              <span class="info-label">Duration:</span>
              <span class="info-value">{{ dispdiagResult.durationSecs }}s</span>
            </div>
            <div class="dispdiag-info-item">
              <span class="info-label">Pass/Fail:</span>
              <span v-if="!dispdiagResult.errors || !dispdiagResult.errors.length" class="pass-badge">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#4CAF50" stroke-width="3"><polyline points="20 6 9 17 4 12"/></svg>
                PASS
              </span>
              <span v-else class="fail-badge">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#F44336" stroke-width="3"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
                FAIL
              </span>
            </div>
            <div class="dispdiag-info-item">
              <span class="info-label">EDID Blocks:</span>
              <span class="info-value" style="color:var(--lenovo-red);font-weight:700">{{ dispdiagResult.edidBlocks }}</span>
            </div>
          </div>

          <!-- Errors -->
          <div v-if="dispdiagResult.errors && dispdiagResult.errors.length" class="dispdiag-errors">
            <div class="sub-title">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#F44336" stroke-width="2">
                <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
              Key Failures / Errors
            </div>
            <ul class="error-list">
              <li v-for="(err, i) in dispdiagResult.errors.slice(0, 15)" :key="i">{{ err }}</li>
            </ul>
          </div>

          <!-- Warnings -->
          <div v-if="dispdiagResult.warnings && dispdiagResult.warnings.length" class="dispdiag-warnings">
            <div class="sub-title">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#FF9800" stroke-width="2">
                <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/>
                <line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/>
              </svg>
              Warnings
            </div>
            <ul class="warn-list">
              <li v-for="(w, i) in dispdiagResult.warnings.slice(0, 10)" :key="i">{{ w }}</li>
            </ul>
          </div>

          <!-- Summary -->
          <div v-if="dispdiagResult.summary" class="dispdiag-summary">
            <div class="sub-title">Summary</div>
            <div class="summary-grid">
              <div class="summary-row" v-if="dispdiagResult.summary.driverVersion">
                <span class="grid-key">Driver</span>
                <code class="grid-val">{{ dispdiagResult.summary.driverVersion }}</code>
              </div>
              <div class="summary-row" v-if="dispdiagResult.summary.buildVersion">
                <span class="grid-key">Build</span>
                <code class="grid-val">{{ dispdiagResult.summary.buildVersion }}</code>
              </div>
              <div class="summary-row" v-if="dispdiagResult.summary.datVersion">
                <span class="grid-key">Dat Version</span>
                <code class="grid-val">{{ dispdiagResult.summary.datVersion }}</code>
              </div>
            </div>
          </div>

          <!-- Brightness Info -->
          <div v-if="dispdiagResult.brightnessInfo && dispdiagResult.brightnessInfo.length" class="dispdiag-brightness">
            <div class="sub-title">Brightness Info</div>
            <div class="brightness-list">
              <code v-for="(b, i) in dispdiagResult.brightnessInfo.slice(0, 10)" :key="i" class="brightness-line">{{ b }}</code>
            </div>
          </div>

          <!-- Raw Content Preview -->
          <div v-if="dispdiagResult.fileContent" class="dispdiag-content">
            <div class="sub-title">Raw Output Preview</div>
            <pre class="dispdiag-preview">{{ dispdiagResult.fileContent.substring(0, 2000) }}
{{ dispdiagResult.fileContent.length > 2000 ? '...(truncated)' : '' }}</pre>
          </div>

          <!-- Actions -->
          <div class="dispdiag-actions">
            <button class="btn-dispdiag-action" @click="openDispdiagFolder">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
              </svg>
              Open Output Folder
            </button>
            <button v-if="dispdiagResult.outputPath" class="btn-dispdiag-action" @click="exportDispdiagJSON">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="7 10 12 15 17 10"/>
                <line x1="12" y1="15" x2="12" y2="3"/>
              </svg>
              Export JSON
            </button>
          </div>
        </div>
      </div>

    </div><!-- /SA Log section -->

    <!-- Log Analysis Tab -->
    <div v-if="activeAIType === 'log'" class="log-section">
      <!-- Log Files Card - Password Protected -->
      <div class="card advanced-section">
        <div class="card-header advanced-toggle" @click="toggleLogSection">
          <span class="card-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:8px">
              <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
              <polyline points="14 2 14 8 20 8"/>
              <line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/>
            </svg>
            Log Files
          </span>
          <svg :class="['chevron', { open: logSectionExpanded }]" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="6 9 12 15 18 9"/></svg>
        </div>

        <!-- Password Prompt -->
        <div v-if="logSectionExpanded && !logUnlocked" class="password-prompt">
          <div class="password-row">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 8px; flex-shrink:0;">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
              <path d="M7 11V7a5 5 0 0110 0v4"/>
            </svg>
            <input
              type="password"
              class="password-input"
              v-model="logPassword"
              placeholder="Enter password to unlock"
              @keydown.enter="unlockLogs"
              ref="logPasswordInput"
            />
            <button class="btn-unlock" @click="unlockLogs">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0110 0v4"/>
              </svg>
              Unlock
            </button>
          </div>
          <div v-if="logPasswordError" class="password-error">{{ logPasswordError }}</div>
        </div>

        <!-- Unlocked Content -->
        <div v-if="logSectionExpanded && logUnlocked">
          <div style="display:flex;justify-content:flex-end;padding:0 0 8px 0;">
            <button class="btn-sm" @click="loadLogs">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="23 4 23 10 17 10"/><path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"/>
              </svg>
              Refresh
            </button>
          </div>
          <div v-if="logFiles.length" class="log-file-list">
            <div v-for="f in logFiles" :key="f.name" class="log-file-item"
              :class="{ active: selectedLog === f.name }" @click="selectLog(f.name)">
              <span class="log-name">{{ f.name }}</span>
              <span class="log-meta">{{ formatSize(f.size) }}  {{ f.modTime }}</span>
            </div>
          </div>
          <div v-else class="empty-hint">No log files found in {{ logDir }}</div>
        </div>
      </div>

      <!-- AI Analysis Card -->
      <div class="card">
        <div class="card-header">
          <span class="card-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:8px">
              <circle cx="12" cy="12" r="10"/><path d="M12 8v4l3 3"/>
            </svg>
            AI Analysis
          </span>
          <button class="btn-analyze" @click="runAnalysis" :disabled="analyzing">
            <svg v-if="!analyzing" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="5 3 19 12 5 21 5 3"/>
            </svg>
            <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="spinning">
              <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
            </svg>
            {{ analyzing ? 'Analyzing...' : 'Analyze Latest Log' }}
          </button>
        </div>
        <div v-if="summary" class="summary-grid">
          <div class="summary-panel">
            <div class="panel-title">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/>
              </svg>
              Workload Levels
            </div>
            <div class="log-lines">
              <div v-for="(l,i) in summary.workloadLevels" :key="i" class="log-line">{{ l }}</div>
              <div v-if="!summary.workloadLevels?.length" class="no-data">No data</div>
            </div>
          </div>
          <div class="summary-panel">
            <div class="panel-title">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/>
              </svg>
              Turbo Events
            </div>
            <div class="log-lines">
              <div v-for="(l,i) in summary.turboEvents" :key="i" class="log-line">{{ l }}</div>
              <div v-if="!summary.turboEvents?.length" class="no-data">No data</div>
            </div>
          </div>
          <div class="summary-panel">
            <div class="panel-title">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="2" y="3" width="20" height="14" rx="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/>
              </svg>
              Foreground Apps
            </div>
            <div class="log-lines">
              <div v-for="(l,i) in summary.appEvents" :key="i" class="log-line">{{ l }}</div>
              <div v-if="!summary.appEvents?.length" class="no-data">No data</div>
            </div>
          </div>
          <div class="summary-panel panel-errors">
            <div class="panel-title error-title">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
              Errors / Warnings
            </div>
            <div class="log-lines">
              <div v-for="(l,i) in summary.errors" :key="i" class="log-line error-line">{{ l }}</div>
              <div v-if="!summary.errors?.length" class="no-data ok-hint">OK - No errors found</div>
            </div>
          </div>
        </div>
        <div v-else-if="!analyzing" class="empty-hint">Click "Analyze Latest Log" to start</div>
        <div v-else class="analyzing-hint">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="spinning">
            <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
          </svg>
          Reading and parsing log...
        </div>
      </div>

    </div><!-- /log-section -->

    <!-- Toolkit Tab -->
    <div v-if="activeAIType === 'toolkit'" class="toolkit-section">
      
      <!-- Header Card -->
      <div class="card">
        <div class="card-header">
          <span class="card-title">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right:8px">
              <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>
            </svg>
            System Tools - One-Click Install
          </span>
          <button class="btn-sm" @click="openToolkitFolder">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
            </svg>
            Open Folder
          </button>
        </div>
        <div class="toolkit-path">Install Directory: {{ toolkitInstallDir }}</div>
      </div>

      <!-- Category Filter -->
      <div class="category-filter">
        <button :class="['cat-btn', { active: toolkitCategory === 'all' }]" @click="toolkitCategory = 'all'">All ({{ toolkitTools.length }})</button>
        <button :class="['cat-btn', { active: toolkitCategory === 'monitor' }]" @click="toolkitCategory = 'monitor'">Monitor</button>
        <button :class="['cat-btn', { active: toolkitCategory === 'system' }]" @click="toolkitCategory = 'system'">System Info</button>
        <button :class="['cat-btn', { active: toolkitCategory === 'benchmark' }]" @click="toolkitCategory = 'benchmark'">Benchmark</button>
        <button :class="['cat-btn', { active: toolkitCategory === 'diagnostic' }]" @click="toolkitCategory = 'diagnostic'">Diagnostic</button>
      </div>

      <!-- Tools Grid -->
      <div class="tools-grid">
        <div v-for="tool in filteredTools" :key="tool.id" class="tool-card" :class="{ installed: toolInstallStatus[tool.id]?.installed }">
          <div class="tool-header">
            <div class="tool-icon" :class="tool.category">
              <svg v-if="tool.category === 'monitor'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
              </svg>
              <svg v-else-if="tool.category === 'system'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="4" y="4" width="16" height="16" rx="2"/><rect x="9" y="9" width="6" height="6"/>
              </svg>
              <svg v-else-if="tool.category === 'benchmark'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
              </svg>
              <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
            </div>
            <div class="tool-info">
              <span class="tool-name">{{ tool.name }}</span>
              <span class="tool-vendor">{{ tool.vendor }}</span>
            </div>
            <span v-if="toolInstallStatus[tool.id]?.installed" class="installed-badge">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                <polyline points="20 6 9 17 4 12"/>
              </svg>
              Installed
            </span>
          </div>
          
          <div class="tool-desc">{{ tool.description }}</div>
          
          <div class="tool-meta">
            <span class="tool-version">v{{ tool.version }}</span>
            <span class="tool-size">{{ tool.sizeMb }} MB</span>
            <span class="tool-winget" v-if="tool.wingetId">已安装 winget</span>
          </div>

          <!-- Installing Status -->
          <div v-if="isToolBusy(tool.id)" class="tool-installing">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="spinning">
              <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
            </svg>
            <span>{{ toolProgress[tool.id]?.message || 'Installing...' }}</span>
          </div>

          <!-- Error Message -->
          <div v-if="toolProgress[tool.id]?.status === 'error'" class="tool-error">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
            {{ toolProgress[tool.id].message }}
          </div>

          <!-- Success Message -->
          <div v-if="toolProgress[tool.id]?.status === 'completed'" class="tool-success">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="20 6 9 17 4 12"/>
            </svg>
            {{ toolProgress[tool.id].message }}
          </div>

          <div class="tool-actions">
            <button v-if="!toolInstallStatus[tool.id]?.installed && !isToolBusy(tool.id)" class="btn-install" @click="installTool(tool.id)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                <polyline points="7 10 12 15 17 10"/>
                <line x1="12" y1="15" x2="12" y2="3"/>
              </svg>
              Install & Run
            </button>
            <button v-if="isToolBusy(tool.id)" class="btn-installing" disabled>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="spinning">
                <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
              </svg>
              Installing...
            </button>
            <button v-if="toolInstallStatus[tool.id]?.installed" class="btn-run" @click="runTool(tool.id)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor" stroke="none">
                <polygon points="5 3 19 12 5 21 5 3"/>
              </svg>
              Run
            </button>
            <a :href="tool.website" target="_blank" class="btn-link">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
                <polyline points="15 3 21 3 21 9"/>
                <line x1="10" y1="14" x2="21" y2="3"/>
              </svg>
              Website
            </a>
          </div>
        </div>
      </div>

      <!-- Quick Install All -->
      <div class="card" v-if="!allToolsInstalled">
        <div class="card-header">
          <span class="card-title">Quick Actions</span>
        </div>
        <div class="quick-actions">
          <button class="btn-batch" @click="installEssentials" :disabled="batchInstalling">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
              <polyline points="7 10 12 15 17 10"/>
              <line x1="12" y1="15" x2="12" y2="3"/>
            </svg>
            {{ batchInstalling ? 'Installing...' : 'Install Essential Tools (HWiNFO, CPU-Z, GPU-Z)' }}
          </button>
        </div>
      </div>

    </div><!-- /toolkit-section -->

    <!-- MSR Tab -->
    <div v-if="activeAIType === 'msr'" class="msr-section">
      <MSR :theme="theme" />
    </div>

    <!-- WMI Tab -->
    <div v-if="activeAIType === 'wmi'" class="wmi-section">
      <WMI :theme="theme" />
    </div>

    <!-- PPM Driver Tab -->
    <div v-if="activeAIType === 'ppm'" class="ppm-section">
      <PPMDriver :theme="theme" />
    </div>

  </div>
</template>

<script>
import { OpenFolder, OpenETLFileDialog, OpenETLInWPA } from '../../wailsjs/go/main/App'
import MSR from './MSR.vue'
import WMI from './WMI.vue'
import PPMDriver from './PPMDriver.vue'

export default {
  name: 'AIAnalysis',
  props: {
    theme: { type: String, default: 'dark' },
  },
  components: {
    MSR,
    WMI,
    PPMDriver,
  },
  beforeDestroy() {
    this._capturePollTimer && clearInterval(this._capturePollTimer)
  },
  data() {
    return {
      activeAIType: 'log',
      // ETL
      etlProfiles: [],
      selectedProfile: 'GeneralProfile',
      captureDuration: 30,
      captureState: { isCapturing: false, profile: '', outputPath: '', status: '', error: '' },
      captureError: '',
      startingCapture: false,
      stoppingCapture: false,
      traceList: [],
      selectedTrace: '',
      analysisResult: null,
      analyzingTrace: false,
      refreshing: false,
      isElevated: false,
      loadedETLPath: '',
      loadingETL: false,
      // Log
      logDir: 'C:\\ProgramData\\Lenovo\\LenovoDispatcher\\Logs',
      logFiles: [],
      logSectionExpanded: true,
      logUnlocked: false,
      logPassword: '',
      logPasswordError: '',
      selectedLog: null,
      summary: null,
      rawLog: '',
      tailLines: 200,
      analyzing: false,
      // Event Log
      eventLogPresets: [
        { label: 'Last 1h', hours: 1, maxEvents: 200 },
        { label: 'Last 500', hours: 0, maxEvents: 500 },
        { label: 'Last 1000', hours: 0, maxEvents: 1000 },
      ],
      selectedEventLogPreset: 'Last 1000',
      capturingEventLog: false,
      eventLogResult: null,
      analyzingEventLog: false,
      eventLogAnalysis: null,
      // EVTX file
      loadingEVTX: false,
      evtxFilePath: '',
      evtxLevelFilter: '',
      // Dispdiag
      dispdiagDump: false,
      dispdiagDelay: 0,
      capturingDispdiag: false,
      dispdiagResult: null,
      // Toolkit
      toolkitTools: [],
      toolkitInstallDir: 'C:\\LenovoDispatcherToolkit\\Tools',
      toolkitCategory: 'all',
      toolInstallStatus: {},
      toolProgress: {},
      batchInstalling: false,
    }
  },
  async mounted() {
    await this.loadLogs()
  },
  methods: {
    async switchToETL() {
      this.activeAIType = 'etl'
      if (!this.etlProfiles.length) await this.loadETLData()
      // Check elevation status from backend
      try {
        if (window.go?.main?.App) {
          this.isElevated = await window.go.main.App.IsElevated()
        }
      } catch (e) { console.error(e) }
    },
    async loadETLData() {
      try {
        if (window.go?.main?.App) {
          const [profiles, status, traces] = await Promise.all([
            window.go.main.App.GetETLProfiles(),
            window.go.main.App.GetETLCaptureStatus(),
            window.go.main.App.GetETLTraceList(),
          ])
          this.etlProfiles = profiles || []
          this.captureState = status || { isCapturing: false }
          this.traceList = traces || []
        }
      } catch (e) { console.error(e) }
    },
    async startCapture() {
      this.startingCapture = true
      this.captureError = ''
      this._capturePollTimer && clearInterval(this._capturePollTimer)
      try {
        if (window.go?.main?.App) {
          const result = await window.go.main.App.StartETLCapture(this.selectedProfile, this.captureDuration)
          this.captureState = result
          if (result.error) this.captureError = result.error
          else if (!result.isCapturing) this.captureError = result.error || 'Failed to start capture'
          else {
            // Poll backend every 2s to detect when goroutine auto-stops the trace
            this._capturePollTimer = setInterval(async () => {
              try {
                const status = await window.go.main.App.GetETLCaptureStatus()
                if (status && !status.isCapturing) {
                  clearInterval(this._capturePollTimer)
                  this._capturePollTimer = null
                  this.captureState = status
                }
              } catch (e) { /* ignore polling errors */ }
            }, 2000)
          }
        }
      } catch (e) { this.captureError = e.message || String(e) }
      finally { this.startingCapture = false }
    },
    async stopCapture() {
      this.stoppingCapture = true
      this._capturePollTimer && clearInterval(this._capturePollTimer)
      try {
        if (window.go?.main?.App) {
          const traceInfo = await window.go.main.App.StopETLCapture()
          this.captureState = { isCapturing: false }
          if (traceInfo.path) {
            await this.refreshTraceList()
            await this.analyzeTrace(traceInfo.path)
          }
        }
      } catch (e) { this.captureError = e.message || String(e) }
      finally { this.stoppingCapture = false }
    },
    async refreshTraceList() {
      this.refreshing = true
      try {
        if (window.go?.main?.App) this.traceList = await window.go.main.App.GetETLTraceList() || []
      } catch (e) { console.error(e) }
      finally { this.refreshing = false }
    },
    async openTraceFolder() {
      try { await OpenFolder('C:\\Users\\Public\\ETL_Traces') } catch (e) { console.error(e) }
    },
    async analyzeTrace(path) {
      this.analyzingTrace = true
      this.analysisResult = null
      try {
        if (window.go?.main?.App) this.analysisResult = await window.go.main.App.AnalyzeETLFile(path || this.selectedTrace)
      } catch (e) { console.error(e) }
      finally { this.analyzingTrace = false }
    },
    async loadETLFile() {
      this.loadingETL = true
      try {
        const filePath = await OpenETLFileDialog()
        if (filePath) {
          this.loadedETLPath = filePath
          // Auto-analyze after loading
          await this.analyzeTrace(filePath)
        }
      } catch (e) { console.error(e) }
      finally { this.loadingETL = false }
    },
    async openWPA() {
      try {
        if (window.go?.main?.App) {
          await window.go.main.App.OpenETLInWPA(this.analysisResult.traceInfo.path)
        }
      } catch (e) { console.error(e) }
    },
    // ============ System Event Log ============
    async captureEventLog() {
      this.capturingEventLog = true
      this.eventLogResult = null
      try {
        const preset = this.eventLogPresets.find(p => p.label === this.selectedEventLogPreset)
        if (!preset) return
        if (window.go?.main?.App) {
          this.eventLogResult = await window.go.main.App.CaptureSystemEventLog(preset.hours, preset.maxEvents)
        }
      } catch (e) {
        console.error(e)
        this.eventLogResult = { error: e.message || String(e), totalEvents: 0, recentErrors: [] }
      } finally {
        this.capturingEventLog = false
      }
    },
    async exportEventLog() {
      try {
        const preset = this.eventLogPresets.find(p => p.label === this.selectedEventLogPreset)
        if (!preset) return
        if (window.go?.main?.App) {
          await window.go.main.App.ExportSystemEventLog('', preset.hours, preset.maxEvents)
          // Open folder after export
          await OpenFolder('C:\\Users\\Public\\ETL_Traces')
        }
      } catch (e) { console.error(e) }
    },
    formatEventTime(ts) {
      if (!ts) return ''
      // Handle various timestamp formats
      const m = ts.match(/(\d{4}-\d{2}-\d{2}\s?T?\d{2}:\d{2}:\d{2})/)
      return m ? m[1].replace('T', ' ') : ts.substring(0, 19)
    },
    openEventViewer() {
      try {
        if (window.go?.main?.App) {
          window.go.main.App.OpenEventViewer()
        }
      } catch (e) { console.error(e) }
    },
    async analyzeEventLog() {
      if (!this.eventLogResult) return
      this.analyzingEventLog = true
      this.eventLogAnalysis = null
      try {
        if (window.go?.main?.App) {
          this.eventLogAnalysis = await window.go.main.App.AnalyzeSystemEventLog(this.eventLogResult)
        }
      } catch (e) {
        this.eventLogAnalysis = { overallHealth: 'Error', summary: 'Analysis failed: ' + (e.message || String(e)), rootCauses: [], repeatErrors: [], knownIssues: [], correlations: [], recommendations: [], timeline: [] }
      } finally {
        this.analyzingEventLog = false
      }
    },
    async loadEVTX() {
      this.loadingEVTX = true
      try {
        if (window.go?.main?.App) {
          const path = await window.go.main.App.OpenEVTXFileDialog()
          if (!path) { this.loadingEVTX = false; return }
          this.evtxFilePath = path
          this.eventLogResult = await window.go.main.App.LoadEVTXFile(path)
          this.eventLogAnalysis = null
          this.evtxLevelFilter = ''
        }
      } catch (e) {
        console.error(e)
        this.eventLogResult = { error: e.message || String(e), totalEvents: 0, recentErrors: [], recentEvents: [] }
      } finally {
        this.loadingEVTX = false
      }
    },
    async evtxToCSV() {
      if (!this.evtxFilePath) return
      try {
        if (window.go?.main?.App) {
          const csvPath = await window.go.main.App.ExportEVTXToCSV(this.evtxFilePath, '')
          if (csvPath && !csvPath.includes('failed')) {
            alert('CSV exported to: ' + csvPath)
          } else if (csvPath) {
            alert(csvPath)
          }
        }
      } catch (e) { console.error(e); alert('Export failed: ' + (e.message || String(e))) }
    },
    // ============ Dispdiag ============
    async captureDispdiag() {
      this.capturingDispdiag = true
      this.dispdiagResult = null
      try {
        if (window.go?.main?.App) {
          this.dispdiagResult = await window.go.main.App.RunDispdiag('', this.dispdiagDelay, this.dispdiagDump)
        }
      } catch (e) {
        console.error(e)
        this.dispdiagResult = { errors: [e.message || String(e)], summary: {} }
      } finally {
        this.capturingDispdiag = false
      }
    },
    async openDispdiagFolder() {
      try {
        if (window.go?.main?.App) {
          const dir = await window.go.main.App.GetDispdiagOutputDir()
          await OpenFolder(dir)
        }
      } catch (e) { console.error(e) }
    },
    async openDispdiagLog() {
      try {
        if (window.go?.main?.App) {
          const result = await window.go.main.App.OpenDispdiagLog()
          if (result && result.includes('No dispdiag log')) {
            alert(result)
          }
        }
      } catch (e) { console.error(e) }
    },
    async exportDispdiagJSON() {
      try {
        if (window.go?.main?.App) {
          await window.go.main.App.ExportDispdiagResult(this.dispdiagResult, '')
        }
      } catch (e) { console.error(e) }
    },
    async loadLogs() {
      try {
        if (window.go?.main?.App) this.logFiles = await window.go.main.App.GetLogFiles() || []
      } catch (e) { console.error(e) }
    },
    selectLog(name) { this.selectedLog = name },
    toggleLogSection() {
      this.logSectionExpanded = !this.logSectionExpanded
    },
    unlockLogs() {
      if (this.logPassword === 'Lenovo2026') {
        this.logUnlocked = true
        this.logPasswordError = ''
        this.logPassword = ''
        this.loadLogs()
      } else {
        this.logPasswordError = 'Need Dispatcher owner check or contact zhoushang2'
        setTimeout(() => { this.logPasswordError = '' }, 3000)
      }
    },
    async runAnalysis() {
      this.analyzing = true
      this.summary = null
      try {
        if (window.go?.main?.App) this.summary = await window.go.main.App.GetLogSummary()
      } catch (e) { console.error(e) }
      finally { this.analyzing = false }
    },
    async loadRawLog() {
      try {
        if (window.go?.main?.App) this.rawLog = await window.go.main.App.ReadLogTail(this.tailLines)
      } catch (e) { console.error(e) }
    },
    formatSize(bytes) {
      if (bytes < 1024) return bytes + ' B'
      if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
      return (bytes / 1024 / 1024).toFixed(1) + ' MB'
    },
    formatNum(n) {
      if (n >= 1e9) return (n / 1e9).toFixed(1) + 'B'
      if (n >= 1e6) return (n / 1e6).toFixed(1) + 'M'
      if (n >= 1000) return (n / 1000).toFixed(1) + 'K'
      return String(n)
    },
    truncatePath(path) {
      if (!path) return ''
      return path.length <= 65 ? path : '...' + path.slice(-62)
    },
    // ============ Toolkit Methods ============
    async switchToToolkit() {
      this.activeAIType = 'toolkit'
      if (!this.toolkitTools.length) await this.loadToolkitData()
    },
    async loadToolkitData() {
      try {
        if (window.go?.main?.App) {
          this.toolkitTools = await window.go.main.App.GetToolkitTools() || []
          this.toolkitInstallDir = await window.go.main.App.GetToolkitInstallDir() || this.toolkitInstallDir
          await this.refreshToolkitStatus()
        }
      } catch (e) { console.error(e) }
    },
    async refreshToolkitStatus() {
      try {
        if (window.go?.main?.App) {
          const statuses = await window.go.main.App.CheckAllToolkitInstalled() || []
          this.toolInstallStatus = {}
          statuses.forEach(s => {
            this.toolInstallStatus[s.toolId] = s
          })
        }
      } catch (e) { console.error(e) }
    },
    async switchToMSR() {
      this.activeAIType = 'msr'
    },
    async switchToWMI() {
      this.activeAIType = 'wmi'
    },
    async switchToPPM() {
      this.activeAIType = 'ppm'
    },
    async installTool(toolId) {
      if (!window.go?.main?.App) return
      // Set installing state
      this.toolProgress[toolId] = { status: 'installing', progress: 0, message: 'Installing via winget...' }
      
      try {
        // Start installation (runs async in backend)
        await window.go.main.App.InstallToolkitTool(toolId)
        
        // Poll for progress
        const pollProgress = async () => {
          if (!this.toolProgress[toolId] || this.toolProgress[toolId].status === 'completed' || this.toolProgress[toolId].status === 'error') {
            return
          }
          try {
            const progress = await window.go.main.App.GetToolkitProgress(toolId)
            this.toolProgress[toolId] = progress
            if (progress.status === 'installing') {
              setTimeout(pollProgress, 1000)
            } else if (progress.status === 'completed') {
              await this.refreshToolkitStatus()
            }
          } catch (e) {
            console.error('Poll error:', e)
          }
        }
        setTimeout(pollProgress, 500)
      } catch (e) {
        this.toolProgress[toolId] = { status: 'error', message: e.message || String(e) }
      }
    },
    async runTool(toolId) {
      try {
        if (window.go?.main?.App) {
          const result = await window.go.main.App.RunToolkitTool(toolId)
          if (result.startsWith('Error')) {
            alert(result)
          }
        }
      } catch (e) { console.error(e) }
    },
    async openToolkitFolder() {
      try {
        if (window.go?.main?.App) {
          await window.go.main.App.OpenToolkitFolder()
        }
      } catch (e) { console.error(e) }
    },
    isToolBusy(toolId) {
      const p = this.toolProgress[toolId]
      return p && p.status === 'installing'
    },
    async installEssentials() {
      this.batchInstalling = true
      const essentials = ['hwinfo64', 'cpuz', 'gpuz']
      for (const id of essentials) {
        if (!this.toolInstallStatus[id]?.installed) {
          await this.installTool(id)
          // Wait for completion
          let attempts = 0
          while (this.isToolBusy(id) && attempts < 120) {
            await new Promise(r => setTimeout(r, 1000))
            attempts++
          }
        }
      }
      this.batchInstalling = false
    },
  },
  computed: {
    filteredTools() {
      if (this.toolkitCategory === 'all') return this.toolkitTools
      return this.toolkitTools.filter(t => t.category === this.toolkitCategory)
    },
    allToolsInstalled() {
      return this.toolkitTools.every(t => this.toolInstallStatus[t.id]?.installed)
    },
    filteredEvents() {
      if (!this.eventLogResult || !this.eventLogResult.recentEvents) return []
      let events
      if (!this.evtxLevelFilter) {
        // Default: show critical + errors
        events = this.eventLogResult.recentErrors || []
      } else {
        events = this.eventLogResult.recentEvents.filter(e => e.level === this.evtxLevelFilter)
      }
      // Sort: Critical first (by time desc), then others (by time desc)
      const critical = events.filter(e => e.level === 'Critical').sort((a, b) => new Date(b.time) - new Date(a.time))
      const others = events.filter(e => e.level !== 'Critical').sort((a, b) => new Date(b.time) - new Date(a.time))
      return [...critical, ...others]
    }
  }
}
</script>

<style scoped>
.ai-analysis-page { display: flex; flex-direction: column; gap: 16px; }

.ai-tabs {
  display: flex; gap: 8px;
}
.ai-tab {
  display: flex; align-items: center; gap: 6px;
  padding: 8px 16px;
  background: var(--bg-card); border: 1px solid var(--border-color);
  border-radius: 8px; color: var(--text-secondary);
  font-size: 13px; font-weight: 500; cursor: pointer;
  transition: var(--transition); font-family: inherit;
}
.ai-tab:hover { border-color: var(--lenovo-red); color: var(--text-primary); }
.ai-tab.active {
  background: linear-gradient(90deg, rgba(230,63,50,0.15) 0%, rgba(230,63,50,0.05) 100%);
  border-color: var(--lenovo-red); color: var(--text-primary);
}

/* ETL */
.etl-section, .log-section, .salog-section { display: flex; flex-direction: column; gap: 16px; }

.elevated-warning {
  display: flex; gap: 12px; padding: 14px 16px;
  background: rgba(255,152,0,0.08); border: 1px solid rgba(255,152,0,0.3);
  border-radius: 8px; color: #FF9800; font-size: 12px; line-height: 1.6;
}
.elevated-warning svg { flex-shrink: 0; margin-top: 2px; }
.elevated-warning strong { display: block; margin-bottom: 4px; font-size: 13px; }
.elevated-warning p { margin: 0; opacity: 0.85; }

.profile-grid {
  display: grid; grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 8px; margin-bottom: 16px;
}
.profile-btn {
  display: flex; flex-direction: column; align-items: center; gap: 4px;
  padding: 12px 8px; background: var(--bg-tertiary);
  border: 1px solid var(--border-color); border-radius: 8px;
  cursor: pointer; transition: var(--transition); text-align: center; font-family: inherit;
}
.profile-btn:hover:not(.disabled) { border-color: var(--lenovo-red); background: rgba(230,63,50,0.06); }
.profile-btn.active { border-color: var(--lenovo-red); background: rgba(230,63,50,0.12); }
.profile-btn.disabled { opacity: 0.5; cursor: not-allowed; }
.profile-icon { font-size: 20px; }
.profile-name { font-size: 12px; font-weight: 600; color: var(--text-primary); }
.profile-desc { font-size: 10px; color: var(--text-tertiary); line-height: 1.3; }

.capture-controls {
  display: flex; gap: 24px; flex-wrap: wrap; margin-bottom: 16px;
  padding: 12px 16px; background: var(--bg-tertiary);
  border-radius: 8px; border: 1px solid var(--border-color);
}
.control-group { display: flex; flex-direction: column; gap: 6px; }
.ctrl-label { font-size: 10px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.8px; color: var(--text-tertiary); }
.duration-buttons { display: flex; gap: 4px; }
.dur-btn {
  padding: 4px 10px; background: var(--bg-card); border: 1px solid var(--border-color);
  border-radius: 5px; color: var(--text-secondary); font-size: 12px; font-weight: 500;
  cursor: pointer; transition: var(--transition); font-family: inherit;
}
.dur-btn.active { background: var(--lenovo-red); border-color: var(--lenovo-red); color: white; }
.dur-btn:disabled { opacity: 0.4; cursor: not-allowed; }
.status-val { font-size: 13px; font-weight: 500; color: var(--text-primary); }
.status-val.recording { color: #4CAF50; }
.status-val.idle { color: var(--text-tertiary); }
.status-val.mono { font-family: 'Consolas','Monaco',monospace; font-size: 11px; }

.capture-actions { display: flex; gap: 8px; flex-wrap: wrap; }
.btn-start, .btn-stop {
  padding: 9px 20px; border: none; border-radius: 8px;
  color: white; font-size: 13px; font-weight: 600; cursor: pointer;
  display: flex; align-items: center; gap: 8px; transition: var(--transition); font-family: inherit;
}
.btn-start { background: linear-gradient(135deg, #4CAF50 0%, #2E7D32 100%); }
.btn-stop { background: linear-gradient(135deg, #F44336 0%, #C62828 100%); }
.btn-start:hover:not(:disabled), .btn-stop:hover:not(:disabled) { opacity: 0.9; transform: translateY(-1px); }
.btn-start:disabled, .btn-stop:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-outline {
  padding: 9px 16px; background: var(--bg-tertiary); border: 1px solid var(--border-color);
  border-radius: 8px; color: var(--text-secondary); font-size: 13px; font-weight: 500;
  cursor: pointer; display: flex; align-items: center; gap: 6px; transition: var(--transition); font-family: inherit;
}
.btn-outline:hover:not(:disabled) { background: var(--bg-card-hover); color: var(--text-primary); border-color: var(--border-light); }
.btn-outline:disabled { opacity: 0.4; cursor: not-allowed; }

.capture-badge {
  display: flex; align-items: center; gap: 6px; padding: 4px 10px;
  border-radius: 12px; font-size: 12px; font-weight: 600;
}
.capture-badge.recording { background: rgba(76,175,80,0.15); color: #4CAF50; border: 1px solid rgba(76,175,80,0.3); }
.rec-dot { width: 8px; height: 8px; border-radius: 50%; background: #4CAF50; animation: blink 1s infinite; }
@keyframes blink { 0%,100%{opacity:1} 50%{opacity:0.3} }

.capture-error {
  display: flex; align-items: flex-start; gap: 8px; margin-top: 12px;
  padding: 10px 14px; background: rgba(244,67,54,0.08);
  border: 1px solid rgba(244,67,54,0.2); border-radius: 6px;
  color: #F44336; font-size: 12px; line-height: 1.5;
}
.capture-error svg { flex-shrink: 0; margin-top: 1px; }

/* ===== Load External ETL ===== */
.load-etl-card { border: 1px dashed var(--border-color); }
.etl-load-body { padding: 16px 20px; }
.etl-load-hint { margin: 0 0 12px 0; font-size: 12px; color: var(--text-tertiary); }
.etl-load-row { display: flex; align-items: center; gap: 12px; }
.loaded-path {
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 10px;
  padding: 4px 10px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  color: var(--accent-green);
  word-break: break-all;
  flex: 1;
}

/* ===== System Event Log ===== */
.event-log-card { border: 1px dashed rgba(76,175,80,0.3); }
.event-log-body { padding: 16px 20px; }
.event-log-hint { margin: 0 0 12px 0; font-size: 12px; color: var(--text-tertiary); }
.event-log-controls { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
.event-log-presets { display: flex; gap: 4px; }
.preset-btn {
  padding: 4px 10px; background: var(--bg-tertiary); border: 1px solid var(--border-color);
  border-radius: 5px; color: var(--text-secondary); font-size: 12px; font-weight: 500;
  cursor: pointer; transition: var(--transition); font-family: inherit;
}
.preset-btn.active { background: var(--lenovo-red); border-color: var(--lenovo-red); color: white; }
.preset-btn:hover:not(.active) { border-color: var(--lenovo-red); }
.btn-capture-eventlog {
  padding: 8px 16px; border-radius: 8px; color: white; font-size: 13px; font-weight: 600;
  cursor: pointer; display: flex; align-items: center; gap: 6px;
  transition: var(--transition); font-family: inherit;
  background: linear-gradient(135deg, #4CAF50 0%, #2E7D32 100%);
  border: none;
}
.btn-capture-eventlog:hover:not(:disabled) { opacity: 0.9; transform: translateY(-1px); }
.btn-capture-eventlog:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-export-eventlog {
  padding: 8px 16px; border: 1px solid var(--border-color); border-radius: 8px;
  background: var(--bg-tertiary); color: var(--text-secondary); font-size: 13px; font-weight: 500;
  cursor: pointer; display: flex; align-items: center; gap: 6px;
  transition: var(--transition); font-family: inherit;
}
.btn-export-eventlog:hover { background: var(--bg-card-hover); color: var(--text-primary); border-color: var(--border-light); }
.btn-open-eventviewer {
  padding: 8px 16px; border: 1px solid var(--border-color); border-radius: 8px;
  background: var(--bg-tertiary); color: var(--text-secondary); font-size: 13px; font-weight: 500;
  cursor: pointer; display: flex; align-items: center; gap: 6px;
  transition: var(--transition); font-family: inherit;
}
.btn-open-eventviewer:hover { background: var(--bg-card-hover); color: var(--text-primary); border-color: var(--lenovo-red); }
.btn-load-evtx {
  padding: 8px 16px; border-radius: 8px; color: white; font-size: 13px; font-weight: 600;
  cursor: pointer; display: flex; align-items: center; gap: 6px;
  transition: var(--transition); font-family: inherit;
  background: linear-gradient(135deg, #FF9800 0%, #F57C00 100%); border: none;
}
.btn-load-evtx:hover:not(:disabled) { opacity: 0.9; transform: translateY(-1px); }
.btn-load-evtx:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-evtx-to-csv {
  padding: 8px 16px; border-radius: 8px; font-size: 13px; font-weight: 600;
  cursor: pointer; display: flex; align-items: center; gap: 6px;
  transition: var(--transition); font-family: inherit;
  background: transparent; border: 1px solid var(--border-color); color: var(--text-secondary);
}
.btn-evtx-to-csv:hover:not(:disabled) { border-color: #FF9800; color: #FF9800; background: rgba(255,152,0,0.06); }
.btn-evtx-to-csv:disabled { opacity: 0.4; cursor: not-allowed; }
.event-level-filter {
  display: flex; align-items: center; gap: 8px; padding: 8px 0 0 0; flex-wrap: wrap;
}
.filter-label { font-size: 12px; color: var(--text-tertiary); font-weight: 600; }
.level-btn {
  padding: 4px 12px; border-radius: 6px; font-size: 12px; font-weight: 600;
  cursor: pointer; transition: var(--transition); font-family: inherit;
  border: 1px solid var(--border-color); background: transparent; color: var(--text-secondary);
}
.level-btn.active.level-critical { background: #F44336; color: white; border-color: #F44336; }
.level-btn.active.level-error { background: #FF5722; color: white; border-color: #FF5722; }
.level-btn.active.level-warning { background: #FF9800; color: white; border-color: #FF9800; }
.level-btn.active.level-info { background: #2196F3; color: white; border-color: #2196F3; }
.level-btn.level-clear { background: transparent; border-color: var(--text-tertiary); color: var(--text-tertiary); }
.level-btn.level-clear:hover { border-color: var(--lenovo-red); color: var(--lenovo-red); }
.table-count {
  margin-left: 8px; font-size: 11px; font-weight: 700; background: var(--bg-tertiary);
  padding: 2px 8px; border-radius: 10px; color: var(--lenovo-red);
}

.event-log-results { padding: 0 20px 16px 20px; display: flex; flex-direction: column; gap: 12px; }
.event-log-summary {
  display: flex; gap: 8px; flex-wrap: wrap;
}
.event-stat {
  display: flex; flex-direction: column; align-items: center; gap: 2px;
  padding: 8px 16px; border-radius: 8px; min-width: 80px;
}
.event-stat.critical { background: rgba(244,67,54,0.1); border: 1px solid rgba(244,67,54,0.2); }
.event-stat.error { background: rgba(255,152,0,0.08); border: 1px solid rgba(255,152,0,0.2); }
.event-stat.warning { background: rgba(255,193,7,0.08); border: 1px solid rgba(255,193,7,0.2); }
.event-stat.info { background: rgba(33,150,243,0.08); border: 1px solid rgba(33,150,243,0.2); }
.event-stat.total { background: var(--bg-tertiary); border: 1px solid var(--border-color); }
.stat-count { font-size: 22px; font-weight: 700; font-family: 'Consolas','Monaco',monospace; color: var(--text-primary); }
.event-stat.critical .stat-count { color: #F44336; }
.event-stat.error .stat-count { color: #FF9800; }
.event-stat.warning .stat-count { color: #FFC107; }
.event-stat.info .stat-count { color: #2196F3; }
.stat-label { font-size: 10px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.5px; color: var(--text-tertiary); }

.event-error-table { }
.table-title {
  display: flex; align-items: center; gap: 6px; font-size: 12px; font-weight: 700;
  color: var(--text-secondary); margin-bottom: 8px; text-transform: uppercase; letter-spacing: 0.5px;
}
.table-scroll { overflow-y: auto; border-radius: 6px; border: 1px solid var(--border-color); }
.event-table { width: 100%; border-collapse: collapse; font-size: 11px; }
.event-table th {
  position: sticky; top: 0; background: var(--bg-tertiary);
  padding: 6px 8px; text-align: left; font-size: 10px; font-weight: 700;
  text-transform: uppercase; letter-spacing: 0.5px; color: var(--text-tertiary);
  border-bottom: 1px solid var(--border-color);
}
.event-table td { padding: 5px 8px; border-bottom: 1px solid var(--border-color); color: var(--text-secondary); }
.event-table tr:last-child td { border-bottom: none; }
.event-table tr:hover td { background: var(--bg-card-hover); }
.cell-time { white-space: nowrap; font-family: 'Consolas','Monaco',monospace; font-size: 10px; color: var(--text-tertiary); }
.cell-source { font-family: 'Consolas','Monaco',monospace; font-size: 11px; color: var(--text-primary); white-space: nowrap; }
.cell-id { font-family: 'Consolas','Monaco',monospace; font-size: 11px; color: var(--lenovo-red); font-weight: 600; text-align: center; }
.cell-msg { font-size: 11px; color: var(--text-secondary); max-width: 400px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.cell-eventdata { font-size: 11px; color: var(--text-tertiary); max-width: 300px; white-space: pre-line; word-break: break-word; line-height: 1.5; }

.level-badge {
  padding: 1px 6px; border-radius: 3px; font-size: 9px; font-weight: 700; text-transform: uppercase; white-space: nowrap;
}
.level-critical, .level-严重 { background: rgba(244,67,54,0.15); color: #F44336; }
.level-error, .level-错误 { background: rgba(255,152,0,0.15); color: #FF9800; }
.level-warning, .level-警告 { background: rgba(255,193,7,0.12); color: #FFC107; }
.level-information, .level-信息, .level-verbose { background: rgba(33,150,243,0.1); color: #2196F3; }
.row-critical td {
  color: #F44336 !important;
  font-weight: 700 !important;
  background: rgba(244,67,54,0.06);
}
.event-table th,
.event-table td {
  border-right: 1px solid rgba(255,255,255,0.08);
}
.event-table th:last-child,
.event-table td:last-child {
  border-right: none;
}

.event-providers { }
.provider-list { display: flex; flex-wrap: wrap; gap: 4px; }
.provider-row {
  display: flex; align-items: center; gap: 6px; padding: 4px 10px;
  background: var(--bg-tertiary); border: 1px solid var(--border-color);
  border-radius: 5px; font-size: 11px;
}
.provider-name { color: var(--text-primary); font-family: 'Consolas','Monaco',monospace; }

/* ===== Event Log Analysis ===== */
.eventlog-analysis-card { border: 1px solid rgba(156,39,176,0.2); background: linear-gradient(135deg, rgba(156,39,176,0.02) 0%, transparent 100%); }
.eventlog-analysis-card .card-header { display: flex; justify-content: space-between; align-items: center; }
.btn-analyze-eventlog { padding: 6px 16px; border-radius: 6px; font-size: 13px; font-weight: 600; background: var(--lenovo-red); color: #fff; border: none; cursor: pointer; display: flex; align-items: center; gap: 6px; transition: all 0.2s; }
.btn-analyze-eventlog:hover:not(:disabled) { opacity: 0.9; transform: translateY(-1px); }
.btn-analyze-eventlog:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-analyze-eventlog svg { width: 14px; height: 14px; }
.eventlog-analysis-body { padding: 16px 20px; display: flex; flex-direction: column; gap: 14px; }
.analysis-hint { padding: 20px; text-align: center; color: var(--text-tertiary); font-size: 13px; }

.health-badge { display: flex; align-items: center; gap: 12px; padding: 14px 16px; border-radius: 8px; border: 1px solid; }
.health-badge.good { background: rgba(76,175,80,0.08); border-color: rgba(76,175,80,0.2); }
.health-badge.warning { background: rgba(255,152,0,0.08); border-color: rgba(255,152,0,0.2); }
.health-badge.critical { background: rgba(244,67,54,0.08); border-color: rgba(244,67,54,0.2); }
.health-label { font-size: 16px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.5px; }
.health-badge.good .health-label { color: #4CAF50; }
.health-badge.warning .health-label { color: #FF9800; }
.health-badge.critical .health-label { color: #F44336; }
.health-summary { flex: 1; font-size: 13px; color: var(--text-secondary); }

.analysis-section { display: flex; flex-direction: column; gap: 8px; }
.analysis-section .sub-title { display: flex; align-items: center; gap: 8px; font-size: 13px; font-weight: 600; color: var(--text-primary); margin-bottom: 4px; }

.root-cause-list { display: flex; flex-direction: column; gap: 6px; }
.root-cause-item { display: flex; align-items: center; gap: 10px; padding: 10px 12px; border-radius: 6px; border: 1px solid var(--border-color); background: var(--bg-card); }
.rc-category { font-weight: 600; color: var(--text-primary); min-width: 100px; }
.rc-detail { flex: 1; font-size: 12px; color: var(--text-secondary); }
.rc-severity { font-size: 11px; padding: 2px 8px; border-radius: 10px; font-weight: 600; text-transform: uppercase; }
.rc-severity.critical { background: rgba(244,67,54,0.1); color: #F44336; }
.rc-severity.warning { background: rgba(255,152,0,0.1); color: #FF9800; }
.rc-severity.info { background: rgba(33,150,243,0.1); color: #2196F3; }

.repeat-list { display: flex; flex-direction: column; gap: 6px; }
.repeat-item { padding: 10px 12px; border-radius: 6px; border: 1px solid var(--border-color); background: var(--bg-card); }
.repeat-header { display: flex; align-items: center; gap: 8px; margin-bottom: 4px; }
.repeat-source { font-family: 'Consolas','Monaco',monospace; font-size: 12px; color: var(--lenovo-red); }
.repeat-count { background: var(--bg-tertiary); padding: 2px 8px; border-radius: 10px; font-size: 11px; font-weight: 600; color: var(--text-primary); }
.repeat-time { font-size: 11px; color: var(--text-tertiary); margin-left: auto; }
.repeat-msg { font-size: 11px; color: var(--text-secondary); padding-left: 4px; }

.known-issue-list { display: flex; flex-direction: column; gap: 8px; }
.known-issue-item { padding: 12px; border-radius: 6px; border: 1px solid rgba(255,152,0,0.3); background: rgba(255,152,0,0.05); }
.ki-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 6px; }
.ki-name { font-weight: 600; color: var(--text-primary); }
.ki-confidence { font-size: 11px; padding: 2px 10px; border-radius: 10px; font-weight: 600; }
.ki-confidence.high { background: rgba(244,67,54,0.1); color: #F44336; }
.ki-confidence.mid { background: rgba(255,152,0,0.1); color: #FF9800; }
.ki-confidence.low { background: rgba(33,150,243,0.1); color: #2196F3; }
.ki-detail { font-size: 12px; color: var(--text-secondary); line-height: 1.5; }

.correlation-list { display: flex; flex-direction: column; gap: 6px; }
.correlation-item { padding: 8px 12px; border-radius: 6px; border: 1px solid var(--border-color); background: var(--bg-card); font-size: 12px; color: var(--text-secondary); }

.timeline-list { display: flex; flex-direction: column; gap: 4px; max-height: 300px; overflow-y: auto; }
.timeline-item { display: flex; align-items: center; gap: 8px; padding: 6px 10px; border-radius: 4px; font-size: 11px; background: var(--bg-card); border-left: 3px solid; }
.timeline-item.critical { border-left-color: #F44336; background: rgba(244,67,54,0.03); }
.timeline-item.error { border-left-color: #FF9800; background: rgba(255,152,0,0.02); }
.tl-time { color: var(--text-tertiary); min-width: 85px; font-family: 'Consolas','Monaco',monospace; }
.tl-level { font-weight: 600; text-transform: uppercase; min-width: 55px; }
.tl-level.critical { color: #F44336; }
.tl-level.error { color: #FF9800; }
.tl-source { color: var(--lenovo-red); }
.tl-msg { flex: 1; color: var(--text-secondary); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

.recommendation-list { margin: 0; padding-left: 20px; display: flex; flex-direction: column; gap: 6px; }
.recommendation-list li { font-size: 13px; color: var(--text-secondary); line-height: 1.5; }

.provider-count { color: var(--lenovo-red); font-weight: 700; }

/* ===== Dispdiag Capture ===== */
.dispdiag-card { border: 1px dashed rgba(33,150,243,0.3); }
.dispdiag-body { padding: 16px 20px; }
.dispdiag-hint { margin: 0 0 12px 0; font-size: 12px; color: var(--text-tertiary); }
.dispdiag-hint code { background: var(--bg-tertiary); padding: 1px 5px; border-radius: 3px; font-size: 11px; color: var(--lenovo-red); }
.dispdiag-controls { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
.dispdiag-options { display: flex; align-items: center; gap: 16px; }
.dispdiag-opt {
  display: flex; align-items: center; gap: 6px; font-size: 12px; color: var(--text-secondary);
  cursor: pointer; user-select: none;
}
.dispdiag-opt input[type="checkbox"] { accent-color: var(--lenovo-red); cursor: pointer; }
.delay-input {
  width: 42px; background: var(--bg-tertiary); border: 1px solid var(--border-color);
  border-radius: 4px; color: var(--text-primary); padding: 2px 4px; font-size: 12px;
  font-family: 'Consolas','Monaco',monospace; text-align: center;
}
.btn-capture-dispdiag {
  padding: 8px 16px; border-radius: 8px; color: white; font-size: 13px; font-weight: 600;
  cursor: pointer; display: flex; align-items: center; gap: 6px;
  transition: var(--transition); font-family: inherit;
  background: linear-gradient(135deg, #2196F3 0%, #1565C0 100%);
  border: none;
}
.btn-capture-dispdiag:hover:not(:disabled) { opacity: 0.9; transform: translateY(-1px); }
.btn-capture-dispdiag:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-open-dispdiag {
  padding: 8px 16px; border-radius: 8px; font-size: 13px; font-weight: 600;
  cursor: pointer; display: flex; align-items: center; gap: 6px;
  transition: var(--transition); font-family: inherit;
  background: transparent; border: 1px solid var(--border-color); color: var(--text-secondary);
}
.btn-open-dispdiag:hover { border-color: #2196F3; color: #2196F3; background: rgba(33,150,243,0.06); transform: translateY(-1px); }

.dispdiag-results { padding: 0 20px 16px 20px; display: flex; flex-direction: column; gap: 12px; }
.dispdiag-info-bar {
  display: flex; gap: 12px; flex-wrap: wrap; align-items: center;
  padding: 10px 14px; background: var(--bg-tertiary); border-radius: 8px;
  border: 1px solid var(--border-color);
}
.dispdiag-info-item { display: flex; align-items: center; gap: 4px; font-size: 12px; }
.info-label { color: var(--text-tertiary); font-weight: 600; text-transform: uppercase; letter-spacing: 0.5px; font-size: 10px; }
.info-value { color: var(--text-primary); font-size: 12px; }
.info-value code { font-family: 'Consolas','Monaco',monospace; font-size: 11px; }
.pass-badge {
  display: flex; align-items: center; gap: 4px; color: #4CAF50; font-weight: 700;
  font-size: 12px; background: rgba(76,175,80,0.1); padding: 2px 8px; border-radius: 4px;
}
.fail-badge {
  display: flex; align-items: center; gap: 4px; color: #F44336; font-weight: 700;
  font-size: 12px; background: rgba(244,67,54,0.1); padding: 2px 8px; border-radius: 4px;
}

.sub-title {
  display: flex; align-items: center; gap: 6px; font-size: 12px; font-weight: 700;
  color: var(--text-secondary); margin-bottom: 6px; text-transform: uppercase; letter-spacing: 0.5px;
}
.dispdiag-errors { }
.error-list, .warn-list {
  margin: 0; padding: 0 0 0 18px; font-size: 11px; color: var(--text-secondary);
  max-height: 200px; overflow-y: auto;
}
.error-list li { color: #F44336; margin-bottom: 2px; font-family: 'Consolas','Monaco',monospace; }
.warn-list li { color: #FF9800; margin-bottom: 2px; font-family: 'Consolas','Monaco',monospace; }

.dispdiag-summary { }
.summary-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); gap: 6px; }
.summary-row {
  display: flex; align-items: center; gap: 8px; padding: 6px 10px;
  background: var(--bg-tertiary); border: 1px solid var(--border-color);
  border-radius: 5px;
}
.grid-key { color: var(--text-tertiary); font-size: 10px; text-transform: uppercase; font-weight: 600; letter-spacing: 0.5px; }
.grid-val { color: var(--text-primary); font-family: 'Consolas','Monaco',monospace; font-size: 11px; }

.dispdiag-brightness { }
.brightness-list { display: flex; flex-direction: column; gap: 3px; }
.brightness-line {
  display: block; padding: 4px 8px; background: var(--bg-tertiary); border-radius: 4px;
  font-size: 11px; color: var(--text-secondary); font-family: 'Consolas','Monaco',monospace;
}

.dispdiag-content { }
.dispdiag-preview {
  margin: 0; padding: 10px 12px; background: #1a1a2e; border-radius: 6px;
  border: 1px solid var(--border-color); font-size: 10px; line-height: 1.4;
  color: #a0a0b0; font-family: 'Consolas','Monaco',monospace;
  max-height: 250px; overflow: auto; white-space: pre-wrap; word-break: break-all;
}

.dispdiag-actions { display: flex; gap: 8px; flex-wrap: wrap; }
.btn-dispdiag-action {
  padding: 8px 14px; border: 1px solid var(--border-color); border-radius: 8px;
  background: var(--bg-tertiary); color: var(--text-secondary); font-size: 12px; font-weight: 500;
  cursor: pointer; display: flex; align-items: center; gap: 6px;
  transition: var(--transition); font-family: inherit;
}
.btn-dispdiag-action:hover { background: var(--bg-card-hover); color: var(--text-primary); border-color: var(--lenovo-red); }

/* ===== Progress Steps ===== */
.step-list { display: flex; flex-direction: column; gap: 6px; }
.step-row {
  display: flex; align-items: center; gap: 10px;
  padding: 8px 12px; border-radius: 6px;
  font-size: 12px; background: var(--bg-tertiary);
  border-left: 3px solid var(--border-color);
}
.step-row.step-done { border-left-color: #4CAF50; background: rgba(76,175,80,0.04); }
.step-row.step-error { border-left-color: #F44336; background: rgba(244,67,54,0.04); }
.step-icon { flex-shrink: 0; width: 18px; display: flex; align-items: center; }
.step-name { font-weight: 600; white-space: nowrap; color: var(--text-primary); }
.step-detail {
  color: var(--text-tertiary); font-family: 'Consolas','Monaco',monospace;
  font-size: 11px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
  flex: 1; min-width: 0;
}

/* ===== Event Types ===== */
.event-list { display: flex; flex-direction: column; gap: 3px; }
.event-row {
  display: flex; align-items: center; gap: 8px;
  padding: 4px 0; border-bottom: 1px solid var(--border-color); font-size: 11px;
}
.event-row:last-child { border-bottom: none; }
.event-cat {
  padding: 1px 6px; border-radius: 3px; font-size: 9px; font-weight: 700;
  text-transform: uppercase; white-space: nowrap;
}
.cat-process { background: rgba(33,150,243,0.15); color: #2196F3; }
.cat-disk { background: rgba(156,39,176,0.15); color: #9C27B0; }
.cat-network { background: rgba(0,188,212,0.15); color: #00BCD4; }
.cat-power { background: rgba(255,152,0,0.15); color: #FF9800; }
.cat-gpu { background: rgba(76,175,80,0.15); color: #4CAF50; }
.cat-memory { background: rgba(233,30,99,0.15); color: #E91E63; }
.cat-thread { background: rgba(158,158,158,0.15); color: #9E9E9E; }
.cat-other { background: var(--bg-tertiary); color: var(--text-tertiary); }
.event-name { color: var(--text-primary); font-family: 'Consolas','Monaco',monospace; flex: 1; }
.event-count { color: var(--text-secondary); font-weight: 600; white-space: nowrap; }

/* ===== Issues ===== */
.card-issues { border-color: rgba(244,67,54,0.2); background: rgba(244,67,54,0.03); }
.issue-list { display: flex; flex-direction: column; gap: 4px; }
.issue-row { display: flex; justify-content: space-between; align-items: center; padding: 5px 0; border-bottom: 1px solid rgba(244,67,54,0.1); font-size: 12px; }
.issue-row:last-child { border-bottom: none; }
.issue-kw { font-family: 'Consolas','Monaco',monospace; color: #F44336; font-weight: 600; text-transform: uppercase; }
.issue-count { color: var(--text-secondary); font-weight: 600; }

/* ===== WPA Card ===== */
.wpa-card { border: 1px dashed rgba(230,63,50,0.3); }
.wpa-hint { padding: 8px 16px 16px 16px; font-size: 12px; color: var(--text-tertiary); }
.btn-wpa {
  padding: 7px 14px;
  background: linear-gradient(135deg, #7B1FA2 0%, #4A148C 100%);
  border: none; border-radius: 6px; color: white; font-size: 12px; font-weight: 600;
  cursor: pointer; display: flex; align-items: center; gap: 6px;
  transition: var(--transition); font-family: inherit;
}
.btn-wpa:hover { opacity: 0.9; transform: translateY(-1px); }

/* Original analysis grid styles below */
.analysis-results { display: flex; flex-direction: column; gap: 12px; }
.result-summary-banner {
  display: flex; align-items: center; gap: 10px; padding: 12px 16px;
  background: linear-gradient(90deg, rgba(230,63,50,0.08) 0%, rgba(230,63,50,0.03) 100%);
  border: 1px solid rgba(230,63,50,0.2); border-radius: 8px;
  color: var(--text-secondary); font-size: 13px; font-family: 'Consolas','Monaco',monospace;
}
.analysis-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(240px, 1fr)); gap: 12px; }
.analysis-card {
  background: var(--bg-card); border: 1px solid var(--border-color);
  border-radius: 8px; padding: 14px; display: flex; flex-direction: column; gap: 6px;
}
.card-title-sm {
  display: flex; align-items: center; gap: 6px; font-size: 12px; font-weight: 700;
  text-transform: uppercase; letter-spacing: 0.8px; color: var(--lenovo-red);
  margin-bottom: 6px; padding-bottom: 8px; border-bottom: 1px solid var(--border-color);
}
.metric-row { display: flex; justify-content: space-between; align-items: center; }
.metric-label { font-size: 12px; color: var(--text-secondary); }
.metric-value { font-size: 12px; font-weight: 600; font-family: 'Consolas','Monaco',monospace; color: var(--text-primary); }
.metric-value.highlight { color: var(--lenovo-red); }
.metric-value.warning { color: #FF9800; }
.proc-list { margin-top: 4px; }
.proc-header { font-size: 10px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.8px; color: var(--text-tertiary); margin-bottom: 4px; }
.proc-row { display: flex; justify-content: space-between; align-items: center; padding: 3px 0; border-bottom: 1px solid var(--border-color); }
.proc-row:last-child { border-bottom: none; }
.proc-name { font-size: 11px; font-family: 'Consolas','Monaco',monospace; color: var(--text-secondary); }
.proc-pct { font-size: 11px; font-weight: 600; font-family: 'Consolas','Monaco',monospace; color: var(--text-primary); }
.proc-pct.warning { color: #FF9800; }
.no-data-sm { font-size: 11px; color: var(--text-tertiary); font-style: italic; }

.analyzing-overlay {
  display: flex; align-items: center; justify-content: center; gap: 12px;
  padding: 32px; background: var(--bg-card); border: 1px solid var(--border-color);
  border-radius: 8px; color: var(--text-secondary); font-size: 13px;
}

/* Log */
.log-file-list { display: flex; flex-direction: column; gap: 4px; max-height: 220px; overflow-y: auto; }
.log-file-item { display: flex; justify-content: space-between; align-items: center; padding: 8px 12px; border-radius: 6px; cursor: pointer; transition: background 0.15s; border: 1px solid transparent; }
.log-file-item:hover { background: var(--bg-card-hover); }
.log-file-item.active { background: rgba(230,63,50,0.08); border-color: rgba(230,63,50,0.25); }
.log-name { font-size: 13px; font-weight: 500; color: var(--text-primary); font-family: 'Consolas','Monaco',monospace; }
.log-meta { font-size: 11px; color: var(--text-tertiary); white-space: nowrap; margin-left: 12px; }

/* Log Files - Password Protected */
.advanced-section { border-color: rgba(245, 158, 11, 0.3); }
.advanced-toggle { cursor: pointer; user-select: none; }
.advanced-toggle .chevron { transition: transform 0.2s; margin-left: auto; }
.advanced-toggle .chevron.open { transform: rotate(180deg); }
.password-prompt { padding: 16px 0 0 0; }
.password-row { display: flex; align-items: center; gap: 8px; }
.password-input { flex: 1; padding: 8px 12px; border: 1px solid var(--border-color); border-radius: 6px; background: var(--bg-primary); color: var(--text-primary); font-size: 13px; outline: none; transition: border-color 0.2s; }
.password-input:focus { border-color: #F59E0B; }
.password-input::placeholder { color: var(--text-tertiary); font-family: inherit; }
.btn-unlock { display: flex; align-items: center; gap: 4px; padding: 8px 16px; border: 1px solid #F59E0B; border-radius: 6px; background: rgba(245,158,11,0.1); color: #F59E0B; font-size: 13px; font-weight: 600; cursor: pointer; transition: all 0.2s; white-space: nowrap; }
.btn-unlock:hover { background: rgba(245,158,11,0.2); border-color: #F59E0B; }
.password-error { margin-top: 8px; font-size: 12px; color: #EF4444; }

.summary-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-top: 4px; }
.summary-panel { background: var(--bg-tertiary); border-radius: 8px; padding: 14px; border: 1px solid var(--border-color); }
.panel-errors { border-color: rgba(244,67,54,0.2); background: rgba(244,67,54,0.04); }
.panel-title { font-size: 11px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.8px; color: var(--text-secondary); margin-bottom: 10px; display: flex; align-items: center; gap: 6px; }
.error-title { color: #F44336; }
.log-lines { display: flex; flex-direction: column; gap: 4px; max-height: 160px; overflow-y: auto; }
.log-line { font-size: 11px; font-family: 'Consolas','Monaco',monospace; color: var(--text-secondary); line-height: 1.5; word-break: break-all; padding: 2px 0; border-bottom: 1px solid var(--border-color); }
.log-line:last-child { border-bottom: none; }
.error-line { color: #F44336; }
.no-data { font-size: 12px; color: var(--text-tertiary); font-style: italic; }
.ok-hint { color: var(--accent-green); font-style: normal; font-weight: 600; }

.raw-log-box { background: var(--bg-tertiary); border-radius: 8px; padding: 16px; max-height: 400px; overflow-y: auto; border: 1px solid var(--border-color); }
.raw-log-box pre { font-family: 'Consolas','Monaco',monospace; font-size: 11px; color: var(--text-secondary); line-height: 1.6; white-space: pre-wrap; word-break: break-all; margin: 0; }

.btn-sm { padding: 6px 12px; background: var(--bg-tertiary); border: 1px solid var(--border-color); border-radius: 6px; color: var(--text-secondary); font-size: 12px; font-weight: 500; cursor: pointer; display: flex; align-items: center; gap: 6px; transition: var(--transition); font-family: inherit; }
.btn-sm:hover { background: var(--bg-card-hover); color: var(--text-primary); border-color: var(--border-light); }
.btn-analyze { padding: 7px 16px; background: linear-gradient(135deg, var(--lenovo-red) 0%, var(--lenovo-red-dark) 100%); border: none; border-radius: 6px; color: white; font-size: 12px; font-weight: 600; cursor: pointer; display: flex; align-items: center; gap: 6px; transition: var(--transition); font-family: inherit; }
.btn-analyze:hover:not(:disabled) { opacity: 0.9; transform: translateY(-1px); }
.btn-analyze:disabled { opacity: 0.6; cursor: not-allowed; }
.header-controls { display: flex; align-items: center; gap: 8px; }
.lines-select { padding: 5px 8px; background: var(--bg-tertiary); border: 1px solid var(--border-color); border-radius: 6px; color: var(--text-primary); font-size: 12px; font-family: inherit; cursor: pointer; }
.empty-hint { padding: 24px; text-align: center; color: var(--text-tertiary); font-size: 13px; }
.analyzing-hint { display: flex; align-items: center; justify-content: center; gap: 10px; padding: 32px; color: var(--text-secondary); font-size: 13px; }
@keyframes spin { to { transform: rotate(360deg); } }
.spinning { animation: spin 0.8s linear infinite; }

/* Toolkit Styles */
.toolkit-section, .ppm-section { display: flex; flex-direction: column; gap: 16px; }
.toolkit-path { font-size: 12px; color: var(--text-tertiary); padding: 8px 12px; background: var(--bg-tertiary); border-radius: 6px; font-family: 'Consolas', monospace; }

.category-filter { display: flex; gap: 8px; flex-wrap: wrap; }
.cat-btn {
  padding: 6px 14px; background: var(--bg-card); border: 1px solid var(--border-color);
  border-radius: 16px; color: var(--text-secondary); font-size: 12px; font-weight: 500;
  cursor: pointer; transition: var(--transition); font-family: inherit;
}
.cat-btn:hover { border-color: var(--lenovo-red); color: var(--text-primary); }
.cat-btn.active { background: var(--lenovo-red); border-color: var(--lenovo-red); color: white; }

.tools-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 12px; }
.tool-card {
  background: var(--bg-card); border: 1px solid var(--border-color);
  border-radius: 10px; padding: 16px; display: flex; flex-direction: column; gap: 10px;
  transition: var(--transition);
}
.tool-card:hover { border-color: var(--border-light); background: var(--bg-card-hover); }
.tool-card.installed { border-color: rgba(76, 175, 80, 0.3); }

.tool-header { display: flex; align-items: center; gap: 12px; }
.tool-icon {
  width: 40px; height: 40px; border-radius: 10px;
  display: flex; align-items: center; justify-content: center;
  background: var(--bg-tertiary); color: var(--text-secondary);
}
.tool-icon.monitor { background: rgba(33, 150, 243, 0.15); color: #2196F3; }
.tool-icon.system { background: rgba(156, 39, 176, 0.15); color: #9C27B0; }
.tool-icon.benchmark { background: rgba(255, 152, 0, 0.15); color: #FF9800; }
.tool-icon.diagnostic { background: rgba(76, 175, 80, 0.15); color: #4CAF50; }

.tool-info { flex: 1; display: flex; flex-direction: column; gap: 2px; }
.tool-name { font-size: 14px; font-weight: 600; color: var(--text-primary); }
.tool-vendor { font-size: 11px; color: var(--text-tertiary); }

.installed-badge {
  display: flex; align-items: center; gap: 4px;
  padding: 3px 8px; background: rgba(76, 175, 80, 0.15);
  border-radius: 10px; color: #4CAF50; font-size: 11px; font-weight: 600;
}

.tool-desc { font-size: 12px; color: var(--text-secondary); line-height: 1.5; }
.tool-meta { display: flex; gap: 12px; font-size: 11px; color: var(--text-tertiary); }
.tool-version { font-family: 'Consolas', monospace; }
.tool-winget { background: rgba(33, 150, 243, 0.15); color: #2196F3; padding: 2px 6px; border-radius: 4px; }

.tool-installing { display: flex; align-items: center; gap: 8px; padding: 8px 10px; background: rgba(255, 152, 0, 0.08); border-radius: 6px; color: #FF9800; font-size: 12px; }

.tool-error, .tool-success { display: flex; align-items: center; gap: 6px; padding: 8px 10px; border-radius: 6px; font-size: 11px; }
.tool-error { background: rgba(244, 67, 54, 0.08); color: #F44336; border: 1px solid rgba(244, 67, 54, 0.2); }
.tool-success { background: rgba(76, 175, 80, 0.08); color: #4CAF50; border: 1px solid rgba(76, 175, 80, 0.2); }

.tool-actions { display: flex; gap: 8px; margin-top: 4px; }
.btn-install, .btn-installing {
  flex: 1; padding: 8px 12px; border: none; border-radius: 6px;
  font-size: 12px; font-weight: 600; cursor: pointer;
  display: flex; align-items: center; justify-content: center; gap: 6px;
  transition: var(--transition); font-family: inherit;
}
.btn-install { background: linear-gradient(135deg, var(--lenovo-red) 0%, var(--lenovo-red-dark) 100%); color: white; }
.btn-install:hover { opacity: 0.9; transform: translateY(-1px); }
.btn-installing { background: var(--bg-tertiary); color: var(--text-secondary); cursor: not-allowed; }
.btn-run {
  flex: 1; padding: 8px 12px; border: none; border-radius: 6px;
  background: linear-gradient(135deg, #4CAF50 0%, #2E7D32 100%);
  color: white; font-size: 12px; font-weight: 600; cursor: pointer;
  display: flex; align-items: center; justify-content: center; gap: 6px;
  transition: var(--transition); font-family: inherit;
}
.btn-run:hover { opacity: 0.9; transform: translateY(-1px); }
.btn-link {
  padding: 8px 12px; background: var(--bg-tertiary); border: 1px solid var(--border-color);
  border-radius: 6px; color: var(--text-secondary); font-size: 12px; font-weight: 500;
  cursor: pointer; display: flex; align-items: center; gap: 6px;
  transition: var(--transition); font-family: inherit; text-decoration: none;
}
.btn-link:hover { background: var(--bg-card-hover); color: var(--text-primary); border-color: var(--border-light); }

.quick-actions { display: flex; gap: 12px; flex-wrap: wrap; }
.btn-batch {
  padding: 10px 20px; background: linear-gradient(135deg, var(--lenovo-red) 0%, var(--lenovo-red-dark) 100%);
  border: none; border-radius: 8px; color: white; font-size: 13px; font-weight: 600;
  cursor: pointer; display: flex; align-items: center; gap: 8px;
  transition: var(--transition); font-family: inherit;
}
.btn-batch:hover:not(:disabled) { opacity: 0.9; transform: translateY(-1px); }
.btn-batch:disabled { opacity: 0.6; cursor: not-allowed; }
</style>
