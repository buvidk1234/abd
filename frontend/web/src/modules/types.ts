export type PagedR<T> = {
  data: T[]
  total: number
  page: number
  pageSize: number
}

export type PagedParams = {
  page?: number
  pageSize?: number
}
