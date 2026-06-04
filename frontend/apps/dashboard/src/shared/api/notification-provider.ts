import type { NotificationProvider } from '@refinedev/core'
import { toast } from 'sonner'

/**
 * sonner 기반 refine 알림 프로바이더 (headless 구성이므로 직접 제공).
 */
export const notificationProvider: NotificationProvider = {
  open: ({ key, message, description, type }) => {
    if (type === 'success') {
      toast.success(message, { id: key, description })
      return
    }
    if (type === 'error') {
      toast.error(message, { id: key, description })
      return
    }
    toast(message, { id: key, description })
  },
  close: (key) => {
    toast.dismiss(key)
  },
}
