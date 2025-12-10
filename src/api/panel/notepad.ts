import { get, post } from '@/utils/request'

export function getNotepadContent() {
    return get<any>({
        url: '/panel/notepad/get',
    })
}

export function saveNotepadContent(data: { content: string }) {
    return post<any>({
        url: '/panel/notepad/save',
        data,
    })
}

export function uploadNotepadFile(data: FormData) {
    return post<any>({
        url: '/panel/notepad/upload',
        data,
        headers: {
            'Content-Type': 'multipart/form-data',
        },
    })
}
