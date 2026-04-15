import service from '@/utils/request'

export const grantContactQuota = (data) => {
  return service({
    url: '/house/contactQuota/grant',
    method: 'post',
    data
  })
}

export const getContactQuotaList = (data) => {
  return service({
    url: '/house/contactQuota/list',
    method: 'post',
    data
  })
}
