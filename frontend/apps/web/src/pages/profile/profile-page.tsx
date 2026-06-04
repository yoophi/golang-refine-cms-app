import { type FormEvent, useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'

import { updateMe, useAuthStore, type UpdateProfileInput } from '@/entities/user'
import {
  Button,
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  Input,
  Label,
} from '@/shared/ui'

export function ProfilePage() {
  const navigate = useNavigate()
  const user = useAuthStore((s) => s.user)
  const setUser = useAuthStore((s) => s.setUser)
  const clear = useAuthStore((s) => s.clear)

  const [name, setName] = useState(user?.name ?? '')
  const [email, setEmail] = useState(user?.email ?? '')
  const [password, setPassword] = useState('')
  const [currentPassword, setCurrentPassword] = useState('')

  const mutation = useMutation({
    mutationFn: updateMe,
    onSuccess: (data) => {
      setUser(data)
      setPassword('')
      setCurrentPassword('')
    },
  })

  const onSubmit = (event: FormEvent) => {
    event.preventDefault()
    const input: UpdateProfileInput = { name, email }
    if (password) {
      input.password = password
      input.currentPassword = currentPassword
    }
    mutation.mutate(input)
  }

  const onLogout = () => {
    clear()
    navigate('/')
  }

  return (
    <div className="mx-auto max-w-md space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold tracking-tight">내 정보</h1>
        <Button variant="outline" size="sm" onClick={onLogout}>
          로그아웃
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>프로필 수정</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSubmit} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="name">이름</Label>
              <Input
                id="name"
                value={name}
                onChange={(event) => setName(event.target.value)}
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="email">이메일</Label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(event) => setEmail(event.target.value)}
                required
              />
            </div>

            <div className="border-t pt-4">
              <p className="text-muted-foreground mb-3 text-sm">
                비밀번호를 변경하려면 아래를 입력하세요.
              </p>
              <div className="space-y-2">
                <Label htmlFor="currentPassword">현재 비밀번호</Label>
                <Input
                  id="currentPassword"
                  type="password"
                  autoComplete="current-password"
                  value={currentPassword}
                  onChange={(event) => setCurrentPassword(event.target.value)}
                />
              </div>
              <div className="mt-3 space-y-2">
                <Label htmlFor="newPassword">새 비밀번호</Label>
                <Input
                  id="newPassword"
                  type="password"
                  autoComplete="new-password"
                  value={password}
                  onChange={(event) => setPassword(event.target.value)}
                />
              </div>
            </div>

            <div className="flex items-center gap-3">
              <Button type="submit" disabled={mutation.isPending}>
                저장
              </Button>
              {mutation.isError ? (
                <span className="text-destructive text-sm">
                  저장에 실패했습니다.
                </span>
              ) : null}
              {mutation.isSuccess ? (
                <span className="text-muted-foreground text-sm">
                  저장되었습니다.
                </span>
              ) : null}
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
