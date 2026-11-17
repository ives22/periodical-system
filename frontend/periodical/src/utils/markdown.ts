/**
 * 轻量级零依赖 GFM (GitHub Flavored Markdown) 渲染引擎
 * 支持：GFM 表格 (| a | b |)、围栏代码块 (```)、标题 (#~####)、引用块 (>)、有序/无序列表、行内代码 (`)、加粗 (**)、分割线 (---)
 */

function escapeHtml(str: string): string {
    return str
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
}

function parseInline(text: string): string {
    let out = escapeHtml(text)
    // 行内代码 `code`
    out = out.replace(/`([^`\n]+)`/g, '<code class="md-inline-code">$1</code>')
    // 加粗 **bold**
    out = out.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
    // 斜体 *italic*
    out = out.replace(/(^|[^*])\*([^*\n]+)\*([^*]|$)/g, '$1<em>$2</em>$3')
    // 书名号高亮《期刊名》
    out = out.replace(/《([^》\n]+)》/g, '<span class="md-book-tag">《$1》</span>')
    return out
}

function isTableSeparatorLine(line: string): boolean {
    const trimmed = line.trim()
    if (!trimmed.includes('|')) return false
    const cells = trimmed
        .replace(/^\|/, '')
        .replace(/\|$/, '')
        .split('|')
        .map((c) => c.trim())
    if (cells.length === 0) return false
    return cells.every((c) => /^:?-{2,}:?$/.test(c))
}

function splitTableCells(line: string): string[] {
    return line
        .trim()
        .replace(/^\|/, '')
        .replace(/\|$/, '')
        .split('|')
        .map((c) => c.trim())
}

export function renderGFMMarkdown(raw: string): string {
    if (!raw || !raw.trim()) return ''

    // 1. 提取围栏代码块 ```lang ... ``` 占位保护
    const codeBlocks: string[] = []
    const textWithoutCode = raw.replace(/```([a-zA-Z0-9_-]*)\n([\s\S]*?)```/g, (_match, lang, code) => {
        const idx = codeBlocks.length
        const langLabel = (lang || 'text').toUpperCase()
        const html = `<div class="md-code-block"><div class="md-code-header"><span>${escapeHtml(
            langLabel,
        )}</span></div><pre><code>${escapeHtml(code.replace(/\n$/, ''))}</code></pre></div>`
        codeBlocks.push(html)
        return `@@MD_CODE_BLOCK_${idx}@@`
    })

    const lines = textWithoutCode.replace(/\r\n/g, '\n').split('\n')
    const htmlBlocks: string[] = []
    let i = 0

    while (i < lines.length) {
        const line = lines[i]
        const trimmed = line.trim()

        // 空行跳过
        if (!trimmed) {
            i++
            continue
        }

        // 代码块占位符直接还原
        const codeMatch = trimmed.match(/^@@MD_CODE_BLOCK_(\d+)@@$/)
        if (codeMatch) {
            htmlBlocks.push(codeBlocks[Number(codeMatch[1])] || '')
            i++
            continue
        }

        // 2. 识别 GFM 表格：当前行含 | 且下一行是 |---|---| 分隔线
        if (trimmed.includes('|') && i + 1 < lines.length && isTableSeparatorLine(lines[i + 1])) {
            const headers = splitTableCells(trimmed)
            const aligns = splitTableCells(lines[i + 1]).map((sep) => {
                if (sep.startsWith(':') && sep.endsWith(':')) return 'center'
                if (sep.endsWith(':')) return 'right'
                return 'left'
            })
            i += 2

            const rows: string[][] = []
            while (i < lines.length && lines[i].trim().includes('|') && lines[i].trim() !== '') {
                rows.push(splitTableCells(lines[i]))
                i++
            }

            let tableHtml = '<div class="md-table-wrap"><table class="md-table"><thead><tr>'
            headers.forEach((h, colIdx) => {
                const align = aligns[colIdx] || 'left'
                tableHtml += `<th style="text-align:${align}">${parseInline(h)}</th>`
            })
            tableHtml += '</tr></thead><tbody>'

            rows.forEach((row) => {
                tableHtml += '<tr>'
                headers.forEach((_h, colIdx) => {
                    const align = aligns[colIdx] || 'left'
                    const cellVal = row[colIdx] !== undefined ? row[colIdx] : ''
                    tableHtml += `<td style="text-align:${align}">${parseInline(cellVal)}</td>`
                })
                tableHtml += '</tr>'
            })
            tableHtml += '</tbody></table></div>'
            htmlBlocks.push(tableHtml)
            continue
        }

        // 3. 水平分割线 --- 或 ***
        if (/^(-{3,}|\*{3,})$/.test(trimmed)) {
            htmlBlocks.push('<hr class="md-hr" />')
            i++
            continue
        }

        // 4. 分级标题 # ~ ####
        const headingMatch = trimmed.match(/^(#{1,4})\s+(.+)$/)
        if (headingMatch) {
            const level = headingMatch[1].length
            htmlBlocks.push(`<h${level} class="md-h${level}">${parseInline(headingMatch[2])}</h${level}>`)
            i++
            continue
        }

        // 5. 引用块 >
        if (trimmed.startsWith('>')) {
            const quoteLines: string[] = []
            while (i < lines.length && lines[i].trim().startsWith('>')) {
                quoteLines.push(lines[i].trim().replace(/^>\s?/, ''))
                i++
            }
            htmlBlocks.push(
                `<blockquote class="md-quote">${quoteLines.map((q) => `<p>${parseInline(q)}</p>`).join('')}</blockquote>`,
            )
            continue
        }

        // 6. 无序列表 - / * / •
        if (/^[-*•]\s+/.test(trimmed)) {
            const items: string[] = []
            while (i < lines.length && /^[-*•]\s+/.test(lines[i].trim())) {
                items.push(lines[i].trim().replace(/^[-*•]\s+/, ''))
                i++
            }
            htmlBlocks.push(
                `<ul class="md-ul">${items.map((it) => `<li>${parseInline(it)}</li>`).join('')}</ul>`,
            )
            continue
        }

        // 7. 有序列表 1. / 2.
        if (/^\d+\.\s+/.test(trimmed)) {
            const items: string[] = []
            while (i < lines.length && /^\d+\.\s+/.test(lines[i].trim())) {
                items.push(lines[i].trim().replace(/^\d+\.\s+/, ''))
                i++
            }
            htmlBlocks.push(
                `<ol class="md-ol">${items.map((it) => `<li>${parseInline(it)}</li>`).join('')}</ol>`,
            )
            continue
        }

        // 8. 普通段落（合并连续非空普通行）
        const paraLines: string[] = []
        while (
            i < lines.length &&
            lines[i].trim() !== '' &&
            !lines[i].trim().startsWith('#') &&
            !lines[i].trim().startsWith('>') &&
            !/^[-*•]\s+/.test(lines[i].trim()) &&
            !/^\d+\.\s+/.test(lines[i].trim()) &&
            !/^@@MD_CODE_BLOCK_\d+@@$/.test(lines[i].trim()) &&
            !(lines[i].trim().includes('|') && i + 1 < lines.length && isTableSeparatorLine(lines[i + 1]))
        ) {
            paraLines.push(parseInline(lines[i]))
            i++
        }
        if (paraLines.length > 0) {
            htmlBlocks.push(`<p class="md-p">${paraLines.join('<br/>')}</p>`)
        } else {
            i++
        }
    }

    return `<div class="md-body">${htmlBlocks.join('')}</div>`
}
