import service from '@/utils/request'

export const getxiaoquList = (data) => {
  return service({
    url: '/xiaoqu/list',
    method: 'post',
    data: data
  })
}

export const getDictBuildingInfo = (id) => {
  return service({
    url: '/dictBuilding/findDictBuilding',
    method: 'get',
    params: { ID: id }
  })
}

export const editXiaoqu = (data) => {
  return service({
    url: '/api/xiaoqu/edit',
    method: 'post',
    data: data
  })
}

export const addDictBuilding = (data) => {
  return service({
    url: '/dictBuilding/createDictBuilding',
    method: 'post',
    data: data
  })
}
