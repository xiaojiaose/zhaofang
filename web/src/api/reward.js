import service from '@/utils/request'

export const getRewardList = (data) => {
  return service({
    url: '/house/reward/list',
    method: 'post',
    data
  })
}

export const rewardAction = (data) => {
  return service({
    url: '/house/reward/action',
    method: 'post',
    data
  })
}

export const rewardBatchAction = async (ids, action) => {
  const results = await Promise.all(
    ids.map((id) =>
      rewardAction({
        id,
        action
      })
    )
  )
  return results
}
