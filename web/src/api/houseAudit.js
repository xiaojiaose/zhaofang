import service from '@/utils/request'

export const getHouseAuditList = (data) => {
  return service({
    url: '/house/landlordContactView/list',
    method: 'post',
    data: data
  })
}

export const houseAuditAction = (data) => {
  return service({
    url: '/api/houseAudit/action',
    method: 'post',
    data: data
  })
}

export const houseAuditBatchAction = (data) => {
  return service({
    url: '/house/landlordContactView/action',
    method: 'post',
    data: data
  })
}
