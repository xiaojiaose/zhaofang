import service from '@/utils/request'

export const getStatis = (params) => {
  return service({
    url: '/house/statis/view',
    method: 'get',
    params
  })
}

export const getVisitList = (data) => {
  return service({
    url: '/house/statis/visit',
    method: 'post',
    data
  })
}
