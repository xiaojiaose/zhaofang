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

export const addOrEditDictBuilding = (data) => {
  return service({
    url: '/xiaoqu/edit',
    method: 'post',
    data: data
  })
}

export const getXiaoquDictTree = (xiaoquId) => {
  return service({
    url: '/xiaoqu/dict/tree',
    method: 'get',
    params: { id: xiaoquId }
  })
}

export const upsertXiaoquDict = (data) => {
  return service({
    url: '/xiaoqu/dict/upsert',
    method: 'post',
    data: data
  })
}

export const deleteXiaoquDict = (data) => {
  return service({
    url: '/xiaoqu/dict/delete',
    method: 'post',
    data: data
  })
}
