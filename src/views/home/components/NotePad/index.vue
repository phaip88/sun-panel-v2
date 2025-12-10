<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { SvgIcon } from '@/components/common'
import { useMessage } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { useDraggable, useDebounceFn } from '@vueuse/core'
import { getNotepadContent, saveNotepadContent, uploadNotepadFile } from '@/api/panel/notepad'

defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'update:visible', visible: boolean): void
}>()

const { t } = useI18n()
const message = useMessage()
const editorRef = ref<HTMLDivElement | null>(null)
const notepadRef = ref<HTMLElement | null>(null)
const headerRef = ref<HTMLElement | null>(null)

// 窗口初始位置
const { x, y } = useDraggable(notepadRef, {
  initialValue: { x: window.innerWidth - 370, y: 80 },
  handle: headerRef
})

// 保存内容（防抖）
const saveContent = useDebounceFn(async () => {
    if (editorRef.value) {
        try {
            await saveNotepadContent({ content: editorRef.value.innerHTML })
        } catch (error) {
            console.error('Save notepad error:', error)
        }
    }
}, 1000)

// 插入文件/图片
const insertFileLink = (fileInfo: { name: string, type: string, url: string }) => {
    if (!editorRef.value) return
    
    let htmlFragment = ''
    let fullUrl = fileInfo.url
    
    if (fileInfo.type.startsWith('image/')) {
        htmlFragment = `<div><img class="note-image" src="${fullUrl}" alt="${fileInfo.name}" /></div>`
    } else {
        htmlFragment = `&nbsp;<a href="${fullUrl}" target="_blank" class="file-attachment" contenteditable="false" title="${t('notepad.clickToDownload')}">📁&nbsp;${fileInfo.name}</a>&nbsp;`
    }
    
    editorRef.value.insertAdjacentHTML('beforeend', htmlFragment)
    saveContent() // 触发保存
    
    nextTick(() => {
        if (editorRef.value) {
           editorRef.value.scrollTop = editorRef.value.scrollHeight
        }
    })
}

// 处理拖拽上传
const handleDrop = async (e: DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    
    if (e.dataTransfer && e.dataTransfer.files) {
        const files = Array.from(e.dataTransfer.files)
        for (const file of files) {
            const formData = new FormData()
            formData.append('file', file)
            try {
                const res = await uploadNotepadFile(formData)
                if (res.code === 0) {
                    const data = res.data
                    insertFileLink({
                        name: data.name || file.name,
                        type: data.type || file.type,
                        url: data.url
                    })
                } else {
                     message.error(t('notepad.saveFailed'))
                }
            } catch (e) {
                message.error(t('notepad.saveFailed'))
            }
        }
    }
}

const handleDragOver = (e: DragEvent) => {
    e.preventDefault()
}

// 处理点击
const handleClick = (e: MouseEvent) => {
   // 默认处理
}

// 初始化
onMounted(async () => {
    try {
        const res = await getNotepadContent()
        if (res.code === 0 && editorRef.value) {
            editorRef.value.innerHTML = res.data.content || ''
        }
    } catch (e) {
        console.error('Failed to load notepad', e)
    }
})

// 处理关闭
const close = () => {
    emit('update:visible', false)
}
</script>

<template>
  <!-- 遮罩层，点击关闭 -->
  <div v-show="visible" class="fixed inset-0 z-[100] bg-transparent" @click="close"></div>

  <!-- 便签主体 -->
  <transition name="note-fade">
    <div
        v-show="visible"
        ref="notepadRef"
        class="fixed z-[101] w-[350px] h-[35vh] flex flex-col shadow-xl rounded-lg overflow-hidden border border-amber-200"
        :style="{ left: x + 'px', top: y + 'px' }"
        @click.stop
        @dragover="handleDragOver"
        @drop="handleDrop"
    >
      <!-- 头部 -->
      <div ref="headerRef" class="h-8 bg-[#fef3c7] flex justify-between items-center px-2 border-b border-[#feebc8] select-none cursor-move">
         <div class="flex items-center text-amber-800 text-sm font-bold">
            <SvgIcon icon="note" class="mr-1" />
            <span>{{ t('notepad.title') }}</span>
         </div>
         <div class="hover:bg-amber-200 rounded p-0.5 cursor-pointer text-amber-900" @click="close">
            <SvgIcon icon="mdi:close" />
         </div>
      </div>

      <!-- 编辑区 -->
      <div class="flex-1 bg-[#fffbeb] relative overflow-hidden">
         <!-- ContentEditable Div -->
         <div
            ref="editorRef"
            contenteditable="true"
            class="w-full h-full p-3 outline-none overflow-y-auto text-sm text-gray-800 break-words font-sans leading-relaxed"
            :data-placeholder="t('notepad.placeholder')"
            @input="saveContent"
            @click="handleClick"
            spellcheck="false"
         ></div>
      </div>
    </div>
  </transition>
</template>

<style scoped>
/* 模拟便签纸质感 */
.note-fade-enter-active,
.note-fade-leave-active {
  transition: all 0.2s ease;
}

.note-fade-enter-from,
.note-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px) scale(0.95);
}

/* 编辑器样式 */
:deep(.file-attachment) {
    display: inline-flex;
    align-items: center;
    background-color: #fff7ed;
    border: 1px solid #fed7aa;
    border-radius: 4px;
    padding: 0 4px;
    margin: 0 2px;
    font-size: 0.85em;
    color: #c2410c; /* amber-700 */
    cursor: pointer;
    user-select: none;
    transition: all 0.2s;
    text-decoration: none; /* 移除下划线 */
}

:deep(.file-attachment:hover) {
    background-color: #ffedd5;
    border-color: #fdba74;
}

:deep(.note-image) {
    max-width: 100%;
    max-height: 150px;
    border-radius: 4px;
    margin: 4px 0;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    display: block; /* 独占一行 */
    cursor: default;
}

/* 占位符效果 */
div[contenteditable]:empty::before {
  content: attr(data-placeholder);
  color: #9ca3af;
  pointer-events: none;
  font-style: italic;
}
</style>
