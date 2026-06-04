import { create } from 'zustand'
import { persist } from 'zustand/middleware'

/**
 * 비로그인 독자의 댓글 작성자 정보를 브라우저에 영구 저장한다.
 * 게시글마다 이름/이메일을 다시 입력하지 않도록 댓글 폼을 자동 채운다.
 */
interface ReaderState {
  authorName: string
  authorEmail: string
  setAuthor: (info: { authorName: string; authorEmail: string }) => void
}

export const useReaderStore = create<ReaderState>()(
  persist(
    (set) => ({
      authorName: '',
      authorEmail: '',
      setAuthor: (info) => set(info),
    }),
    { name: 'cms.web.reader' },
  ),
)
