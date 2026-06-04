import { Fragment, type ReactNode } from 'react'

export interface DescriptionListItem {
  label: string
  value: ReactNode
}

/** 상세(show) 화면의 라벨-값 정의 목록을 렌더링하는 공용 컴포넌트. */
export function DescriptionList({ items }: { items: DescriptionListItem[] }) {
  return (
    <dl className="grid grid-cols-[8rem_1fr] gap-y-3 text-sm">
      {items.map((item) => (
        <Fragment key={item.label}>
          <dt className="text-muted-foreground">{item.label}</dt>
          <dd>{item.value}</dd>
        </Fragment>
      ))}
    </dl>
  )
}
