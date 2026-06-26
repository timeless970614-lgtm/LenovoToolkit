<template>
  <div class="settings-page">

    <!-- Appearance Card -->
    <div class="card">
      <div class="card-header">
        <span class="card-title">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 8px;">
            <circle cx="12" cy="12" r="5"/>
            <path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/>
          </svg>
          {{ t.appearance }}
        </span>
      </div>

      <div class="settings-list">
        <!-- Theme -->
        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.theme }}</div>
            <div class="setting-desc">{{ t.themeDesc }}</div>
          </div>
          <div class="setting-control">
            <div class="theme-grid">
              <button :class="['theme-btn', theme === 'dark' ? 'active' : '']" @click="setTheme('dark')">
                <span class="theme-preview dark"></span>
                <span class="theme-name">{{ t.dark }}</span>
              </button>
              <button :class="['theme-btn', theme === 'light' ? 'active' : '']" @click="setTheme('light')">
                <span class="theme-preview light"></span>
                <span class="theme-name">{{ t.light }}</span>
              </button>
            </div>
          </div>
        </div>

        <!-- Language -->
        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.language }}</div>
            <div class="setting-desc">{{ t.languageDesc }}</div>
          </div>
          <div class="setting-control">
            <div class="toggle-group">
              <button :class="['toggle-btn', lang === 'en' ? 'active' : '']" @click="setLang('en')">
                <span class="lang-flag">🇺🇸</span>
                English
              </button>
              <button :class="['toggle-btn', lang === 'zh' ? 'active' : '']" @click="setLang('zh')">
                <span class="lang-flag">🇨🇳</span>
                中文
              </button>
            </div>
          </div>
        </div>

        <!-- Animation -->
        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.animations }}</div>
            <div class="setting-desc">{{ t.animationsDesc }}</div>
          </div>
          <div class="setting-control">
            <div :class="['switch', animations ? 'on' : '']" @click="animations = !animations">
              <div class="switch-thumb"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- General Settings Card -->
    <div class="card">
      <div class="card-header">
        <span class="card-title">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 8px;">
            <circle cx="12" cy="12" r="3"/>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>
          </svg>
          {{ t.general }}
        </span>
      </div>

      <div class="settings-list">
        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.refreshInterval }}</div>
            <div class="setting-desc">{{ t.refreshIntervalDesc }}</div>
          </div>
          <div class="setting-control">
            <div class="interval-group">
              <button
                v-for="opt in intervalOptions"
                :key="opt.value"
                :class="['interval-btn', String(pollInterval) === opt.value ? 'active' : '']"
                @click="setRefreshInterval(opt.value)"
              >{{ opt.label }}</button>
            </div>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.autoStart }}</div>
            <div class="setting-desc">{{ t.autoStartDesc }}</div>
          </div>
          <div class="setting-control">
            <div :class="['switch', autoStart ? 'on' : '']" @click="autoStart = !autoStart">
              <div class="switch-thumb"></div>
            </div>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.minimizeToTray }}</div>
            <div class="setting-desc">{{ t.minimizeToTrayDesc }}</div>
          </div>
          <div class="setting-control">
            <div :class="['switch', minimizeToTray ? 'on' : '']" @click="minimizeToTray = !minimizeToTray">
              <div class="switch-thumb"></div>
            </div>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.closeConfirm }}</div>
            <div class="setting-desc">{{ t.closeConfirmDesc }}</div>
          </div>
          <div class="setting-control">
            <div :class="['switch', closeConfirm ? 'on' : '']" @click="closeConfirm = !closeConfirm">
              <div class="switch-thumb"></div>
            </div>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.trayBattery }}</div>
            <div class="setting-desc">{{ t.trayBatteryDesc }}</div>
          </div>
          <div class="setting-control">
            <div :class="['switch', trayBattery ? 'on' : '']" @click="trayBattery = !trayBattery">
              <div class="switch-thumb"></div>
            </div>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.tempUnit }}</div>
            <div class="setting-desc">{{ t.tempUnitDesc }}</div>
          </div>
          <div class="setting-control">
            <div class="toggle-group">
              <button :class="['toggle-btn', tempUnit === 'c' ? 'active' : '']" @click="tempUnit = 'c'">
                °C
              </button>
              <button :class="['toggle-btn', tempUnit === 'f' ? 'active' : '']" @click="tempUnit = 'f'">
                °F
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Performance Card -->
    <div class="card">
      <div class="card-header">
        <span class="card-title">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 8px;">
            <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/>
          </svg>
          {{ t.performance }}
        </span>
      </div>

      <div class="settings-list">
        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.powerPlan }}</div>
            <div class="setting-desc">{{ t.powerPlanDesc }}</div>
          </div>
          <div class="setting-control">
            <select v-model="powerPlan" class="select-control">
              <option value="balanced">{{ t.balanced }}</option>
              <option value="high">{{ t.highPerformance }}</option>
              <option value="saver">{{ t.powerSaver }}</option>
            </select>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.gpuMode }}</div>
            <div class="setting-desc">{{ t.gpuModeDesc }}</div>
          </div>
          <div class="setting-control">
            <div class="toggle-group">
              <button :class="['toggle-btn', gpuMode === 'auto' ? 'active' : '']" @click="gpuMode = 'auto'">
                {{ t.auto }}
              </button>
              <button :class="['toggle-btn', gpuMode === 'integrated' ? 'active' : '']" @click="gpuMode = 'integrated'">
                {{ t.integrated }}
              </button>
              <button :class="['toggle-btn', gpuMode === 'discrete' ? 'active' : '']" @click="gpuMode = 'discrete'">
                {{ t.discrete }}
              </button>
            </div>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.cpuBoost }}</div>
            <div class="setting-desc">{{ t.cpuBoostDesc }}</div>
          </div>
          <div class="setting-control">
            <div :class="['switch', cpuBoost ? 'on' : '']" @click="cpuBoost = !cpuBoost">
              <div class="switch-thumb"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Notifications Card -->
    <div class="card">
      <div class="card-header">
        <span class="card-title">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 8px;">
            <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/>
            <path d="M13.73 21a2 2 0 0 1-3.46 0"/>
          </svg>
          {{ t.notifications }}
        </span>
      </div>

      <div class="settings-list">
        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.notifyModeChange }}</div>
            <div class="setting-desc">{{ t.notifyModeChangeDesc }}</div>
          </div>
          <div class="setting-control">
            <div :class="['switch', notifyModeChange ? 'on' : '']" @click="notifyModeChange = !notifyModeChange">
              <div class="switch-thumb"></div>
            </div>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.notifyBattery }}</div>
            <div class="setting-desc">{{ t.notifyBatteryDesc }}</div>
          </div>
          <div class="setting-control">
            <div :class="['switch', notifyBattery ? 'on' : '']" @click="notifyBattery = !notifyBattery">
              <div class="switch-thumb"></div>
            </div>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.notifyTemperature }}</div>
            <div class="setting-desc">{{ t.notifyTemperatureDesc }}</div>
          </div>
          <div class="setting-control">
            <div :class="['switch', notifyTemperature ? 'on' : '']" @click="notifyTemperature = !notifyTemperature">
              <div class="switch-thumb"></div>
            </div>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.soundNotify }}</div>
            <div class="setting-desc">{{ t.soundNotifyDesc }}</div>
          </div>
          <div class="setting-control">
            <div :class="['switch', soundNotify ? 'on' : '']" @click="soundNotify = !soundNotify">
              <div class="switch-thumb"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Data Management Card -->
    <div class="card">
      <div class="card-header">
        <span class="card-title">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 8px;">
            <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
            <polyline points="14 2 14 8 20 8"/>
            <line x1="16" y1="13" x2="8" y2="13"/>
            <line x1="16" y1="17" x2="8" y2="17"/>
            <polyline points="10 9 9 9 8 9"/>
          </svg>
          {{ t.dataManagement }}
        </span>
      </div>

      <div class="settings-list">
        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.cacheSize }}</div>
            <div class="setting-desc">{{ t.cacheSizeDesc }}</div>
          </div>
          <div class="setting-control">
            <button class="btn btn-secondary btn-sm" @click="clearCache">
              {{ t.clearCache }}
            </button>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.exportLogs }}</div>
            <div class="setting-desc">{{ t.exportLogsDesc }}</div>
          </div>
          <div class="setting-control">
            <button class="btn btn-secondary btn-sm" @click="exportLogs">
              {{ t.export }}
            </button>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.logLevel }}</div>
            <div class="setting-desc">{{ t.logLevelDesc }}</div>
          </div>
          <div class="setting-control">
            <select v-model="logLevel" class="select-control">
              <option value="error">{{ t.error }}</option>
              <option value="warn">{{ t.warn }}</option>
              <option value="info">{{ t.info }}</option>
              <option value="debug">{{ t.debug }}</option>
            </select>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.resetAll }}</div>
            <div class="setting-desc">{{ t.resetAllDesc }}</div>
          </div>
          <div class="setting-control">
            <button class="btn btn-danger btn-sm" @click="resetAll">
              {{ t.reset }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Advanced Card -->
    <div class="card">
      <div class="card-header">
        <span class="card-title">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 8px;">
            <path d="M12 20h9"/>
            <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/>
          </svg>
          {{ t.advanced }}
        </span>
      </div>

      <div class="settings-list">
        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.devMode }}</div>
            <div class="setting-desc">{{ t.devModeDesc }}</div>
          </div>
          <div class="setting-control">
            <div :class="['switch', devMode ? 'on' : '']" @click="devMode = !devMode">
              <div class="switch-thumb"></div>
            </div>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.debugLogging }}</div>
            <div class="setting-desc">{{ t.debugLoggingDesc }}</div>
          </div>
          <div class="setting-control">
            <div :class="['switch', debugLogging ? 'on' : '']" @click="debugLogging = !debugLogging">
              <div class="switch-thumb"></div>
            </div>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.hardwarePolling }}</div>
            <div class="setting-desc">{{ t.hardwarePollingDesc }}</div>
          </div>
          <div class="setting-control">
            <div class="toggle-group">
              <button :class="['toggle-btn', hardwarePolling === 'fast' ? 'active' : '']" @click="hardwarePolling = 'fast'">
                {{ t.fast }}
              </button>
              <button :class="['toggle-btn', hardwarePolling === 'normal' ? 'active' : '']" @click="hardwarePolling = 'normal'">
                {{ t.normal }}
              </button>
              <button :class="['toggle-btn', hardwarePolling === 'slow' ? 'active' : '']" @click="hardwarePolling = 'slow'">
                {{ t.slow }}
              </button>
            </div>
          </div>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.autoUpdate }}</div>
            <div class="setting-desc">{{ t.autoUpdateDesc }}</div>
          </div>
          <div class="setting-control">
            <div :class="['switch', autoUpdate ? 'on' : '']" @click="autoUpdate = !autoUpdate">
              <div class="switch-thumb"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- About Card -->
    <div class="card">
      <div class="card-header">
        <span class="card-title">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 8px;">
            <circle cx="12" cy="12" r="10"/>
            <path d="M12 16v-4M12 8h.01"/>
          </svg>
          {{ t.aboutTitle }}
        </span>
      </div>
      <div class="settings-list">
        <div class="setting-row about-row">
          <div class="about-logo">
            <svg width="36" height="36" viewBox="0 0 24 24" fill="none" stroke="var(--lenovo-red)" stroke-width="1.5">
              <rect x="2" y="3" width="20" height="18" rx="2"/>
              <path d="M8 21h8M12 3v18"/>
            </svg>
          </div>
          <div class="about-info">
            <div class="about-name">Lenovo Toolkit</div>
            <div class="about-version">v1.0.22</div>
            <div class="about-copyright">© 2026 Lenovo. All rights reserved.</div>
          </div>
        </div>
        <div class="setting-row">
          <div class="setting-info">
            <div class="setting-name">{{ t.checkUpdate }}</div>
            <div class="setting-desc">{{ t.checkUpdateDesc }}</div>
          </div>
          <div class="setting-control">
            <button class="btn btn-secondary btn-sm" @click="checkUpdate">
              {{ t.check }}
            </button>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script>
const i18n = {
  en: {
    appearance: 'Appearance',
    theme: 'Theme',
    themeDesc: 'Choose your preferred color scheme',
    dark: 'Dark',
    light: 'Light',
    language: 'Language',
    languageDesc: 'Switch between English and Chinese',
    animations: 'Animations',
    animationsDesc: 'Enable smooth transitions and animations',
    general: 'General Settings',
    refreshInterval: 'Auto Refresh Interval',
    refreshIntervalDesc: 'How often to refresh system data across all pages',
    autoStart: 'Start with Windows',
    autoStartDesc: 'Launch toolkit automatically on startup',
    minimizeToTray: 'Minimize to Tray',
    minimizeToTrayDesc: 'Keep running in system tray when closed',
    closeConfirm: 'Close Confirmation',
    closeConfirmDesc: 'Show confirmation dialog when closing the app',
    trayBattery: 'Tray Battery Info',
    trayBatteryDesc: 'Display battery status in system tray icon',
    tempUnit: 'Temperature Unit',
    tempUnitDesc: 'Choose temperature display unit',
    performance: 'Performance',
    powerPlan: 'Power Plan',
    powerPlanDesc: 'Set system power management mode',
    balanced: 'Balanced',
    highPerformance: 'High Performance',
    powerSaver: 'Power Saver',
    gpuMode: 'GPU Mode',
    gpuModeDesc: 'Control discrete GPU behavior',
    auto: 'Auto',
    integrated: 'Integrated Only',
    discrete: 'Discrete Only',
    cpuBoost: 'CPU Boost',
    cpuBoostDesc: 'Allow CPU to exceed base frequency',
    notifications: 'Notifications',
    notifyModeChange: 'Mode Change Alert',
    notifyModeChangeDesc: 'Notify when performance mode changes',
    notifyBattery: 'Battery Alert',
    notifyBatteryDesc: 'Alert when battery is low',
    notifyTemperature: 'Temperature Alert',
    notifyTemperatureDesc: 'Alert when temperature is high',
    soundNotify: 'Sound Notifications',
    soundNotifyDesc: 'Play sound for notifications',
    dataManagement: 'Data Management',
    cacheSize: 'Cache Size',
    cacheSizeDesc: 'Clear cached data to free up space',
    clearCache: 'Clear Cache',
    exportLogs: 'Export Logs',
    exportLogsDesc: 'Export application logs for troubleshooting',
    export: 'Export',
    logLevel: 'Log Level',
    logLevelDesc: 'Set the verbosity of application logs',
    error: 'Error',
    warn: 'Warning',
    info: 'Info',
    debug: 'Debug',
    resetAll: 'Reset All Settings',
    resetAllDesc: 'Restore all settings to default values',
    reset: 'Reset',
    advanced: 'Advanced',
    devMode: 'Developer Mode',
    devModeDesc: 'Enable developer tools and debug options',
    debugLogging: 'Debug Logging',
    debugLoggingDesc: 'Log detailed debug information',
    hardwarePolling: 'Hardware Polling Rate',
    hardwarePollingDesc: 'Control frequency of hardware sensor reads',
    fast: 'Fast',
    normal: 'Normal',
    slow: 'Slow',
    autoUpdate: 'Auto Update',
    autoUpdateDesc: 'Automatically check for and install updates',
    aboutTitle: 'About',
    checkUpdate: 'Check for Updates',
    checkUpdateDesc: 'Check if a newer version is available',
    check: 'Check',
    aboutCopyright: '© 2026 Lenovo. All rights reserved.',
  },
  zh: {
    appearance: '外观',
    theme: '主题',
    themeDesc: '选择你偏好的配色方案',
    dark: '深色',
    light: '浅色',
    language: '语言',
    languageDesc: '切换中英文界面',
    animations: '动画效果',
    animationsDesc: '启用平滑过渡和动画效果',
    general: '通用设置',
    refreshInterval: '自动刷新间隔',
    refreshIntervalDesc: '所有页面的系统数据刷新频率',
    autoStart: '开机自启',
    autoStartDesc: '系统启动时自动运行工具箱',
    minimizeToTray: '最小化到托盘',
    minimizeToTrayDesc: '关闭窗口时保持在系统托盘运行',
    closeConfirm: '关闭确认',
    closeConfirmDesc: '关闭应用时显示确认对话框',
    trayBattery: '托盘电量显示',
    trayBatteryDesc: '在系统托盘图标中显示电池状态',
    tempUnit: '温度单位',
    tempUnitDesc: '选择温度显示单位',
    performance: '性能',
    powerPlan: '电源计划',
    powerPlanDesc: '设置系统电源管理模式',
    balanced: '平衡',
    highPerformance: '高性能',
    powerSaver: '节能',
    gpuMode: '显卡模式',
    gpuModeDesc: '控制独立显卡行为',
    auto: '自动',
    integrated: '仅集成显卡',
    discrete: '仅独立显卡',
    cpuBoost: 'CPU 加速',
    cpuBoostDesc: '允许 CPU 超过基准频率',
    notifications: '通知',
    notifyModeChange: '模式切换提醒',
    notifyModeChangeDesc: '性能模式变更时发送通知',
    notifyBattery: '电池提醒',
    notifyBatteryDesc: '电池电量低时发送提醒',
    notifyTemperature: '温度提醒',
    notifyTemperatureDesc: '温度过高时发送提醒',
    soundNotify: '声音通知',
    soundNotifyDesc: '通知时播放提示音',
    dataManagement: '数据管理',
    cacheSize: '缓存大小',
    cacheSizeDesc: '清除缓存数据以释放空间',
    clearCache: '清除缓存',
    exportLogs: '导出日志',
    exportLogsDesc: '导出应用日志用于故障排查',
    export: '导出',
    logLevel: '日志级别',
    logLevelDesc: '设置应用日志的详细程度',
    error: '错误',
    warn: '警告',
    info: '信息',
    debug: '调试',
    resetAll: '重置所有设置',
    resetAllDesc: '将所有设置恢复为默认值',
    reset: '重置',
    advanced: '高级',
    devMode: '开发者模式',
    devModeDesc: '启用开发者工具和调试选项',
    debugLogging: '调试日志',
    debugLoggingDesc: '记录详细的调试信息',
    hardwarePolling: '硬件轮询频率',
    hardwarePollingDesc: '控制硬件传感器读取频率',
    fast: '快速',
    normal: '正常',
    slow: '慢速',
    autoUpdate: '自动更新',
    autoUpdateDesc: '自动检查并安装更新',
    aboutTitle: '关于',
    checkUpdate: '检查更新',
    checkUpdateDesc: '检查是否有新版本可用',
    check: '检查',
    aboutCopyright: '© 2026 Lenovo。保留所有权利。',
  }
}

export default {
  name: 'Settings',
  props: {
    theme: {
      type: String,
      default: 'dark'
    },
    lang: {
      type: String,
      default: 'en'
    },
    pollInterval: {
      type: Number,
      default: 5000
    }
  },
  emits: ['update:theme', 'update:lang', 'update:poll-interval'],
  data() {
    return {
      autoStart: false,
      minimizeToTray: true,
      closeConfirm: true,
      trayBattery: false,
      tempUnit: 'c',
      animations: true,
      powerPlan: 'balanced',
      gpuMode: 'auto',
      cpuBoost: true,
      notifyModeChange: true,
      notifyBattery: true,
      notifyTemperature: false,
      soundNotify: true,
      logLevel: 'info',
      devMode: false,
      debugLogging: false,
      hardwarePolling: 'normal',
      autoUpdate: true,
    }
  },
  computed: {
    t() {
      return i18n[this.lang] || i18n.en
    },
    intervalOptions() {
      return [
        { value: '2000', label: this.lang === 'zh' ? '2秒' : '2s' },
        { value: '5000', label: this.lang === 'zh' ? '5秒' : '5s' },
        { value: '10000', label: this.lang === 'zh' ? '10秒' : '10s' },
        { value: '30000', label: this.lang === 'zh' ? '30秒' : '30s' },
      ]
    }
  },
  methods: {
    setTheme(newTheme) {
      this.$emit('update:theme', newTheme)
    },
    setLang(newLang) {
      this.$emit('update:lang', newLang)
    },
    setRefreshInterval(ms) {
      this.$emit('update:poll-interval', Number(ms))
    },
    clearCache() {
      if (confirm(this.lang === 'zh' ? '确定要清除缓存吗？' : 'Are you sure you want to clear the cache?')) {
        localStorage.clear()
        alert(this.lang === 'zh' ? '缓存已清除' : 'Cache cleared')
      }
    },
    exportLogs() {
      alert(this.lang === 'zh' ? '日志导出功能开发中' : 'Log export feature is under development')
    },
    resetAll() {
      const msg = this.lang === 'zh' ? '确定要将所有设置恢复为默认值吗？' : 'Are you sure you want to reset all settings to default?'
      if (confirm(msg)) {
        this.autoStart = false
        this.minimizeToTray = true
        this.closeConfirm = true
        this.trayBattery = false
        this.tempUnit = 'c'
        this.animations = true
        this.powerPlan = 'balanced'
        this.gpuMode = 'auto'
        this.cpuBoost = true
        this.notifyModeChange = true
        this.notifyBattery = true
        this.notifyTemperature = false
        this.soundNotify = true
        this.logLevel = 'info'
        this.devMode = false
        this.debugLogging = false
        this.hardwarePolling = 'normal'
        this.autoUpdate = true
        alert(this.lang === 'zh' ? '所有设置已重置' : 'All settings have been reset')
      }
    },
    checkUpdate() {
      alert(this.lang === 'zh' ? '已是最新版本' : 'You are up to date')
    }
  }
}
</script>

<style scoped>
.settings-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.settings-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-radius: var(--radius-md);
  transition: var(--transition);
}

.setting-row:hover {
  background: var(--bg-tertiary);
}

.setting-info {
  flex: 1;
  min-width: 0;
}

.setting-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 3px;
}

.setting-desc {
  font-size: 12px;
  color: var(--text-tertiary);
  line-height: 1.4;
}

.setting-control {
  flex-shrink: 0;
  margin-left: 24px;
}

.toggle-group {
  display: flex;
  background: var(--bg-tertiary);
  border-radius: var(--radius-md);
  padding: 3px;
  gap: 2px;
}

.toggle-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 7px 14px;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  background: transparent;
  color: var(--text-secondary);
  transition: var(--transition);
  font-family: inherit;
}

.toggle-btn.active {
  background: var(--lenovo-red);
  color: white;
  box-shadow: 0 2px 8px rgba(230, 63, 50, 0.3);
}

.lang-flag {
  font-size: 14px;
  line-height: 1;
}

/* Interval button group */
.interval-group {
  display: flex;
  background: var(--bg-tertiary);
  border-radius: var(--radius-md);
  padding: 3px;
  gap: 2px;
}

.interval-btn {
  padding: 7px 16px;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  background: transparent;
  color: var(--text-secondary);
  transition: var(--transition);
  font-family: inherit;
}

.interval-btn:hover {
  color: var(--text-primary);
}

.interval-btn.active {
  background: var(--lenovo-red);
  color: white;
  box-shadow: 0 2px 8px rgba(230, 63, 50, 0.3);
}

.switch {
  width: 48px;
  height: 26px;
  background: var(--bg-tertiary);
  border: 2px solid var(--border-color);
  border-radius: 13px;
  cursor: pointer;
  transition: var(--transition);
  position: relative;
}

.switch.on {
  background: var(--lenovo-red);
  border-color: var(--lenovo-red);
}

.switch-thumb {
  width: 18px;
  height: 18px;
  background: white;
  border-radius: 50%;
  position: absolute;
  top: 2px;
  left: 2px;
  transition: var(--transition);
  box-shadow: 0 2px 4px rgba(0,0,0,0.3);
}

.switch.on .switch-thumb {
  left: 24px;
}

/* Theme Grid */
.theme-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.theme-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  border: 2px solid var(--border-color);
  border-radius: var(--radius-md);
  background: var(--bg-tertiary);
  cursor: pointer;
  transition: var(--transition);
}

.theme-btn:hover {
  border-color: var(--border-light);
}

.theme-btn.active {
  border-color: var(--lenovo-red);
  background: rgba(230, 63, 50, 0.1);
}

.theme-preview {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-md);
  border: 2px solid var(--border-color);
}

.theme-preview.dark {
  background: linear-gradient(135deg, #0D0D0D 0%, #232323 100%);
}

.theme-preview.light {
  background: linear-gradient(135deg, #F5F5F5 0%, #FFFFFF 100%);
}

.theme-name {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary);
}

.theme-btn.active .theme-name {
  color: var(--lenovo-red);
}

/* Select control */
.select-control {
  padding: 7px 12px;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  background: var(--bg-tertiary);
  color: var(--text-primary);
  font-size: 13px;
  font-family: inherit;
  cursor: pointer;
  outline: none;
  min-width: 140px;
}

.select-control:focus {
  border-color: var(--lenovo-red);
}

/* Buttons */
.btn {
  padding: 7px 16px;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  font-family: inherit;
  transition: var(--transition);
}

.btn-sm {
  padding: 6px 14px;
  font-size: 12px;
}

.btn-secondary {
  background: var(--bg-tertiary);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.btn-secondary:hover {
  background: var(--bg-hover);
  border-color: var(--border-light);
}

.btn-danger {
  background: rgba(248, 113, 113, 0.15);
  color: #f87171;
  border: 1px solid rgba(248, 113, 113, 0.3);
}

.btn-danger:hover {
  background: rgba(248, 113, 113, 0.25);
}

/* About row */
.about-row {
  gap: 16px;
  cursor: default;
}

.about-row:hover {
  background: transparent;
}

.about-logo {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.about-info {
  flex: 1;
}

.about-name {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.about-version {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 2px;
}

.about-copyright {
  font-size: 11px;
  color: var(--text-muted);
  margin-top: 4px;
}
</style>
