/**
 * 跨环境剪贴板复制工具
 * - 在 HTTPS / localhost 安全上下文下优先使用 navigator.clipboard.writeText
 * - 在 HTTP 公网 IP 等非安全上下文下自动无缝降级为隐藏 textarea + document.execCommand('copy')
 * - 自动挂载至当前激活的 el-dialog / el-drawer 容器内，规避焦点陷阱 (Focus Trap) 抢夺选区问题
 */
export async function copyTextToClipboard(text: string): Promise<boolean> {
    const content = String(text ?? '')
    if (!content) return false

    // 1. 安全上下文优先使用现代 Clipboard API
    if (typeof navigator !== 'undefined' && navigator.clipboard && window.isSecureContext) {
        try {
            await navigator.clipboard.writeText(content)
            return true
        } catch {
            // 若权限受限或异常，继续走下方同步 DOM 复制兜底
        }
    }

    // 2. HTTP 非安全上下文兜底（无前置 await，保持在点击事件同步调用栈内执行）
    if (typeof document === 'undefined') return false

    let textarea: HTMLTextAreaElement | null = null
    let container: HTMLElement = document.body
    const previousActive = document.activeElement as HTMLElement | null

    try {
        textarea = document.createElement('textarea')
        textarea.value = content
        textarea.setAttribute('readonly', 'readonly')
        textarea.setAttribute('aria-hidden', 'true')
        textarea.style.position = 'fixed'
        textarea.style.top = '0'
        textarea.style.left = '-9999px'
        textarea.style.width = '1px'
        textarea.style.height = '1px'
        textarea.style.padding = '0'
        textarea.style.border = 'none'
        textarea.style.outline = 'none'
        textarea.style.boxShadow = 'none'
        textarea.style.background = 'transparent'
        textarea.style.opacity = '0'
        textarea.style.pointerEvents = 'none'

        const trapHost =
            (previousActive?.closest?.('.el-dialog, .el-drawer, [role="dialog"]') as HTMLElement | null) ||
            (Array.from(
                document.querySelectorAll<HTMLElement>(
                    '.el-overlay:not([style*="display: none"]) .el-dialog, .el-overlay:not([style*="display: none"]) .el-drawer',
                ),
            ).pop() ??
                null)
        if (trapHost) {
            container = trapHost
        }

        container.appendChild(textarea)
        textarea.focus({ preventScroll: true })
        textarea.select()
        textarea.setSelectionRange(0, textarea.value.length)

        const ok = document.execCommand('copy')
        return Boolean(ok)
    } catch (err) {
        console.error('[clipboard] copy failed:', err)
        return false
    } finally {
        if (textarea && textarea.parentNode) {
            textarea.parentNode.removeChild(textarea)
        }
        if (previousActive && typeof previousActive.focus === 'function') {
            try {
                previousActive.focus({ preventScroll: true })
            } catch {
                // ignore
            }
        }
    }
}
