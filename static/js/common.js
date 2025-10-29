/**
 * CMDB系统通用JavaScript函数库
 */

// 全局配置
const CMDB_CONFIG = {
    API_BASE_URL: '/api',
    REFRESH_INTERVAL: 60000, // 默认刷新间隔1分钟
    MAX_PAGE_SIZE: 100
};

/**
 * 初始化侧边栏交互
 */
function initSidebar() {
    // 侧边栏切换按钮
    const sidebarToggle = document.getElementById('sidebar-toggle');
    const sidebar = document.querySelector('.sidebar');
    const mainContent = document.querySelector('.main-content');
    
    if (sidebarToggle && sidebar && mainContent) {
        sidebarToggle.addEventListener('click', function() {
            sidebar.classList.toggle('open');
            mainContent.classList.toggle('shrink');
        });
    }
    
    // 下拉菜单切换
    const dropdownToggles = document.querySelectorAll('.nav-link.dropdown-toggle');
    dropdownToggles.forEach(toggle => {
        toggle.addEventListener('click', function(e) {
            e.preventDefault();
            const parent = this.closest('.has-dropdown');
            parent.classList.toggle('open');
        });
    });
}

/**
 * 初始化退出登录功能
 */
function initLogout() {
    const logoutBtn = document.getElementById('logout-btn');
    if (logoutBtn) {
        logoutBtn.addEventListener('click', function(e) {
            e.preventDefault();
            if (confirm('确定要退出登录吗？')) {
                logout();
            }
        });
    }
}

/**
 * 执行退出登录操作
 */
function logout() {
    fetch(`${CMDB_CONFIG.API_BASE_URL}/auth/logout`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => {
        window.location.href = '/login';
    })
    .catch(error => {
        console.error('退出登录失败:', error);
        window.location.href = '/login';
    });
}

/**
 * 显示通知消息
 * @param {string} message - 通知内容
 * @param {string} type - 通知类型 (success, error, warning, info)
 * @param {number} duration - 显示时长(毫秒)
 */
function showNotification(message, type = 'info', duration = 3000) {
    // 创建通知元素
    const notification = document.createElement('div');
    notification.className = `alert alert-${type} notification fade-in`;
    notification.style.position = 'fixed';
    notification.style.top = '20px';
    notification.style.right = '20px';
    notification.style.zIndex = '9999';
    notification.style.padding = '12px 20px';
    notification.style.boxShadow = '0 4px 12px rgba(0, 0, 0, 0.15)';
    notification.style.maxWidth = '400px';
    notification.textContent = message;
    
    // 添加关闭按钮
    const closeBtn = document.createElement('button');
    closeBtn.type = 'button';
    closeBtn.className = 'close';
    closeBtn.style.float = 'right';
    closeBtn.style.marginLeft = '10px';
    closeBtn.style.background = 'none';
    closeBtn.style.border = 'none';
    closeBtn.style.color = 'inherit';
    closeBtn.style.fontSize = '18px';
    closeBtn.style.cursor = 'pointer';
    closeBtn.innerHTML = '&times;';
    closeBtn.addEventListener('click', () => {
        document.body.removeChild(notification);
    });
    
    notification.appendChild(closeBtn);
    document.body.appendChild(notification);
    
    // 自动关闭
    setTimeout(() => {
        notification.style.opacity = '0';
        notification.style.transition = 'opacity 0.3s ease';
        setTimeout(() => {
            if (document.body.contains(notification)) {
                document.body.removeChild(notification);
            }
        }, 300);
    }, duration);
}

/**
 * 显示加载状态
 * @param {HTMLElement} element - 目标元素
 * @param {boolean} show - 是否显示加载状态
 */
function showLoading(element, show = true) {
    if (!element) return;
    
    if (show) {
        // 保存原始内容
        element.dataset.originalContent = element.innerHTML;
        
        // 创建加载元素
        const loadingDiv = document.createElement('div');
        loadingDiv.className = 'loading-container text-center py-4';
        loadingDiv.innerHTML = `
            <div class="loading"></div>
            <p class="mt-2">加载中...</p>
        `;
        
        element.innerHTML = '';
        element.appendChild(loadingDiv);
    } else {
        // 恢复原始内容
        if (element.dataset.originalContent) {
            element.innerHTML = element.dataset.originalContent;
            delete element.dataset.originalContent;
        }
    }
}

/**
 * 表格排序功能
 * @param {HTMLElement} table - 表格元素
 * @param {number} columnIndex - 排序列索引
 * @param {boolean} isNumeric - 是否为数字排序
 */
function sortTable(table, columnIndex, isNumeric = false) {
    const tbody = table.querySelector('tbody');
    const rows = Array.from(tbody.querySelectorAll('tr'));
    const sortAsc = table.dataset.sortAsc === 'true';
    
    // 排序
    rows.sort((a, b) => {
        let aValue = a.cells[columnIndex].textContent.trim();
        let bValue = b.cells[columnIndex].textContent.trim();
        
        if (isNumeric) {
            aValue = parseFloat(aValue) || 0;
            bValue = parseFloat(bValue) || 0;
        }
        
        if (sortAsc) {
            return aValue > bValue ? 1 : -1;
        } else {
            return aValue < bValue ? 1 : -1;
        }
    });
    
    // 重新排列行
    rows.forEach(row => tbody.appendChild(row));
    
    // 更新排序状态
    table.dataset.sortAsc = (!sortAsc).toString();
    
    // 更新排序图标
    const headers = table.querySelectorAll('thead th');
    headers.forEach((header, index) => {
        header.innerHTML = header.innerHTML.replace(/ <i\s+class="bi\s+bi-sort-.+?<\/i>/g, '');
        if (index === columnIndex) {
            const icon = sortAsc ? 'bi-sort-down' : 'bi-sort-up';
            header.innerHTML += ` <i class="bi ${icon}"></i>`;
        }
    });
}

/**
 * 格式化日期时间
 * @param {string|Date} date - 日期对象或字符串
 * @param {string} format - 格式化模板
 * @returns {string} 格式化后的日期字符串
 */
function formatDateTime(date, format = 'YYYY-MM-DD HH:mm:ss') {
    if (!date) return '';
    
    const d = typeof date === 'string' ? new Date(date) : date;
    
    if (isNaN(d.getTime())) return '';
    
    const year = d.getFullYear();
    const month = String(d.getMonth() + 1).padStart(2, '0');
    const day = String(d.getDate()).padStart(2, '0');
    const hours = String(d.getHours()).padStart(2, '0');
    const minutes = String(d.getMinutes()).padStart(2, '0');
    const seconds = String(d.getSeconds()).padStart(2, '0');
    
    return format
        .replace('YYYY', year)
        .replace('MM', month)
        .replace('DD', day)
        .replace('HH', hours)
        .replace('mm', minutes)
        .replace('ss', seconds);
}

/**
 * 验证IP地址格式
 * @param {string} ip - IP地址字符串
 * @returns {boolean} 是否为有效IP地址
 */
function isValidIP(ip) {
    if (!ip) return false;
    
    const ipv4Regex = /^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
    return ipv4Regex.test(ip);
}

/**
 * 深拷贝对象
 * @param {*} obj - 要拷贝的对象
 * @returns {*} 拷贝后的对象
 */
function deepClone(obj) {
    if (obj === null || typeof obj !== 'object') return obj;
    if (obj instanceof Date) return new Date(obj.getTime());
    if (obj instanceof Array) return obj.map(item => deepClone(item));
    
    const clonedObj = {};
    for (const key in obj) {
        if (obj.hasOwnProperty(key)) {
            clonedObj[key] = deepClone(obj[key]);
        }
    }
    return clonedObj;
}

/**
 * 获取URL参数
 * @param {string} name - 参数名
 * @returns {string|null} 参数值
 */
function getURLParam(name) {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get(name);
}

/**
 * 设置URL参数
 * @param {string} name - 参数名
 * @param {string} value - 参数值
 */
function setURLParam(name, value) {
    const urlParams = new URLSearchParams(window.location.search);
    urlParams.set(name, value);
    window.history.replaceState({}, '', `${window.location.pathname}?${urlParams.toString()}`);
}

/**
 * 移除URL参数
 * @param {string} name - 参数名
 */
function removeURLParam(name) {
    const urlParams = new URLSearchParams(window.location.search);
    urlParams.delete(name);
    window.history.replaceState({}, '', `${window.location.pathname}?${urlParams.toString()}`);
}

/**
 * 自动调整文本框高度
 * @param {HTMLElement} textarea - 文本框元素
 */
function autoResizeTextarea(textarea) {
    textarea.style.height = 'auto';
    textarea.style.height = (textarea.scrollHeight) + 'px';
}

/**
 * 防抖函数
 * @param {Function} func - 要防抖的函数
 * @param {number} wait - 等待时间(毫秒)
 * @returns {Function} 防抖后的函数
 */
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

/**
 * 节流函数
 * @param {Function} func - 要节流的函数
 * @param {number} limit - 时间限制(毫秒)
 * @returns {Function} 节流后的函数
 */
function throttle(func, limit) {
    let inThrottle;
    return function(...args) {
        if (!inThrottle) {
            func.apply(this, args);
            inThrottle = true;
            setTimeout(() => inThrottle = false, limit);
        }
    };
}

/**
 * 检查用户权限
 * @param {string} permission - 权限名称
 * @returns {Promise<boolean>} 是否有权限
 */
async function checkPermission(permission) {
    try {
        const response = await fetch(`${CMDB_CONFIG.API_BASE_URL}/permissions/check`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ permission })
        });
        const data = await response.json();
        return data.hasPermission || false;
    } catch (error) {
        console.error('检查权限失败:', error);
        return false;
    }
}

/**
 * 初始化数据表格（通用）
 * @param {string} tableId - 表格ID
 * @param {Object} options - 配置选项
 */
function initDataTable(tableId, options = {}) {
    const table = document.getElementById(tableId);
    if (!table) return;
    
    // 配置默认值
    const defaultOptions = {
        sortable: true,
        searchable: true,
        pagination: true
    };
    
    const config = { ...defaultOptions, ...options };
    
    // 初始化排序
    if (config.sortable) {
        const headers = table.querySelectorAll('thead th[data-sortable="true"]');
        headers.forEach((header, index) => {
            header.style.cursor = 'pointer';
            header.addEventListener('click', () => {
                const isNumeric = header.dataset.numeric === 'true';
                sortTable(table, index, isNumeric);
            });
        });
    }
    
    // 更多功能可根据需要扩展
}

/**
 * 初始化所有全局功能
 */
function initGlobalFunctions() {
    // 初始化侧边栏
    initSidebar();
    
    // 初始化退出登录
    initLogout();
    
    // 初始化自动调整文本框
    const textareas = document.querySelectorAll('textarea.auto-resize');
    textareas.forEach(textarea => {
        textarea.addEventListener('input', () => autoResizeTextarea(textarea));
        autoResizeTextarea(textarea); // 初始调整
    });
    
    // 禁止右键菜单（可选）
    if (window.location.pathname !== '/login') {
        document.addEventListener('contextmenu', e => {
            if (!e.target.closest('textarea') && !e.target.closest('input')) {
                e.preventDefault();
            }
        });
    }
    
    // 监听网络状态
    window.addEventListener('online', () => {
        showNotification('网络连接已恢复', 'success');
    });
    
    window.addEventListener('offline', () => {
        showNotification('网络连接已断开', 'error');
    });
}

// 页面加载完成后初始化
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initGlobalFunctions);
} else {
    initGlobalFunctions();
}

// 暴露全局函数
window.CMDB = {
    showNotification,
    showLoading,
    sortTable,
    formatDateTime,
    isValidIP,
    deepClone,
    getURLParam,
    setURLParam,
    removeURLParam,
    autoResizeTextarea,
    debounce,
    throttle,
    checkPermission,
    initDataTable,
    logout
};